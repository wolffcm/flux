package functions

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/influxdata/flux"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/plan"
	"github.com/influxdata/flux/semantic"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const FromSQLKind = "fromSQL"

type FromSQLOpSpec struct {
	DriverName     string `json:"driverName,omitempty"`
	DataSourceName string `json:"dataSourceName,omitempty"`
	Query          string `json:"query,omitempty"`
}

var fromSQLSignature = semantic.FunctionSignature{
	Params: map[string]semantic.Type{
		"driverName":     semantic.String,
		"dataSourceName": semantic.String,
		"query":          semantic.String,
	},
	ReturnType: flux.TableObjectType,
}

func init() {
	flux.RegisterFunction(FromSQLKind, createFromSQLOpSpec, fromSQLSignature)
	flux.RegisterOpSpec(FromSQLKind, newFromSQLOp)
	plan.RegisterProcedureSpec(FromSQLKind, newFromSQLProcedure, FromSQLKind)
	execute.RegisterSource(FromSQLKind, createFromSQLSource)
}

func createFromSQLOpSpec(args flux.Arguments, a *flux.Administration) (flux.OperationSpec, error) {
	spec := new(FromSQLOpSpec)

	if driverName, ok, err := args.GetString("driverName"); err != nil {
		return nil, err
	} else if ok {
		spec.DriverName = driverName
	}

	if dataSourceName, ok, err := args.GetString("dataSourceName"); err != nil {
		return nil, err
	} else if ok {
		spec.DataSourceName = dataSourceName
	}

	if query, ok, err := args.GetString("query"); err != nil {
		return nil, err
	} else if ok {
		spec.Query = query
	}

	if spec.DriverName == "" && spec.DataSourceName == "" && spec.Query == "" {
		return nil, errors.New("must specify driverName, dataSourceName, and query")
	}
	return spec, nil
}

func newFromSQLOp() flux.OperationSpec {
	return new(FromSQLOpSpec)
}

func (s *FromSQLOpSpec) Kind() flux.OperationKind {
	return FromSQLKind
}

type FromSQLProcedureSpec struct {
	DriverName     string
	DataSourceName string
	Query          string
}

func newFromSQLProcedure(qs flux.OperationSpec, pa plan.Administration) (plan.ProcedureSpec, error) {
	spec, ok := qs.(*FromSQLOpSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", qs)
	}

	return &FromSQLProcedureSpec{
		DriverName:     spec.DriverName,
		DataSourceName: spec.DataSourceName,
		Query:          spec.Query,
	}, nil
}

func (s *FromSQLProcedureSpec) Kind() plan.ProcedureKind {
	return FromSQLKind
}

func (s *FromSQLProcedureSpec) Copy() plan.ProcedureSpec {
	ns := new(FromSQLProcedureSpec)
	ns.DriverName = s.DriverName
	ns.DataSourceName = s.DataSourceName
	ns.Query = s.Query
	return ns
}

func createFromSQLSource(prSpec plan.ProcedureSpec, dsid execute.DatasetID, administration execute.Administration) (execute.Source, error) {
	spec, ok := prSpec.(*FromSQLProcedureSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", prSpec)
	}

	if spec.DriverName != "postgres" {
		return nil, fmt.Errorf("sql driver %s not supported", spec.DriverName)
	}

	SQLIterator := SQLIterator{id: dsid, spec: spec, administration: administration}

	return &SQLIterator, nil
}

type SQLIterator struct {
	id             execute.DatasetID
	data           flux.Result
	ts             []execute.Transformation
	administration *flux.Administration
	spec           *FromSQLProcedureSpec
	db             *sql.DB
}

func (c *SQLIterator) Connect() error {
	db, err := sql.Open(c.spec.DriverName, c.spec.DataSourceName)
	if err != nil {
		return err
	}
	c.db = db

	return nil
}

// if fetch gets table, return true; if only fetching one time, then return false
func (c *SQLIterator) Fetch() (bool, error) {
	rows, err := c.db.Query(c.spec.Query)
	if err != nil {
		return false, err
	}

	spew.Dump(rows)

	return false, nil
}

func (c *SQLIterator) Decode() (flux.Table, error) {
	// TODO implement sql package with decoder
	// decoder := sql.NewResultDecoder(sql.ResultDecoderConfig{})
	// result, err := decoder.Decode(strings.NewReader(rows))
	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
}

func (c *SQLIterator) Do(f func(flux.Table) error) error {
	c.Connect()

	more, err := c.Fetch()
	if err != nil {
		return err
	}
	for more {
		tbl, err := c.Decode()
		if err != nil {
			return err
		}
		f(tbl)
		more, err = c.Fetch()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *SQLIterator) AddTransformation(t execute.Transformation) {
	c.ts = append(c.ts, t)
}

func (c *SQLIterator) Run(ctx context.Context) {
	var err error
	var max execute.Time
	maxSet := false
	err = c.Do(func(tbl flux.Table) error {
		for _, t := range c.ts {
			err := t.Process(c.id, tbl)
			if err != nil {
				return err
			}
			if idx := execute.ColIdx(execute.DefaultStopColLabel, tbl.Key().Cols()); idx >= 0 {
				if stop := tbl.Key().ValueTime(idx); !maxSet || stop > max {
					max = stop
					maxSet = true
				}
			}
		}
		return nil
	})
	if err != nil {
		goto FINISH
	}

	if maxSet {
		for _, t := range c.ts {
			t.UpdateWatermark(c.id, max)
		}
	}

FINISH:
	for _, t := range c.ts {
		t.Finish(c.id, err)
	}
}

package functions

import (
	"context"
	"fmt"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/ebpf"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/plan"
	"github.com/influxdata/flux/semantic"
)

const FromEbpfKind = "fromEbpf"

type FromEbpfOpSpec struct {
	File string
}

var fromEbpfSignature = semantic.FunctionSignature{
	Params: map[string]semantic.Type{
		"file": semantic.String,
	},
	ReturnType: flux.TableObjectType,
}

func init() {
	flux.RegisterFunction(FromEbpfKind, createFromEbpfOpSpec, fromEbpfSignature)
	flux.RegisterOpSpec(FromEbpfKind, newFromEbpfOp)
	plan.RegisterProcedureSpec(FromEbpfKind, newFromEbpfProcedure, FromEbpfKind)
	execute.RegisterSource(FromEbpfKind, createFromEbpfSource)
}

func createFromEbpfOpSpec(args flux.Arguments, a *flux.Administration) (flux.OperationSpec, error) {
	spec := new(FromEbpfOpSpec)

	if t, ok, err := args.GetString("file"); err != nil {
		return nil, err
	}

	return spec, nil
}

func newFromEbpfOp() flux.OperationSpec {
	return new(FromEbpfOpSpec)
}

func (s *FromEbpfOpSpec) Kind() flux.OperationKind {
	return FromEbpfKind
}

type FromEbpfProcedureSpec struct {
	File string
}

func newFromEbpfProcedure(qs flux.OperationSpec, pa plan.Administration) (plan.ProcedureSpec, error) {
	// TODO: copy over data from the OpSpec to the ProcedureSpec
	spec, ok := qs.(*FromEbpfOpSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", qs)
	}

	return &FromEbpfProcedureSpec{
		File: spec.File,
	}, nil
}

func (s *FromEbpfProcedureSpec) Kind() plan.ProcedureKind {
	return FromEbpfKind
}

func (s *FromEbpfProcedureSpec) Copy() plan.ProcedureSpec {
	ns := new(FromEbpfProcedureSpec)

	return ns
}

func createFromEbpfSource(prSpec plan.ProcedureSpec, dsid execute.DatasetID, a execute.Administration) (execute.Source, error) {
	spec, ok := prSpec.(*FromEbpfProcedureSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", prSpec)
	}

	s := ebpf.NewSource(a.Allocator())
	s.File = spec.File

	return &EbpfIterator{
		source: s,
		id:     dsid,
	}, nil
}

func (c *EbpfIterator) Connect() error {
	return nil
}
func (c *EbpfIterator) Fetch() (bool, error) {
	return false, nil
}
func (c *EbpfIterator) Decode() (flux.Table, error) {
	return nil, nil
}

type EbpfIterator struct {
	source *generate.Source
	id     execute.DatasetID
	ts     []execute.Transformation
}

func (c *EbpfIterator) Do(f func(flux.Table) error) error {
	c.source.Connect()

	more, err := c.source.Fetch()
	if err != nil {
		return err
	}
	for more {
		tbl, err := c.source.Decode()
		if err != nil {
			return err
		}
		f(tbl)
		more, err = c.source.Fetch()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *EbpfIterator) AddTransformation(t execute.Transformation) {
	c.ts = append(c.ts, t)
}

func (c *EbpfIterator) Run(ctx context.Context) {
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

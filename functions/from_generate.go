package functions

import (
	"context"
	"fmt"
	"time"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/generate"
	"github.com/influxdata/flux/plan"
	"github.com/influxdata/flux/semantic"
)

const FromGenerateKind = "fromGenerate"

type FromGenerateOpSpec struct {
	Start time.Time
	Stop  time.Time
	Count int64
}

var fromGenerateSignature = semantic.FunctionSignature{
	Params: map[string]semantic.Type{
		"start": semantic.Time,
		"stop":  semantic.Time,
		"count": semantic.Int,
	},
	ReturnType: flux.TableObjectType,
}

func init() {
	flux.RegisterFunction(FromGenerateKind, createFromGenerateOpSpec, fromGenerateSignature)
	flux.RegisterOpSpec(FromGenerateKind, newFromGenerateOp)
	plan.RegisterProcedureSpec(FromGenerateKind, newFromGenerateProcedure, FromGenerateKind)
	execute.RegisterSource(FromGenerateKind, createFromGenerateSource)
}

func createFromGenerateOpSpec(args flux.Arguments, a *flux.Administration) (flux.OperationSpec, error) {
	spec := new(FromGenerateOpSpec)

	// TODO:  read in arguments of your custom function
	if t, ok, err := args.GetTime("start"); err != nil {
		return nil, err
	} else if ok {
		spec.Start = t.Time(time.Now())
	}

	if t, ok, err := args.GetTime("stop"); err != nil {
		return nil, err
	} else if ok {
		spec.Stop = t.Time(time.Now())
	}

	if i, ok, err := args.GetInt("count"); err != nil {
		return nil, err
	} else if ok {
		spec.Count = i
	}

	return spec, nil
}

func newFromGenerateOp() flux.OperationSpec {
	return new(FromGenerateOpSpec)
}

func (s *FromGenerateOpSpec) Kind() flux.OperationKind {
	return FromGenerateKind
}

type FromGenerateProcedureSpec struct {
	Start time.Time
	Stop  time.Time
	Count int64
}

func newFromGenerateProcedure(qs flux.OperationSpec, pa plan.Administration) (plan.ProcedureSpec, error) {
	// TODO: copy over data from the OpSpec to the ProcedureSpec
	spec, ok := qs.(*FromGenerateOpSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", qs)
	}

	return &FromGenerateProcedureSpec{
		Count: spec.Count,
		Start: spec.Start,
		Stop:  spec.Stop,
	}, nil
}

func (s *FromGenerateProcedureSpec) Kind() plan.ProcedureKind {
	return FromGenerateKind
}

func (s *FromGenerateProcedureSpec) Copy() plan.ProcedureSpec {
	ns := new(FromGenerateProcedureSpec)

	return ns
}

func createFromGenerateSource(prSpec plan.ProcedureSpec, dsid execute.DatasetID, a execute.Administration) (execute.Source, error) {
	// TODO: copy over info from the ProcedureSpec you need to run your source..
	spec, ok := prSpec.(*FromGenerateProcedureSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", prSpec)
	}

	s := generate.NewSource(a.Allocator())
	s.Start = spec.Start
	s.Stop = spec.Stop
	s.Count = spec.Count

	return &GenerateIterator{
		source: s,
		id:     dsid,
	}, nil
}

func (c *GenerateIterator) Connect() error {
	return nil
}
func (c *GenerateIterator) Fetch() (bool, error) {
	return false, nil
}
func (c *GenerateIterator) Decode() (flux.Table, error) {
	return nil, nil
}

type GenerateIterator struct {
	// TODO: add fields you need to connect, fetch, etc.

	//source execute.Source
	source *generate.Source
	id     execute.DatasetID
	ts     []execute.Transformation
}

func (c *GenerateIterator) Do(f func(flux.Table) error) error {
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

func (c *GenerateIterator) AddTransformation(t execute.Transformation) {
	c.ts = append(c.ts, t)
}

func (c *GenerateIterator) Run(ctx context.Context) {
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

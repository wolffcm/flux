package functions

import (

	"context"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/plan"
	"github.com/influxdata/flux/semantic"

)

const FromTemplateKind = "fromTemplate"

type FromTemplateOpSpec struct {

}

var fromTemplateSignature = semantic.FunctionSignature{
	Params: map[string]semantic.Type{
		// TODO: indiccate the arguments and their types.
	},
	ReturnType: flux.TableObjectType,
}

func init() {
	flux.RegisterFunction(FromTemplateKind, createFromTemplateOpSpec, fromTemplateSignature)
	flux.RegisterOpSpec(FromTemplateKind, newFromTemplateOp)
	plan.RegisterProcedureSpec(FromTemplateKind, newFromTemplateProcedure, FromTemplateKind)
	execute.RegisterSource(FromTemplateKind, createFromTemplateSource)
}

func createFromTemplateOpSpec(args flux.Arguments, a *flux.Administration) (flux.OperationSpec, error) {
	spec := new(FromTemplateOpSpec)

	// TODO:  read in arguments of your custom function

	return spec, nil
}

func newFromTemplateOp() flux.OperationSpec {
	return new(FromTemplateOpSpec)
}

func (s *FromTemplateOpSpec) Kind() flux.OperationKind {
	return FromTemplateKind
}

type FromTemplateProcedureSpec struct {

}

func newFromTemplateProcedure(qs flux.OperationSpec, pa plan.Administration) (plan.ProcedureSpec, error) {
	// TODO: copy over data from the OpSpec to the ProcedureSpec
	//spec, ok := qs.(*FromTemplateOpSpec)
	//if !ok {
	//	return nil, fmt.Errorf("invalid spec type %T", qs)
	//}

	return &FromTemplateProcedureSpec{

	}, nil
}

func (s *FromTemplateProcedureSpec) Kind() plan.ProcedureKind {
	return FromTemplateKind
}

func (s *FromTemplateProcedureSpec) Copy() plan.ProcedureSpec {
	ns := new(FromTemplateProcedureSpec)

	return ns
}

func createFromTemplateSource(prSpec plan.ProcedureSpec, dsid execute.DatasetID, a execute.Administration) (execute.Source, error) {
	// TODO: copy over info from the ProcedureSpec you need to run your source..
	//spec, ok := prSpec.(*FromTemplateProcedureSpec)
	//if !ok {
	//	return nil, fmt.Errorf("invalid spec type %T", prSpec)
	//}

	return &TemplateIterator{}, nil
}

func (c *TemplateIterator) Connect() error {
	return nil
}
func (c *TemplateIterator) Fetch() (bool, error) {
	return false, nil
}
func (c *TemplateIterator) Decode() (flux.Table, error) {
	return nil, nil
}


type TemplateIterator struct {
	// TODO: add fields you need to connect, fetch, etc.

	source execute.Source
	id   execute.DatasetID
	ts   []execute.Transformation
}

func (c *TemplateIterator) Do(f func(flux.Table) error) error {
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



func (c *TemplateIterator) AddTransformation(t execute.Transformation) {
	c.ts = append(c.ts, t)
}

func (c *TemplateIterator) Run(ctx context.Context) {
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

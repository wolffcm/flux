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
	CSV  string `json:"csv"`
	File string `json:"file"`
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
	//spec, ok := prSpec.(*FromTemplateProcedureSpec)
	//if !ok {
	//	return nil, fmt.Errorf("invalid spec type %T", prSpec)
	//}

	return &TemplateIterator{}, nil
}

type TemplateIterator struct {

}

func (c *TemplateIterator) Connect() error {
return nil
}
func (c *TemplateIterator) Fetch() (bool, error) {
return false, nil
}
func (c *TemplateIterator) Decode() flux.Table {
return nil
}

func (c *TemplateIterator) AddTransformation(t execute.Transformation) {
	//c.ts = append(c.ts, t)
}

func (c *TemplateIterator) Run(ctx context.Context) {
//	var err error
//	var max execute.Time
//	maxSet := false
//	err = c.data.Tables().Do(func(tbl flux.Table) error {
//		for _, t := range c.ts {
//			err := t.Process(c.id, tbl)
//			if err != nil {
//				return err
//			}
//			if idx := execute.ColIdx(execute.DefaultStopColLabel, tbl.Key().Cols()); idx >= 0 {
//				if stop := tbl.Key().ValueTime(idx); !maxSet || stop > max {
//					max = stop
//					maxSet = true
//				}
//			}
//		}
//		return nil
//	})
//	if err != nil {
//		goto FINISH
//	}
//
//	if maxSet {
//		for _, t := range c.ts {
//			t.UpdateWatermark(c.id, max)
//		}
//	}
//
//FINISH:
//	for _, t := range c.ts {
//		t.Finish(c.id, err)
//	}
}

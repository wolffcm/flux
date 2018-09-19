package functions

import (
	"fmt"
	"time"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/compiler"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/generate"
	"github.com/influxdata/flux/interpreter"
	"github.com/influxdata/flux/plan"
	"github.com/influxdata/flux/semantic"
)

const FromGenerateKind = "fromGenerate"

type FromGenerateOpSpec struct {
	Start time.Time                    `json:"start"`
	Stop  time.Time                    `json:"stop"`
	Count int64                        `json:"count"`
	Fn    *semantic.FunctionExpression `json:"fn"`
}

var fromGenerateSignature = semantic.FunctionSignature{
	Params: map[string]semantic.Type{
		"start": semantic.Time,
		"stop":  semantic.Time,
		"count": semantic.Int,
		"fn":    semantic.Function,
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

	if t, err := args.GetRequiredTime("start"); err != nil {
		return nil, err
	} else {
		spec.Start = t.Time(time.Now())
	}

	if t, err := args.GetRequiredTime("stop"); err != nil {
		return nil, err
	} else {
		spec.Stop = t.Time(time.Now())
	}

	if i, err := args.GetRequiredInt("count"); err != nil {
		return nil, err
	} else {
		spec.Count = i
	}

	if f, err := args.GetRequiredFunction("fn"); err != nil {
		return nil, err
	} else {
		fn, err := interpreter.ResolveFunction(f)
		if err != nil {
			return nil, err
		}
		spec.Fn = fn
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
	Param string
	Fn    compiler.Func
}

func newFromGenerateProcedure(qs flux.OperationSpec, pa plan.Administration) (plan.ProcedureSpec, error) {
	// TODO: copy over data from the OpSpec to the ProcedureSpec
	spec, ok := qs.(*FromGenerateOpSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", qs)
	}

	fn, param, err := compileFnParam(spec.Fn, semantic.Int, semantic.Int)
	if err != nil {
		return nil, err
	}
	return &FromGenerateProcedureSpec{
		Count: spec.Count,
		Start: spec.Start,
		Stop:  spec.Stop,
		Param: param,
		Fn:    fn,
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
	s.Param = spec.Param
	s.Fn = spec.Fn

	return CreateFromSourceIterator(s, dsid, a)
}

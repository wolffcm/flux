package functions

import (
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

	return CreateFromSourceIterator(s, dsid, a)
}


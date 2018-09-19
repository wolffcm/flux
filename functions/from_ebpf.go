package functions

import (
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

	if f, ok, err := args.GetString("file"); err != nil {
		return nil, err
	} else if ok {
		spec.File = f
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

	return CreateFromSourceIterator(s, dsid, a)
}

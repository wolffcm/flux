package functions

import (
	"fmt"
	"io/ioutil"
	"os"

	"context"

	"github.com/influxdata/flux"
	// "github.com/influxdata/flux/csv"

	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/plan"
	"github.com/influxdata/flux/semantic"
	"github.com/pkg/errors"
)

const FromBPFKind = "fromBPF"

type FromBPFOpSpec struct {
	File string `json:"file"`
}

var fromBPFSignature = semantic.FunctionSignature{
	Params: map[string]semantic.Type{
		"file": semantic.String,
	},
	ReturnType: flux.TableObjectType,
}

func init() {
	flux.RegisterFunction(FromBPFKind, createFromBPFOpSpec, fromBPFSignature)
	flux.RegisterOpSpec(FromBPFKind, newFromBPFOp)
	plan.RegisterProcedureSpec(FromBPFKind, newFromBPFProcedure, FromBPFKind)
	execute.RegisterSource(FromBPFKind, createFromBPFSource)
}

func createFromBPFOpSpec(args flux.Arguments, a *flux.Administration) (flux.OperationSpec, error) {
	spec := new(FromBPFOpSpec)

	if file, ok, err := args.GetString("file"); err != nil {
		return nil, err
	} else if ok {
		spec.File = file
	}

	if spec.File == "" {
		return nil, errors.New("must provide filename containing c source")
	}

	if _, err := os.Stat(spec.File); err != nil {
		return nil, errors.Wrap(err, "failed to stat c source file: ")
	}

	return spec, nil
}

func newFromBPFOp() flux.OperationSpec {
	return new(FromBPFOpSpec)
}

func (s *FromBPFOpSpec) Kind() flux.OperationKind {
	return FromBPFKind
}

func newFromBPFProcedure(qs flux.OperationSpec, pa plan.Administration) (plan.ProcedureSpec, error) {
	spec, ok := qs.(*FromBPFOpSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", qs)
	}

	return &FromBPFProcedureSpec{
		File: spec.File,
	}, nil
}

func (s *FromBPFOpSpec) Kind() plan.ProcedureKind {
	return FromBPFKind
}

func (s *FromBPFOpSpec) Copy() plan.ProcedureSpec {
	ns := new(FromBPFOpSpec)
	ns.File = s.File
	return ns
}

func createFromBPFSource(prSpec plan.ProcedureSpec, dsid execute.DatasetID, a execute.Administration) (execute.Source, error) {
	spec, ok := prSpec.(*FromBPFOpSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", prSpec)
	}

	bpfBytes, err := ioutil.ReadFile(spec.File)
	if err != nil {
		return nil, err
	}
	bpfSource := string(bpfBytes)

	// TODO: bpf goes heres
	// decoder := csv.NewResultDecoder(csv.ResultDecoderConfig{})
	// result, err := decoder.Decode(strings.NewReader(bpfSource))
	// if err != nil {
	// 	return nil, err
	// }
	csvSource := BPFSource{id: dsid, data: result}

	return &csvSource, nil
}

type BPFSource struct {
	id   execute.DatasetID
	data flux.Result
	ts   []execute.Transformation
}

func (b *BPFSource) AddTransformation(t execute.Transformation) {
	b.ts = append(b.ts, t)
}

func (b *BPFSource) Run(ctx context.Context) {
	var err error
	var max execute.Time
	maxSet := false
	err = b.data.Tables().Do(func(tbl flux.Table) error {
		for _, t := range b.ts {
			err := t.Process(b.id, tbl)
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
		for _, t := range b.ts {
			t.UpdateWatermark(b.id, max)
		}
	}

FINISH:
	for _, t := range b.ts {
		t.Finish(b.id, err)
	}
}

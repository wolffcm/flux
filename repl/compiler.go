package repl

import (
	"context"

	"github.com/wolffcm/flux"
	"github.com/wolffcm/flux/lang"
	"github.com/wolffcm/flux/plan"
)

// CompilerType specific to the Flux REPL
const CompilerType = "REPL"

// Compiler specific to the Flux REPL
type Compiler struct {
	Spec *flux.Spec `json:"spec"`
}

func (c Compiler) Compile(ctx context.Context) (flux.Program, error) {
	planner := plan.PlannerBuilder{}.Build()
	ps, err := planner.Plan(c.Spec)
	if err != nil {
		return nil, err
	}

	return &lang.Program{
		PlanSpec: ps,
	}, err
}

func (c Compiler) CompilerType() flux.CompilerType {
	return CompilerType
}

func AddCompilerToMappings(mappings flux.CompilerMappings) error {
	return mappings.Add(CompilerType, func() flux.Compiler {
		return new(Compiler)
	})
}

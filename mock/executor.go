package mock

import (
	"context"

	"github.com/wolffcm/flux"
	"github.com/wolffcm/flux/execute"
	"github.com/wolffcm/flux/memory"
	"github.com/wolffcm/flux/plan"
)

var _ execute.Executor = (*Executor)(nil)

var NoMetadata <-chan flux.Metadata

// Executor is a mock implementation of an execute.Executor.
type Executor struct {
	ExecuteFn func(ctx context.Context, p *plan.Spec, a *memory.Allocator) (map[string]flux.Result, <-chan flux.Metadata, error)
}

// NewExecutor returns a mock Executor where its methods will return zero values.
func NewExecutor() *Executor {
	return &Executor{
		ExecuteFn: func(context.Context, *plan.Spec, *memory.Allocator) (map[string]flux.Result, <-chan flux.Metadata, error) {
			return nil, NoMetadata, nil
		},
	}
}

func (e *Executor) Execute(ctx context.Context, p *plan.Spec, a *memory.Allocator) (map[string]flux.Result, <-chan flux.Metadata, error) {
	return e.ExecuteFn(ctx, p, a)
}

func init() {
	noMetaCh := make(chan flux.Metadata)
	close(noMetaCh)
	NoMetadata = noMetaCh
}

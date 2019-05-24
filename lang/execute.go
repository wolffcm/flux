package lang

import (
	"context"
	"log"
	"time"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/internal/spec"
	"github.com/influxdata/flux/interpreter"
	"github.com/influxdata/flux/memory"
	"github.com/influxdata/flux/values"
	"go.uber.org/zap"
)

type CreateContextAwareValue func(e StreamEvaluator) values.Value

var executionAwareValues = make(map[string]CreateContextAwareValue)

func RegisterContextAwareValue(name string, v CreateContextAwareValue) bool {
	_, replace := executionAwareValues[name]
	executionAwareValues[name] = v
	return replace
}

func BindContextAwareValues(prelude interpreter.Scope, e StreamEvaluator) {
	for name, f := range executionAwareValues {
		prelude.Set(name, f(e))
	}
}

// StreamEvaluator evaluates an intermediate result of computation (a stream, or table object)
// and returns its result.
type StreamEvaluator interface {
	Eval(to *flux.TableObject) (flux.Query, error)
}

// executor is a StreamEvaluator.
type executor struct {
	opts   *compileOptions
	ctx    context.Context
	now    time.Time
	alloc  *memory.Allocator
	deps   execute.Dependencies
	logger *zap.Logger
}

func (e *executor) Eval(to *flux.TableObject) (flux.Query, error) {
	p, err := CompileTableObject(to, e.now)
	if err != nil {
		return nil, err
	}
	p.opts = e.opts
	p.SetExecutorDependencies(e.deps)
	p.SetLogger(e.logger)
	return p.Start(e.ctx, e.alloc)
}

// TableObjectCompiler compiles a TableObject into an executable flux.Program.
// It is not added to CompilerMappings and it is not serializable, because
// it is impossible to use it outside of the context of an ongoing execution of a program.
type TableObjectCompiler struct {
	Tables *flux.TableObject
	Now    time.Time
}

func (c *TableObjectCompiler) Compile(ctx context.Context) (flux.Program, error) {
	// Ignore context, it will be provided upon Program Start.
	return CompileTableObject(c.Tables, c.Now)
}

func (*TableObjectCompiler) CompilerType() flux.CompilerType {
	panic("TableObjectCompiler is not associated with a CompilerType")
}

// CompileTableObject evaluates a TableObject and produces a flux.Program.
// `now` parameter must be non-zero, that is the default now time should be set before compiling.
func CompileTableObject(to *flux.TableObject, now time.Time, opts ...CompileOption) (*Program, error) {
	o := applyOptions(opts...)
	s := spec.FromTableObject(to, now)
	if o.verbose {
		log.Println("Query Spec: ", flux.Formatted(s, flux.FmtJSON))
	}
	ps, err := buildPlan(s, o)
	if err != nil {
		return nil, err
	}
	return &Program{
		opts:     o,
		PlanSpec: ps,
	}, nil
}

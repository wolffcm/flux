package llvm

import (
	"fmt"
	"github.com/influxdata/flux"
	"github.com/influxdata/flux/ast"
	"github.com/influxdata/flux/semantic"
	"github.com/llvm-mirror/llvm/bindings/go/llvm"
	"runtime/debug"
	"sort"
)

const (
	target   = "asmjs-unknown-emscripten"
	mainFunc = "flux_main"
)

func Build(astPkg *ast.Package) (mod llvm.Module, err error) {
	defer func() {
		if e := recover(); e != nil {
			mod = llvm.Module{}
			err = fmt.Errorf("panic: %v\nstack:\n%v", e, string(debug.Stack()))
		}
	}()

	pkg, ts, err := toSemantic(astPkg)
	if err != nil {
		return llvm.Module{}, err
	}
	v := &builder{
		typeSol:              ts,
		b:                    llvm.NewBuilder(),
		condStates:           make(map[*semantic.ConditionalExpression]condState),
		symTab: newSymbolTable(),
	}
	mod = llvm.NewModule("flux_module")

	// Declare builtins
	for _, bi := range builtins {
		llvm.AddFunction(mod, bi.name, bi.typ)
	}

	// Create top-level function
	main := llvm.FunctionType(llvm.VoidType(), []llvm.Type{}, false)
	llvm.AddFunction(mod, mainFunc, main)
	mainFunc := mod.NamedFunction(mainFunc)
	v.module = mod
	v.currentFn = mainFunc
	block := llvm.AddBasicBlock(mainFunc, "entry")

	v.b.SetInsertPointAtEnd(block)

	// Define global strings
	for name, str := range globalStrings {
		v.b.CreateGlobalStringPtr(str, name)
	}

	if err := v.Walk(pkg); err != nil {
		return llvm.Module{}, fmt.Errorf("could not generate IR: %v", v.err)
	}
	v.b.CreateRetVoid()

	if err := llvm.VerifyModule(mod, llvm.ReturnStatusAction); err != nil {
		return llvm.Module{}, fmt.Errorf("error verifying module: %v", err.Error())
	}

	mod.SetTarget(target)

	return mod, nil
}

func toSemantic(astPkg *ast.Package) (semantic.Node, semantic.TypeSolution, error) {
	semPkg, err := semantic.New(astPkg)
	if err != nil {
		return nil, nil, err
	}

	// Sort arguments in each call expression, to avoid having to do it
	// later every time we visit a call.
	sortCallParams(semPkg)

	extern := &semantic.Extern{
		Block: &semantic.ExternBlock{
			Node: semPkg,
		},
	}
	extern.Assignments = []*semantic.ExternalVariableAssignment{
		{
			Identifier: &semantic.Identifier{
				Name: "println",
			},
			ExternType: semantic.NewFunctionPolyType(semantic.FunctionPolySignature{
				Parameters: map[string]semantic.PolyType{"v": semantic.Tvar(1)},
				Required:   semantic.LabelSet([]string{"v"}),
				Return:     semantic.Int,
			}),
		},
	}
	ts, err := semantic.InferTypes(extern, flux.StdLib())
	if err != nil {
		return nil, nil, err
	}

	return extern, ts, nil
}

func sortCallParams(semPkg semantic.Node) {
	semantic.Walk(paramSortingVisitor{}, semPkg)
}

type paramSortingVisitor struct{}

func (v paramSortingVisitor) Visit(node semantic.Node) semantic.Visitor {
	return v
}

func (paramSortingVisitor) Done(node semantic.Node) {
	if ce, ok := node.(*semantic.CallExpression); ok {
		args := ce.Arguments.Properties
		sort.Slice(args, func(i, j int) bool {
			return args[i].Key.Key() < args[j].Key.Key()
		})
	} else if fe, ok := node.(*semantic.FunctionExpression); ok {
		params := fe.Block.Parameters.List
		sort.Slice(params, func(i, j int) bool {
			return params[i].Key.Name < params[j].Key.Name
		})
	}
}



package llvm

import (
	"errors"
	"fmt"
	"github.com/influxdata/flux"
	"runtime/debug"
	"sort"

	"github.com/influxdata/flux/ast"
	"github.com/influxdata/flux/semantic"
	"github.com/llvm-mirror/llvm/bindings/go/llvm"
)

const (
	target   = "asmjs-unknown-emscripten"
	mainFunc = "flux_main"

	printlnI64Fmt = "println_i64_fmt"
	printlnStrFmt = "println_str_fmt"
)

var i8PtrTy llvm.Type

// TODO(cwolff): move these to their own file
var builtins map[string]builtinInfo
var globalStrings map[string]string

func init() {
	i8PtrTy = llvm.PointerType(llvm.Int8Type(), 0)

	builtins = map[string]builtinInfo{
		"println": {
			name: "printf",
			typ: llvm.FunctionType(
				llvm.Int32Type(),
				[]llvm.Type{
					llvm.PointerType(llvm.Int8Type(), 0),
				},
				true,
			),
			nargs: 2,
			pushArgs: func(b *builder, ce *semantic.CallExpression) error {
				fluxArg := ce.Arguments.Properties[0].Value
				var format llvm.Value
				typ, err := b.ts.TypeOf(fluxArg)
				if err != nil {
					// If the type if polymorphic, just assume int64 for now
					typ = semantic.Int
				}
				switch typ {
				case semantic.Int:
					format = b.m.NamedGlobal(printlnI64Fmt)
				case semantic.String:
					format = b.m.NamedGlobal(printlnStrFmt)
				default:
					return errors.New("unsupported type to println: " + typ.Nature().String())
				}
				cast := b.b.CreatePointerCast(format, i8PtrTy, "")
				b.push(cast)

				semantic.Walk(b, fluxArg)

				return nil
			},
		},
	}
	globalStrings = map[string]string{
		printlnI64Fmt: "%lld\n",
		printlnStrFmt: "%s\n",
	}

}

type builtinInfo struct {
	name     string
	typ      llvm.Type
	nargs    int
	pushArgs func(b *builder, ce *semantic.CallExpression) error
}

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
		ts:                   ts,
		b:                    llvm.NewBuilder(),
		names:                make(map[string]llvm.Value),
		condStates:           make(map[*semantic.ConditionalExpression]condState),
		callStates:           make(map[*semantic.CallExpression]builtinInfo),
		builtinReverseLookup: make(map[llvm.Value]builtinInfo),
	}
	mod = llvm.NewModule("flux_module")

	// Declare builtins
	for _, bi := range builtins {
		llvm.AddFunction(mod, bi.name, bi.typ)
		fn := mod.NamedFunction(bi.name)
		v.builtinReverseLookup[fn] = bi
	}

	// Create top-level function
	main := llvm.FunctionType(llvm.VoidType(), []llvm.Type{}, false)
	llvm.AddFunction(mod, mainFunc, main)
	mainFunc := mod.NamedFunction(mainFunc)
	v.m = mod
	v.f = mainFunc
	block := llvm.AddBasicBlock(mainFunc, "entry")

	v.b.SetInsertPointAtEnd(block)

	// Define global strings
	for name, str := range globalStrings {
		v.b.CreateGlobalStringPtr(str, name)
	}

	semantic.Walk(v, pkg)
	if v.err != nil {
		return llvm.Module{}, fmt.Errorf("coult not generate IR: %v", v.err)
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

type builder struct {
	ts     semantic.TypeSolution
	m      llvm.Module
	f      llvm.Value
	values []llvm.Value
	b      llvm.Builder
	names  map[string]llvm.Value
	idCtr  int64

	builtinReverseLookup map[llvm.Value]builtinInfo

	err error

	condStates map[*semantic.ConditionalExpression]condState
	callStates map[*semantic.CallExpression]builtinInfo
}

type condState struct {
	before              llvm.BasicBlock
	consEntry, consExit llvm.BasicBlock
	altEntry, altExit   llvm.BasicBlock
	after               llvm.BasicBlock
}

func (b *builder) newID() int64 {
	v := b.idCtr
	b.idCtr++
	return v
}

func (b *builder) push(v llvm.Value) {
	b.values = append(b.values, v)
}

func (b *builder) pop() llvm.Value {
	v := b.values[len(b.values)-1]
	b.values = b.values[:len(b.values)-1]
	return v
}

func (b *builder) peek() llvm.Value {
	return b.values[len(b.values)-1]
}

func (b *builder) Visit(node semantic.Node) semantic.Visitor {
	if b.err != nil {
		return nil
	}
	switch n := node.(type) {
	case *semantic.ConditionalExpression:

		// Generate code for test, leave register on stack
		semantic.Walk(b, n.Test)

		cs := condState{
			before: b.b.GetInsertBlock(),
			after:  llvm.AddBasicBlock(b.f, fmt.Sprintf("merge%d", b.newID())),
		}

		cs.consEntry = llvm.AddBasicBlock(b.f, fmt.Sprintf("true%d", b.newID()))
		b.b.SetInsertPointAtEnd(cs.consEntry)
		semantic.Walk(b, n.Consequent)
		b.b.CreateBr(cs.after)
		cs.consExit = b.b.GetInsertBlock()

		cs.altEntry = llvm.AddBasicBlock(b.f, fmt.Sprintf("false%d", b.newID()))
		b.b.SetInsertPointAtEnd(cs.altEntry)
		semantic.Walk(b, n.Alternate)
		b.b.CreateBr(cs.after)
		cs.altExit = b.b.GetInsertBlock()

		cs.after.MoveAfter(cs.altExit)

		b.b.SetInsertPointAtEnd(cs.before)

		b.condStates[n] = cs
		// We already recursed into all children, so return nil.
		return nil
	case *semantic.CallExpression:
		semantic.Walk(b, n.Callee)

		if bi, ok := b.builtinReverseLookup[b.peek()]; ok && bi.pushArgs != nil {
			b.callStates[n] = bi
			if err := bi.pushArgs(b, n); err != nil {
				b.err = err
				return nil
			}
		} else {
			pushArgs(b, n)
		}

		if n.Pipe != nil {
			panic("pipe expression unsupported")
		}

		return nil
	case *semantic.FunctionExpression:
		v := buildFunctionExpression(b, n)
		b.push(v)
		return nil
	}
	return b
}

func (b *builder) Done(node semantic.Node) {
	if b.err != nil {
		return
	}
	switch n := node.(type) {
	case *semantic.NativeVariableAssignment:
		v := b.pop()
		b.names[n.Identifier.Name] = b.b.CreateAlloca(v.Type(), n.Identifier.Name)
		b.b.CreateStore(v, b.names[n.Identifier.Name])
	case *semantic.ExpressionStatement:
		b.pop()
	case *semantic.IdentifierExpression:
		if v, ok := b.names[n.Name]; ok {
			lv := b.b.CreateLoad(v, "")
			b.push(lv)
		} else {
			// Must be a call to a pre-defined function
			bi, ok := builtins[n.Name]
			if !ok {
				b.err = errors.New("Undefined identifier: " + n.Name)
			}
			callee := b.m.NamedFunction(bi.name)
			b.push(callee)
		}
	case *semantic.BinaryExpression:
		op2 := b.pop()
		op1 := b.pop()
		var v llvm.Value
		switch n.Operator {
		case ast.AdditionOperator:
			v = b.b.CreateAdd(op1, op2, "")
		case ast.SubtractionOperator:
			v = b.b.CreateSub(op1, op2, "")
		case ast.MultiplicationOperator:
			v = b.b.CreateMul(op1, op2, "")
		case ast.DivisionOperator:
			v = b.b.CreateSDiv(op1, op2, "")
		case ast.EqualOperator:
			v = b.b.CreateICmp(llvm.IntEQ, op1, op2, "")
		case ast.GreaterThanOperator:
			v = b.b.CreateICmp(llvm.IntSGT, op1, op2, "")
		case ast.LessThanOperator:
			v = b.b.CreateICmp(llvm.IntSLT, op1, op2, "")
		case ast.GreaterThanEqualOperator:
			v = b.b.CreateICmp(llvm.IntSGE, op1, op2, "")
		case ast.LessThanEqualOperator:
			v = b.b.CreateICmp(llvm.IntSLE, op1, op2, "")
		default:
			panic("unsupported binary operand")
		}
		b.push(v)
	case *semantic.ConditionalExpression:
		cs := b.condStates[n]
		alt := b.pop()
		cons := b.pop()
		t := b.pop()
		b.b.CreateCondBr(t, cs.consEntry, cs.altEntry)

		b.b.SetInsertPointAtEnd(cs.after)
		phi := b.b.CreatePHI(cons.Type(), "")
		phi.AddIncoming([]llvm.Value{cons, alt}, []llvm.BasicBlock{cs.consExit, cs.altExit})
		b.push(phi)

		delete(b.condStates, n)
	case *semantic.Identifier:
		// Do nothing, parent will generate appropriate code for context.
	case *semantic.IntegerLiteral:
		v := llvm.ConstInt(llvm.Int64Type(), uint64(n.Value), false)
		b.push(v)
	case *semantic.StringLiteral:
		lit := b.b.CreateGlobalStringPtr(n.Value, "str")
		v := b.b.CreatePointerCast(lit, i8PtrTy, "")
		b.push(v)
	case *semantic.CallExpression:
		var nargs int
		if bi, ok := b.callStates[n]; ok {
			nargs = bi.nargs
		} else {
			nargs = len(n.Arguments.Properties)
		}

		args := make([]llvm.Value, nargs)
		for i := nargs - 1; i >= 0; i-- {
			args[i] = b.pop()
		}
		callee := b.pop()
		v := b.b.CreateCall(callee, args, "")
		b.push(v)
		delete(b.callStates, n)
	case *semantic.FunctionExpression:
		// Do nothing; function expression is on top of stack
	case *semantic.ExternalVariableAssignment:
		if _, ok := builtins[n.Identifier.Name]; !ok {
			b.err = errors.New("undefined extern: " + n.Identifier.Name)
		}
	case *semantic.ExternBlock:
	case *semantic.Extern:
	case *semantic.File:
	case *semantic.Package:
	default:
		panic("unsupported node: " + node.NodeType())
	}
}

func pushArgs(b *builder, ce *semantic.CallExpression) {
	args := ce.Arguments.Properties
	sortedArgs := make([]*semantic.Property, len(args))
	copy(sortedArgs, args)
	sort.Slice(sortedArgs, func(i, j int) bool {
		return sortedArgs[i].Key.Key() < sortedArgs[j].Key.Key()
	})

	for _, a := range sortedArgs {
		semantic.Walk(b, a.Value)
	}
}

func buildFunctionExpression(b *builder, fe *semantic.FunctionExpression) llvm.Value {
	if fe.Defaults != nil && len(fe.Defaults.Properties) > 0 {
		panic("default arguments not supported")
	}

	fty := buildFunctionType(b, fe)
	fn := llvm.AddFunction(b.m, "fun", fty)
	entry := llvm.AddBasicBlock(fn, "entry")

	caller := b.f
	callerNames := b.names
	callerBlock := b.b.GetInsertBlock()

	defer func() {
		b.f = caller
		b.names = callerNames
		b.b.SetInsertPointAtEnd(callerBlock)
	}()

	b.f = fn
	b.names = make(map[string]llvm.Value)
	b.b.SetInsertPointAtEnd(entry)

	// TODO(cwolff): should sort these to get a deterministic order

	for i, param := range fn.Params() {
		name := fe.Block.Parameters.List[i].Key.Name
		v := b.b.CreateAlloca(llvm.Int64Type(), name)
		b.b.CreateStore(param, v)
		b.names[name] = v
	}

	// For now assume that body is just a simple expression
	e := fe.Block.Body.(semantic.Expression)
	semantic.Walk(b, e)
	v := b.pop()
	b.b.CreateRet(v)

	return fn
}

func buildFunctionType(b *builder, fe *semantic.FunctionExpression) llvm.Type {
	// For now, assume all inputs are int64 and output is int64
	rty := llvm.Int64Type()

	ptys := make([]llvm.Type, len(fe.Block.Parameters.List))
	for i := range ptys {
		ptys[i] = llvm.Int64Type()
	}
	return llvm.FunctionType(rty, ptys, false)
}

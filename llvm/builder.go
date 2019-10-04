package llvm

import (
	"errors"
	"fmt"
	"github.com/influxdata/flux"
	"github.com/influxdata/flux/ast"
	"github.com/influxdata/flux/semantic"
	"github.com/llvm-mirror/llvm/bindings/go/llvm"
	"runtime/debug"
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
		ts:                   ts,
		b:                    llvm.NewBuilder(),
		names:                make(map[string]llvm.Value),
		condStates:           make(map[*semantic.ConditionalExpression]condState),
		builtinReverseLookup: make(map[llvm.Value]builtinInfo),
		env: make(map[string]semantic.Expression),
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

var i8PtrTy llvm.Type

func init() {
	i8PtrTy = llvm.PointerType(llvm.Int8Type(), 0)
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
	env map[string]semantic.Expression
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
	case *semantic.NativeVariableAssignment:
		b.env[n.Identifier.Name] = n.Init
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
		return b.buildCallExpression(n)
	case *semantic.FunctionExpression:
		// This builds the function expression and leaves on the stack.
		return b.buildFunctionExpression(n)
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

		var nat semantic.Nature
		if ty, err := b.ts.TypeOf(n.Left); err != nil || ty == nil {
			// probably we have unbound type variables.
			// choose a concrete llvm type based on the polytype.
			ty, err := b.ts.PolyTypeOf(n)
			if err != nil {
				b.err = err
				return
			} else if ty == nil {
				b.err = errors.New("type of " + n.NodeType() + " node is nil")
			}
			nat, err = b.natureFromPolyType(ty)
			if err != nil {
				b.err = err
				return
			}
		} else {
			nat = ty.Nature()
		}

		var v llvm.Value
		var err error
		switch nat {
		case semantic.Int:
			v, err = b.genBinaryIntInsn(n, op1, op2)
		case semantic.Float:
			v, err = b.genBinaryFloatInsn(n, op1, op2)
		default:
			err = errors.New("unable to get type of " + nat.String())
		}
		if err != nil {
			b.err = err
			return
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
	case *semantic.FloatLiteral:
		v := llvm.ConstFloat(llvm.DoubleType(), n.Value)
		b.push(v)
	case *semantic.StringLiteral:
		lit := b.b.CreateGlobalStringPtr(n.Value, "str")
		v := b.b.CreatePointerCast(lit, i8PtrTy, "")
		b.push(v)
	case *semantic.CallExpression:
		// Do nothing; call expression is on top of stack
	case *semantic.FunctionExpression:
		// Do nothing; function expression is on top of stack
	case *semantic.ExternalVariableAssignment:
		if _, ok := builtins[n.Identifier.Name]; !ok {
			b.err = errors.New("undefined extern: " + n.Identifier.Name)
		}
	case *semantic.ReturnStatement:
		v := b.pop()
		b.b.CreateRet(v)
	case *semantic.ExternBlock:
	case *semantic.Extern:
	case *semantic.File:
	case *semantic.Package:
	default:
		b.err = errors.New("unsupported node: " + node.NodeType())
		return
	}
}


func (b *builder) natureFromPolyType(ty semantic.PolyType) (semantic.Nature, error) {
	return semantic.Int, nil
}

func (b *builder) genBinaryFloatInsn(node *semantic.BinaryExpression, op1, op2 llvm.Value) (llvm.Value, error) {
	var v llvm.Value
	switch node.Operator {
	case ast.AdditionOperator:
		v = b.b.CreateFAdd(op1, op2, "")
	case ast.SubtractionOperator:
		v = b.b.CreateFSub(op1, op2, "")
	case ast.MultiplicationOperator:
		v = b.b.CreateFMul(op1, op2, "")
	case ast.DivisionOperator:
		v = b.b.CreateFDiv(op1, op2, "")
	case ast.EqualOperator:
		v = b.b.CreateFCmp(llvm.FloatOEQ, op1, op2, "")
	case ast.GreaterThanOperator:
		v = b.b.CreateFCmp(llvm.FloatOGT, op1, op2, "")
	case ast.LessThanOperator:
		v = b.b.CreateFCmp(llvm.FloatOLT, op1, op2, "")
	case ast.GreaterThanEqualOperator:
		v = b.b.CreateFCmp(llvm.FloatOGE, op1, op2, "")
	case ast.LessThanEqualOperator:
		v = b.b.CreateFCmp(llvm.FloatOLE, op1, op2, "")
	default:
		return llvm.Value{}, errors.New("unsupported binary operand")
	}

	return v, nil
}

func (b *builder) genBinaryIntInsn(node *semantic.BinaryExpression, op1, op2 llvm.Value) (llvm.Value, error) {
	var v llvm.Value
	switch node.Operator {
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
		return llvm.Value{}, errors.New("unsupported binary operand")
	}

	return v, nil
}

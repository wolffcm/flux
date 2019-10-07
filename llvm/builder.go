package llvm

import (
	"errors"
	"fmt"
	"github.com/influxdata/flux/ast"
	"github.com/influxdata/flux/semantic"
	"github.com/llvm-mirror/llvm/bindings/go/llvm"
)

type builder struct {
	err error

	idCtr  int64

	typeSol    semantic.TypeSolution
	module     llvm.Module
	b          llvm.Builder

	currentFn  llvm.Value
	valueStack []llvm.Value
	condStates map[*semantic.ConditionalExpression]condState
	symTab *symbolTable
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
	b.valueStack = append(b.valueStack, v)
}

func (b *builder) pop() llvm.Value {
	v := b.valueStack[len(b.valueStack)-1]
	b.valueStack = b.valueStack[:len(b.valueStack)-1]
	return v
}

func (b *builder) peek() llvm.Value {
	return b.valueStack[len(b.valueStack)-1]
}

func (b *builder) Walk(node semantic.Node) error {
	semantic.Walk(b, node)
	if b.err != nil {
		return b.err
	}
	return nil
}

func (b *builder) Visit(node semantic.Node) semantic.Visitor {
	if b.err != nil {
		return nil
	}
	switch n := node.(type) {
	case *semantic.NativeVariableAssignment:
		err := b.symTab.addEntry(n.Identifier.Name, n.Init, nil)
		if err != nil {
			b.err = err
			return nil
		}
	case *semantic.ConditionalExpression:

		// Generate code for test, leave register on stack
		if err := b.Walk(n.Test); err != nil {
			return nil
		}

		cs := condState{
			before: b.b.GetInsertBlock(),
			after:  llvm.AddBasicBlock(b.currentFn, fmt.Sprintf("merge%d", b.newID())),
		}

		cs.consEntry = llvm.AddBasicBlock(b.currentFn, fmt.Sprintf("true%d", b.newID()))
		b.b.SetInsertPointAtEnd(cs.consEntry)
		if err := b.Walk(n.Consequent); err != nil {
			return nil
		}
		b.b.CreateBr(cs.after)
		cs.consExit = b.b.GetInsertBlock()

		cs.altEntry = llvm.AddBasicBlock(b.currentFn, fmt.Sprintf("false%d", b.newID()))
		b.b.SetInsertPointAtEnd(cs.altEntry)
		if err := b.Walk(n.Alternate); err != nil {
			return nil
		}
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
		name := n.Identifier.Name
		v := b.pop()
		typ := v.Type()
		alloca := b.b.CreateAlloca(typ, name)
		b.b.CreateStore(v, alloca)
		if err := b.symTab.addEntry(name, n.Init, &alloca); err != nil {
			b.err = err
			return
		}
	case *semantic.ExpressionStatement:
		b.pop()
	case *semantic.IdentifierExpression:
		v, err := b.symTab.getSingleValue(n.Name)
		if err != nil && err != symbolNotFound {
			b.err = err
			return
		} else if err == symbolNotFound {
			// Must be a call to a pre-defined function
			// TODO(cwolff): add predfined symbols to table
			bi, ok := builtins[n.Name]
			if !ok {
				b.err = errors.New("Undefined identifier: " + n.Name)
			}
			callee := b.module.NamedFunction(bi.name)
			b.push(callee)
		} else {
			lv := b.b.CreateLoad(v, "")
			b.push(lv)
		}
	case *semantic.BinaryExpression:
		op2 := b.pop()
		op1 := b.pop()

		var v llvm.Value
		var err error
		switch t := op1.Type(); t {
		case llvmIntType:
			v, err = b.genBinaryIntInsn(n, op1, op2)
		case llvmFloatType:
			v, err = b.genBinaryFloatInsn(n, op1, op2)
		case llvmStringType:
			v, err = b.genBinaryStringInsn(n, op1, op2)
		default:
			err = errors.New("unable to gen binary insn for type of " + t.String())
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
		v := llvm.ConstInt(llvmIntType, uint64(n.Value), false)
		b.push(v)
	case *semantic.FloatLiteral:
		v := llvm.ConstFloat(llvmFloatType, n.Value)
		b.push(v)
	case *semantic.StringLiteral:
		lit := b.b.CreateGlobalStringPtr(n.Value, "str")
		v := b.b.CreatePointerCast(lit, llvmStringType, "")
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

func (b *builder) genBinaryFloatInsn(node *semantic.BinaryExpression, op1, op2 llvm.Value) (llvm.Value, error) {
	var v llvm.Value
	switch o := node.Operator; o {
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
		return llvm.Value{}, errors.New("unsupported binary operand for float: " + o.String())
	}

	return v, nil
}

func (b *builder) genBinaryIntInsn(node *semantic.BinaryExpression, op1, op2 llvm.Value) (llvm.Value, error) {
	var v llvm.Value
	switch o := node.Operator; o {
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
		return llvm.Value{}, errors.New("unsupported binary operand for int: " + o.String())
	}

	return v, nil
}

func (b *builder) genBinaryStringInsn(node *semantic.BinaryExpression, op1, op2 llvm.Value) (llvm.Value, error) {
	var v llvm.Value
	switch o := node.Operator; o {
	case ast.AdditionOperator:
		s1Len := b.b.CreateCall(b.module.NamedFunction("strlen"), []llvm.Value{op1}, "s1_len")
		s2Len := b.b.CreateCall(b.module.NamedFunction("strlen"), []llvm.Value{op2}, "s2_len")

		allocLen := b.b.CreateAdd(s1Len, s2Len, "sum_len")
		allocLen = b.b.CreateAdd(allocLen, llvm.ConstInt(llvmSizeType, 1, false), "alloc_len")

		// This allocation is never freed, and is a memory leak in the generated code.
		resultBuf := b.b.CreateCall(b.module.NamedFunction("malloc"), []llvm.Value{allocLen}, "res_buf")

		b.b.CreateCall(b.module.NamedFunction("strcpy"), []llvm.Value{resultBuf, op1}, "strcpy")
		v = b.b.CreateCall(b.module.NamedFunction("strcat"), []llvm.Value{resultBuf, op2}, "strcat")
	default:
		return llvm.Value{}, errors.New("unsupported binary operand for string: " + o.String())
	}

	return v, nil
}

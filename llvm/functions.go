package llvm

import (
	"errors"
	"fmt"
	"github.com/influxdata/flux/semantic"
	"github.com/llvm-mirror/llvm/bindings/go/llvm"
)

func (b *builder) buildCallExpression(callExpr *semantic.CallExpression) semantic.Visitor {
	if callExpr.Pipe != nil {
		b.err = errors.New("pipe expression unsupported")
	}

	// Generate code for the callee.
	// It might be an identifier, or something else.
	semantic.Walk(b, callExpr.Callee)
	callee := b.pop()

	// Determine if this is a call to a builtin or a function defined in Flux.
	if bi, ok := b.builtinReverseLookup[callee]; ok {
		return b.buildBuiltinCallExpression(bi, callExpr, callee)
	} else {
		return b.buildFluxCallExpression(callExpr, callee)
	}
}

func (b *builder) buildBuiltinCallExpression(bi builtinInfo, callExpr *semantic.CallExpression, callee llvm.Value) semantic.Visitor {
	// this is a builtin
	llvmArgs, err := bi.getLLVMArgs(b, callExpr.Arguments)
	if err != nil {
		b.err = err
		return nil
	}
	v := b.b.CreateCall(callee, llvmArgs, "")
	b.push(v)
	return nil
}

// Function expressions in Flux will be assigned general types with
// type variables.  At the callsite we will know the actual types
// required and can generate the corresponding code.

func (b *builder) buildFluxCallExpression(callExpr *semantic.CallExpression, callee llvm.Value) semantic.Visitor {
	fluxCalleeType, err := b.ts.PolyTypeOf(callExpr.Callee)
	if err != nil {
		b.err = err
		return nil
	}
	llvmCalleeType, _ := polyTypeToLLVMType(fluxCalleeType, true)
	fmt.Println("llvm call expr callee type: ", llvmCalleeType.String())

	llvmDefType := callee.Type().ElementType()

	if llvmDefType != llvmCalleeType {
		b.err = fmt.Errorf("call needs specialization; definition type is %v, callsite type is %v", llvmDefType, llvmCalleeType)
		return nil
	}

	args := callExpr.Arguments.Properties
	llvmArgs := make([]llvm.Value, len(args))
	for i, a := range args {
		semantic.Walk(b, a.Value)
		llvmArgs[i] = b.pop()
	}

	v := b.b.CreateCall(callee, llvmArgs, "")
	b.push(v)
	return nil
}

func (b *builder) buildFunctionExpression(fe *semantic.FunctionExpression) semantic.Visitor {
	if fe.Defaults != nil && len(fe.Defaults.Properties) > 0 {
		b.err = errors.New("default arguments not supported")
		return nil
	}

	fty, err := b.getLLVMType(fe, true)
	if err != nil {
		b.err = err
		return nil
	}

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

	// The code generator expects identifiers to have addresses, so generate
	// local variables to hold the arguments.
	llvmParamTypes := fty.ParamTypes()
	for i, param := range fn.Params() {
		name := fe.Block.Parameters.List[i].Key.Name
		v := b.b.CreateAlloca(llvmParamTypes[i], name)
		b.b.CreateStore(param, v)
		b.names[name] = v
	}

	if e, ok := fe.Block.Body.(semantic.Expression); ok {
		semantic.Walk(b, e)
		v := b.pop()
		b.b.CreateRet(v)
	} else {
		block := fe.Block.Body.(*semantic.Block)
		for _, stmt := range block.Body {
			semantic.Walk(b, stmt)
		}
	}

	b.push(fn)

	return nil
}

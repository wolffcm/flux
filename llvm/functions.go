package llvm

import (
	"errors"
	"fmt"
	"github.com/influxdata/flux/semantic"
	"github.com/llvm-mirror/llvm/bindings/go/llvm"
	"sort"
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

func (b *builder) buildFluxCallExpression(callExpr *semantic.CallExpression, callee llvm.Value) semantic.Visitor {
	// Check to see if we need to generate a specialization here.
	calleeType, err := b.ts.PolyTypeOf(callExpr.Callee)
	if err != nil {
		b.err = err
		return nil
	}

	var origType semantic.PolyType
	if ie, ok := callExpr.Callee.(*semantic.IdentifierExpression); ok {
		rhs, ok := b.env[ie.Name]
		if !ok {
			b.err = errors.New("use of undefined var")
			return nil
		}
		origType, err = b.ts.PolyTypeOf(rhs)
		if err != nil {
			b.err = err
			return nil
		}
	} else {
		origType = calleeType
	}

	if str, ok := calleeType.(fmt.Stringer); ok {
		fmt.Println("Entering call, type of callee: ", str)
	}

	if str, ok := origType.(fmt.Stringer); ok {
		fmt.Println("Entering call, type of orig fn expr: ", str)
	}

	if calleeType.Equal(origType) {
		fmt.Println("types are equal")
	} else {
		fmt.Println("types are not equal")
	}

	args := callExpr.Arguments.Properties
	sortedArgs := make([]*semantic.Property, len(args))
	copy(sortedArgs, args)
	sort.Slice(sortedArgs, func(i, j int) bool {
		return sortedArgs[i].Key.Key() < sortedArgs[j].Key.Key()
	})

	llvmArgs := make([]llvm.Value, len(args))
	for i, a := range sortedArgs {
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

	// Note: sorting the parameter list modifies it, but we can't make a copy
	// or we'll invalidate the type solution.
	paramList := fe.Block.Parameters.List
	sort.Slice(paramList, func(i, j int) bool {
		return paramList[i].Key.Name < paramList[j].Key.Name
	})

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

	// The code generator expects identifiers to have addresses, so generate
	// local variables to hold the arguments.
	for i, param := range fn.Params() {
		name := fe.Block.Parameters.List[i].Key.Name
		v := b.b.CreateAlloca(llvm.Int64Type(), name)
		b.b.CreateStore(param, v)
		b.names[name] = v
	}

	// For now assume that body is just a simple expression
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

func buildFunctionType(b *builder, fe *semantic.FunctionExpression) llvm.Type {
	// For now, assume all inputs are int64 and output is int64
	rty := llvm.Int64Type()

	ptys := make([]llvm.Type, len(fe.Block.Parameters.List))
	for i := range ptys {
		ptys[i] = llvm.Int64Type()
	}
	return llvm.FunctionType(rty, ptys, false)
}

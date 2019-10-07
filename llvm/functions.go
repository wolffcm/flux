package llvm

import (
	"errors"
	"fmt"
	"github.com/influxdata/flux/semantic"
	"github.com/llvm-mirror/llvm/bindings/go/llvm"
	"log"
)

func (b *builder) buildCallExpression(callExpr *semantic.CallExpression) semantic.Visitor {
	if callExpr.Pipe != nil {
		b.err = errors.New("pipe expression unsupported")
	}

	var name string
	if id, ok := callExpr.Callee.(*semantic.IdentifierExpression); ok {
		name = id.Name
	}

	//// Generate code for the callee.
	//// It might be an identifier, or something else.
	//if err := b.Walk(callExpr.Callee); err != nil {
	//	return nil
	//}
	//callee := b.pop()

	// Determine if this is a call to a builtin or a function defined in Flux.
	if bi, ok := builtins[name]; ok {
		return b.buildBuiltinCallExpression(bi, callExpr)
	} else {
		return b.buildFluxCallExpression(callExpr)
	}
}

func (b *builder) buildBuiltinCallExpression(bi builtinInfo, callExpr *semantic.CallExpression) semantic.Visitor {
	// Generate code for the callee.
	if err := b.Walk(callExpr.Callee); err != nil {
		return nil
	}
	callee := b.pop()

	// this is a builtin
	llvmArgs, err := bi.getLLVMArgs(b, callExpr.Arguments)
	if err != nil {
		b.err = err
		return nil
	}
	log.Printf("generating call to builtin %q with type %s", bi.name, callee.Type())
	v := b.b.CreateCall(callee, llvmArgs, "")
	b.push(v)
	return nil
}

func (b *builder) buildFluxCallExpression(callExpr *semantic.CallExpression) semantic.Visitor {
	var calleeName string
	if id, ok := callExpr.Callee.(*semantic.IdentifierExpression); ok {
		calleeName = id.Name
	}

	var llvmCallee llvm.Value
	if calleeName != "" {
		// See if this function expression type differs from the callee type.
		se := b.symTab.getEntry(calleeName)
		if se == nil {
			b.err = fmt.Errorf("could not find callee name %q", calleeName)
			return nil
		}
		fe := se.fluxExpr
		fnExprType, err := b.typeSol.PolyTypeOf(fe)
		if err != nil {
			b.err = err
			return nil
		}

		llvmFnType, err := polyTypeToLLVMType(fnExprType, true)
		if err != nil {
			b.err = err
			return nil
		}

		llvmFnType = llvm.PointerType(llvm.PointerType(llvmFnType, 0), 0)
		fn := b.symTab.getSpecialization(calleeName, llvmFnType)
		if fn == nil {
			fn, err := b.specializeFunction(callExpr)
			if err != nil {
				b.err = err
				return nil
			}
			llvmCallee = fn
		} else {
			llvmCallee = *fn
		}
	} else {
		// callee must be a function literal
		err := b.Walk(callExpr.Callee)
		if err != nil {
			b.err = err
			return nil
		}
		llvmCallee = b.pop()
	}

	args := callExpr.Arguments.Properties
	llvmArgs := make([]llvm.Value, len(args))
	for i, a := range args {
		if err := b.Walk(a.Value); err != nil {
			return nil
		}
		llvmArgs[i] = b.pop()
	}

	fn := b.b.CreateLoad(llvmCallee, calleeName)
	v := b.b.CreateCall(fn, llvmArgs, "")
	b.push(v)
	return nil
}

func (b *builder) buildFunctionExpression(fe *semantic.FunctionExpression) semantic.Visitor {
	if fe.Defaults != nil && len(fe.Defaults.Properties) > 0 {
		b.err = errors.New("default arguments not supported")
		return nil
	}

	name := b.symTab.findName(fe)
	if name == "" {
		name = "fn"
	}

	fty, err := b.getLLVMType(fe, true)
	if err != nil {
		b.err = err
		return nil
	}

	fn := llvm.AddFunction(b.module, name, fty)
	entry := llvm.AddBasicBlock(fn, "entry")

	caller := b.currentFn
	callerSymTab := b.symTab
	callerBlock := b.b.GetInsertBlock()

	defer func() {
		b.currentFn = caller
		b.symTab = callerSymTab
		b.b.SetInsertPointAtEnd(callerBlock)
	}()

	b.currentFn = fn
	b.symTab = newSymbolTable()
	b.b.SetInsertPointAtEnd(entry)

	// The code generator expects identifiers to have addresses, so generate
	// local variables to hold the arguments.
	llvmParamTypes := fty.ParamTypes()
	for i, param := range fn.Params() {
		fluxParam := fe.Block.Parameters.List[i]
		name := fluxParam.Key.Name
		v := b.b.CreateAlloca(llvmParamTypes[i], name)
		b.b.CreateStore(param, v)

		if err := b.symTab.addEntry(name, fluxParam, &v); err != nil {
			b.err = err
			return nil
		}
	}

	if e, ok := fe.Block.Body.(semantic.Expression); ok {
		if err := b.Walk(e); err != nil {
			return nil
		}
		v := b.pop()
		b.b.CreateRet(v)
	} else {
		block := fe.Block.Body.(*semantic.Block)
		for _, stmt := range block.Body {
			if err := b.Walk(stmt); err != nil {
				return nil
			}
		}
	}

	b.push(fn)

	return nil
}

func (b *builder) specializeFunction(ce *semantic.CallExpression) (llvm.Value, error) {
	callee := ce.Callee
	var defFn *semantic.FunctionExpression
	id, ok := ce.Callee.(*semantic.IdentifierExpression)
	if ! ok {
		// When can this happen?
		return llvm.Value{}, errors.New("could not find defined function")
	}

	se := b.symTab.getEntry(id.Name)
	defFn = se.fluxExpr.(*semantic.FunctionExpression)

	// Update type solution to reflect call arguments
	origTypeSol := b.typeSol
	defer func() {
		b.typeSol = origTypeSol
	}()
	b.typeSol = b.typeSol.FreshSolution()

	fnExprType, err := b.typeSol.PolyTypeOf(defFn)
	if err != nil {
		return llvm.Value{}, err
	}

	calleeType, err := b.typeSol.PolyTypeOf(callee)
	if err != nil {
		return llvm.Value{}, err
	}

	if err := b.typeSol.AddConstraint(fnExprType,  calleeType); err != nil {
		return llvm.Value{}, err
	}
	// Regenerate function expression with new type solution
	if b.buildFunctionExpression(defFn); err != nil {
		return llvm.Value{}, err
	}

	fn := b.pop()
	err = b.symTab.addEntry(id.Name, defFn, &fn)
	if err != nil {
		return llvm.Value{}, err
	}

	return fn, nil
}
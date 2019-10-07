package llvm

import (
	"errors"
	"fmt"
	"sort"

	"github.com/influxdata/flux/semantic"
	"github.com/llvm-mirror/llvm/bindings/go/llvm"
)

var (
	llvmStringType = llvm.PointerType(llvm.Int8Type(), 0)
	llvmIntType    = llvm.Int64Type()
	llvmFloatType  = llvm.DoubleType()
	llvmSizeType   = llvm.Int32Type()
)

type fnSigger interface {
	Signature() semantic.FunctionPolySignature
}

func (b *builder) getLLVMType(node semantic.Node, allowTypeVars bool) (llvm.Type, error) {
	pt, err := b.typeSol.PolyTypeOf(node)
	if err != nil {
		return llvm.Type{}, err
	}

	return polyTypeToLLVMType(pt, allowTypeVars)
}

func polyTypeToLLVMType(polyType semantic.PolyType, allowTypeVars bool) (llvm.Type, error) {
	if fs, ok := polyType.(fnSigger); ok {
		sig := fs.Signature()
		lrt, err := polyTypeToLLVMType(sig.Return, allowTypeVars)
		if err != nil {
			return llvm.Type{}, err
		}

		fpts := paramTypes(sig.Parameters)
		lpts := make([]llvm.Type, len(fpts))
		for i, fpt := range fpts {
			lpt, err := polyTypeToLLVMType(fpt, allowTypeVars)
			if err != nil {
				return llvm.Type{}, err
			}
			lpts[i] = lpt
		}

		return llvm.FunctionType(lrt, lpts, false), nil
	}

	if mt, ok := polyType.MonoType(); ok {
		switch nat := mt.Nature(); nat {
		case semantic.Int:
			return llvmIntType, nil
		case semantic.Float:
			return llvmFloatType, nil
		case semantic.String:
			return llvmStringType, nil
		default:
			return llvm.Type{}, errors.New("unsupported nature: " + nat.String())
		}
	}

	if _, ok := polyType.(semantic.Tvar); ok {
		if allowTypeVars {
			// For now, just treat any tvar as a stand-in for integer.
			return llvmIntType, nil
		}
		return llvm.Type{}, errors.New("cannot translate tvar to LLVM type")
	}

	return llvm.Type{}, fmt.Errorf("could not determine LLVM type for %#v : %T", polyType, polyType)
}

func paramTypes(params map[string]semantic.PolyType) []semantic.PolyType {
	names := make([]string, 0, len(params))
	for k, _ := range params {
		names = append(names, k)
	}
	sort.Strings(names)

	pts := make([]semantic.PolyType, len(names))
	for i, name := range names {
		pts[i] = params[name]
	}

	return pts
}

package llvm

import (
	"fmt"

	"github.com/influxdata/flux/semantic"
	"github.com/llvm-mirror/llvm/bindings/go/llvm"
)

const (
	printlnI64Fmt    = "println_i64_fmt"
	printlnStrFmt    = "println_str_fmt"
	printlnDoubleFmt = "println_double_fmt"
)

type builtinInfo struct {
	name        string
	typ         llvm.Type
	getLLVMArgs func(b *builder, fluxArgs *semantic.ObjectExpression) ([]llvm.Value, error)
}

var builtins map[string]builtinInfo
var globalStrings map[string]string

func init() {
	builtins = map[string]builtinInfo{
		"println": {
			name: "printf",
			typ: llvm.FunctionType(
				llvm.Int32Type(),
				[]llvm.Type{llvmStringType},
				true,
			),
			getLLVMArgs: func(b *builder, fluxArgs *semantic.ObjectExpression) ([]llvm.Value, error) {
				llvmArgs := make([]llvm.Value, 2)
				fluxArg := fluxArgs.Properties[0].Value

				var format llvm.Value
				typ, err := b.ts.TypeOf(fluxArg)
				if err != nil {
					// If the type if polymorphic, just assume int64 for now
					typ = semantic.Int
				}
				if typ == nil {
					return nil, fmt.Errorf("could not get type for %v", fluxArg)
				}
				switch typ {
				case semantic.Int:
					format = b.m.NamedGlobal(printlnI64Fmt)
				case semantic.String:
					format = b.m.NamedGlobal(printlnStrFmt)
				case semantic.Float:
					format = b.m.NamedGlobal(printlnDoubleFmt)
				default:
					return nil, fmt.Errorf("unsupported type to println: %v", typ)
				}
				cast := b.b.CreatePointerCast(format, llvmStringType, "")
				llvmArgs[0] = cast

				if err := b.Walk(fluxArg); err != nil {
					return nil, err
				}
				llvmArgs[1] = b.pop()

				return llvmArgs, nil
			},
		},
		"strcat": {
			name: "strcat",
			typ: llvm.FunctionType(
				llvmStringType,
				[]llvm.Type{
					llvmStringType,
					llvmStringType,
				},
				false,
			),
		},
		"strlen": {
			name: "strlen",
			typ: llvm.FunctionType(
				llvmSizeType,
				[]llvm.Type{llvmStringType},
				false,
			),
		},
		"strcpy": {
			name: "strcpy",
			typ: llvm.FunctionType(
				llvmStringType,
				[]llvm.Type{llvmStringType, llvmStringType},
				false,
			),
		},
		"malloc": {
			name: "malloc",
			typ: llvm.FunctionType(
				llvmStringType,
				[]llvm.Type{llvmSizeType},
				false,
			),
		},
	}
	globalStrings = map[string]string{
		printlnI64Fmt:    "%lld\n",
		printlnStrFmt:    "%s\n",
		printlnDoubleFmt: "%f\n",
	}
}

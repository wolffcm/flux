package llvm

import (
	"github.com/google/go-cmp/cmp"
	"github.com/influxdata/flux/semantic"
	"testing"
)

func TestPolyTypeToLLVMType(t *testing.T) {
	type testcase struct {
		name          string
		polyType      semantic.PolyType
		allowTypeVars bool
		want          string
	}

	tcs := []testcase{
		{
			name:     "int",
			polyType: semantic.Int,
			want:     "IntegerType(64 bits)",
		},
		{
			name:     "float",
			polyType: semantic.Float,
			want:     "DoubleType",
		},
		{
			name:     "string",
			polyType: semantic.String,
			want:     "PointerType(IntegerType(8 bits))",
		},
		{
			name: "fn",
			polyType: semantic.NewFunctionPolyType(semantic.FunctionPolySignature{
				Return: semantic.Float,
				Parameters: map[string]semantic.PolyType{
					"a": semantic.Float,
					"b": semantic.Float,
				},
			}),
			want: "FunctionType(DoubleType, DoubleType):DoubleType",
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			typ, err := polyTypeToLLVMType(tc.polyType, tc.allowTypeVars)
			if err != nil {
				t.Fatal(err)
			}
			got := typ.String()
			if tc.want != got {
				t.Fatalf("got unexpected LLVM type; -want/+got:\n%s", cmp.Diff(tc.want, got))
			}
		})
	}
}

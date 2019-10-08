package llvm_test

import (
	"strings"
	"testing"

	"github.com/influxdata/flux"
	_ "github.com/influxdata/flux/builtin"
	"github.com/influxdata/flux/llvm"
)

func TestBuilder(t *testing.T) {
	for _, tc := range []struct {
		name string
		flux string
		want []string
		err  string
	}{
		{
			name: "simple",
			flux: `x = 0
x`,
			want: []string{`
entry:
  %x = alloca i64
  store i64 0, i64* %x
  %0 = load i64, i64* %x
  ret void
`},
		},
		{
			name: "add",
			flux: `x = 10
y = x + 1
y
`,
			want: []string{`
  %x = alloca i64
  store i64 10, i64* %x
  %0 = load i64, i64* %x
  %1 = add i64 %0, 1
  %y = alloca i64
  store i64 %1, i64* %y
  %2 = load i64, i64* %y
  ret void
`},
		},
		{
			name: "conditional",
			flux: `x = 10
y = if x > 9 then x * 10 else x * 100
y`,
			want: []string{`
entry:
  %x = alloca i64
  store i64 10, i64* %x
  %0 = load i64, i64* %x
  %1 = icmp sgt i64 %0, 9
  br i1 %1, label %true1, label %false2

true1:                                            ; preds = %entry
  %2 = load i64, i64* %x
  %3 = mul i64 %2, 10
  br label %merge0

false2:                                           ; preds = %entry
  %4 = load i64, i64* %x
  %5 = mul i64 %4, 100
  br label %merge0

merge0:                                           ; preds = %false2, %true1
  %6 = phi i64 [ %3, %true1 ], [ %5, %false2 ]
  %y = alloca i64
  store i64 %6, i64* %y
  %7 = load i64, i64* %y
  ret void
`},
		},
		{
			name: "nested conditional",
			flux: `x = 10
y = if x < 1024 then 
      if x < 512 
        then x 
        else x * 10
      else
        if x < 2048 
          then x * 100
          else x * 1000
y`,
			want: []string{`
entry:
  %x = alloca i64
  store i64 10, i64* %x
  %0 = load i64, i64* %x
  %1 = icmp slt i64 %0, 1024
  br i1 %1, label %true1, label %false5

true1:                                            ; preds = %entry
  %2 = load i64, i64* %x
  %3 = icmp slt i64 %2, 512
  br i1 %3, label %true3, label %false4

true3:                                            ; preds = %true1
  %4 = load i64, i64* %x
  br label %merge2

false4:                                           ; preds = %true1
  %5 = load i64, i64* %x
  %6 = mul i64 %5, 10
  br label %merge2

merge2:                                           ; preds = %false4, %true3
  %7 = phi i64 [ %4, %true3 ], [ %6, %false4 ]
  br label %merge0

false5:                                           ; preds = %entry
  %8 = load i64, i64* %x
  %9 = icmp slt i64 %8, 2048
  br i1 %9, label %true7, label %false8

true7:                                            ; preds = %false5
  %10 = load i64, i64* %x
  %11 = mul i64 %10, 100
  br label %merge6

false8:                                           ; preds = %false5
  %12 = load i64, i64* %x
  %13 = mul i64 %12, 1000
  br label %merge6

merge6:                                           ; preds = %false8, %true7
  %14 = phi i64 [ %11, %true7 ], [ %13, %false8 ]
  br label %merge0

merge0:                                           ; preds = %merge6, %merge2
  %15 = phi i64 [ %7, %merge2 ], [ %14, %merge6 ]
  %y = alloca i64
  store i64 %15, i64* %y
  %16 = load i64, i64* %y
  ret void
`},
		},
		{
			name: "println",
			flux: `x = 17
println(v: x)
x + 1
`,
			want: []string{`
entry:
  %x = alloca i64
  store i64 17, i64* %x
  %0 = load i64, i64* %x
  %1 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([6 x i8], [6 x i8]* @println_i64_fmt, i32 0, i32 0), i64 %0)
  %2 = load i64, i64* %x
  %3 = add i64 %2, 1
  ret void
`},
		},
		{
			name: "hello_world",
			flux: `println(v: "hello world")`,
			want: []string{`
entry:
  %0 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([4 x i8], [4 x i8]* @println_str_fmt, i32 0, i32 0), i8* getelementptr inbounds ([12 x i8], [12 x i8]* @str, i32 0, i32 0))
  ret void
`},
		},
		{
			name: "string assign",
			flux: `foo = "hello world"
println(v: foo)`,
			want: []string{`
entry:
  %foo = alloca i8*
  store i8* getelementptr inbounds ([12 x i8], [12 x i8]* @str, i32 0, i32 0), i8** %foo
  %0 = load i8*, i8** %foo
  %1 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([4 x i8], [4 x i8]* @println_str_fmt, i32 0, i32 0), i8* %0)
  ret void
`},
		},
		{
			name: "define function",
			flux: `f = (n) => n * n
x = f(n: 3)
println(v: x)
`,
			want: []string{`
define void @flux_main() {
entry:
  %f = alloca i64 (i64)*
  store i64 (i64)* @f, i64 (i64)** %f
  %f1 = load i64 (i64)*, i64 (i64)** %f
  %0 = call i64 %f1(i64 3)
  %x = alloca i64
  store i64 %0, i64* %x
  %1 = load i64, i64* %x
  %2 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([6 x i8], [6 x i8]* @println_i64_fmt, i32 0, i32 0), i64 %1)
  ret void
}

define i64 @f(i64) {
entry:
  %n = alloca i64
  store i64 %0, i64* %n
  %1 = load i64, i64* %n
  %2 = load i64, i64* %n
  %3 = mul i64 %1, %2
  ret i64 %3
}`},
		},
		{
			name: "polymorphic function",
			flux: `ident = (v) => v
ident(v: 1)
ident(v: 20.0)
ident(v: "foo")
`,
			want: []string{`
define void @flux_main() {
entry:
  %ident = alloca i64 (i64)*
  store i64 (i64)* @ident, i64 (i64)** %ident
  %ident1 = load i64 (i64)*, i64 (i64)** %ident
  %0 = call i64 %ident1(i64 1)
  %ident2 = alloca double (double)*
  store double (double)* @ident.1, double (double)** %ident2
  %ident3 = load double (double)*, double (double)** %ident2
  %1 = call double %ident3(double 2.000000e+01)
  %ident4 = alloca i8* (i8*)*
  store i8* (i8*)* @ident.2, i8* (i8*)** %ident4
  %ident5 = load i8* (i8*)*, i8* (i8*)** %ident4
  %2 = call i8* %ident5(i8* getelementptr inbounds ([4 x i8], [4 x i8]* @str, i32 0, i32 0))
  ret void
}

define i64 @ident(i64) {
entry:
  %v = alloca i64
  store i64 %0, i64* %v
  %1 = load i64, i64* %v
  ret i64 %1
}

define double @ident.1(double) {
entry:
  %v = alloca double
  store double %0, double* %v
  %1 = load double, double* %v
  ret double %1
}

define i8* @ident.2(i8*) {
entry:
  %v = alloca i8*
  store i8* %0, i8** %v
  %1 = load i8*, i8** %v
  ret i8* %1
}`,
			},
		},
		{
			name: "concatenate strings",
			flux: `"foo" + "bar"`,
			want: []string{
				`
  %s1_len = call i32 @strlen(i8* getelementptr inbounds ([4 x i8], [4 x i8]* @str, i32 0, i32 0))
  %s2_len = call i32 @strlen(i8* getelementptr inbounds ([4 x i8], [4 x i8]* @str.1, i32 0, i32 0))
  %sum_len = add i32 %s1_len, %s2_len
  %alloc_len = add i32 %sum_len, 1
  %res_buf = call i8* @malloc(i32 %alloc_len)
  %strcpy = call i8* @strcpy(i8* %res_buf, i8* getelementptr inbounds ([4 x i8], [4 x i8]* @str, i32 0, i32 0))
  %strcat = call i8* @strcat(i8* %res_buf, i8* getelementptr inbounds ([4 x i8], [4 x i8]* @str.1, i32 0, i32 0))`,
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			astPkg, err := flux.Parse(tc.flux)
			if err != nil {
				t.Fatal(err)
			}

			llvmMod, err := llvm.Build(astPkg)
			if tc.err == "" && err != nil {
				t.Fatal(err)
			} else if tc.err != "" && err == nil {
				t.Fatalf("expected test to fail with %q but it passed", tc.err)
			} else if tc.err != "" && err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("expected test to fail with %q but got %q", tc.err, err.Error())
				}
				return
			}

			llvmText := llvmMod.String()

			for _, want := range tc.want {
				if !strings.Contains(llvmText, want) {
					t.Fatalf("did not get expected llvm IR; want:\n%s\ngot:\n%s", tc.want, llvmText)
				}
			}
		})
	}
}

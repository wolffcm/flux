package experimental_test

import (
	"context"
	"testing"

	"github.com/wolffcm/flux"
	"github.com/wolffcm/flux/codes"
	"github.com/wolffcm/flux/dependencies/dependenciestest"
	"github.com/wolffcm/flux/internal/errors"
	"github.com/wolffcm/flux/semantic"
	"github.com/wolffcm/flux/values"
)

func addFail(scope values.Scope) {
	scope.Set("fail", values.NewFunction(
		"fail",
		semantic.NewFunctionPolyType(semantic.FunctionPolySignature{
			Return: semantic.Bool,
		}),
		func(ctx context.Context, args values.Object) (values.Value, error) {
			return nil, errors.New(codes.Aborted, "fail")
		},
		false,
	))
}

func TestObjectKeys(t *testing.T) {
	script := `
import "experimental"

o = {a: 1, b: 2, c: 3}
experimental.objectKeys(o: o) == ["a", "b", "c"] or fail()
`
	ctx := dependenciestest.Default().Inject(context.Background())
	if _, _, err := flux.Eval(ctx, script, addFail); err != nil {
		t.Fatal("evaluation of objectKeys failed: ", err)
	}
}

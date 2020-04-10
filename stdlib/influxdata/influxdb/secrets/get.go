package secrets

import (
	"context"

	"github.com/wolffcm/flux"
	"github.com/wolffcm/flux/codes"
	"github.com/wolffcm/flux/internal/errors"
	"github.com/wolffcm/flux/interpreter"
	"github.com/wolffcm/flux/semantic"
	"github.com/wolffcm/flux/values"
)

const GetKind = "get"

func init() {
	flux.RegisterPackageValue("influxdata/influxdb/secrets", GetKind, GetFunc)
}

// GetFunc is a function that calls Get.
var GetFunc = makeGetFunc()

func makeGetFunc() values.Function {
	sig := semantic.NewFunctionPolyType(semantic.FunctionPolySignature{
		Parameters: map[string]semantic.PolyType{
			"key": semantic.String,
		},
		Required: semantic.LabelSet{"key"},
		Return:   semantic.String,
	})
	return values.NewFunction("get", sig, Get, false)
}

// Get retrieves the secret key identifier for a given secret.
func Get(ctx context.Context, args values.Object) (values.Value, error) {
	fargs := interpreter.NewArguments(args)
	key, err := fargs.GetRequiredString("key")
	if err != nil {
		return nil, err
	}

	ss, err := flux.GetDependencies(ctx).SecretService()
	if err != nil {
		return nil, errors.Wrapf(err, codes.Inherit, "cannot retrieve secret %q", key)
	}

	value, err := ss.LoadSecret(ctx, key)
	if err != nil {
		return nil, err
	}
	return values.NewString(value), nil
}

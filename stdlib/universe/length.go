package universe

import (
	"context"

	"github.com/wolffcm/flux/codes"
	"github.com/wolffcm/flux/internal/errors"
	"github.com/wolffcm/flux/interpreter"
	"github.com/wolffcm/flux/runtime"
	"github.com/wolffcm/flux/semantic"
	"github.com/wolffcm/flux/values"
)

// MakeLengthFunc create the "length()" function.
//
// Length will return the length of the given arr array.
func MakeLengthFunc() values.Function {
	return values.NewFunction(
		"length",
		runtime.MustLookupBuiltinType("universe", "length"),
		func(ctx context.Context, args values.Object) (values.Value, error) {
			a := interpreter.NewArguments(args)
			v, err := a.GetRequired("arr")
			if err != nil {
				return nil, err
			} else if got := v.Type().Nature(); got != semantic.Array {
				return nil, errors.Newf(codes.Invalid, "arr must be an array, got %s", got)
			}
			l := v.Array().Len()
			return values.NewInt(int64(l)), nil
		}, false,
	)
}

func init() {
	runtime.RegisterPackageValue("universe", "length", MakeLengthFunc())
}

package system

import (
	"context"
	"time"

	"github.com/wolffcm/flux/runtime"
	"github.com/wolffcm/flux/values"
)

var systemTimeFuncName = "time"

func init() {
	runtime.RegisterPackageValue("system", systemTimeFuncName, values.NewFunction(
		systemTimeFuncName,
		runtime.MustLookupBuiltinType("system", systemTimeFuncName),
		func(ctx context.Context, args values.Object) (values.Value, error) {
			return values.NewTime(values.ConvertTime(time.Now().UTC())), nil
		},
		false,
	))
}

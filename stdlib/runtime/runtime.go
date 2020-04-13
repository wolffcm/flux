package runtime

import (
	"github.com/wolffcm/flux/codes"
	"github.com/wolffcm/flux/internal/errors"
	"github.com/wolffcm/flux/runtime"
	"github.com/wolffcm/flux/values"
)

const versionFuncName = "version"

var errBuildInfoNotPresent = errors.New(codes.NotFound, "build info is not present")

func init() {
	runtime.RegisterPackageValue("runtime", versionFuncName, values.NewFunction(
		versionFuncName,
		runtime.MustLookupBuiltinType("runtime", versionFuncName),
		Version,
		false,
	))
}

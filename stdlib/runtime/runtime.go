package runtime

import (
	"github.com/wolffcm/flux"
	"github.com/wolffcm/flux/codes"
	"github.com/wolffcm/flux/internal/errors"
	"github.com/wolffcm/flux/semantic"
	"github.com/wolffcm/flux/values"
)

const versionFuncName = "version"

var errBuildInfoNotPresent = errors.New(codes.NotFound, "build info is not present")

func init() {
	flux.RegisterPackageValue("runtime", versionFuncName, values.NewFunction(
		versionFuncName,
		semantic.NewFunctionPolyType(semantic.FunctionPolySignature{
			Return: semantic.String,
		}),
		Version,
		false,
	))
}

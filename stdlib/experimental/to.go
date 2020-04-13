package experimental

import (
	"github.com/wolffcm/flux"
	"github.com/wolffcm/flux/runtime"
)

// ToKind is the kind for the experimental `to` flux function
const ExperimentalToKind = "experimental-to"

var ToSignature = runtime.MustLookupBuiltinType("experimental", "to")

func init() {
	runtime.RegisterPackageValue("experimental", "to", flux.MustValue(flux.FunctionValueWithSideEffect("to", nil, ToSignature)))
}

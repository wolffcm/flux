package influxdb

import (
	"github.com/wolffcm/flux"
	"github.com/wolffcm/flux/runtime"
)

// ToKind is the kind for the `to` flux function
const ToKind = "to"

var ToSignature = runtime.MustLookupBuiltinType("influxdata/influxdb", "to")

func init() {
	runtime.RegisterPackageValue("influxdata/influxdb", ToKind, flux.MustValue(flux.FunctionValueWithSideEffect(ToKind, nil, ToSignature)))
}

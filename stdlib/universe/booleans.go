package universe

import (
	"github.com/wolffcm/flux/runtime"
	"github.com/wolffcm/flux/values"
)

func init() {
	runtime.RegisterPackageValue("universe", "true", values.NewBool(true))
	runtime.RegisterPackageValue("universe", "false", values.NewBool(false))
}

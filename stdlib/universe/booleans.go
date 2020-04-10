package universe

import (
	"github.com/wolffcm/flux"
	"github.com/wolffcm/flux/values"
)

func init() {
	flux.RegisterPackageValue("universe", "true", values.NewBool(true))
	flux.RegisterPackageValue("universe", "false", values.NewBool(false))
}

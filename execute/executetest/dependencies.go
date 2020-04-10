package executetest

import (
	"github.com/wolffcm/flux"
	"github.com/wolffcm/flux/dependencies/dependenciestest"
)

func NewTestExecuteDependencies() flux.Dependencies {
	return dependenciestest.Default()
}

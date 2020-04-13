package langtest

import (
	"github.com/wolffcm/flux/lang"
	"github.com/wolffcm/flux/memory"
)

func DefaultExecutionDependencies() lang.ExecutionDependencies {
	return lang.ExecutionDependencies{
		Allocator: new(memory.Allocator),
	}
}

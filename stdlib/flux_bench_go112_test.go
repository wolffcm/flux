// +build !go1.13

package stdlib_test

import (
	"testing"

	"github.com/wolffcm/flux"
)

func reportStatistics(b *testing.B, stats flux.Statistics) {
	// Not supported in go 1.12.
}

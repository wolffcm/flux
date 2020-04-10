package values_test

import (
	"testing"

	"github.com/wolffcm/flux/semantic"
	"github.com/wolffcm/flux/values"
)

func TestArrayEqual(t *testing.T) {
	r := values.NewArray(semantic.Int)
	r.Append(values.NewInt(1))
	l := values.NewArray(semantic.Int)
	l.Append(values.NewInt(1))

	if !l.Equal(r) {
		t.Fatal("expected arrays to be equal")
	}

	l.Set(0, values.NewInt(2))
	if l.Equal(r) {
		t.Fatal("expected arrays to be unequal")
	}

	r.Set(0, values.NewInt(2))
	if !l.Equal(r) {
		t.Fatal("expected objects to be equal")
	}
	l.Append(values.NewInt(1))
	if l.Equal(r) {
		t.Fatal("expected objects to be unequal")
	}
}

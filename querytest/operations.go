package querytest

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/wolffcm/flux"
	"github.com/wolffcm/flux/semantic/semantictest"
	"github.com/wolffcm/flux/stdlib/universe"
)

func OperationMarshalingTestHelper(t *testing.T, data []byte, expOp *flux.Operation) {
	t.Helper()

	opts := append(
		semantictest.CmpOptions,
		cmp.AllowUnexported(universe.JoinOpSpec{}),
		cmpopts.IgnoreUnexported(universe.JoinOpSpec{}),
	)

	// Ensure we can properly unmarshal a spec
	gotOp := new(flux.Operation)
	if err := json.Unmarshal(data, gotOp); err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(gotOp, expOp, opts...) {
		t.Errorf("unexpected operation -want/+got %s", cmp.Diff(expOp, gotOp, opts...))
	}

	// Marshal the spec and ensure we can unmarshal it again.
	data, err := json.Marshal(expOp)
	if err != nil {
		t.Fatal(err)
	}
	gotOp = new(flux.Operation)
	if err := json.Unmarshal(data, gotOp); err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(gotOp, expOp, opts...) {
		t.Errorf("unexpected operation after marshalling -want/+got %s", cmp.Diff(expOp, gotOp, opts...))
	}
}

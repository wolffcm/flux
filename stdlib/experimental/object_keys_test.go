package experimental_test

import (
	"context"
	"testing"

	"github.com/wolffcm/flux/dependencies/dependenciestest"
	"github.com/wolffcm/flux/runtime"
)

func TestObjectKeys(t *testing.T) {
	script := `
import "experimental"
import "internal/testutil"

o = {a: 1, b: 2, c: 3}
experimental.objectKeys(o: o) == ["a", "b", "c"] or testutil.fail()
`
	ctx := dependenciestest.Default().Inject(context.Background())
	if _, _, err := runtime.Eval(ctx, script); err != nil {
		t.Fatal("evaluation of objectKeys failed: ", err)
	}
}

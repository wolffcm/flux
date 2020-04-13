package interptest

import (
	"context"

	"github.com/wolffcm/flux/codes"
	"github.com/wolffcm/flux/internal/errors"
	"github.com/wolffcm/flux/interpreter"
	"github.com/wolffcm/flux/runtime"
	"github.com/wolffcm/flux/values"
)

func Eval(ctx context.Context, itrp *interpreter.Interpreter, scope values.Scope, importer interpreter.Importer, src string) ([]interpreter.SideEffect, error) {
	node, err := runtime.AnalyzeSource(src)
	if err != nil {
		return nil, errors.Wrap(err, codes.Inherit, "could not analyze program")
	}
	return itrp.Eval(ctx, node, scope, importer)
}

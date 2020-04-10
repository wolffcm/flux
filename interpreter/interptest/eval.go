package interptest

import (
	"context"

	"github.com/wolffcm/flux/ast"
	"github.com/wolffcm/flux/interpreter"
	"github.com/wolffcm/flux/parser"
	"github.com/wolffcm/flux/semantic"
	"github.com/wolffcm/flux/values"
)

func Eval(ctx context.Context, itrp *interpreter.Interpreter, scope values.Scope, importer interpreter.Importer, src string) ([]interpreter.SideEffect, error) {
	pkg := parser.ParseSource(src)
	if ast.Check(pkg) > 0 {
		return nil, ast.GetError(pkg)
	}
	node, err := semantic.New(pkg)
	if err != nil {
		return nil, err
	}
	return itrp.Eval(ctx, node, scope, importer)
}

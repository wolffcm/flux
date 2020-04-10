// +build !libflux

package parser

import (
	"fmt"

	"github.com/wolffcm/flux/ast"
	"github.com/wolffcm/flux/internal/parser"
	"github.com/wolffcm/flux/internal/token"
)

func parseFile(f *token.File, src []byte) (*ast.File, error) {
	if useRustParser() {
		panic(fmt.Sprintf(`%v set to %q but this is not a libflux build`, fluxParserTypeEnvVar, parserTypeRust))
	}
	return parser.ParseFile(f, src), nil
}

func isLibfluxBuild() bool {
	return false
}

package fluxfmt

import (
	"fmt"

	"github.com/influxdata/flux/parser"
)

func Format(query string) (string, error) {
	program, err := parser.NewAST(query)
	if err != nil {
		return "", fmt.Errorf("could not parse query: %v", err)
	}

	return program.String(), nil
}

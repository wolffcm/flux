package execute

import (
	"io"

	"github.com/wolffcm/flux"
	"github.com/wolffcm/flux/values"
)

type RowReader interface {
	Next() bool
	GetNextRow() ([]values.Value, error)
	ColumnNames() []string
	ColumnTypes() []flux.ColType
	SetColumns([]interface{})
	io.Closer
}

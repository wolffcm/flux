package llvm

import (
	"errors"
	"fmt"
	"github.com/influxdata/flux/semantic"
	"github.com/llvm-mirror/llvm/bindings/go/llvm"
)

type symbolEntry struct {
	fluxExpr semantic.Node
	llvmValues map[llvm.Type]llvm.Value
}

type symbolTable struct {
	entries map[string]symbolEntry
}

func newSymbolTable() *symbolTable{
	return &symbolTable{
		entries: make(map[string]symbolEntry),
	}
}

func (st *symbolTable) addEntry(name string, fluxExpr semantic.Node, value *llvm.Value) error {
	if se, ok := st.entries[name]; ok {
		if fluxExpr != nil && fluxExpr != se.fluxExpr {
			return fmt.Errorf("found differing flux exprs for symbol %q", name)
		}

		if value != nil {
			vty := value.Type()
			if _, ok := se.llvmValues[vty]; ok {
				return fmt.Errorf("found existing value with type %s for symbol %q", vty, name)
			}

			se.llvmValues[vty] = *value
		}
	} else {
		se.fluxExpr = fluxExpr
		se.llvmValues = make(map[llvm.Type]llvm.Value)
		if value != nil {
			se.llvmValues[value.Type()] = *value
		}
		st.entries[name] = se
	}

	return nil
}

var symbolNotFound = errors.New("could not find symbol")

func (st *symbolTable) getSingleValue(name string) (llvm.Value, error) {
	se, ok := st.entries[name]
	if !ok {
		return llvm.Value{}, symbolNotFound
	}

	if len(se.llvmValues) != 1 {
		return llvm.Value{}, fmt.Errorf("symbol table has %d values for %q", len(se.llvmValues), name)
	}

	var val llvm.Value
	for _, v := range se.llvmValues {
		val = v
	}

	return val, nil
}

func (st *symbolTable) getEntry(name string) *symbolEntry {
	se, ok := st.entries[name]
	if !ok {
		return nil
	}

	return &se
}

func (st *symbolTable) getSpecialization(name string, lty llvm.Type) *llvm.Value {
	se, ok := st.entries[name]
	if ! ok {
		return nil
	}

	v, ok := se.llvmValues[lty]
	if !ok {
		return nil
	}

	return &v
}

func (st *symbolTable) findName(n semantic.Node) string {
	for name, se := range st.entries {
		if se.fluxExpr == n {
			return name
		}
	}
	return ""
}
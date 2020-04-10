package types

import (
	"sort"
	"strings"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/wolffcm/flux/codes"
	"github.com/wolffcm/flux/internal/errors"
	"github.com/wolffcm/flux/semantic/internal/fbsemantic"
)

// PolyType represents a polytype.  This struct is a thin wrapper around
// Go code generated by the FlatBuffers compiler.
type PolyType struct {
	fb *fbsemantic.PolyType
}

// NewPolyType returns a new polytype given a flatbuffers polytype.
func NewPolyType(fb *fbsemantic.PolyType) (*PolyType, error) {
	if fb == nil {
		return nil, errors.New(codes.Internal, "got nil fbsemantic.polytype")
	}
	return &PolyType{fb: fb}, nil
}

// NumVars returns the number of type variables in this polytype.
func (pt *PolyType) NumVars() int {
	return pt.fb.VarsLength()
}

// Var returns the type variable at ordinal position i.
func (pt *PolyType) Var(i int) (*fbsemantic.Var, error) {
	if i < 0 || i >= pt.NumVars() {
		return nil, errors.Newf(codes.Internal, "request for polytype var out of bounds: %v in %v", i, pt.NumVars())
	}
	v := new(fbsemantic.Var)
	if !pt.fb.Vars(v, i) {
		return nil, errors.Newf(codes.Internal, "missing var")
	}
	return v, nil
}

// NumConstraints returns the number of kind constraints in this polytype.
func (pt *PolyType) NumConstraints() int {
	return pt.fb.ConsLength()
}

// Constraint returns the constraint at ordinal position i.
func (pt *PolyType) Constraint(i int) (*fbsemantic.Constraint, error) {
	if i < 0 || i >= pt.NumConstraints() {
		return nil, errors.Newf(codes.Internal, "request for constraint out of bounds: %v in %v", i, pt.NumConstraints())
	}
	c := new(fbsemantic.Constraint)
	if !pt.fb.Cons(c, i) {
		return nil, errors.Newf(codes.Internal, "missing constraint")
	}
	return c, nil

}

// SortedConstraints returns the constraints for this polytype sorted by type variable and constraint kind.
func (pt *PolyType) SortedConstraints() ([]*fbsemantic.Constraint, error) {
	ncs := pt.NumConstraints()
	cs := make([]*fbsemantic.Constraint, ncs)
	for i := 0; i < ncs; i++ {
		c, err := pt.Constraint(i)
		if err != nil {
			return nil, err
		}
		cs[i] = c
	}
	sort.Slice(cs, func(i, j int) bool {
		tvi, tvj := cs[i].Tvar(nil).I(), cs[j].Tvar(nil).I()
		if tvi == tvj {
			return cs[i].Kind() < cs[j].Kind()
		}
		return tvi < tvj
	})
	return cs, nil
}

// Expr returns the monotype expression for this polytype.
func (pt *PolyType) Expr() (*MonoType, error) {
	tbl := new(flatbuffers.Table)
	if !pt.fb.Expr(tbl) {
		return nil, errors.New(codes.Internal, "missing a polytype expr")
	}

	return NewMonoType(tbl, pt.fb.ExprType())
}

func (pt PolyType) SortedVars() ([]*fbsemantic.Var, error) {
	nvars := pt.NumVars()
	vars := make([]*fbsemantic.Var, nvars)
	for i := 0; i < nvars; i++ {
		arg, err := pt.Var(i)
		if err != nil {
			return nil, err
		}
		vars[i] = arg
	}
	sort.Slice(vars, func(i, j int) bool {
		return vars[i].I() < vars[j].I()
	})
	return vars, nil

}

// String returns a string representation for this polytype.
func (pt *PolyType) String() string {
	var sb strings.Builder

	sb.WriteString("forall [")
	needComma := false
	svars, err := pt.SortedVars()
	if err != nil {
		return "<" + err.Error() + ">"
	}
	for _, v := range svars {
		if needComma {
			sb.WriteString(", ")
		} else {
			needComma = true
		}
		mt := monoTypeFromVar(v)
		sb.WriteString(mt.String())
	}
	sb.WriteString("] ")

	needWhere := true
	cs, err := pt.SortedConstraints()
	if err != nil {
		return "<" + err.Error() + ">"
	}
	for i := 0; i < len(cs); i++ {
		cons := cs[i]
		tv := cons.Tvar(nil)
		k := cons.Kind()

		if needWhere {
			sb.WriteString("where ")
			needWhere = false
		}
		mtv := monoTypeFromVar(tv)
		sb.WriteString(mtv.String())
		sb.WriteString(": ")
		sb.WriteString(fbsemantic.EnumNamesKind[k])

		if i < pt.NumConstraints()-1 {
			sb.WriteString(", ")
		} else {
			sb.WriteString(" ")
		}
	}

	mt, err := pt.Expr()
	if err != nil {
		return "<" + err.Error() + ">"
	}
	sb.WriteString(mt.String())

	return sb.String()
}

// GetCanonicalMapping returns a map of type variable numbers to
// canonicalized numbers that start from 0.
// Tests that do type inference will have type variables that are sensitive
// to changes in the standard library, this helps to solve that problem.
func (pt *PolyType) GetCanonicalMapping() (map[uint64]int, error) {
	tvm := make(map[uint64]int)
	counter := 0

	svars, err := pt.SortedVars()
	if err != nil {
		return nil, err
	}
	for _, v := range svars {
		updateTVarMap(&counter, tvm, v.I())
	}

	mt, err := pt.Expr()
	if err != nil {
		return nil, err
	}
	if err := mt.getCanonicalMapping(&counter, tvm); err != nil {
		return nil, err
	}

	return tvm, nil
}

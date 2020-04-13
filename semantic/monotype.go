package semantic

import (
	"fmt"
	"sort"
	"strings"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/wolffcm/flux/codes"
	"github.com/wolffcm/flux/internal/errors"
	"github.com/wolffcm/flux/internal/fbsemantic"
)

type fbTabler interface {
	Init(buf []byte, i flatbuffers.UOffsetT)
	Table() flatbuffers.Table
}

// MonoType represents a monotype.  This struct is a thin wrapper around
// Go code generated by the FlatBuffers compiler.
type MonoType struct {
	mt  fbsemantic.MonoType
	tbl fbTabler
}

// NewMonoType constructs a new monotype from a FlatBuffers table and the given kind of monotype.
func NewMonoType(tbl flatbuffers.Table, t fbsemantic.MonoType) (MonoType, error) {
	var tbler fbTabler
	switch t {
	case fbsemantic.MonoTypeNONE:
		return MonoType{}, nil
	case fbsemantic.MonoTypeBasic:
		tbler = new(fbsemantic.Basic)
	case fbsemantic.MonoTypeVar:
		tbler = new(fbsemantic.Var)
	case fbsemantic.MonoTypeArr:
		tbler = new(fbsemantic.Arr)
	case fbsemantic.MonoTypeRow:
		tbler = new(fbsemantic.Row)
	case fbsemantic.MonoTypeFun:
		tbler = new(fbsemantic.Fun)
	default:
		return MonoType{}, errors.Newf(codes.Internal, "unknown type (%v)", t)
	}
	tbler.Init(tbl.Bytes, tbl.Pos)
	return MonoType{mt: t, tbl: tbler}, nil
}

func (mt MonoType) Nature() Nature {
	switch mt.mt {
	case fbsemantic.MonoTypeBasic:
		t, _ := mt.Basic()
		switch t {
		case fbsemantic.TypeBool:
			return Bool
		case fbsemantic.TypeInt:
			return Int
		case fbsemantic.TypeUint:
			return UInt
		case fbsemantic.TypeFloat:
			return Float
		case fbsemantic.TypeString:
			return String
		case fbsemantic.TypeDuration:
			return Duration
		case fbsemantic.TypeTime:
			return Time
		case fbsemantic.TypeRegexp:
			return Regexp
		case fbsemantic.TypeBytes:
			return Bytes
		default:
			return Invalid
		}
	case fbsemantic.MonoTypeArr:
		return Array
	case fbsemantic.MonoTypeRow:
		return Object
	case fbsemantic.MonoTypeFun:
		return Function
	case fbsemantic.MonoTypeNONE,
		fbsemantic.MonoTypeVar:
		fallthrough
	default:
		return Invalid
	}
}

// Kind specifies a particular kind of monotype.
type Kind fbsemantic.MonoType

const (
	Unknown = Kind(fbsemantic.MonoTypeNONE)
	Basic   = Kind(fbsemantic.MonoTypeBasic)
	Var     = Kind(fbsemantic.MonoTypeVar)
	Arr     = Kind(fbsemantic.MonoTypeArr)
	Row     = Kind(fbsemantic.MonoTypeRow)
	Fun     = Kind(fbsemantic.MonoTypeFun)
)

// Kind returns what kind of monotype the receiver is.
func (mt MonoType) Kind() Kind {
	return Kind(mt.mt)
}

var (
	BasicBool     = newBasicType(fbsemantic.TypeBool)
	BasicInt      = newBasicType(fbsemantic.TypeInt)
	BasicUint     = newBasicType(fbsemantic.TypeUint)
	BasicFloat    = newBasicType(fbsemantic.TypeFloat)
	BasicString   = newBasicType(fbsemantic.TypeString)
	BasicDuration = newBasicType(fbsemantic.TypeDuration)
	BasicTime     = newBasicType(fbsemantic.TypeTime)
	BasicRegexp   = newBasicType(fbsemantic.TypeRegexp)
	BasicBytes    = newBasicType(fbsemantic.TypeBytes)
)

func getBasic(tbl fbTabler) (*fbsemantic.Basic, error) {
	b, ok := tbl.(*fbsemantic.Basic)
	if !ok {
		return nil, errors.New(codes.Internal, "MonoType is not a basic type")
	}
	return b, nil
}

// Basic returns the basic type for this monotype if it is a basic type,
// and an error otherwise.
func (mt MonoType) Basic() (fbsemantic.Type, error) {
	b, err := getBasic(mt.tbl)
	if err != nil {
		return fbsemantic.TypeBool, err
	}
	return b.T(), nil
}

func getVar(tbl fbTabler) (*fbsemantic.Var, error) {
	v, ok := tbl.(*fbsemantic.Var)
	if !ok {
		return nil, errors.New(codes.Internal, "MonoType is not a type var")
	}
	return v, nil

}

// VarNum returns the type variable number if this monotype is a type variable,
// and an error otherwise.
func (mt MonoType) VarNum() (uint64, error) {
	v, err := getVar(mt.tbl)
	if err != nil {
		return 0, err
	}
	return v.I(), nil
}

func monoTypeFromVar(v *fbsemantic.Var) MonoType {
	return MonoType{
		mt:  fbsemantic.MonoTypeVar,
		tbl: v,
	}
}

func getFun(tbl fbTabler) (*fbsemantic.Fun, error) {
	f, ok := tbl.(*fbsemantic.Fun)
	if !ok {
		return nil, errors.New(codes.Internal, "MonoType is not a function")
	}
	return f, nil
}

// NumArguments returns the number of arguments if this monotype is a function,
// and an error otherwise.
func (mt MonoType) NumArguments() (int, error) {
	f, err := getFun(mt.tbl)
	if err != nil {
		return 0, err
	}
	return f.ArgsLength(), nil
}

// Argument returns the argument give an ordinal position if this monotype is a function,
// and an error otherwise.
func (mt MonoType) Argument(i int) (*Argument, error) {
	f, err := getFun(mt.tbl)
	if err != nil {
		return nil, err
	}
	if i < 0 || i >= f.ArgsLength() {
		return nil, errors.Newf(codes.Internal, "request for out-of-bounds argument: %v of %v", i, f.ArgsLength())
	}
	a := new(fbsemantic.Argument)
	if !f.Args(a, i) {
		return nil, errors.New(codes.Internal, "missing argument")
	}
	return newArgument(a)
}

// SortedArguments returns a slice of function arguments,
// sorted by argument name, if this monotype is a function.
func (mt MonoType) SortedArguments() ([]*Argument, error) {
	nargs, err := mt.NumArguments()
	if err != nil {
		return nil, err
	}
	args := make([]*Argument, nargs)
	for i := 0; i < nargs; i++ {
		arg, err := mt.Argument(i)
		if err != nil {
			return nil, err
		}
		args[i] = arg
	}
	sort.Slice(args, func(i, j int) bool {
		return string(args[i].Name()) < string(args[j].Name())
	})
	return args, nil
}

func (mt MonoType) ReturnType() (MonoType, error) {
	f, ok := mt.tbl.(*fbsemantic.Fun)
	if !ok {
		return MonoType{}, errors.New(codes.Internal, "ReturnType() called on non-function MonoType")
	}
	var tbl flatbuffers.Table
	if !f.Retn(&tbl) {
		return MonoType{}, errors.New(codes.Internal, "missing return type")
	}
	return NewMonoType(tbl, f.RetnType())
}

func getArr(tbl fbTabler) (*fbsemantic.Arr, error) {
	arr, ok := tbl.(*fbsemantic.Arr)
	if !ok {
		return nil, errors.New(codes.Internal, "MonoType is not an array")
	}
	return arr, nil
}

// ElemType returns the element type if this monotype is an array, and an error otherise.
func (mt MonoType) ElemType() (MonoType, error) {
	arr, err := getArr(mt.tbl)
	if err != nil {
		return MonoType{}, err
	}
	var tbl flatbuffers.Table
	if !arr.T(&tbl) {
		return MonoType{}, errors.New(codes.Internal, "missing array type")
	}
	return NewMonoType(tbl, arr.TType())
}

func getRow(tbl fbTabler) (*fbsemantic.Row, error) {
	row, ok := tbl.(*fbsemantic.Row)
	if !ok {
		return nil, errors.New(codes.Internal, "MonoType is not a row")
	}
	return row, nil

}

// NumProperties returns the number of properties if this monotype is a row, and an error otherwise.
func (mt MonoType) NumProperties() (int, error) {
	row, err := getRow(mt.tbl)
	if err != nil {
		return 0, err
	}
	return row.PropsLength(), nil
}

// RowProperty returns a property given its ordinal position if this monotype is a row, and an error otherwise.
func (mt MonoType) RowProperty(i int) (*RowProperty, error) {
	row, err := getRow(mt.tbl)
	if err != nil {
		return nil, err
	}
	if i < 0 || i >= row.PropsLength() {
		return nil, errors.Newf(codes.Internal, "request for out-of-bounds property: %v of %v", i, row.PropsLength())
	}
	p := new(fbsemantic.Prop)
	if !row.Props(p, i) {
		return nil, errors.New(codes.Internal, "missing property")
	}
	return &RowProperty{fb: p}, nil
}

// SortedProperties returns the properties for a Row monotype, sorted by
// key.  It's possible that there are duplicate keys with different types,
// in this case, this function preserves their order.
func (mt MonoType) SortedProperties() ([]*RowProperty, error) {
	nps, err := mt.NumProperties()
	if err != nil {
		return nil, err
	}
	ps := make([]*RowProperty, nps)
	for i := 0; i < nps; i++ {
		ps[i], err = mt.RowProperty(i)
		if err != nil {
			return nil, err
		}
	}
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].Name() == ps[j].Name() {
			return i < j
		}
		return ps[i].Name() < ps[j].Name()
	})
	return ps, nil
}

// Extends returns the extending type variable if this monotype is a row, and an error otherwise.
// If the type is a row but does not extend anything a false is returned.
func (mt MonoType) Extends() (MonoType, bool, error) {
	row, err := getRow(mt.tbl)
	if err != nil {
		return MonoType{}, false, err
	}
	v := row.Extends(nil)
	if v == nil {
		return MonoType{}, false, nil
	}
	return monoTypeFromVar(v), true, nil
}

// Argument represents a function argument.
type Argument struct {
	*fbsemantic.Argument
}

func newArgument(fb *fbsemantic.Argument) (*Argument, error) {
	if fb == nil {
		return nil, errors.Newf(codes.Internal, "nil argument")
	}
	return &Argument{Argument: fb}, nil
}

// TypeOf returns the type of the function argument.
func (a *Argument) TypeOf() (MonoType, error) {
	var tbl flatbuffers.Table
	if !a.T(&tbl) {
		return MonoType{}, errors.New(codes.Internal, "missing argument type")
	}
	argTy, err := NewMonoType(tbl, a.TType())
	if err != nil {
		return MonoType{}, err
	}
	return argTy, nil
}

// Property represents a property of a row.
type RowProperty struct {
	fb *fbsemantic.Prop
}

// Name returns the name of the property.
func (p *RowProperty) Name() string {
	return string(p.fb.K())
}

// TypeOf returns the type of the property.
func (p *RowProperty) TypeOf() (MonoType, error) {
	var tbl flatbuffers.Table
	if !p.fb.V(&tbl) {
		return MonoType{}, nil
	}
	return NewMonoType(tbl, p.fb.VType())
}

// String returns a string representation of this monotype.
func (mt MonoType) String() string {
	return mt.string(nil)
}

// CanonicalString returns a string representation of this monotype
// where the tvar numbers are contiguous and indexed starting at zero.
func (mt MonoType) CanonicalString() string {
	ctr := uint64(0)
	m := make(map[uint64]uint64)
	if err := mt.getCanonicalMapping(&ctr, m); err != nil {
		return "<" + err.Error() + ">"
	}
	return mt.string(m)
}

func (mt MonoType) getCanonicalMapping(counter *uint64, tvm map[uint64]uint64) error {
	switch tk := mt.Kind(); tk {
	case Var:
		tv, err := mt.VarNum()
		if err != nil {
			return err
		}
		updateTVarMap(counter, tvm, tv)
	case Arr:
		et, err := mt.ElemType()
		if err != nil {
			return err
		}
		if err := et.getCanonicalMapping(counter, tvm); err != nil {
			return err
		}
	case Row:
		props, err := mt.SortedProperties()
		if err != nil {
			return err
		}
		for _, p := range props {
			pt, err := p.TypeOf()
			if err != nil {
				return err
			}
			if err := pt.getCanonicalMapping(counter, tvm); err != nil {
				return err
			}
		}
		evar, ok, err := mt.Extends()
		if err != nil {
			return err
		} else if ok {
			if err := evar.getCanonicalMapping(counter, tvm); err != nil {
				return err
			}
		}
	case Fun:
		args, err := mt.SortedArguments()
		if err != nil {
			return err
		}
		for _, arg := range args {
			at, err := arg.TypeOf()
			if err != nil {
				return err
			}
			if err := at.getCanonicalMapping(counter, tvm); err != nil {
				return err
			}
		}
		rt, err := mt.ReturnType()
		if err != nil {
			return err
		}
		if err := rt.getCanonicalMapping(counter, tvm); err != nil {
			return err
		}
	}

	return nil
}

func (mt MonoType) string(m map[uint64]uint64) string {
	if mt.tbl == nil {
		return "null"
	}
	switch tk := mt.Kind(); tk {
	case Unknown:
		return "<" + fbsemantic.EnumNamesMonoType[fbsemantic.MonoType(tk)] + ">"
	case Basic:
		b, err := mt.Basic()
		if err != nil {
			return "<" + err.Error() + ">"
		}
		return strings.ToLower(fbsemantic.EnumNamesType[byte(b)])
	case Var:
		i, err := mt.VarNum()
		if err != nil {
			return "<" + err.Error() + ">"
		}
		if m != nil {
			var ok bool
			if i, ok = m[i]; !ok {
				return "<could not find var num in map>"
			}
		}
		return fmt.Sprintf("t%d", i)
	case Arr:
		et, err := mt.ElemType()
		if err != nil {
			return "<" + err.Error() + ">"
		}
		return "[" + et.string(m) + "]"
	case Row:
		var sb strings.Builder
		sb.WriteString("{")
		sprops, err := mt.SortedProperties()
		if err != nil {
			return "<" + err.Error() + ">"
		}
		needBar := false
		for _, prop := range sprops {
			if needBar {
				sb.WriteString(" | ")
			} else {
				needBar = true
			}
			sb.WriteString(prop.Name() + ": ")
			ty, err := prop.TypeOf()
			if err != nil {
				return "<" + err.Error() + ">"
			}
			sb.WriteString(ty.string(m))
		}
		extends, ok, err := mt.Extends()
		if err != nil {
			return "<" + err.Error() + ">"
		}
		if ok {
			if needBar {
				sb.WriteString(" | ")
			}
			sb.WriteString(extends.string(m))
		}
		sb.WriteString("}")
		return sb.String()
	case Fun:
		var sb strings.Builder
		sb.WriteString("(")
		needComma := false
		sargs, err := mt.SortedArguments()
		if err != nil {
			return "<" + err.Error() + ">"
		}
		for _, arg := range sargs {
			if needComma {
				sb.WriteString(", ")
			} else {
				needComma = true
			}
			if arg.Optional() {
				sb.WriteString("?")
			} else if arg.Pipe() {
				sb.WriteString("<-")
			}
			sb.WriteString(string(arg.Name()) + ": ")
			argTyp, err := arg.TypeOf()
			if err != nil {
				return "<" + err.Error() + ">"
			}
			sb.WriteString(argTyp.string(m))
		}
		sb.WriteString(") -> ")
		rt, err := mt.ReturnType()
		if err != nil {
			return "<" + err.Error() + ">"
		}
		sb.WriteString(rt.string(m))
		return sb.String()
	default:
		return "<" + fmt.Sprintf("unknown monotype (%v)", tk) + ">"
	}
}

func (l MonoType) Equal(r MonoType) bool {
	return l.String() == r.String()
}

func newBasicType(t fbsemantic.Type) MonoType {
	builder := flatbuffers.NewBuilder(16)
	offset := buildBasicType(builder, t)
	builder.Finish(offset)

	buf := builder.FinishedBytes()
	basic := fbsemantic.GetRootAsBasic(buf, 0)
	mt, err := NewMonoType(basic.Table(), fbsemantic.MonoTypeBasic)
	if err != nil {
		panic(err)
	}
	return mt
}

// NewArrayType will construct a new Array MonoType
// where the inner element for the array is elemType.
func NewArrayType(elemType MonoType) MonoType {
	builder := flatbuffers.NewBuilder(32)
	offset := buildArrayType(builder, elemType)
	builder.Finish(offset)

	buf := builder.FinishedBytes()
	arr := fbsemantic.GetRootAsArr(buf, 0)
	mt, err := NewMonoType(arr.Table(), fbsemantic.MonoTypeArr)
	if err != nil {
		panic(err)
	}
	return mt
}

type ArgumentType struct {
	Name     []byte
	Type     MonoType
	Pipe     bool
	Optional bool
}

// NewFunctionType will construct a new Function MonoType
// that has a return value that matches retn and arguments
// for each of the values in ArgumentType.
func NewFunctionType(retn MonoType, args []ArgumentType) MonoType {
	builder := flatbuffers.NewBuilder(64)
	offset := buildFunctionType(builder, retn, args)
	builder.Finish(offset)

	buf := builder.FinishedBytes()
	fun := fbsemantic.GetRootAsFun(buf, 0)
	mt, err := NewMonoType(fun.Table(), fbsemantic.MonoTypeFun)
	if err != nil {
		panic(err)
	}
	return mt
}

type PropertyType struct {
	Key   []byte
	Value MonoType
}

// NewObjectType will construct a new Object MonoType with
// the properties in properties.
//
// The MonoType will be constructed with the properties in the
// same order as they appear in the array.
func NewObjectType(properties []PropertyType) MonoType {
	builder := flatbuffers.NewBuilder(64)
	offset := buildObjectType(builder, properties, nil)
	builder.Finish(offset)

	buf := builder.FinishedBytes()
	row := fbsemantic.GetRootAsRow(buf, 0)
	mt, err := NewMonoType(row.Table(), fbsemantic.MonoTypeRow)
	if err != nil {
		panic(err)
	}
	return mt
}

// copyMonoType will reconstruct the type contained within the
// MonoType for the new builder. When building a new buffer,
// flatbuffers cannot reference data in another buffer and the
// flatbuffers types contain references to offsets that are no
// longer valid when copied to a new buffer.
//
// This method will access the existing types so it can correctly
// rebuild an already constructed MonoType inside of another buffer.
func copyMonoType(builder *flatbuffers.Builder, t MonoType) flatbuffers.UOffsetT {
	if t.mt == fbsemantic.MonoTypeNONE {
		return 0
	}

	table := t.tbl.Table()
	switch t.mt {
	case fbsemantic.MonoTypeNONE:
		panic("monotype type not set")
	case fbsemantic.MonoTypeBasic:
		var basic fbsemantic.Basic
		basic.Init(table.Bytes, table.Pos)
		return buildBasicType(builder, basic.T())
	case fbsemantic.MonoTypeVar:
		var tv fbsemantic.Var
		tv.Init(table.Bytes, table.Pos)
		return buildVarType(builder, tv.I())
	case fbsemantic.MonoTypeArr:
		var arr fbsemantic.Arr
		arr.Init(table.Bytes, table.Pos)

		elem := monoTypeFromFunc(arr.T, arr.TType())
		return buildArrayType(builder, elem)
	case fbsemantic.MonoTypeRow:
		var row fbsemantic.Row
		row.Init(table.Bytes, table.Pos)

		properties := make([]PropertyType, row.PropsLength())
		for i := 0; i < len(properties); i++ {
			var prop fbsemantic.Prop
			row.Props(&prop, i)
			properties[i] = PropertyType{
				Key:   prop.K(),
				Value: monoTypeFromFunc(prop.V, prop.VType()),
			}
		}
		extends := row.Extends(nil)
		return buildObjectType(builder, properties, extends)
	case fbsemantic.MonoTypeFun:
		var fun fbsemantic.Fun
		fun.Init(table.Bytes, table.Pos)

		args := make([]ArgumentType, fun.ArgsLength())
		for i := 0; i < len(args); i++ {
			var arg fbsemantic.Argument
			fun.Args(&arg, i)
			args[i] = ArgumentType{
				Name:     arg.Name(),
				Type:     monoTypeFromFunc(arg.T, arg.TType()),
				Pipe:     arg.Pipe(),
				Optional: arg.Optional(),
			}
		}
		retn := monoTypeFromFunc(fun.Retn, fun.RetnType())
		return buildFunctionType(builder, retn, args)
	default:
		panic(fmt.Sprintf("unknown monotype (%v)", t.mt))
	}
}

// monoTypeFromFunc will initialize a MonoType using the table
// initialized from the function. If the property does not exist,
// this will panic.
func monoTypeFromFunc(fn func(obj *flatbuffers.Table) bool, t fbsemantic.MonoType) MonoType {
	var table flatbuffers.Table
	if !fn(&table) {
		return MonoType{}
	}
	mt, err := NewMonoType(table, t)
	if err != nil {
		panic(err)
	}
	return mt
}

// buildBasicType will construct a basic type in the builder
// and return the offset for the type.
func buildBasicType(builder *flatbuffers.Builder, t fbsemantic.Type) flatbuffers.UOffsetT {
	fbsemantic.BasicStart(builder)
	fbsemantic.BasicAddT(builder, t)
	return fbsemantic.BasicEnd(builder)
}

// buildVarType will construct a var type in the builder
// and return the offset for the type.
func buildVarType(builder *flatbuffers.Builder, i uint64) flatbuffers.UOffsetT {
	fbsemantic.VarStart(builder)
	fbsemantic.VarAddI(builder, i)
	return fbsemantic.VarEnd(builder)
}

// buildArrayType will construct an arr type in the builder
// and return the offset for the type.
func buildArrayType(builder *flatbuffers.Builder, elemType MonoType) flatbuffers.UOffsetT {
	offset := copyMonoType(builder, elemType)
	fbsemantic.ArrStart(builder)
	fbsemantic.ArrAddTType(builder, elemType.mt)
	fbsemantic.ArrAddT(builder, offset)
	return fbsemantic.ArrEnd(builder)
}

// buildFunctionType will construct a fun type in the builder
// and return the offset for the type.
func buildFunctionType(builder *flatbuffers.Builder, retn MonoType, args []ArgumentType) flatbuffers.UOffsetT {
	retnOffset := copyMonoType(builder, retn)
	argsOffsets := make([]flatbuffers.UOffsetT, len(args))
	for i, arg := range args {
		nOffset := builder.CreateByteString(arg.Name)
		tOffset := copyMonoType(builder, arg.Type)
		fbsemantic.ArgumentStart(builder)
		fbsemantic.ArgumentAddName(builder, nOffset)
		fbsemantic.ArgumentAddTType(builder, arg.Type.mt)
		fbsemantic.ArgumentAddT(builder, tOffset)
		argsOffsets[i] = fbsemantic.ArgumentEnd(builder)
	}

	fbsemantic.FunStartArgsVector(builder, len(argsOffsets))
	for i := len(argsOffsets) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(argsOffsets[i])
	}
	argsOffset := builder.EndVector(len(argsOffsets))

	fbsemantic.FunStart(builder)
	fbsemantic.FunAddRetnType(builder, retn.mt)
	fbsemantic.FunAddRetn(builder, retnOffset)
	fbsemantic.FunAddArgs(builder, argsOffset)
	return fbsemantic.FunEnd(builder)
}

// buildObjectType will construct a row type in the builder
// and return the offset for the type.
func buildObjectType(builder *flatbuffers.Builder, properties []PropertyType, extends *fbsemantic.Var) flatbuffers.UOffsetT {
	propOffsets := make([]flatbuffers.UOffsetT, len(properties))
	for i, p := range properties {
		kOffset := builder.CreateByteString(p.Key)
		vOffset := copyMonoType(builder, p.Value)
		fbsemantic.PropStart(builder)
		fbsemantic.PropAddK(builder, kOffset)
		if p.Value.mt != fbsemantic.MonoTypeNONE {
			fbsemantic.PropAddVType(builder, p.Value.mt)
			fbsemantic.PropAddV(builder, vOffset)
		}
		propOffsets[i] = fbsemantic.PropEnd(builder)
	}

	var extendsOffset flatbuffers.UOffsetT
	if extends != nil {
		fbsemantic.VarStart(builder)
		fbsemantic.VarAddI(builder, extends.I())
		extendsOffset = fbsemantic.VarEnd(builder)
	}

	fbsemantic.RowStartPropsVector(builder, len(propOffsets))
	for i := len(propOffsets) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(propOffsets[i])
	}
	props := builder.EndVector(len(propOffsets))
	fbsemantic.RowStart(builder)
	fbsemantic.RowAddProps(builder, props)
	if extends != nil {
		fbsemantic.RowAddExtends(builder, extendsOffset)
	}
	return fbsemantic.RowEnd(builder)
}

func updateTVarMap(counter *uint64, m map[uint64]uint64, tv uint64) {
	if _, ok := m[tv]; ok {
		return
	}
	m[tv] = *counter
	*counter++
}

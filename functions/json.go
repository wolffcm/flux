package functions

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/semantic"
	"github.com/influxdata/flux/values"
	"github.com/pkg/errors"
)

func init() {
	flux.RegisterBuiltInValue("parseJSON", &parseJSON{})
}

type parseJSON struct {
}

var parseJSONType = semantic.NewFunctionType(semantic.FunctionSignature{
	Params:     map[string]semantic.Type{"str": semantic.String},
	ReturnType: semantic.Object,
})

func (c *parseJSON) Type() semantic.Type {
	return parseJSONType
}
func (c *parseJSON) Str() string {
	panic(values.UnexpectedKind(semantic.Function, semantic.String))
}
func (c *parseJSON) Int() int64 {
	panic(values.UnexpectedKind(semantic.Function, semantic.Int))
}
func (c *parseJSON) UInt() uint64 {
	panic(values.UnexpectedKind(semantic.Function, semantic.UInt))
}
func (c *parseJSON) Float() float64 {
	panic(values.UnexpectedKind(semantic.Function, semantic.Float))
}
func (c *parseJSON) Bool() bool {
	panic(values.UnexpectedKind(semantic.Function, semantic.Bool))
}
func (c *parseJSON) Time() values.Time {
	panic(values.UnexpectedKind(semantic.Function, semantic.Time))
}
func (c *parseJSON) Duration() values.Duration {
	panic(values.UnexpectedKind(semantic.Function, semantic.Duration))
}
func (c *parseJSON) Regexp() *regexp.Regexp {
	panic(values.UnexpectedKind(semantic.Function, semantic.Regexp))
}
func (c *parseJSON) Array() values.Array {
	panic(values.UnexpectedKind(semantic.Function, semantic.Array))
}
func (c *parseJSON) Object() values.Object {
	panic(values.UnexpectedKind(semantic.Function, semantic.Object))
}
func (c *parseJSON) Function() values.Function {
	return c
}
func (c *parseJSON) Equal(rhs values.Value) bool {
	if c.Type() != rhs.Type() {
		return false
	}
	f, ok := rhs.(*parseJSON)
	return ok && (c == f)
}
func (c *parseJSON) HasSideEffect() bool {
	return false
}

func (c *parseJSON) Call(args values.Object) (values.Value, error) {
	str, ok := args.Get("str")
	if !ok {
		return nil, errors.New(`missing argument "str"`)
	}
	if str.Type() != semantic.String {
		return nil, fmt.Errorf("\"str\" parameter must be a string, got %v", str.Type())
	}
	tmpl, ok := args.Get("template")
	if !ok {
		return nil, errors.New(`missing argument "template"`)
	}
	if tmpl.Type().Kind() != semantic.Object {
		return nil, fmt.Errorf("\"template\" parameter must be an object, got %v", tmpl.Type())
	}

	raw := make(map[string]interface{})
	if err := json.Unmarshal([]byte(str.Str()), &raw); err != nil {
		return nil, err
	}
	o, err := constructValue(tmpl.Object(), raw)
	if err != nil {
		return nil, err
	}
	return o, nil
}

// constructValue creates a new value from raw using the type of the existing value.
func constructValue(existing values.Value, raw interface{}) (values.Value, error) {
	switch k := existing.Type().Kind(); k {
	case semantic.Object:
		obj := existing.Object()
		// TODO(nathanielc): Uncomment these to copy the value instead.
		//nobj := values.NewObject()
		o, ok := raw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("cannot convert %T to object", raw)
		}
		var err error
		obj.Range(func(k string, existing values.Value) {
			if err != nil {
				return
			}
			raw, ok := o[k]
			if !ok {
				// Copy existing
				//nobj.Set(k, existing)
				return
			}
			var v values.Value
			v, err = constructValue(existing, raw)
			if err != nil {
				return
			}
			//nobj.Set(k, v)
			obj.Set(k, v)
		})
		//return nobj, nil
		return obj, nil
	case semantic.Array:
		//TODO(nathanielc): Handle construcing arrays from the existing type information.
		// Should the raw array completly replace the existing array or should they some how be merged.
		// Is an existing array required?
		return nil, errors.New("array parsing not supported")
	default:
		return values.NewValue(raw, k)
	}
}

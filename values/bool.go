package values

import "github.com/wolffcm/flux/semantic"

var (
	trueValue Value = value{
		t: semantic.Bool,
		v: true,
	}
	falseValue Value = value{
		t: semantic.Bool,
		v: false,
	}
)

func NewBool(v bool) Value {
	if v {
		return trueValue
	} else {
		return falseValue
	}
}

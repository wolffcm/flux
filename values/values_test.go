package values_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/wolffcm/flux/semantic"
	"github.com/wolffcm/flux/values"
)

func TestNew(t *testing.T) {
	for _, tt := range []struct {
		v    interface{}
		want values.Value
	}{
		{v: "abc", want: values.NewString("abc")},
		{v: int64(4), want: values.NewInt(4)},
		{v: uint64(4), want: values.NewUInt(4)},
		{v: float64(6.0), want: values.NewFloat(6.0)},
		{v: true, want: values.NewBool(true)},
		{v: values.Time(1000), want: values.NewTime(values.Time(1000))},
		{v: values.ConvertDuration(1), want: values.NewDuration(values.ConvertDuration(1))},
		{v: regexp.MustCompile(`.+`), want: values.NewRegexp(regexp.MustCompile(`.+`))},
	} {
		t.Run(fmt.Sprint(tt.want.Type()), func(t *testing.T) {
			if want, got := tt.want, values.New(tt.v); !want.Equal(got) {
				t.Fatalf("unexpected value -want/+got\n\t- %s\n\t+ %s", want, got)
			}
		})
	}
}

func TestNewNull(t *testing.T) {
	v := values.NewNull(semantic.BasicString)
	if want, got := true, v.IsNull(); want != got {
		t.Fatalf("unexpected value -want/+got\n\t- %v\n\t+ %v", want, got)
	}
}

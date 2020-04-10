package universe_test

import (
	"testing"

	"github.com/wolffcm/flux"
	"github.com/wolffcm/flux/execute"
	"github.com/wolffcm/flux/execute/executetest"
	"github.com/wolffcm/flux/querytest"
	"github.com/wolffcm/flux/stdlib/universe"
)

func TestRelativeStrengthIndex_Marshaling(t *testing.T) {
	data := []byte(`{"id":"relativeStrengthIndex","kind":"relativeStrengthIndex","spec":{"n":1}}`)
	op := &flux.Operation{
		ID: "relativeStrengthIndex",
		Spec: &universe.RelativeStrengthIndexOpSpec{
			N: 1,
		},
	}
	querytest.OperationMarshalingTestHelper(t, data, op)
}

func TestRelativeStrengthIndex_PassThrough(t *testing.T) {
	executetest.TransformationPassThroughTestHelper(t, func(d execute.Dataset, c execute.TableBuilderCache) execute.Transformation {
		s := universe.NewRelativeStrengthIndexTransformation(
			d,
			c,
			&universe.RelativeStrengthIndexProcedureSpec{},
		)
		return s
	})
}

func TestRelativeStrengthIndex_Process(t *testing.T) {
	testCases := []struct {
		name    string
		spec    *universe.RelativeStrengthIndexProcedureSpec
		data    []flux.Table
		want    []*executetest.Table
		wantErr error
	}{
		{
			name: "float",
			spec: &universe.RelativeStrengthIndexProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       10,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), float64(1)},
					{execute.Time(2), float64(2)},
					{execute.Time(3), float64(3)},
					{execute.Time(4), float64(4)},
					{execute.Time(5), float64(5)},
					{execute.Time(6), float64(6)},
					{execute.Time(7), float64(7)},
					{execute.Time(8), float64(8)},
					{execute.Time(9), float64(9)},
					{execute.Time(10), float64(10)},
					{execute.Time(11), float64(11)},
					{execute.Time(12), float64(12)},
					{execute.Time(13), float64(13)},
					{execute.Time(14), float64(14)},
					{execute.Time(15), float64(15)},
					{execute.Time(16), float64(14)},
					{execute.Time(17), float64(13)},
					{execute.Time(18), float64(12)},
					{execute.Time(19), float64(11)},
					{execute.Time(20), float64(10)},
					{execute.Time(21), float64(9)},
					{execute.Time(22), float64(8)},
					{execute.Time(23), float64(7)},
					{execute.Time(24), float64(6)},
					{execute.Time(25), float64(5)},
					{execute.Time(26), float64(4)},
					{execute.Time(27), float64(3)},
					{execute.Time(28), float64(2)},
					{execute.Time(29), float64(1)},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(11), float64(100)},
					{execute.Time(12), float64(100)},
					{execute.Time(13), float64(100)},
					{execute.Time(14), float64(100)},
					{execute.Time(15), float64(100)},
					{execute.Time(16), float64(90)},
					{execute.Time(17), float64(81)},
					{execute.Time(18), float64(72.9)},
					{execute.Time(19), float64(65.61)},
					{execute.Time(20), float64(59.04900000000001)},
					{execute.Time(21), float64(53.144099999999995)},
					{execute.Time(22), float64(47.82969000000001)},
					{execute.Time(23), float64(43.046721)},
					{execute.Time(24), float64(38.74204890000001)},
					{execute.Time(25), float64(34.86784401000001)},
					{execute.Time(26), float64(31.381059609000005)},
					{execute.Time(27), float64(28.242953648100013)},
					{execute.Time(28), float64(25.418658283290014)},
					{execute.Time(29), float64(22.876792454961006)},
				},
			}},
		},
		{
			name: "float with chunks",
			spec: &universe.RelativeStrengthIndexProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       10,
			},
			data: []flux.Table{&executetest.RowWiseTable{
				Table: &executetest.Table{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), float64(1)},
						{execute.Time(2), float64(2)},
						{execute.Time(3), float64(3)},
						{execute.Time(4), float64(4)},
						{execute.Time(5), float64(5)},
						{execute.Time(6), float64(6)},
						{execute.Time(7), float64(7)},
						{execute.Time(8), float64(8)},
						{execute.Time(9), float64(9)},
						{execute.Time(10), float64(10)},
						{execute.Time(11), float64(11)},
						{execute.Time(12), float64(12)},
						{execute.Time(13), float64(13)},
						{execute.Time(14), float64(14)},
						{execute.Time(15), float64(15)},
						{execute.Time(16), float64(14)},
						{execute.Time(17), float64(13)},
						{execute.Time(18), float64(12)},
						{execute.Time(19), float64(11)},
						{execute.Time(20), float64(10)},
						{execute.Time(21), float64(9)},
						{execute.Time(22), float64(8)},
						{execute.Time(23), float64(7)},
						{execute.Time(24), float64(6)},
						{execute.Time(25), float64(5)},
						{execute.Time(26), float64(4)},
						{execute.Time(27), float64(3)},
						{execute.Time(28), float64(2)},
						{execute.Time(29), float64(1)},
					},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(11), float64(100)},
					{execute.Time(12), float64(100)},
					{execute.Time(13), float64(100)},
					{execute.Time(14), float64(100)},
					{execute.Time(15), float64(100)},
					{execute.Time(16), float64(90)},
					{execute.Time(17), float64(81)},
					{execute.Time(18), float64(72.9)},
					{execute.Time(19), float64(65.61)},
					{execute.Time(20), float64(59.04900000000001)},
					{execute.Time(21), float64(53.144099999999995)},
					{execute.Time(22), float64(47.82969000000001)},
					{execute.Time(23), float64(43.046721)},
					{execute.Time(24), float64(38.74204890000001)},
					{execute.Time(25), float64(34.86784401000001)},
					{execute.Time(26), float64(31.381059609000005)},
					{execute.Time(27), float64(28.242953648100013)},
					{execute.Time(28), float64(25.418658283290014)},
					{execute.Time(29), float64(22.876792454961006)},
				},
			}},
		},
		{
			name: "int",
			spec: &universe.RelativeStrengthIndexProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       10,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TInt},
				},
				Data: [][]interface{}{
					{execute.Time(1), int64(1)},
					{execute.Time(2), int64(2)},
					{execute.Time(3), int64(3)},
					{execute.Time(4), int64(4)},
					{execute.Time(5), int64(5)},
					{execute.Time(6), int64(6)},
					{execute.Time(7), int64(7)},
					{execute.Time(8), int64(8)},
					{execute.Time(9), int64(9)},
					{execute.Time(10), int64(10)},
					{execute.Time(11), int64(11)},
					{execute.Time(12), int64(12)},
					{execute.Time(13), int64(13)},
					{execute.Time(14), int64(14)},
					{execute.Time(15), int64(15)},
					{execute.Time(16), int64(14)},
					{execute.Time(17), int64(13)},
					{execute.Time(18), int64(12)},
					{execute.Time(19), int64(11)},
					{execute.Time(20), int64(10)},
					{execute.Time(21), int64(9)},
					{execute.Time(22), int64(8)},
					{execute.Time(23), int64(7)},
					{execute.Time(24), int64(6)},
					{execute.Time(25), int64(5)},
					{execute.Time(26), int64(4)},
					{execute.Time(27), int64(3)},
					{execute.Time(28), int64(2)},
					{execute.Time(29), int64(1)},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(11), float64(100)},
					{execute.Time(12), float64(100)},
					{execute.Time(13), float64(100)},
					{execute.Time(14), float64(100)},
					{execute.Time(15), float64(100)},
					{execute.Time(16), float64(90)},
					{execute.Time(17), float64(81)},
					{execute.Time(18), float64(72.9)},
					{execute.Time(19), float64(65.61)},
					{execute.Time(20), float64(59.04900000000001)},
					{execute.Time(21), float64(53.144099999999995)},
					{execute.Time(22), float64(47.82969000000001)},
					{execute.Time(23), float64(43.046721)},
					{execute.Time(24), float64(38.74204890000001)},
					{execute.Time(25), float64(34.86784401000001)},
					{execute.Time(26), float64(31.381059609000005)},
					{execute.Time(27), float64(28.242953648100013)},
					{execute.Time(28), float64(25.418658283290014)},
					{execute.Time(29), float64(22.876792454961006)},
				},
			}},
		},
		{
			name: "int with chunks",
			spec: &universe.RelativeStrengthIndexProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       10,
			},
			data: []flux.Table{&executetest.RowWiseTable{
				Table: &executetest.Table{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TInt},
					},
					Data: [][]interface{}{
						{execute.Time(1), int64(1)},
						{execute.Time(2), int64(2)},
						{execute.Time(3), int64(3)},
						{execute.Time(4), int64(4)},
						{execute.Time(5), int64(5)},
						{execute.Time(6), int64(6)},
						{execute.Time(7), int64(7)},
						{execute.Time(8), int64(8)},
						{execute.Time(9), int64(9)},
						{execute.Time(10), int64(10)},
						{execute.Time(11), int64(11)},
						{execute.Time(12), int64(12)},
						{execute.Time(13), int64(13)},
						{execute.Time(14), int64(14)},
						{execute.Time(15), int64(15)},
						{execute.Time(16), int64(14)},
						{execute.Time(17), int64(13)},
						{execute.Time(18), int64(12)},
						{execute.Time(19), int64(11)},
						{execute.Time(20), int64(10)},
						{execute.Time(21), int64(9)},
						{execute.Time(22), int64(8)},
						{execute.Time(23), int64(7)},
						{execute.Time(24), int64(6)},
						{execute.Time(25), int64(5)},
						{execute.Time(26), int64(4)},
						{execute.Time(27), int64(3)},
						{execute.Time(28), int64(2)},
						{execute.Time(29), int64(1)},
					},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(11), float64(100)},
					{execute.Time(12), float64(100)},
					{execute.Time(13), float64(100)},
					{execute.Time(14), float64(100)},
					{execute.Time(15), float64(100)},
					{execute.Time(16), float64(90)},
					{execute.Time(17), float64(81)},
					{execute.Time(18), float64(72.9)},
					{execute.Time(19), float64(65.61)},
					{execute.Time(20), float64(59.04900000000001)},
					{execute.Time(21), float64(53.144099999999995)},
					{execute.Time(22), float64(47.82969000000001)},
					{execute.Time(23), float64(43.046721)},
					{execute.Time(24), float64(38.74204890000001)},
					{execute.Time(25), float64(34.86784401000001)},
					{execute.Time(26), float64(31.381059609000005)},
					{execute.Time(27), float64(28.242953648100013)},
					{execute.Time(28), float64(25.418658283290014)},
					{execute.Time(29), float64(22.876792454961006)},
				},
			}},
		},
		{
			name: "uint",
			spec: &universe.RelativeStrengthIndexProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       10,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TUInt},
				},
				Data: [][]interface{}{
					{execute.Time(1), uint64(1)},
					{execute.Time(2), uint64(2)},
					{execute.Time(3), uint64(3)},
					{execute.Time(4), uint64(4)},
					{execute.Time(5), uint64(5)},
					{execute.Time(6), uint64(6)},
					{execute.Time(7), uint64(7)},
					{execute.Time(8), uint64(8)},
					{execute.Time(9), uint64(9)},
					{execute.Time(10), uint64(10)},
					{execute.Time(11), uint64(11)},
					{execute.Time(12), uint64(12)},
					{execute.Time(13), uint64(13)},
					{execute.Time(14), uint64(14)},
					{execute.Time(15), uint64(15)},
					{execute.Time(16), uint64(14)},
					{execute.Time(17), uint64(13)},
					{execute.Time(18), uint64(12)},
					{execute.Time(19), uint64(11)},
					{execute.Time(20), uint64(10)},
					{execute.Time(21), uint64(9)},
					{execute.Time(22), uint64(8)},
					{execute.Time(23), uint64(7)},
					{execute.Time(24), uint64(6)},
					{execute.Time(25), uint64(5)},
					{execute.Time(26), uint64(4)},
					{execute.Time(27), uint64(3)},
					{execute.Time(28), uint64(2)},
					{execute.Time(29), uint64(1)},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(11), float64(100)},
					{execute.Time(12), float64(100)},
					{execute.Time(13), float64(100)},
					{execute.Time(14), float64(100)},
					{execute.Time(15), float64(100)},
					{execute.Time(16), float64(90)},
					{execute.Time(17), float64(81)},
					{execute.Time(18), float64(72.9)},
					{execute.Time(19), float64(65.61)},
					{execute.Time(20), float64(59.04900000000001)},
					{execute.Time(21), float64(53.144099999999995)},
					{execute.Time(22), float64(47.82969000000001)},
					{execute.Time(23), float64(43.046721)},
					{execute.Time(24), float64(38.74204890000001)},
					{execute.Time(25), float64(34.86784401000001)},
					{execute.Time(26), float64(31.381059609000005)},
					{execute.Time(27), float64(28.242953648100013)},
					{execute.Time(28), float64(25.418658283290014)},
					{execute.Time(29), float64(22.876792454961006)},
				},
			}},
		},
		{
			name: "uint with chunks",
			spec: &universe.RelativeStrengthIndexProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       10,
			},
			data: []flux.Table{&executetest.RowWiseTable{
				Table: &executetest.Table{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TUInt},
					},
					Data: [][]interface{}{
						{execute.Time(1), uint64(1)},
						{execute.Time(2), uint64(2)},
						{execute.Time(3), uint64(3)},
						{execute.Time(4), uint64(4)},
						{execute.Time(5), uint64(5)},
						{execute.Time(6), uint64(6)},
						{execute.Time(7), uint64(7)},
						{execute.Time(8), uint64(8)},
						{execute.Time(9), uint64(9)},
						{execute.Time(10), uint64(10)},
						{execute.Time(11), uint64(11)},
						{execute.Time(12), uint64(12)},
						{execute.Time(13), uint64(13)},
						{execute.Time(14), uint64(14)},
						{execute.Time(15), uint64(15)},
						{execute.Time(16), uint64(14)},
						{execute.Time(17), uint64(13)},
						{execute.Time(18), uint64(12)},
						{execute.Time(19), uint64(11)},
						{execute.Time(20), uint64(10)},
						{execute.Time(21), uint64(9)},
						{execute.Time(22), uint64(8)},
						{execute.Time(23), uint64(7)},
						{execute.Time(24), uint64(6)},
						{execute.Time(25), uint64(5)},
						{execute.Time(26), uint64(4)},
						{execute.Time(27), uint64(3)},
						{execute.Time(28), uint64(2)},
						{execute.Time(29), uint64(1)},
					},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(11), float64(100)},
					{execute.Time(12), float64(100)},
					{execute.Time(13), float64(100)},
					{execute.Time(14), float64(100)},
					{execute.Time(15), float64(100)},
					{execute.Time(16), float64(90)},
					{execute.Time(17), float64(81)},
					{execute.Time(18), float64(72.9)},
					{execute.Time(19), float64(65.61)},
					{execute.Time(20), float64(59.04900000000001)},
					{execute.Time(21), float64(53.144099999999995)},
					{execute.Time(22), float64(47.82969000000001)},
					{execute.Time(23), float64(43.046721)},
					{execute.Time(24), float64(38.74204890000001)},
					{execute.Time(25), float64(34.86784401000001)},
					{execute.Time(26), float64(31.381059609000005)},
					{execute.Time(27), float64(28.242953648100013)},
					{execute.Time(28), float64(25.418658283290014)},
					{execute.Time(29), float64(22.876792454961006)},
				},
			}},
		},
		{
			name: "pass through",
			spec: &universe.RelativeStrengthIndexProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       10,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TInt},
					{Label: "a", Type: flux.TBool},
					{Label: "b", Type: flux.TString},
					{Label: "c", Type: flux.TInt},
				},
				Data: [][]interface{}{
					{execute.Time(1), int64(1), true, "a", int64(1)},
					{execute.Time(2), int64(2), false, "b", int64(2)},
					{execute.Time(3), int64(3), false, "c", int64(3)},
					{execute.Time(4), int64(4), true, "d", int64(4)},
					{execute.Time(5), int64(5), false, "e", int64(5)},
					{execute.Time(6), int64(6), true, "f", int64(6)},
					{execute.Time(7), int64(7), true, "g", int64(7)},
					{execute.Time(8), int64(8), true, "h", int64(8)},
					{execute.Time(9), int64(9), false, "i", int64(9)},
					{execute.Time(10), int64(10), false, "j", int64(10)},
					{execute.Time(11), int64(11), true, "k", int64(11)},
					{execute.Time(12), int64(12), false, "l", int64(12)},
					{execute.Time(13), int64(13), false, "m", int64(13)},
					{execute.Time(14), int64(14), false, "n", int64(14)},
					{execute.Time(15), int64(15), false, "o", int64(15)},
					{execute.Time(16), int64(14), true, "p", int64(14)},
					{execute.Time(17), int64(13), true, "q", int64(13)},
					{execute.Time(18), int64(12), true, "r", int64(12)},
					{execute.Time(19), int64(11), false, "s", int64(11)},
					{execute.Time(20), int64(10), false, "t", int64(10)},
					{execute.Time(21), int64(9), false, "u", int64(9)},
					{execute.Time(22), int64(8), true, "v", int64(8)},
					{execute.Time(23), int64(7), false, "w", int64(7)},
					{execute.Time(24), int64(6), false, "x", int64(6)},
					{execute.Time(25), int64(5), true, "y", int64(5)},
					{execute.Time(26), int64(4), true, "z", int64(4)},
					{execute.Time(27), int64(3), false, "aa", int64(3)},
					{execute.Time(28), int64(2), true, "ab", int64(2)},
					{execute.Time(29), int64(1), false, "ac", int64(1)},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
					{Label: "a", Type: flux.TBool},
					{Label: "b", Type: flux.TString},
					{Label: "c", Type: flux.TInt},
				},
				Data: [][]interface{}{
					{execute.Time(11), float64(100), true, "k", int64(11)},
					{execute.Time(12), float64(100), false, "l", int64(12)},
					{execute.Time(13), float64(100), false, "m", int64(13)},
					{execute.Time(14), float64(100), false, "n", int64(14)},
					{execute.Time(15), float64(100), false, "o", int64(15)},
					{execute.Time(16), float64(90), true, "p", int64(14)},
					{execute.Time(17), float64(81), true, "q", int64(13)},
					{execute.Time(18), float64(72.9), true, "r", int64(12)},
					{execute.Time(19), float64(65.61), false, "s", int64(11)},
					{execute.Time(20), float64(59.04900000000001), false, "t", int64(10)},
					{execute.Time(21), float64(53.144099999999995), false, "u", int64(9)},
					{execute.Time(22), float64(47.82969000000001), true, "v", int64(8)},
					{execute.Time(23), float64(43.046721), false, "w", int64(7)},
					{execute.Time(24), float64(38.74204890000001), false, "x", int64(6)},
					{execute.Time(25), float64(34.86784401000001), true, "y", int64(5)},
					{execute.Time(26), float64(31.381059609000005), true, "z", int64(4)},
					{execute.Time(27), float64(28.242953648100013), false, "aa", int64(3)},
					{execute.Time(28), float64(25.418658283290014), true, "ab", int64(2)},
					{execute.Time(29), float64(22.876792454961006), false, "ac", int64(1)},
				},
			}},
		},
		{
			name: "pass through with chunks",
			spec: &universe.RelativeStrengthIndexProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       10,
			},
			data: []flux.Table{&executetest.RowWiseTable{
				Table: &executetest.Table{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TInt},
						{Label: "a", Type: flux.TBool},
						{Label: "b", Type: flux.TString},
						{Label: "c", Type: flux.TInt},
					},
					Data: [][]interface{}{
						{execute.Time(1), int64(1), true, "a", int64(1)},
						{execute.Time(2), int64(2), false, "b", int64(2)},
						{execute.Time(3), int64(3), false, "c", int64(3)},
						{execute.Time(4), int64(4), true, "d", int64(4)},
						{execute.Time(5), int64(5), false, "e", int64(5)},
						{execute.Time(6), int64(6), true, "f", int64(6)},
						{execute.Time(7), int64(7), true, "g", int64(7)},
						{execute.Time(8), int64(8), true, "h", int64(8)},
						{execute.Time(9), int64(9), false, "i", int64(9)},
						{execute.Time(10), int64(10), false, "j", int64(10)},
						{execute.Time(11), int64(11), true, "k", int64(11)},
						{execute.Time(12), int64(12), false, "l", int64(12)},
						{execute.Time(13), int64(13), false, "m", int64(13)},
						{execute.Time(14), int64(14), false, "n", int64(14)},
						{execute.Time(15), int64(15), false, "o", int64(15)},
						{execute.Time(16), int64(14), true, "p", int64(14)},
						{execute.Time(17), int64(13), true, "q", int64(13)},
						{execute.Time(18), int64(12), true, "r", int64(12)},
						{execute.Time(19), int64(11), false, "s", int64(11)},
						{execute.Time(20), int64(10), false, "t", int64(10)},
						{execute.Time(21), int64(9), false, "u", int64(9)},
						{execute.Time(22), int64(8), true, "v", int64(8)},
						{execute.Time(23), int64(7), false, "w", int64(7)},
						{execute.Time(24), int64(6), false, "x", int64(6)},
						{execute.Time(25), int64(5), true, "y", int64(5)},
						{execute.Time(26), int64(4), true, "z", int64(4)},
						{execute.Time(27), int64(3), false, "aa", int64(3)},
						{execute.Time(28), int64(2), true, "ab", int64(2)},
						{execute.Time(29), int64(1), false, "ac", int64(1)},
					},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
					{Label: "a", Type: flux.TBool},
					{Label: "b", Type: flux.TString},
					{Label: "c", Type: flux.TInt},
				},
				Data: [][]interface{}{
					{execute.Time(11), float64(100), true, "k", int64(11)},
					{execute.Time(12), float64(100), false, "l", int64(12)},
					{execute.Time(13), float64(100), false, "m", int64(13)},
					{execute.Time(14), float64(100), false, "n", int64(14)},
					{execute.Time(15), float64(100), false, "o", int64(15)},
					{execute.Time(16), float64(90), true, "p", int64(14)},
					{execute.Time(17), float64(81), true, "q", int64(13)},
					{execute.Time(18), float64(72.9), true, "r", int64(12)},
					{execute.Time(19), float64(65.61), false, "s", int64(11)},
					{execute.Time(20), float64(59.04900000000001), false, "t", int64(10)},
					{execute.Time(21), float64(53.144099999999995), false, "u", int64(9)},
					{execute.Time(22), float64(47.82969000000001), true, "v", int64(8)},
					{execute.Time(23), float64(43.046721), false, "w", int64(7)},
					{execute.Time(24), float64(38.74204890000001), false, "x", int64(6)},
					{execute.Time(25), float64(34.86784401000001), true, "y", int64(5)},
					{execute.Time(26), float64(31.381059609000005), true, "z", int64(4)},
					{execute.Time(27), float64(28.242953648100013), false, "aa", int64(3)},
					{execute.Time(28), float64(25.418658283290014), true, "ab", int64(2)},
					{execute.Time(29), float64(22.876792454961006), false, "ac", int64(1)},
				},
			}},
		},
		{
			name: "nulls",
			spec: &universe.RelativeStrengthIndexProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       10,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TInt},
				},
				Data: [][]interface{}{
					{execute.Time(0), nil},
					{execute.Time(1), int64(1)},
					{execute.Time(2), int64(2)},
					{execute.Time(3), int64(3)},
					{execute.Time(4), int64(4)},
					{execute.Time(5), int64(5)},
					{execute.Time(6), int64(6)},
					{execute.Time(7), int64(7)},
					{execute.Time(8), int64(8)},
					{execute.Time(9), int64(9)},
					{execute.Time(10), int64(10)},
					{execute.Time(11), int64(11)},
					{execute.Time(12), int64(12)},
					{execute.Time(13), nil},
					{execute.Time(14), int64(13)},
					{execute.Time(15), int64(14)},
					{execute.Time(16), int64(15)},
					{execute.Time(17), int64(14)},
					{execute.Time(18), int64(13)},
					{execute.Time(19), int64(12)},
					{execute.Time(20), int64(11)},
					{execute.Time(21), int64(10)},
					{execute.Time(22), int64(9)},
					{execute.Time(23), int64(8)},
					{execute.Time(24), int64(7)},
					{execute.Time(25), int64(6)},
					{execute.Time(26), nil},
					{execute.Time(27), int64(5)},
					{execute.Time(28), int64(4)},
					{execute.Time(29), int64(3)},
					{execute.Time(30), int64(2)},
					{execute.Time(31), int64(1)},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(10), float64(100)},
					{execute.Time(11), float64(100)},
					{execute.Time(12), float64(100)},
					{execute.Time(13), float64(100)},
					{execute.Time(14), float64(100)},
					{execute.Time(15), float64(100)},
					{execute.Time(16), float64(100)},
					{execute.Time(17), float64(90)},
					{execute.Time(18), float64(81)},
					{execute.Time(19), float64(72.9)},
					{execute.Time(20), float64(65.61)},
					{execute.Time(21), float64(59.04900000000001)},
					{execute.Time(22), float64(53.144099999999995)},
					{execute.Time(23), float64(47.82969000000001)},
					{execute.Time(24), float64(43.046721)},
					{execute.Time(25), float64(38.74204890000001)},
					{execute.Time(26), float64(38.74204890000001)},
					{execute.Time(27), float64(34.86784401000001)},
					{execute.Time(28), float64(31.381059609000005)},
					{execute.Time(29), float64(28.242953648100013)},
					{execute.Time(30), float64(25.418658283290014)},
					{execute.Time(31), float64(22.876792454961006)},
				},
			}},
		},
		{
			name: "nulls with chunks",
			spec: &universe.RelativeStrengthIndexProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       10,
			},
			data: []flux.Table{&executetest.RowWiseTable{
				Table: &executetest.Table{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TInt},
					},
					Data: [][]interface{}{
						{execute.Time(0), nil},
						{execute.Time(1), int64(1)},
						{execute.Time(2), int64(2)},
						{execute.Time(3), int64(3)},
						{execute.Time(4), int64(4)},
						{execute.Time(5), int64(5)},
						{execute.Time(6), int64(6)},
						{execute.Time(7), int64(7)},
						{execute.Time(8), int64(8)},
						{execute.Time(9), int64(9)},
						{execute.Time(10), int64(10)},
						{execute.Time(11), int64(11)},
						{execute.Time(12), int64(12)},
						{execute.Time(13), nil},
						{execute.Time(14), int64(13)},
						{execute.Time(15), int64(14)},
						{execute.Time(16), int64(15)},
						{execute.Time(17), int64(14)},
						{execute.Time(18), int64(13)},
						{execute.Time(19), int64(12)},
						{execute.Time(20), int64(11)},
						{execute.Time(21), int64(10)},
						{execute.Time(22), int64(9)},
						{execute.Time(23), int64(8)},
						{execute.Time(24), int64(7)},
						{execute.Time(25), int64(6)},
						{execute.Time(26), nil},
						{execute.Time(27), int64(5)},
						{execute.Time(28), int64(4)},
						{execute.Time(29), int64(3)},
						{execute.Time(30), int64(2)},
						{execute.Time(31), int64(1)},
					},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(10), float64(100)},
					{execute.Time(11), float64(100)},
					{execute.Time(12), float64(100)},
					{execute.Time(13), float64(100)},
					{execute.Time(14), float64(100)},
					{execute.Time(15), float64(100)},
					{execute.Time(16), float64(100)},
					{execute.Time(17), float64(90)},
					{execute.Time(18), float64(81)},
					{execute.Time(19), float64(72.9)},
					{execute.Time(20), float64(65.61)},
					{execute.Time(21), float64(59.04900000000001)},
					{execute.Time(22), float64(53.144099999999995)},
					{execute.Time(23), float64(47.82969000000001)},
					{execute.Time(24), float64(43.046721)},
					{execute.Time(25), float64(38.74204890000001)},
					{execute.Time(26), float64(38.74204890000001)},
					{execute.Time(27), float64(34.86784401000001)},
					{execute.Time(28), float64(31.381059609000005)},
					{execute.Time(29), float64(28.242953648100013)},
					{execute.Time(30), float64(25.418658283290014)},
					{execute.Time(31), float64(22.876792454961006)},
				},
			}},
		},
		{
			name: "less rows than period",
			spec: &universe.RelativeStrengthIndexProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       6,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TInt},
				},
				Data: [][]interface{}{
					{execute.Time(1), int64(1)},
					{execute.Time(2), int64(2)},
					{execute.Time(3), int64(3)},
					{execute.Time(4), int64(4)},
					{execute.Time(5), int64(5)},
					{execute.Time(6), int64(4)},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(6), float64(83.33333333333334)},
				},
			}},
		},
		{
			name: "less rows than period with chunks",
			spec: &universe.RelativeStrengthIndexProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       6,
			},
			data: []flux.Table{&executetest.RowWiseTable{
				Table: &executetest.Table{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TInt},
					},
					Data: [][]interface{}{
						{execute.Time(1), int64(1)},
						{execute.Time(2), int64(2)},
						{execute.Time(3), int64(3)},
						{execute.Time(4), int64(4)},
						{execute.Time(5), int64(5)},
						{execute.Time(6), int64(4)},
					},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(6), float64(83.33333333333334)},
				},
			}},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			executetest.ProcessTestHelper(
				t,
				tc.data,
				tc.want,
				tc.wantErr,
				func(d execute.Dataset, c execute.TableBuilderCache) execute.Transformation {
					return universe.NewRelativeStrengthIndexTransformation(d, c, tc.spec)
				})
		})
	}
}

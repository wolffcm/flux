package generate

import (
	"time"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/compiler"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/values"
)

type Source struct {
	done  bool
	Start time.Time
	Stop  time.Time
	Count int64
	alloc *execute.Allocator
	Fn    compiler.Func
	Param string
}

func NewSource(a *execute.Allocator) *Source {
	return &Source{alloc: a}
}

func (s *Source) Connect() error {
	return nil
}

func (s *Source) Fetch() (bool, error) {
	return !s.done, nil
}

func (s *Source) Decode() (flux.Table, error) {
	defer func() {
		s.done = true
	}()
	ks := []flux.ColMeta{
		flux.ColMeta{
			Label: "_start",
			Type:  flux.TTime,
		},
		flux.ColMeta{
			Label: "_stop",
			Type:  flux.TTime,
		},
	}
	vs := []values.Value{
		values.NewTimeValue(values.ConvertTime(s.Start)),
		values.NewTimeValue(values.ConvertTime(s.Stop)),
	}
	groupKey := execute.NewGroupKey(ks, vs)
	b := execute.NewColListTableBuilder(groupKey, s.alloc)

	cols := []flux.ColMeta{
		flux.ColMeta{
			Label: "_time",
			Type:  flux.TTime,
		},
		flux.ColMeta{
			Label: "_value",
			Type:  flux.TInt,
		},
	}

	for _, col := range cols {
		b.AddCol(col)
	}

	cols = b.Cols()

	colIndex := map[string]int{}
	for i, col := range cols {
		colIndex[col.Label] = i
	}

	deltaT := s.Stop.Sub(s.Start) / time.Duration(s.Count)
	timeIdx := execute.ColIdx("_time", cols)
	valueIdx := execute.ColIdx("_value", cols)
	for i := 0; i < int(s.Count); i++ {
		b.AppendTime(timeIdx, values.ConvertTime(s.Start.Add(time.Duration(i)*deltaT)))
		scope := map[string]values.Value{s.Param: values.NewIntValue(int64(i))}
		v, err := s.Fn.EvalInt(scope)
		if err != nil {
			return nil, err
		}
		b.AppendInt(valueIdx, v)
	}

	return b.Table()
}

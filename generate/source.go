package generate

import (
	"time"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/values"
)

type Source struct {
	called bool
	Start  time.Time
	Stop   time.Time
	Count  int64
	alloc  *execute.Allocator
}

func NewSource(a *execute.Allocator) *Source {
	return &Source{alloc: a}
}

func (s *Source) Connect() error {
	return nil
}

func (s *Source) Fetch() (bool, error) {
	return !s.called, nil
}

func (s *Source) Decode() (flux.Table, error) {
	defer func() {
		s.called = true
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
	for i := 0; i < int(s.Count); i++ {
		b.AppendTime(colIndex["_time"], values.ConvertTime(s.Start.Add(time.Duration(i)*deltaT)))
		b.AppendInt(colIndex["_value"], int64(i))
	}

	return b.Table()
}

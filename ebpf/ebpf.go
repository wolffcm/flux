package ebpf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/values"

	bpf "github.com/iovisor/gobpf/bcc"
)

type Source struct {
	called   bool
	File     string
	Duration int64
	alloc    *execute.Allocator
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

	groupKey := execute.NewGroupKey([]flux.ColMeta{}, []values.Value{})
	b := execute.NewColListTableBuilder(groupKey, s.alloc)

	cols := []flux.ColMeta{
		flux.ColMeta{
			Label: "_time",
			Type:  flux.TTime,
		},
		flux.ColMeta{
			Label: "_value",
			Type:  flux.TString,
		},
	}

	for _, col := range cols {
		b.AddCol(col)
	}

	err := s.buildBpfResult(b)
	if err != nil {
		return nil, err
	}

	return b.Table()
}

type readlineEvent struct {
	Pid uint32
	Str [80]byte
}

func (s *Source) buildBpfResult(b *execute.ColListTableBuilder) error {
	if s.Duration < 1 {
		s.Duration = 20
	}
	sourceBytes, err := ioutil.ReadFile(s.File)
	if err != nil {
		return fmt.Errorf("Failed to read file: %s", err)
	}

	m := bpf.NewModule(string(sourceBytes), []string{})
	defer m.Close()

	readlineUretprobe, err := m.LoadUprobe("get_return_value")
	if err != nil {
		return fmt.Errorf("Failed to load get_return_value: %s", err)
	}

	err = m.AttachUretprobe("/bin/bash", "readline", readlineUretprobe, -1)
	if err != nil {
		return fmt.Errorf("Failed to attach return_value: %s", err)
	}

	table := bpf.NewTable(m.TableId("readline_events"), m)

	channel := make(chan []byte)

	perfMap, err := bpf.InitPerfMap(table, channel)
	if err != nil {
		return fmt.Errorf("Failed to init perf map: %s", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	colIndex := map[string]int{}
	for i, col := range b.Cols() {
		colIndex[col.Label] = i
	}

	go func() {
		var event readlineEvent
		for {
			data := <-channel
			err := binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &event)
			if err != nil {
				fmt.Printf("failed to decode received data: %s\n", err)
				continue
			}
			// Convert C string (null-terminated) to Go string
			comm := string(event.Str[:bytes.IndexByte(event.Str[:], 0)])
			b.AppendTime(colIndex["_time"], values.ConvertTime(time.Now()))
			b.AppendString(colIndex["_value"], comm)
		}
	}()
	fmt.Println("THIKNGA", s.Duration)

	perfMap.Start()
	<-time.After(time.Second * time.Duration(s.Duration))
	perfMap.Stop()

	return nil
}

package ebpf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/values"

	bpf "github.com/iovisor/gobpf/bcc"
)

type Source struct {
	called bool
	File   string
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
			Label: "_file",
			Type:  flux.TString,
		},
	}
	vs := []values.Value{
		values.NewStringValue(s.File),
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
			Type:  flux.TString,
		},
	}

	for _, col := range cols {
		b.AddCol(col)
	}

	buildBpfResult(b)

	return b.Table()
}

// temporary program here
const source string = `
	#include <uapi/linux/ptrace.h>
	
	struct readline_event_t {
	        u32 pid;
	        char str[80];
	} __attribute__((packed));
	
	BPF_PERF_OUTPUT(readline_events);
	
	int get_return_value(struct pt_regs *ctx) {
	        struct readline_event_t event = {};
	        u32 pid;
	        if (!PT_REGS_RC(ctx))
	                return 0;
	        pid = bpf_get_current_pid_tgid();
	        event.pid = pid;
	        bpf_probe_read(&event.str, sizeof(event.str), (void *)PT_REGS_RC(ctx));
	        readline_events.perf_submit(ctx, &event, sizeof(event));
	
	        return 0;
	}
	`

type readlineEvent struct {
	Pid uint32
	Str [80]byte
}

func buildBpfResult(b *execute.ColListTableBuilder) {
	m := bpf.NewModule(source, []string{})
	defer m.Close()

	readlineUretprobe, err := m.LoadUprobe("get_return_value")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load get_return_value: %s\n", err)
		os.Exit(1)
	}

	err = m.AttachUretprobe("/bin/bash", "readline", readlineUretprobe, -1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to attach return_value: %s\n", err)
		os.Exit(1)
	}

	table := bpf.NewTable(m.TableId("readline_events"), m)

	channel := make(chan []byte)

	perfMap, err := bpf.InitPerfMap(table, channel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init perf map: %s\n", err)
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	colIndex := map[string]int{}
	for i, col := range b.Cols() {
		colIndex[col.Label] = i
	}

	fmt.Printf("%10s\t%s\n", "PID", "COMMAND")
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

	perfMap.Start()
	<-time.After(time.Second * 10)
	perfMap.Stop()
}

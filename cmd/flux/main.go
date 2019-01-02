package main

import (
	"context"
	"log"
	"math"

	"github.com/influxdata/flux"
	_ "github.com/influxdata/flux/builtin"
	"github.com/influxdata/flux/control"
	"github.com/influxdata/flux/repl"
)

func main() {
	querier := NewQuerier()
	r := repl.New(querier)
	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}

type Querier struct {
	c *control.Controller
}

func (q *Querier) Query(ctx context.Context, c flux.Compiler) (flux.ResultIterator, error) {
	qry, err := q.c.Query(ctx, c)
	if err != nil {
		return nil, err
	}
	results := flux.NewResultIteratorFromQuery(qry)
	return results, nil
}

func NewQuerier() *Querier {
	config := control.Config{
		ConcurrencyQuota: 1,
		MemoryBytesQuota: math.MaxInt64,
	}

	c := control.New(config)

	return &Querier{
		c: c,
	}
}

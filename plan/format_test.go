package plan_test

import (
	"fmt"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/wolffcm/flux/execute/executetest"
	"github.com/wolffcm/flux/interpreter"
	"github.com/wolffcm/flux/plan"
	"github.com/wolffcm/flux/plan/plantest"
	"github.com/wolffcm/flux/stdlib/influxdata/influxdb"
	"github.com/wolffcm/flux/stdlib/universe"
)

func TestFormatted(t *testing.T) {
	fromSpec := &influxdb.FromProcedureSpec{
		Bucket: influxdb.NameOrID{Name: "my-bucket"},
	}

	// (r) => r._value > 5.0
	filterSpec := &universe.FilterProcedureSpec{
		Fn: interpreter.ResolvedFunction{
			Fn: executetest.FunctionExpression(t, `(r) => r._value > 5.0`),
		},
	}

	type testcase struct {
		name string
		plan *plantest.PlanSpec
		want string
	}

	tcs := []testcase{
		{
			name: "from |> filter",
			plan: &plantest.PlanSpec{
				Nodes: []plan.Node{
					plan.CreateLogicalNode("from", fromSpec),
					plan.CreateLogicalNode("filter", filterSpec),
				},
				Edges: [][2]int{
					{0, 1},
				},
			},
			want: `digraph {
  from
  filter
  // r._value > 5.000000

  from -> filter
}
`,
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ps := plantest.CreatePlanSpec(tc.plan)
			got := fmt.Sprintf("%v", plan.Formatted(ps, plan.WithDetails()))
			if tc.want != got {
				t.Fatalf("unexpected output: -want/+got:\n%v", diff.LineDiff(tc.want, got))
			}
		})
	}
}

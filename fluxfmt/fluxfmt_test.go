package fluxfmt_test

import (
	"testing"

	"github.com/influxdata/flux/fluxfmt"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		query    string
		expected string
	}{
		{
			query: `from(bucket: "greetings") |> range(start: -1m) |> filter(fn: (r) => r._measurement == "howdy")`,
			expected: `from(bucket: "greetings")
  |> range(start: -1m)
  |> filter(fn: (r) => r._measurement == "howdy")`,
		},
	}

	for _, tt := range tests {
		actual, err := fluxfmt.Format(tt.query)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if actual != tt.expected {
			t.Fatalf("expected:\n```\n%v\n```\nactual:\n```\n%v\n```\n", tt.expected, actual)
		}
	}
}

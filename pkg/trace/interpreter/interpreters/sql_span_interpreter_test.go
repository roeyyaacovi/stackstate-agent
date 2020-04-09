package interpreters

import (
	"github.com/StackVista/stackstate-agent/pkg/trace/interpreter/config"
	"github.com/StackVista/stackstate-agent/pkg/trace/pb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSQLSpanInterpreter(t *testing.T) {
	sqlInterpreter := MakeSQLSpanInterpreter(config.DefaultInterpreterConfig())
	for _, tc := range []struct {
		testCase    string
		interpreter *SQLSpanInterpreter
		span        pb.Span
		expected    pb.Span
	}{
		{
			testCase:    "Should set span.serviceType to 'database' when no db.type metadata exists",
			interpreter: sqlInterpreter,
			span:        pb.Span{Service: "SpanServiceName"},
			expected:    pb.Span{Service: "SpanServiceName", Meta: map[string]string{"span.serviceType": "database"}},
		},
		{
			testCase:    "Should set span.serviceType to 'postgresql' when the db.type is 'postgresql'",
			interpreter: sqlInterpreter,
			span:        pb.Span{Service: "SpanServiceName", Meta: map[string]string{"db.type": "postgresql"}},
			expected:    pb.Span{Service: "SpanServiceName", Meta: map[string]string{"db.type": "postgresql", "span.serviceType": "postgresql"}},
		},
	} {
		t.Run(tc.testCase, func(t *testing.T) {
			actual := tc.interpreter.Interpret(&tc.span)
			assert.EqualValues(t, tc.expected, *actual)
		})
	}
}

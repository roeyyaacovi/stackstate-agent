package interpreters

import (
	"github.com/StackVista/stackstate-agent/pkg/trace/interpreter/config"
	"github.com/StackVista/stackstate-agent/pkg/trace/pb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessSpanInterpreter(t *testing.T) {
	processInterpreter := MakeProcessSpanInterpreter(config.DefaultInterpreterConfig())
	for _, tc := range []struct {
		testCase    string
		interpreter *ProcessSpanInterpreter
		span        pb.Span
		expected    pb.Span
	}{
		{
			testCase:    "Should set span.serviceType to 'service' when no language metadata exists",
			interpreter: processInterpreter,
			span:        pb.Span{Service: "SpanServiceName"},
			expected:    pb.Span{Service: "SpanServiceName", Meta: map[string]string{"span.serviceType": "service"}},
		},
		{
			testCase:    "Should set span.serviceType to 'process' when an unknown language is detected",
			interpreter: processInterpreter,
			span:        pb.Span{Service: "SpanServiceName", Meta: map[string]string{"language": "unknown"}},
			expected:    pb.Span{Service: "SpanServiceName", Meta: map[string]string{"language": "unknown", "span.serviceType": "process"}},
		},
		{
			testCase:    "Should set span.serviceType to 'java' when the language is 'jvm'",
			interpreter: processInterpreter,
			span:        pb.Span{Service: "SpanServiceName", Meta: map[string]string{"language": "jvm"}},
			expected:    pb.Span{Service: "SpanServiceName", Meta: map[string]string{"language": "jvm", "span.serviceType": "java"}},
		},
	} {
		t.Run(tc.testCase, func(t *testing.T) {
			actual := tc.interpreter.Interpret(&tc.span)
			assert.EqualValues(t, tc.expected, *actual)
		})
	}
}

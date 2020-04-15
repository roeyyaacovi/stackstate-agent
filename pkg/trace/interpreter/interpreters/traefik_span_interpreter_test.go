package interpreters

import (
	"github.com/StackVista/stackstate-agent/pkg/trace/interpreter/config"
	"github.com/StackVista/stackstate-agent/pkg/trace/pb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTraefikSpanInterpreter(t *testing.T) {
	traefikInterpreter := MakeTraefikInterpreter(config.DefaultInterpreterConfig())
	for _, tc := range []struct {
		testCase    string
		interpreter *TraefikInterpreter
		span        pb.Span
		expected    pb.Span
	}{
		{
			testCase:    "Should set span.serviceType to 'traefik' when no span.kind metadata exists",
			interpreter: traefikInterpreter,
			span:        pb.Span{Service: "SpanServiceName"},
			expected:    pb.Span{Service: "SpanServiceName", Meta: map[string]string{"span.serviceType": "traefik"}},
		},
		{
			testCase:    "Should set name and service to 'http.host' when span.kind is 'server'",
			interpreter: traefikInterpreter,
			span:        pb.Span{Service: "SpanServiceName", Meta: map[string]string{"http.host": "hostname", "span.kind": "server"}},
			expected:    pb.Span{Name: "hostname", Service: "hostname", Meta: map[string]string{"http.host": "hostname", "span.kind": "server", "span.serviceType": "traefik"}},
		},
		{
			testCase:    "Should set name and service to 'http.host' when span.kind is 'client'",
			interpreter: traefikInterpreter,
			span:        pb.Span{Service: "SpanServiceName", Meta: map[string]string{"backend.name": "backend-service-name", "span.kind": "client"}},
			expected:    pb.Span{Name: "service-name", Service: "service-name", Meta: map[string]string{"backend.name": "backend-service-name", "span.kind": "client", "span.serviceType": "traefik"}},
		},
	} {
		t.Run(tc.testCase, func(t *testing.T) {
			actual := tc.interpreter.Interpret(&tc.span)
			assert.EqualValues(t, tc.expected, *actual)
		})
	}
}

package interpreters

import (
	"github.com/StackVista/stackstate-agent/pkg/trace/interpreter/config"
	"github.com/StackVista/stackstate-agent/pkg/trace/interpreter/util"
	"github.com/StackVista/stackstate-agent/pkg/trace/pb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessSpanInterpreter(t *testing.T) {
	processInterpreter := MakeProcessSpanInterpreter(config.DefaultInterpreterConfig())
	for _, tc := range []struct {
		testCase    string
		interpreter *ProcessSpanInterpreter
		span        util.SpanWithMeta
		expected    pb.Span
	}{
		{
			testCase:    "Should set span.serviceType to 'service' when no language metadata exists",
			interpreter: processInterpreter,
			span:        util.SpanWithMeta{
				Span: &pb.Span{
					Name: "SpanServiceName",
					Service: "SpanServiceName",
				},
				SpanMetadata: &util.SpanMetadata{
					CreateTime: 1586441095,
					Hostname: "hostname",
					PID: 10,
					Type: "web",
					Kind: "some-kind",
				},
			},
			expected:    pb.Span{
				Name: "SpanServiceName",
				Service: "SpanServiceName",
				Meta: map[string]string{
					"span.serviceInstanceIdentifier":"urn:service-instance:/SpanServiceName:/hostname:10:1586441095",
					"span.serviceType": "service",
				},
			},
		},
		{
			testCase:    "Should set span.serviceType to 'process' when an unknown language is detected",
			interpreter: processInterpreter,
			span:        util.SpanWithMeta{
				Span: &pb.Span{
					Name: "SpanServiceName",
					Service: "SpanServiceName",
					Meta: map[string]string{
						"language": "unknown",
					},
				},
				SpanMetadata: &util.SpanMetadata{
					CreateTime: 1586441095,
					Hostname: "hostname",
					PID: 10,
					Type: "web",
					Kind: "some-kind",
				},
			},
			expected:    pb.Span{
				Name: "SpanServiceName",
				Service: "SpanServiceName",
				Meta: map[string]string{
					"span.serviceInstanceIdentifier":"urn:service-instance:/SpanServiceName:/hostname:10:1586441095",
					"language": "unknown", "span.serviceType": "process",
				},
			},
		},
		{
			testCase:    "Should set span.serviceType to 'java' when the language is 'jvm'",
			interpreter: processInterpreter,
			span:        util.SpanWithMeta{
				Span: &pb.Span{
					Name: "SpanServiceName",
					Service: "SpanServiceName",
					Meta: map[string]string{
						"language": "jvm",
					},
				},
				SpanMetadata: &util.SpanMetadata{
					CreateTime: 1586441095,
					Hostname: "hostname",
					PID: 10,
					Type: "web",
					Kind: "some-kind",
				},
			},
			expected:    pb.Span{
				Name: "SpanServiceName",
				Service: "SpanServiceName",
				Meta: map[string]string{
					"span.serviceInstanceIdentifier":"urn:service-instance:/SpanServiceName:/hostname:10:1586441095",
					"language": "jvm", "span.serviceType": "java",
				},
			},
		},
	} {
		t.Run(tc.testCase, func(t *testing.T) {
			actual := tc.interpreter.Interpret(&tc.span)
			assert.EqualValues(t, tc.expected, *actual)
		})
	}
}

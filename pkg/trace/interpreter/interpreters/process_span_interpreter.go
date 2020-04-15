package interpreters

import (
	"github.com/StackVista/stackstate-agent/pkg/trace/interpreter/config"
	"github.com/StackVista/stackstate-agent/pkg/trace/pb"
)

// ProcessSpanInterpreter sets up the process span interpreter
type ProcessSpanInterpreter struct {
	interpreter
}

// ServiceTypeName returns the default service type
const ServiceTypeName = "service"

// ProcessSpanInterpreterName is the name used for matching this interpreter
const ProcessSpanInterpreterName = "process"

// ProcessTypeName returns the default process type
const ProcessTypeName = "process"

// MakeProcessSpanInterpreter creates an instance of the process span interpreter
func MakeProcessSpanInterpreter(config *config.Config) *ProcessSpanInterpreter {
	return &ProcessSpanInterpreter{interpreter{Config: config}}
}

// Interpret performs the interpretation for the ProcessSpanInterpreter
func (in *ProcessSpanInterpreter) Interpret(span *pb.Span) *pb.Span {
	serviceType := ServiceTypeName

	// no meta, add a empty map
	if span.Meta == nil {
		span.Meta = map[string]string{}
	}

	if language, found := span.Meta["language"]; found {
		serviceType = in.LanguageToComponentType(language)
	}
	span.Meta["span.serviceType"] = serviceType

	return span
}

// LanguageToComponentType converts a trace language to a component type
func (in *ProcessSpanInterpreter) LanguageToComponentType(spanLanguage string) string {
	switch spanLanguage {
	case "jvm":
		return "java"
	default:
		return ProcessTypeName
	}
}

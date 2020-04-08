package interpreter

import "github.com/StackVista/stackstate-agent/pkg/trace/pb"

type ProcessSpanInterpreter struct {
	interpreter
}

func MakeProcessSpanInterpreter(config *Config) *ProcessSpanInterpreter {
	return &ProcessSpanInterpreter{interpreter{Config: config}}
}

func (in *ProcessSpanInterpreter) interpret(span *pb.Span) *pb.Span {
	serviceType := SERVICE_TYPE_NAME
	if language, found := span.Meta["language"]; found {
		serviceType	= in.LanguageToComponentType(language)
	}
	span.Meta["span.serviceType"] = serviceType

	return span
}

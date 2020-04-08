package interpreter

import "github.com/StackVista/stackstate-agent/pkg/trace/pb"

type DefaultSpanInterpreter struct {
	interpreter
}

func MakeDefaultSpanInterpreter(config *Config) *DefaultSpanInterpreter {
	return &DefaultSpanInterpreter{interpreter{Config: config}}
}

func (in *DefaultSpanInterpreter) interpret(span *pb.Span) *pb.Span {
	span.Name = in.ServiceName(span)
	return span
}


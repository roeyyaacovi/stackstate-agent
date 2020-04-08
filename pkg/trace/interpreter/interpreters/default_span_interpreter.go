package interpreters

import (
	"github.com/StackVista/stackstate-agent/pkg/trace/interpreter/config"
	"github.com/StackVista/stackstate-agent/pkg/trace/pb"
)

type DefaultSpanInterpreter struct {
	interpreter
}

func MakeDefaultSpanInterpreter(config *config.Config) *DefaultSpanInterpreter {
	return &DefaultSpanInterpreter{interpreter{Config: config}}
}

func (in *DefaultSpanInterpreter) Interpret(span *pb.Span) *pb.Span {
	span.Name = in.ServiceName(span)
	return span
}

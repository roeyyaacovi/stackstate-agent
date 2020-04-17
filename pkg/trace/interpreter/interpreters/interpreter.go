package interpreters

import (
	"github.com/StackVista/stackstate-agent/pkg/trace/interpreter/config"
	"github.com/StackVista/stackstate-agent/pkg/trace/interpreter/util"
	"github.com/StackVista/stackstate-agent/pkg/trace/pb"
)

// SourceInterpreter provides the interface for the different source interpreters
type SourceInterpreter interface {
	Interpret(span *pb.Span) *pb.Span
}

// TypeInterpreter provides the interface for the different type interpreters
type TypeInterpreter interface {
	Interpret(span *util.SpanWithMeta) *pb.Span
}

type interpreter struct {
	Config *config.Config
}

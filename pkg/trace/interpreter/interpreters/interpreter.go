package interpreters

import (
	"github.com/StackVista/stackstate-agent/pkg/trace/interpreter/config"
	"github.com/StackVista/stackstate-agent/pkg/trace/pb"
)

// Interpreter provides the interface for the different interpreters
type Interpreter interface {
	Interpret(span *pb.Span) *pb.Span
}

type interpreter struct {
	Config *config.Config
}

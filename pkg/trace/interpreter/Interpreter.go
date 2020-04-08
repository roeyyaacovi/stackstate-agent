package interpreter

import (
	"fmt"
	"github.com/StackVista/stackstate-agent/pkg/trace/pb"
	"strings"
)

const SERVICE_TYPE_NAME = "service"
const PROCESS_TYPE_NAME = "process"

type Interpreter interface {
	interpret(span *pb.Span) *pb.Span
}

type interpreter struct {
	Config *Config
}

// Calculates a Service Name for this span given the interpreter config
func (in *interpreter) ServiceName(span *pb.Span) string {
	serviceNameSet := make([]string, 0)
	for _, identifier := range in.Config.ServiceIdentifiers {
		if identifierValue, found := span.Meta[identifier]; found {
			serviceNameSet = append(serviceNameSet, identifierValue)
		}
	}

	if len(serviceNameSet)> 0 {
		return fmt.Sprintf("%s:%s", span.Service, strings.Join(serviceNameSet, ":"))
	} else {
		return span.Service
	}
}

func (in *interpreter) LanguageToComponentType(spanLanguage string) string {
	switch spanLanguage {
	case "jvm":
		return "java"
	default:
		return PROCESS_TYPE_NAME
	}
}

package interpreter

import (
	"github.com/StackVista/stackstate-agent/pkg/trace/config"
	interpreterConfig "github.com/StackVista/stackstate-agent/pkg/trace/interpreter/config"
	"github.com/StackVista/stackstate-agent/pkg/trace/interpreter/interpreters"
	"github.com/StackVista/stackstate-agent/pkg/trace/pb"
)

type SpanInterpreterEngine struct {
	SpanInterpreterEngineContext
	DefaultSpanInterpreter *interpreters.DefaultSpanInterpreter
	Interpreters           map[string]interpreters.Interpreter
}

// MakeSpanInterpreterEngine creates a SpanInterpreterEngine given the config and interpreters
func MakeSpanInterpreterEngine(config *interpreterConfig.Config, ins map[string]interpreters.Interpreter) *SpanInterpreterEngine {
	return &SpanInterpreterEngine{
		SpanInterpreterEngineContext: MakeSpanInterpreterEngineContext(config),
		DefaultSpanInterpreter:       interpreters.MakeDefaultSpanInterpreter(config),
		Interpreters:                 ins,
	}
}

// NewSpanInterpreterEngine creates a SpanInterpreterEngine given the config and the default interpreters
func NewSpanInterpreterEngine(agentConfig *config.AgentConfig) *SpanInterpreterEngine {
	interpreterConfig := agentConfig.InterpreterConfig
	ins := make(map[string]interpreters.Interpreter, 0)
	ins[interpreters.PROCESS_SPAN_INTERPRETER] = interpreters.MakeProcessSpanInterpreter(interpreterConfig)
	ins[interpreters.SQL_SPAN_INTERPRETER] = interpreters.MakeSQLSpanInterpreter(interpreterConfig)
	ins[interpreters.TRAEFIK_SPAN_INTERPRETER] = interpreters.MakeTraefikInterpreter(interpreterConfig)

	return MakeSpanInterpreterEngine(interpreterConfig, ins)
}

// Interprets spans using the configured SpanInterpreterEngine
func (se *SpanInterpreterEngine) Interpret(span *pb.Span) *pb.Span {
	span = se.DefaultSpanInterpreter.Interpret(span)

	meta, err := se.extractSpanMetadata(span)
	// no metadata, let's look for the span's source.
	if err != nil {
		if source, found := span.Meta["source"]; found {
			// interpret the source if we have a interpreter.
			if interpreter, found := se.Interpreters[source]; found {
				span = interpreter.Interpret(span)
			}
		}
	} else {
		// process different span types

		// interpret the type if we have a interpreter, otherwise run it through the process interpreter.
		if interpreter, found := se.Interpreters[meta.Type]; found {
			span = interpreter.Interpret(span)
		} else {
			span = se.Interpreters["process"].Interpret(span)
		}
	}
	// we mutate so we return the "same" span
	return span
}

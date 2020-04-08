package interpreter

import (
	"github.com/StackVista/stackstate-agent/pkg/trace/config"
	"github.com/StackVista/stackstate-agent/pkg/trace/pb"
)

type SpanInterpreterEngine struct {
	SpanInterpreterEngineContext
	DefaultSpanInterpreter DefaultSpanInterpreter
	Interpreters map[string]Interpreter
}

func MakeSpanIntepreterEngine(config *Config, interpreters map[string]Interpreter) *SpanInterpreterEngine {
	return &SpanInterpreterEngine{
		SpanInterpreterEngineContext: MakeSpanInterpreterEngineContext(config),
		Interpreters: interpreters,
	}
}

func NewSpanIntepreterEngine(agentConfig *config.AgentConfig) *SpanInterpreterEngine {
	interpreterConfig := agentConfig.InterpreterConfig
	interpreters := make(map[string]Interpreter, 0)
	interpreters[PROCESS_SPAN_INTERPRETER] = MakeProcessSpanInterpreter(interpreterConfig)
	interpreters[SQL_SPAN_INTERPRETER] = MakeSQLSpanInterpreter(interpreterConfig)
	interpreters[TRAEFIK_SPAN_INTERPRETER] = MakeTraefikInterpreter(interpreterConfig)

	return MakeSpanIntepreterEngine(interpreterConfig, interpreters)
}

func (se *SpanInterpreterEngine) Interpret(span *pb.Span) *pb.Span {
	span = se.DefaultSpanInterpreter.interpret(span)

	meta, err := se.extractSpanMetadata(span)
	// no metadata, let's look for the span's source.
	if err != nil {
		if source, found := span.Meta["source"]; found {
			// interpret the source if we have a interpreter.
			if interpreter, found := se.Interpreters[source]; found {
				span = interpreter.interpret(span)
			}
		}
	} else {
		// process different span types

		// interpret the type if we have a interpreter, otherwise run it through the process interpreter.
		if interpreter, found := se.Interpreters[meta.Type]; found {
			span = interpreter.interpret(span)
		} else {
			span = se.Interpreters["process"].interpret(span)
		}
	}
	// we mutate so we return the "same" span
	return span
}

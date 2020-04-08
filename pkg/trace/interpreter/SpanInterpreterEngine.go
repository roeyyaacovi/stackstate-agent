package interpreter

import "github.com/StackVista/stackstate-agent/pkg/trace/pb"

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

func (se *SpanInterpreterEngine) Interpret(trace pb.Trace) pb.Trace {
	for _, span := range trace {
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
	}
	// we mutate so we return the "same" trace
	return trace
}

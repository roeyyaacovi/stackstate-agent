package interpreter

import "github.com/StackVista/stackstate-agent/pkg/trace/pb"

type SQLSpanInterpreter struct {
	interpreter
}

const SQL_SPAN_INTERPRETER = "sql"

func MakeSQLSpanInterpreter(config *Config) *SQLSpanInterpreter {
	return &SQLSpanInterpreter{interpreter{Config: config}}
}

func (in *SQLSpanInterpreter) interpret(span *pb.Span) *pb.Span {
	dbType := "database"
	if database, found := span.Meta["db.type"]; found {
		dbType	= database
	}
	span.Meta["span.serviceType"] = dbType

	return span
}

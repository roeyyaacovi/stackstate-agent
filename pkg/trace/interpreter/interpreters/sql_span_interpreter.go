package interpreters

import (
	"github.com/StackVista/stackstate-agent/pkg/trace/interpreter/config"
	"github.com/StackVista/stackstate-agent/pkg/trace/pb"
)

type SQLSpanInterpreter struct {
	interpreter
}

const SQL_SPAN_INTERPRETER = "sql"

func MakeSQLSpanInterpreter(config *config.Config) *SQLSpanInterpreter {
	return &SQLSpanInterpreter{interpreter{Config: config}}
}

func (in *SQLSpanInterpreter) Interpret(span *pb.Span) *pb.Span {
	dbType := "database"
	if database, found := span.Meta["db.type"]; found {
		dbType = database
	}
	span.Meta["span.serviceType"] = dbType

	return span
}

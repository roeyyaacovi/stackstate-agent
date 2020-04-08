package interpreters

import (
	"github.com/StackVista/stackstate-agent/pkg/trace/interpreter/config"
	"github.com/StackVista/stackstate-agent/pkg/trace/pb"
	"strings"
)

type TraefikInterpreter struct {
	interpreter
}

const TRAEFIK_SPAN_INTERPRETER = "traefik"

func MakeTraefikInterpreter(config *config.Config) *TraefikInterpreter {
	return &TraefikInterpreter{interpreter{Config: config}}
}

func (in *TraefikInterpreter) Interpret(span *pb.Span) *pb.Span {
	if kind, found := span.Meta["span.kind"]; found {
		switch kind {
		case "server":
			// this is the calling service, take the host as Name and Service
			// e.g. urn:service:/api-service-router.staging.furby.ps
			if host, found := span.Meta["http.host"]; found {
				span.Name = host
				span.Service = host
			}
		case "client":
			// this is the called service, takes the backend.name as identifier
			// e.g. "backend-stackstate-books-app" -> urn:service:/stackstate-books-app
			if backendName, found := span.Meta["backend.name"]; found {
				backendName = strings.TrimPrefix(backendName, "backend-")
				span.Name = backendName
				span.Service = backendName
			}
		}
	}

	span.Meta["span.serviceType"] = "traefik"

	return span
}

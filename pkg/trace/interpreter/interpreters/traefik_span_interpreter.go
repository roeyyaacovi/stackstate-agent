package interpreters

import (
	"fmt"
	"github.com/StackVista/stackstate-agent/pkg/trace/interpreter/config"
	"github.com/StackVista/stackstate-agent/pkg/trace/interpreter/util"
	"github.com/StackVista/stackstate-agent/pkg/trace/pb"
	"net/url"
	"strings"
)

// TraefikInterpreter sets up the default span interpreter
type TraefikInterpreter struct {
	interpreter
}

// TraefikSpanInterpreterSpan is the name used for matching this interpreter
const TraefikSpanInterpreterSpan = "traefik"

// MakeTraefikInterpreter creates an instance of the traefik span interpreter
func MakeTraefikInterpreter(config *config.Config) *TraefikInterpreter {
	return &TraefikInterpreter{interpreter{Config: config}}
}

// Interpret performs the interpretation for the TraefikInterpreter
func (in *TraefikInterpreter) Interpret(span *pb.Span) *pb.Span {

	// no meta, add a empty map
	if span.Meta == nil {
		span.Meta = map[string]string{}
	}

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

			// create the service identifier using the already interpreted name and the meta "http.url"
			if urlString, found := span.Meta["http.url"]; found {
				url, err := url.Parse(urlString)
				if err == nil {
					span.Meta["span.serviceInstanceIdentifier"] = createTraefikServiceInstanceURN(span.Name, url.Host)
				}
			}

		}
	}

	span.Meta["span.serviceType"] = "traefik"

	return span
}

// createTraefikServiceInstanceURN creates the urn identifier for all service instance components
func createTraefikServiceInstanceURN(serviceName string, hostname string) string {
	return fmt.Sprintf("urn:%s:/%s:/%s", util.ServiceInstanceTypeName, serviceName, hostname)
}

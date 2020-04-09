package interpreter

import (
	"fmt"
	"github.com/StackVista/stackstate-agent/pkg/trace/interpreter/config"
	"github.com/StackVista/stackstate-agent/pkg/trace/interpreter/util"
	"github.com/StackVista/stackstate-agent/pkg/trace/pb"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

// SpanInterpreterEngineContext helper functions that is used by the span interpreter engine context
type SpanInterpreterEngineContext interface {
	nanosToMillis(nanos int64) int64
	extractSpanMetadata(span *pb.Span) (*util.SpanMetadata, error)
}

type spanInterpreterEngineContext struct {
	Config *config.Config
}

// MakeSpanInterpreterEngineContext creates a SpanInterpreterEngineContext for config
func MakeSpanInterpreterEngineContext(config *config.Config) SpanInterpreterEngineContext {
	return &spanInterpreterEngineContext{Config: config}
}

func (c *spanInterpreterEngineContext) nanosToMillis(nanos int64) int64 {
	return nanos / int64(time.Millisecond)
}

func (c *spanInterpreterEngineContext) extractSpanMetadata(span *pb.Span) (*util.SpanMetadata, error) {

	var hostname string
	var createTime int64
	var pid int
	var kind string
	var found bool

	if hostname, found = span.Meta[c.Config.ExtractionFields.HostnameField]; !found {
		return nil, createSpanMetadataError(c.Config.ExtractionFields.HostnameField)
	}

	if pidStr, found := span.Meta[c.Config.ExtractionFields.PidField]; found {
		p, err := strconv.Atoi(pidStr)
		if err != nil {
			return nil, err
		}
		pid = p
	} else {
		return nil, createSpanMetadataError(c.Config.ExtractionFields.PidField)
	}

	if kind, found = span.Meta[c.Config.ExtractionFields.KindField]; !found {
		return nil, createSpanMetadataError(c.Config.ExtractionFields.KindField)
	}

	// try to get the create time, otherwise default to span start
	if createTimeStr, found := span.Meta[c.Config.ExtractionFields.CreateTimeField]; found {
		ct, err := strconv.ParseInt(createTimeStr, 10, 64)
		if err != nil {
			return nil, err
		}
		createTime = ct
	} else {
		createTime = c.nanosToMillis(span.Start)
	}

	return &util.SpanMetadata{
		CreateTime: createTime,
		Hostname:   hostname,
		PID:        pid,
		Type:       span.Type,
		Kind:       kind,
	}, nil
}

func createSpanMetadataError(configField string) error {
	return errors.New(fmt.Sprintf("Field [%s] not found in Span", configField))
}

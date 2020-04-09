package util

// SpanMetadata contains the fields of the span meta that we are interested in
type SpanMetadata struct {
	CreateTime int64
	Hostname   string
	PID        int
	Type       string
	Kind       string
}

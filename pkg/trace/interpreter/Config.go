package interpreter

// InterpreterConfig holds the configuration that allows the span interpreter
// to interpret and enrich various span types.
type Config struct {
	ServiceIdentifiers []string         `mapstructure:"service_identifiers"`
	ExtractionFields   ExtractionFields `mapstructure:"extraction_fields"`
}

type ExtractionFields struct {
	CreateTimeField string `mapstructure:"create_time"`
	HostnameField string `mapstructure:"host_name"`
	PidField string `mapstructure:"pid"`
	KindField string `mapstructure:"kind"`
}

func DefaultInterpreterConfig() *Config {
	return &Config{
		ServiceIdentifiers: []string{"db.instance"},
		ExtractionFields: ExtractionFields{
			CreateTimeField: "span.starttime",
			HostnameField: "span.hostname",
			PidField: "span.pid",
			KindField: "span.kind",
		},
	}
}

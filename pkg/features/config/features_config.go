package config

import "time"

// FeaturesConfig contains the configuration to customize the behaviour of the Features functionality.
type FeaturesConfig struct {
	// HTTPRequestTimeoutSecs is the HTTP timeout for POST requests to the StackState backend
	HTTPRequestTimeoutSecs time.Duration
	FeatureRequestTicker   *time.Ticker
	MaxRetries             int
}

// DefaultFeaturesConfig creates a new instance of a FeaturesConfig using default values.
func DefaultFeaturesConfig() FeaturesConfig {
	return FeaturesConfig{
		HTTPRequestTimeoutSecs: 10 * time.Second,
		FeatureRequestTicker:   time.NewTicker(5 * time.Second),
		MaxRetries:             5,
	}
}

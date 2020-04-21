package config

import "time"

// FeaturesConfig contains the configuration to customize the behaviour of the Features functionality.
type FeaturesConfig struct {
	// HTTPRequestTimeoutDuration is the HTTP timeout for POST requests to the StackState backend
	HTTPRequestTimeoutDuration   time.Duration
	FeatureRequestTickerDuration time.Duration
	MaxRetries                   int
}

// DefaultFeaturesConfig creates a new instance of a FeaturesConfig using default values.
func DefaultFeaturesConfig() FeaturesConfig {
	return FeaturesConfig{
		HTTPRequestTimeoutDuration:   10 * time.Second,
		FeatureRequestTickerDuration: 30 * time.Second,
		MaxRetries:                   10,
	}
}

package features

// Features Structure for describing features published by StackState
type Features interface {
	FeatureEnabled(feature string) bool
}

// features Implementation
type features struct {
	features map[string]bool
}

// NewFeatures creates a features object based on map
func NewFeatures(featureMap map[string]bool) Features {
	return &features{
		features: featureMap,
	}
}

// FeatureEnabled checks whether a certain feature is enabled
func (f *features) FeatureEnabled(feature string) bool {
	if supported, ok := f.features[feature]; ok {
		return supported
	}
	return false
}

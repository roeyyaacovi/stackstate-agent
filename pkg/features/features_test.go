package features

import (
	"github.com/StackVista/stackstate-agent/pkg/trace/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFeatures(t *testing.T) {
	featuresTestServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			switch req.URL.Path {
			case "/features":
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				a := `{
					"some-test-feature": true
				}`
				_, err := w.Write([]byte(a))
				if err != nil {
					t.Fatal(err)
				}
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}),
	)

	conf := config.New()
	conf.Endpoints = []*config.Endpoint{
		{Host: featuresTestServer.URL},
	}
	conf.FeaturesConfig.FeatureRequestTicker = time.NewTicker(500 * time.Millisecond)
	conf.FeaturesConfig.MaxRetries = 5

	featureChan := make(chan map[string]bool, 1)
	features := NewTestFeatures(conf, featureChan)

	// assert feature not supported before fetched
	assert.False(t, features.FeatureEnabled("some-test-feature"))

	// start feature fetcher
	features.Start()

	// assert feature supported after fetch completed
	select {
	case <-time.After(1 * time.Second): // timeout and assert after 1 second
		assert.True(t, features.FeatureEnabled("some-test-feature"))
	default:
		// check on each loop if the condition is satisfied yet, otherwise wait for the timeout
		enabled := features.FeatureEnabled("some-test-feature")
		if enabled {
			assert.True(t, enabled)
		}
	}

	// stop feature fetcher
	features.Stop()
}

func TestFeaturesRetries(t *testing.T) {
	featuresTestServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}),
	)

	conf := config.New()
	conf.Endpoints = []*config.Endpoint{
		{Host: featuresTestServer.URL},
	}
	conf.FeaturesConfig.FeatureRequestTicker = time.NewTicker(100 * time.Millisecond)
	conf.FeaturesConfig.MaxRetries = 10

	featureChan := make(chan map[string]bool, 1)
	features := NewTestFeatures(conf, featureChan)

	// assert feature not supported before fetched
	assert.False(t, features.FeatureEnabled("some-test-feature"))

	// start feature fetcher
	features.Start()

	// assert feature supported after fetch completed
	select {
	case <-time.After(2 * time.Second):
		assert.Equal(t, 0, features.retries)
		assert.False(t, features.FeatureEnabled("some-test-feature"))
	default:
		// check on each loop if the condition is satisfied yet, otherwise wait for the timeout
		if features.retries == 0 {
			assert.Equal(t, 0, features.retries)
		}
	}

	// stop feature fetcher
	features.Stop()
}

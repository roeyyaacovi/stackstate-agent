package features

import (
	"github.com/StackVista/stackstate-agent/pkg/trace/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFeaturesWithRetries(t *testing.T) {
	GlobalRetries := 0
	featuresTestServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if GlobalRetries >= 3 {
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
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}),
	)

	conf := config.New()
	conf.Endpoints = []*config.Endpoint{
		{Host: featuresTestServer.URL},
	}
	conf.FeaturesConfig.FeatureRequestTicker = time.NewTicker(500 * time.Millisecond)
	conf.FeaturesConfig.MaxRetries = 10

	featureChan := make(chan map[string]bool, 1)
	features := NewTestFeatures(conf, featureChan)

	// assert feature not supported before fetched
	assert.False(t, features.FeatureEnabled("some-test-feature"))

	// start feature fetcher
	features.Start()

	// assert feature supported after fetch completed
	timeout := time.After(5 * time.Second)
	assert := func() {
		// assert we had at least 3 retries in the test scenario
		assert.True(t, features.retries >= 3)
		// assert that the feature is enabled, so we got the response from the backend
		assert.True(t, features.FeatureEnabled("some-test-feature"))
		// stop feature fetcher
		features.Stop()
	}

	assertLoop:
		for {
		select {
		case <-timeout:
			assert()
			break assertLoop
		default:
			GlobalRetries = features.retries

			// check on each loop if the condition is satisfied yet, otherwise continue until the timeout
			enabled := features.FeatureEnabled("some-test-feature")
			if enabled {
				assert()
				break assertLoop
			}
		}
	}
}

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
			case "/feature":
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
	conf.FeaturesConfig.MaxRetries = 2

	features := NewFeatures(conf)

	// assert feature not supported before fetched
	assert.False(t, features.FeatureEnabled("some-test-feature"))

	// start feature fetcher
	features.Start()

	// assert feature supported after fetch completed
	select {
	case <-conf.FeaturesConfig.FeatureRequestTicker.C:
		assert.True(t, features.FeatureEnabled("some-test-feature"))
	case <-time.After(5 * time.Second):
		assert.True(t, features.FeatureEnabled("some-test-feature"))
	}
}

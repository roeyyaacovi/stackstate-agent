package features

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/StackVista/stackstate-agent/pkg/trace/config"
	"github.com/StackVista/stackstate-agent/pkg/trace/info"
	"github.com/StackVista/stackstate-agent/pkg/trace/watchdog"
	"github.com/StackVista/stackstate-agent/pkg/util/log"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type featureEndpoint struct {
	*config.Endpoint
	client *http.Client
}

type Features struct {
	config *config.AgentConfig
	endpoint *featureEndpoint
	featureTicket *time.Ticker
	featureChan chan map[string]bool
	retries int
	features map[string]bool
}

// NewFeatures returns a Features type given the config
func NewFeatures(conf *config.AgentConfig) *Features {
	endpoint := conf.Endpoints[0]
	client := newClient(conf, false)
	if endpoint.NoProxy {
		client = newClient(conf, true)
	}

	return &Features{
		config: conf,
		endpoint: &featureEndpoint{
			Endpoint: endpoint,
			client: client,
		},
		featureTicket: time.NewTicker(5 * time.Second),
		featureChan: make(chan map[string]bool, 1),
		retries: 5,
	}
}

func (f *Features) Start() {
	go func() {
		defer watchdog.LogOnPanic()
		for {
			select {
			case <-f.featureTicket.C:
				f.getSupportedFeatures()
			case featuresMap := <-f.featureChan:
				// Set the supported features
				f.features = featuresMap
				// Stop polling and close this channel
				f.featureTicket.Stop()
				close(f.featureChan)
			}
		}
	}()
}

func (f *Features) Stop() {
	f.featureTicket.Stop()
	close(f.featureChan)
}
// timeout is the HTTP timeout for POST requests to the StackState backend
const timeout = 10 * time.Second

// getSupportedFeatures returns the features supported by the StackState API
func (f *Features) getSupportedFeatures() {
	f.retries = f.retries -1
	if f.retries == 0 {
		f.featureChan <- map[string]bool{}
	}

	resp, accessErr := f.makeFeatureRequest()

	// Handle error response
	if accessErr != nil {
		// Soo we got a 404, meaning we were able to contact stackstate, but it had no features path. We can publish a result
		if resp != nil {
			log.Info("Found StackState version which does not support feature detection yet")
			return
		}
		// Log
		_ = log.Error(accessErr)
		return
	}

	defer resp.Body.Close()

	// Get byte array
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		_ = log.Errorf("could not decode response body from features: %s", err)
		return
	}
	var data interface{}
	// Parse json
	err = json.Unmarshal(body, &data)
	if err != nil {
		_ = log.Errorf("error unmarshalling features json: %s of body %s", err, body)
		return
	}

	// Validate structure
	featureMap, ok := data.(map[string]interface{})
	if !ok {
		_ = log.Errorf("Json was wrongly formatted, expected map type, got: %s", reflect.TypeOf(data))
		return
	}

	featuresParsed := make(map[string]bool)

	for k, v := range featureMap {
		featureValue, okV := v.(bool)
		if !okV {
			_ = log.Warnf("Json was wrongly formatted, expected boolean type, got: %s, skipping feature %s", reflect.TypeOf(v), k)
		}
		featuresParsed[k] = featureValue
	}

	log.Infof("Server supports features: %s", featuresParsed)
	f.featureChan <- featuresParsed
}

// FeatureEnabled checks whether a certain feature is enabled
func (f *Features) FeatureEnabled(feature string) bool {
	if supported, ok := f.features[feature]; ok {
		return supported
	}
	return false
}

func (f *Features) makeFeatureRequest() (*http.Response, error) {
	url := fmt.Sprintf("%s/features", f.endpoint.Host)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request to %s/features: %s", url, err)
	}

	req.Header.Add("content-encoding", "identity")
	req.Header.Add("sts-api-key", f.endpoint.APIKey)
	req.Header.Add("sts-hostname", f.endpoint.Host)
	req.Header.Add("sts-traceagentversion", info.VersionString())

	resp, err := f.endpoint.client.Do(req)
	if err != nil {
		if isHTTPTimeout(err) {
			return nil, fmt.Errorf("timeout detected on %s, %s", url, err)
		}
		return nil, fmt.Errorf("error submitting payload to %s: %s", url, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		defer resp.Body.Close()
		io.Copy(ioutil.Discard, resp.Body)
		return resp, fmt.Errorf("unexpected response from %s. Status: %s", url, resp.Status)
	}

	return resp, nil
}

// newClient returns a http.Client configured with the Agent options.
func newClient(conf *config.AgentConfig, ignoreProxy bool) *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: conf.SkipSSLValidation},
	}
	if conf.ProxyURL != nil && !ignoreProxy {
		log.Infof("configuring proxy through: %s", conf.ProxyURL.String())
		transport.Proxy = http.ProxyURL(conf.ProxyURL)
	}
	return &http.Client{Timeout: timeout, Transport: transport}
}

// IsTimeout returns true if the error is due to reaching the timeout limit on the http.client
func isHTTPTimeout(err error) bool {
	if netErr, ok := err.(interface {
		Timeout() bool
	}); ok && netErr.Timeout() {
		return true
	} else if strings.Contains(err.Error(), "use of closed network connection") { //To deprecate when using GO > 1.5
		return true
	}
	return false
}

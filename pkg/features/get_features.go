package features

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/StackVista/stackstate-agent/pkg/trace/config"
	"github.com/StackVista/stackstate-agent/pkg/trace/info"
	"github.com/StackVista/stackstate-agent/pkg/util/log"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// timeout is the HTTP timeout for POST requests to the StackState backend
const timeout = 10 * time.Second

func GetSupportedFeatures(conf *config.AgentConfig) (Features, error) {
	resp, accessErr := makeFeatureRequest(conf)

	// Handle error response
	if accessErr != nil {
		// Soo we got a 404, meaning we were able to contact stackstate, but it had no features path. We can publish a result
		if resp != nil {
			log.Info("Found StackState version which does not support feature detection yet")
			return nil, errors.New("")
		}
		// Log
		_ = log.Error(accessErr)
		return nil, accessErr
	}

	defer resp.Body.Close()

	// Get byte array
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		_ = log.Errorf("could not decode response body from features: %s", err)
		return nil, err
	}
	var data interface{}
	// Parse json
	err = json.Unmarshal(body, &data)
	if err != nil {
		_ = log.Errorf("error unmarshalling features json: %s of body %s", err, body)
		return nil, err
	}

	// Validate structure
	featureMap, ok := data.(map[string]interface{})
	if !ok {
		_ = log.Errorf("Json was wrongly formatted, expected map type, got: %s", reflect.TypeOf(data))
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
	return NewFeatures(featuresParsed), nil
}

func makeFeatureRequest(conf *config.AgentConfig) (*http.Response, error) {
	ignoreProxy := true
	endpoint := conf.Endpoints[0]
	client := newClient(conf, !ignoreProxy)
	if endpoint.NoProxy {
		client = newClient(conf, ignoreProxy)
	}
	url := fmt.Sprintf("%s/features", endpoint.Host)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request to %s/features: %s", url, err)
	}

	req.Header.Add("content-encoding", "identity")
	req.Header.Add("sts-api-key", endpoint.APIKey)
	req.Header.Add("sts-hostname", endpoint.Host)
	req.Header.Add("sts-traceagentversion", info.VersionString())

	resp, err := client.Do(req)
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
			DualStack: true,
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

// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package providers

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/StackVista/stackstate-agent/pkg/autodiscovery/integration"
	"github.com/StackVista/stackstate-agent/pkg/config"
)

func TestParseJSONValue(t *testing.T) {
	// empty value
	res, err := parseJSONValue("")
	assert.Nil(t, res)
	assert.NotNil(t, err)

	// value is not a list
	res, err = parseJSONValue("{}")
	assert.Nil(t, res)
	assert.NotNil(t, err)

	// invalid json
	res, err = parseJSONValue("[{]")
	assert.Nil(t, res)
	assert.NotNil(t, err)

	// bad type
	res, err = parseJSONValue("[1, {\"test\": 1}, \"test\"]")
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, "found non JSON object type, value is: '1'", err.Error())

	// valid input
	res, err = parseJSONValue("[{\"test\": 1}, {\"test\": 2}]")
	assert.Nil(t, err)
	assert.NotNil(t, res)
	require.Len(t, res, 2)
	assert.Equal(t, integration.Data("{\"test\":1}"), res[0])
	assert.Equal(t, integration.Data("{\"test\":2}"), res[1])
}

func TestParseCheckNames(t *testing.T) {
	// empty value
	res, err := parseCheckNames("")
	assert.Nil(t, res)
	assert.NotNil(t, err)

	// value is not a list
	res, err = parseCheckNames("{}")
	assert.Nil(t, res)
	assert.NotNil(t, err)

	// invalid json
	res, err = parseCheckNames("[{]")
	assert.Nil(t, res)
	assert.NotNil(t, err)

	// ignore bad type
	res, err = parseCheckNames("[1, {\"test\": 1}, \"test\"]")
	assert.Nil(t, res)
	assert.NotNil(t, err)

	// valid input
	res, err = parseCheckNames("[\"test1\", \"test2\"]")
	assert.Nil(t, err)
	assert.NotNil(t, res)
	require.Len(t, res, 2)
	assert.Equal(t, []string{"test1", "test2"}, res)
}

func TestBuildStoreKey(t *testing.T) {
	res := buildStoreKey()
	assert.Equal(t, "/datadog/check_configs", res)
	res = buildStoreKey("")
	assert.Equal(t, "/datadog/check_configs", res)
	res = buildStoreKey("foo")
	assert.Equal(t, "/datadog/check_configs/foo", res)
	res = buildStoreKey("foo", "bar")
	assert.Equal(t, "/datadog/check_configs/foo/bar", res)
	res = buildStoreKey("foo", "bar", "bazz")
	assert.Equal(t, "/datadog/check_configs/foo/bar/bazz", res)
}

func TestBuildTemplates(t *testing.T) {
	// wrong number of checkNames
	res := buildTemplates("id",
		[]string{"a", "b"},
		[]integration.Data{integration.Data("")},
		[]integration.Data{integration.Data("")})
	assert.Len(t, res, 0)

	res = buildTemplates("id",
		[]string{"a", "b"},
		[]integration.Data{integration.Data("{\"test\": 1}"), integration.Data("{}")},
		[]integration.Data{integration.Data("{}"), integration.Data("{1:2}")})
	require.Len(t, res, 2)

	assert.Len(t, res[0].ADIdentifiers, 1)
	assert.Equal(t, "id", res[0].ADIdentifiers[0])
	assert.Equal(t, res[0].Name, "a")
	assert.Equal(t, res[0].InitConfig, integration.Data("{\"test\": 1}"))
	assert.Equal(t, res[0].Instances, []integration.Data{integration.Data("{}")})

	assert.Len(t, res[1].ADIdentifiers, 1)
	assert.Equal(t, "id", res[1].ADIdentifiers[0])
	assert.Equal(t, res[1].Name, "b")
	assert.Equal(t, res[1].InitConfig, integration.Data("{}"))
	assert.Equal(t, res[1].Instances, []integration.Data{integration.Data("{1:2}")})
}

func TestExtractTemplatesFromMap(t *testing.T) {
	for nb, tc := range []struct {
		source       map[string]string
		adIdentifier string
		prefix       string
		output       []integration.Config
		errs         []error
	}{
		{
			// Nominal case with two templates
			source: map[string]string{
				"prefix.check_names":  "[\"apache\",\"http_check\"]",
				"prefix.init_configs": "[{},{}]",
				"prefix.instances":    "[{\"apache_status_url\":\"http://%%host%%/server-status?auto\"},{\"name\":\"My service\",\"timeout\":1,\"url\":\"http://%%host%%\"}]",
			},
			adIdentifier: "id",
			prefix:       "prefix.",
			output: []integration.Config{
				{
					Name:          "apache",
					Instances:     []integration.Data{integration.Data("{\"apache_status_url\":\"http://%%host%%/server-status?auto\"}")},
					InitConfig:    integration.Data("{}"),
					ADIdentifiers: []string{"id"},
				},
				{
					Name:          "http_check",
					Instances:     []integration.Data{integration.Data("{\"name\":\"My service\",\"timeout\":1,\"url\":\"http://%%host%%\"}")},
					InitConfig:    integration.Data("{}"),
					ADIdentifiers: []string{"id"},
				},
			},
		},
		{
			// Take one, ignore one
			source: map[string]string{
				"prefix.check_names":   "[\"apache\"]",
				"prefix.init_configs":  "[{}]",
				"prefix.instances":     "[{\"apache_status_url\":\"http://%%host%%/server-status?auto\"}]",
				"prefix2.check_names":  "[\"apache\"]",
				"prefix2.init_configs": "[{}]",
				"prefix2.instances":    "[{\"apache_status_url\":\"http://%%host%%/server-status?auto\"}]",
			},
			adIdentifier: "id",
			prefix:       "prefix.",
			output: []integration.Config{
				{
					Name:          "apache",
					Instances:     []integration.Data{integration.Data("{\"apache_status_url\":\"http://%%host%%/server-status?auto\"}")},
					InitConfig:    integration.Data("{}"),
					ADIdentifiers: []string{"id"},
				},
			},
		},
		{
			// Logs config
			source: map[string]string{
				"prefix.logs": "[{\"service\":\"any_service\",\"source\":\"any_source\"}]",
			},
			adIdentifier: "id",
			prefix:       "prefix.",
			output: []integration.Config{
				{
					LogsConfig:    integration.Data("[{\"service\":\"any_service\",\"source\":\"any_source\"}]"),
					ADIdentifiers: []string{"id"},
				},
			},
		},
		{
			// Check + logs
			source: map[string]string{
				"prefix.check_names":  "[\"apache\"]",
				"prefix.init_configs": "[{}]",
				"prefix.instances":    "[{\"apache_status_url\":\"http://%%host%%/server-status?auto\"}]",
				"prefix.logs":         "[{\"service\":\"any_service\",\"source\":\"any_source\"}]",
			},
			adIdentifier: "id",
			prefix:       "prefix.",
			output: []integration.Config{
				{
					Name:          "apache",
					Instances:     []integration.Data{integration.Data("{\"apache_status_url\":\"http://%%host%%/server-status?auto\"}")},
					InitConfig:    integration.Data("{}"),
					ADIdentifiers: []string{"id"},
				},
				{
					LogsConfig:    integration.Data("[{\"service\":\"any_service\",\"source\":\"any_source\"}]"),
					ADIdentifiers: []string{"id"},
				},
			},
		},
		{
			// Missing check_names, silently ignore map
			source: map[string]string{
				"prefix.init_configs": "[{}]",
				"prefix.instances":    "[{\"apache_status_url\":\"http://%%host%%/server-status?auto\"}]",
			},
			adIdentifier: "id",
			prefix:       "prefix.",
			output:       nil,
		},
		{
			// Missing init_configs, error out
			source: map[string]string{
				"prefix.check_names": "[\"apache\"]",
				"prefix.instances":   "[{\"apache_status_url\":\"http://%%host%%/server-status?auto\"}]",
			},
			adIdentifier: "id",
			prefix:       "prefix.",
			output:       nil,
			errs:         []error{errors.New("could not extract checks config: missing init_configs key")},
		},
		{
			// Invalid instances json
			source: map[string]string{
				"prefix.check_names":  "[\"apache\"]",
				"prefix.init_configs": "[{}]",
				"prefix.instances":    "[{\"apache_status_url\" \"http://%%host%%/server-status?auto\"}]",
			},
			adIdentifier: "id",
			prefix:       "prefix.",
			output:       nil,
			errs:         []error{errors.New("could not extract checks config: in instances: Failed to unmarshal JSON")},
		},
		{
			// Invalid logs json
			source: map[string]string{
				"prefix.logs": "{\"service\":\"any_service\",\"source\":\"any_source\"}",
			},
			adIdentifier: "id",
			prefix:       "prefix.",
			output:       nil,
			errs:         []error{errors.New("could not extract logs config: invalid format, expected an array, got: ")},
		},
		{
			// Invalid checks but valid logs
			source: map[string]string{
				"prefix.check_names": "[\"apache\"]",
				"prefix.instances":   "[{\"apache_status_url\":\"http://%%host%%/server-status?auto\"}]",
				"prefix.logs":        "[{\"service\":\"any_service\",\"source\":\"any_source\"}]",
			},
			adIdentifier: "id",
			prefix:       "prefix.",
			errs:         []error{errors.New("could not extract checks config: missing init_configs key")},
			output: []integration.Config{
				{
					LogsConfig:    integration.Data("[{\"service\":\"any_service\",\"source\":\"any_source\"}]"),
					ADIdentifiers: []string{"id"},
				},
			},
		},
		{
			// Invalid checks and invalid logs
			source: map[string]string{
				"prefix.check_names": "[\"apache\"]",
				"prefix.instances":   "[{\"apache_status_url\":\"http://%%host%%/server-status?auto\"}]",
				"prefix.logs":        "{\"service\":\"any_service\",\"source\":\"any_source\"}",
			},
			adIdentifier: "id",
			prefix:       "prefix.",
			errs: []error{
				errors.New("could not extract checks config: missing init_configs key"),
				errors.New("could not extract logs config: invalid format, expected an array, got: "),
			},
			output: nil,
		},
	} {
		t.Run(fmt.Sprintf("case %d: %s", nb, tc.source), func(t *testing.T) {
			assert := assert.New(t)
			configs, errs := extractTemplatesFromMap(tc.adIdentifier, tc.source, tc.prefix)
			assert.EqualValues(tc.output, configs)

			if len(tc.errs) == 0 {
				assert.Equal(0, len(errs))
			} else {
				for i, err := range errs {
					assert.NotNil(err)
					assert.Contains(err.Error(), tc.errs[i].Error())
				}
			}
		})
	}
}

func TestGetPollInterval(t *testing.T) {
	cp := config.ConfigurationProviders{}
	assert.Equal(t, GetPollInterval(cp), 10*time.Second)
	cp = config.ConfigurationProviders{
		PollInterval: "foo",
	}
	assert.Equal(t, GetPollInterval(cp), 10*time.Second)
	cp = config.ConfigurationProviders{
		PollInterval: "1s",
	}
	assert.Equal(t, GetPollInterval(cp), 1*time.Second)
}

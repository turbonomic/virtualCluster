package client

import (
	"net/url"
	"reflect"
	"testing"
)

func TestNewConfigBuilder(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost")
	table := []struct {
		serverAddress         *url.URL
		expectedConfigBuilder *ConfigBuilder
	}{
		{
			serverAddress: baseURL,
			expectedConfigBuilder: &ConfigBuilder{
				serverAddress: baseURL,
			},
		},
	}

	for _, item := range table {
		configBuilder := NewConfigBuilder(item.serverAddress)
		if !reflect.DeepEqual(item.expectedConfigBuilder, configBuilder) {
			t.Errorf("Expect ConfigBuilder %++v, got %++v", item.expectedConfigBuilder, configBuilder)
		}
	}
}

func TestConfigBuilder_BasicAuthentication(t *testing.T) {
	table := []struct {
		username          string
		password          string
		expectedBasicAuth *BasicAuthentication
	}{
		{
			username:          "foo",
			password:          "bar",
			expectedBasicAuth: &BasicAuthentication{"foo", "bar"},
		},
	}

	baseURL, _ := url.Parse("http://localhost")
	for _, item := range table {
		configBuilder := NewConfigBuilder(baseURL).BasicAuthentication(item.username, item.password)
		if !reflect.DeepEqual(item.expectedBasicAuth, configBuilder.basicAuth) {
			t.Errorf("Expect basic authentication %++v, got %++v",
				item.expectedBasicAuth, configBuilder.basicAuth)
		}
	}
}

func TestConfigBuilder_Create(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost")
	table := []struct {
		serverAddress  *url.URL
		apiPath        string
		username       string
		password       string
		expectedConfig *Config
	}{
		{
			serverAddress:  baseURL,
			apiPath:        "path/to/api",
			username:       "foo",
			password:       "bar",
			expectedConfig: &Config{baseURL, "path/to/api", &BasicAuthentication{"foo", "bar"}},
		},
		{
			serverAddress:  baseURL,
			expectedConfig: &Config{baseURL, defaultAPIPath, nil},
		},
	}
	for _, item := range table {
		cb := NewConfigBuilder(item.serverAddress)
		if item.apiPath != "" {
			cb = cb.APIPath(item.apiPath)
		}
		if item.username != "" && item.password != "" {
			cb = cb.BasicAuthentication(item.username, item.password)
		}
		config := cb.Create()
		if !reflect.DeepEqual(item.expectedConfig, config) {
			t.Errorf("Expect config %++v, got %++v",
				item.expectedConfig, config)
		}
	}
}

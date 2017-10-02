package client

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestNewRESTClient(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost")
	expectedRESTClient := &RESTClient{
		http.DefaultClient,
		baseURL,
		"path/to/api",
		nil,
	}
	restClient := NewRESTClient(http.DefaultClient, baseURL, "path/to/api")
	if !reflect.DeepEqual(restClient, expectedRESTClient) {
		t.Errorf("Expected REST client %++v, got %++v", expectedRESTClient, restClient)
	}
}

func TestRESTClient_Verb(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost")
	client := http.DefaultClient
	apiPath := "/path/to/api"
	table := []struct {
		verb            string
		basicAuth       *BasicAuthentication
		expectedRequest *Request
	}{
		{
			"GET",
			nil,
			&Request{
				client:     client,
				verb:       "GET",
				baseURL:    baseURL,
				pathPrefix: apiPath,
			},
		},
		{
			"POST",
			&BasicAuthentication{"foo", "bar"},
			&Request{
				client:     client,
				verb:       "POST",
				baseURL:    baseURL,
				pathPrefix: apiPath,
				basicAuth:  &BasicAuthentication{"foo", "bar"},
			},
		},
	}

	for _, item := range table {
		restClient := NewRESTClient(client, baseURL, apiPath)
		if item.basicAuth != nil {
			restClient.basicAuth = item.basicAuth
		}
		request := restClient.Verb(item.verb)
		if !reflect.DeepEqual(item.expectedRequest, request) {
			t.Errorf("Expected request %++v, got %++v", item.expectedRequest, request)
		}
	}
}

func TestRESTClient_Get(t *testing.T) {
	restClient := buildExampleRESTClient()

	request := restClient.Get()
	if request.verb != "GET" {
		t.Errorf("Expected GET, got %s", request.verb)
	}
}

func TestRESTClient_Delete(t *testing.T) {
	restClient := buildExampleRESTClient()

	request := restClient.Delete()
	if request.verb != "DELETE" {
		t.Errorf("Expected DELETE, got %s", request.verb)
	}
}

func TestRESTClient_Post(t *testing.T) {
	restClient := buildExampleRESTClient()

	request := restClient.Post()
	if request.verb != "POST" {
		t.Errorf("Expected POST, got %s", request.verb)
	}
}

func TestRESTClient_BasicAuthentication(t *testing.T) {
	basicAuth := &BasicAuthentication{"user", "pass"}
	restClient := buildExampleRESTClient().BasicAuthentication(basicAuth)
	if !reflect.DeepEqual(restClient.basicAuth, basicAuth) {
		t.Errorf("Expected basic authentication %++v, got %++v", basicAuth, restClient.basicAuth)
	}
}

func buildExampleRESTClient() *RESTClient{
	baseURL, _ := url.Parse("http://localhost")
	client := http.DefaultClient
	apiPath := "/path/to/api"
	return NewRESTClient(client, baseURL, apiPath)
}

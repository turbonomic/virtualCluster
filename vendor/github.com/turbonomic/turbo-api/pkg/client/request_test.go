package client

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/turbonomic/turbo-api/pkg/api"
)

func TestNewRequest(t *testing.T) {
	// NewRequest(client HTTPClient, verb string, baseURL *url.URL, apiPath string) *Request
	client := http.DefaultClient
	verb := "GET"
	baseURL, _ := url.Parse("http://localhost")
	table := []struct {
		apiPath       string
		expectRequest *Request
	}{
		{
			"",
			&Request{
				client:     client,
				verb:       verb,
				baseURL:    baseURL,
				pathPrefix: "",
			},
		},
		{
			"foo",
			&Request{
				client:     client,
				verb:       verb,
				baseURL:    baseURL,
				pathPrefix: "/foo",
			},
		},
		{
			"/bar",
			&Request{
				client:     client,
				verb:       verb,
				baseURL:    baseURL,
				pathPrefix: "/bar",
			},
		},
	}

	for _, item := range table {
		request := NewRequest(client, verb, baseURL, item.apiPath)
		if !reflect.DeepEqual(request, item.expectRequest) {
			t.Errorf("expected %++v, got %++v", item.expectRequest, request)
		}
	}
}

func TestRequest_BasicAuthentication(t *testing.T) {
	table := []struct {
		username   string
		password   string
		expectAuth *BasicAuthentication
	}{
		{"foo", "31415", &BasicAuthentication{"foo", "31415"}},
		{"bar", "", &BasicAuthentication{"bar", ""}},
	}

	u, _ := url.Parse("http://localhost")
	for _, item := range table {
		basicAuth := &BasicAuthentication{item.username, item.password}
		request := NewRequest(http.DefaultClient, "GET", u, "").BasicAuthentication(basicAuth)
		if !reflect.DeepEqual(request.basicAuth, item.expectAuth) {
			t.Errorf("expected %++v, got %++v", item.expectAuth, request)
		}
	}
}

func TestRequest_Resource(t *testing.T) {
	table := []struct {
		existingError    error
		existingResource api.ResourceType
		resource         api.ResourceType
		expectedResource api.ResourceType
		expectsError     bool
	}{
		{
			existingError: fmt.Errorf("error"),
			expectsError:  true,
		},
		{
			existingResource: api.Resource_Type_Target,
			resource:         api.Resource_Type_External_Target,
			expectsError:     true,
		},
		{
			resource:         api.Resource_Type_Target,
			expectedResource: api.Resource_Type_Target,
			expectsError:     false,
		},
	}

	u, _ := url.Parse("http://localhost")
	for _, item := range table {
		request := NewRequest(http.DefaultClient, "GET", u, "")
		if item.existingError != nil {
			request.err = item.existingError
		}
		if item.existingResource != "" {
			request.resource = item.existingResource
		}
		request = request.Resource(item.resource)
		if request.err != nil != item.expectsError {
			t.Error("Error handling check failed.")
		}
		if !item.expectsError && request.resource != item.expectedResource {
			t.Error("Expected Reource %s, got %s", item.expectedResource, request.resource)
		}
	}
}

func TestRequest_Name(t *testing.T) {
	table := []struct {
		existingError        error
		existingResourceName string
		resourceName         string
		expectedResourceName string
		expectsError         bool
	}{
		{
			existingError: fmt.Errorf("error"),
			expectsError:  true,
		},
		{
			existingResourceName: "target1",
			resourceName:         "target2",
			expectsError:         true,
		},
		{
			resourceName: "",
			expectsError: true,
		},
		{
			resourceName:         "target1",
			expectedResourceName: "target1",
			expectsError:         false,
		},
	}

	u, _ := url.Parse("http://localhost")
	for _, item := range table {
		request := NewRequest(http.DefaultClient, "GET", u, "")
		if item.existingError != nil {
			request.err = item.existingError
		}
		if item.existingResourceName != "" {
			request.resourceName = item.existingResourceName
		}
		request = request.Name(item.resourceName)
		if request.err != nil != item.expectsError {
			t.Error("Error handling check failed.")
		}
		if !item.expectsError && request.resourceName != item.expectedResourceName {
			t.Error("Expected Reource %s, got %s", item.expectedResourceName, request.resourceName)
		}
	}
}

func TestRequest_Param(t *testing.T) {
	u, _ := url.Parse("http://localhost")
	table := []struct {
		name      string
		testVal   string
		expectStr string
	}{
		{"foo", "31415", "http://localhost?foo=31415"},
		{"bar", "42", "http://localhost?bar=42"},
		{"baz", "0", "http://localhost?baz=0"},
	}

	for _, item := range table {
		r := NewRequest(http.DefaultClient, "GET", u, "").Param(item.name, item.testVal)
		if e, a := item.expectStr, r.URL().String(); e != a {
			t.Errorf("expected %v, got %v", e, a)
		}
	}
}

func TestRequest_NameURL(t *testing.T) {
	u, _ := url.Parse("http://localhost")
	tests := []struct {
		name      string
		expectStr string
	}{
		{"bar", "http://localhost/bar"},
		{"foo", "http://localhost/foo"},
	}
	for _, test := range tests {
		r := NewRequest(http.DefaultClient, "GET", u, "").Name(test.name)
		if e, a := test.expectStr, r.URL().String(); e != a {
			t.Errorf("expected %s, got %s", e, a)
		}
	}
}

func TestRequest_ResourceURL(t *testing.T) {
	u, _ := url.Parse("http://localhost")
	tests := []struct {
		resource  api.ResourceType
		expectStr string
	}{
		{api.Resource_Type_Target, u.String() + "/targets"},
		{api.Resource_Type_External_Target, u.String() + "/externaltargets"},
	}
	for _, test := range tests {
		r := NewRequest(http.DefaultClient, "GET", u, "").Resource(test.resource)
		if e, a := test.expectStr, r.URL().String(); e != a {
			t.Errorf("expected %s, got %s", e, a)
		}
	}
}

func TestURLInOrder(t *testing.T) {
	u, _ := url.Parse("http://localhost")
	tests := []struct {
		resource     api.ResourceType
		resourceName string
		parameters   map[string]string
		expectStr    string
	}{
		{
			resource:     api.Resource_Type_Target,
			resourceName: "foo",
			expectStr:    "http://localhost/targets/foo",
		},
		{
			resource: api.Resource_Type_External_Target,
			parameters: map[string]string{
				"foo": "12",
			},
			expectStr: "http://localhost/externaltargets?foo=12"},
	}
	for _, test := range tests {
		r := NewRequest(http.DefaultClient, "GET", u, "").Resource(test.resource).Name(test.resourceName)
		for key, val := range test.parameters {
			r.Param(key, val)
		}
		if e, a := test.expectStr, r.URL().String(); e != a {
			t.Errorf("expected %s, got %s", e, a)
		}
	}
}

// Only test resp == nil for now.
func TestParseHTTPResponse(t *testing.T) {
	table := []struct {
		resp         *http.Response
		expectsError bool
	}{
		{
			nil,
			true,
		},
	}
	for _, item := range table {
		result := parseHTTPResponse(item.resp)
		if item.expectsError && result.err == nil {
			t.Error("Expected error, but got nil in err field in Result")
		}
	}
}

func TestRequest_Header(t *testing.T) {
	table := []struct {
		existingHeader map[string]string
		key            string
		value          string

		existingError error
	}{
		{
			existingHeader: map[string]string{
				"foo_1": "bar_1",
			},
			key:           "foo_key",
			value:         "bar_value",
			existingError: fmt.Errorf("Error"),
		},
		{
			existingHeader: map[string]string{
				"foo_1": "bar_1",
			},
			key:   "foo_key",
			value: "bar_value",
		},
		{
			key:   "foo_key",
			value: "bar_value",
		},
	}
	u, _ := url.Parse("http://localhost")
	for _, item := range table {
		request := NewRequest(http.DefaultClient, "GET", u, "")
		expectedRequest := &Request{
			client:     request.client,
			verb:       request.verb,
			baseURL:    request.baseURL,
			pathPrefix: request.pathPrefix,

			headers: request.headers,
		}
		if item.existingError != nil {
			request.err = item.existingError
			expectedRequest.err = item.existingError
		} else {
			if expectedRequest.headers == nil {
				expectedRequest.headers = map[string]string{}
			}
			expectedRequest.headers[item.key] = item.value
		}
		newRequest := request.Header(item.key, item.value)
		if !reflect.DeepEqual(expectedRequest, newRequest) {
			t.Errorf("Expected %v, got %v", expectedRequest, newRequest)
		}

	}
}

func TestRequest_Data(t *testing.T) {
	table := []struct {
		data        []byte
		existingErr error
	}{
		{
			data:        []byte("Some string"),
			existingErr: fmt.Errorf("Error"),
		},
		{
			data: []byte("Some string"),
		},
	}
	u, _ := url.Parse("http://localhost")
	for _, item := range table {
		request := NewRequest(http.DefaultClient, "GET", u, "")
		expectedRequest := &Request{
			client:     request.client,
			verb:       request.verb,
			baseURL:    request.baseURL,
			pathPrefix: request.pathPrefix,
		}
		if item.existingErr != nil {
			request.err = item.existingErr
			expectedRequest.err = item.existingErr
		} else {
			expectedRequest.data = bytes.NewBuffer(item.data)
		}
		newRequest := request.Data(item.data)
		if !reflect.DeepEqual(expectedRequest, newRequest) {
			t.Errorf("Expected %v, got %v", expectedRequest, newRequest)
		}
	}
}

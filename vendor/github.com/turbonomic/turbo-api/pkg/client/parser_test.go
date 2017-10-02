package client

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/turbonomic/turbo-api/pkg/api"
)

func TestParseAPIErrorDTO(t *testing.T) {
	table := []struct {
		responseType int
		exception    string
		message      string

		expectsErr bool
	}{
		{
			responseType: 400,
			exception:    "internal exception",
			message:      "error happend",

			expectsErr: false,
		},
	}

	for _, item := range table {
		input := fmt.Sprintf("{\"type\": %d, \"exception\":\"%s\", \"message\":\"%s\"}", item.responseType,
			item.exception, item.message)
		errDTO, err := parseAPIErrorDTO(input)
		if err != nil {
			if !item.expectsErr {
				t.Errorf("Unexpected error %s", err)
			}
		} else {
			expectedErrorDTO := &api.APIErrorDTO{
				ResponseType: item.responseType,
				Exception:    item.exception,
				Message:      item.message,
			}
			if !reflect.DeepEqual(expectedErrorDTO, errDTO) {
				t.Errorf("Expected %v, got %v", expectedErrorDTO, errDTO)
			}
		}
	}
}

func TestParseAPIErrorDTOUnMarshallError(t *testing.T) {
	input := "invalid content"
	_, err := parseAPIErrorDTO(input)
	if err == nil {
		t.Error("Expected UnMarshalling error, got not error.")
	}
}

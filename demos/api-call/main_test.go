package main

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockHTTPClient struct {
	statusCode int
	body       string
	err        error
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(m.body)),
	}, nil
}

func TestGetUser(t *testing.T) {
	tests := map[string]struct {
		statusCode int
		body       string
		mockError  error

		expectedName   string
		expectedErrMsg string
	}{
		"success": {
			statusCode:   http.StatusOK,
			body:         `{"name": "John Doe"}`,
			expectedName: "John Doe",
		},
		"http error": {
			mockError:      io.ErrUnexpectedEOF,
			expectedErrMsg: "failed to fetch user",
		},
		"non-OK status code": {
			statusCode:     http.StatusNotFound,
			expectedErrMsg: "unexpected status code: 404",
		},
		"invalid JSON": {
			statusCode:     http.StatusOK,
			body:           "invalid json",
			expectedErrMsg: "failed to parse user",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				statusCode: tt.statusCode,
				body:       tt.body,
				err:        tt.mockError,
			}

			user, err := GetUser(mockClient, 1)

			if tt.expectedErrMsg != "" {
				assert.ErrorContains(t, err, tt.expectedErrMsg)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedName, user.Name)
		})
	}
}

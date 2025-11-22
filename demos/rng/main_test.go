package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockStringGenerator returns a predetermined string
type mockStringGenerator struct {
	value string
}

func (m *mockStringGenerator) Generate() string {
	return m.value
}

func TestGenerateName(t *testing.T) {
	tests := map[string]struct {
		baseName     string
		suffix       string
		expectedName string
	}{
		"simple pod name": {
			baseName:     "my-pod",
			suffix:       "abcde",
			expectedName: "my-pod-abcde",
		},
		"deployment with numbers": {
			baseName:     "nginx-deployment",
			suffix:       "01234",
			expectedName: "nginx-deployment-01234",
		},
		"service with mixed": {
			baseName:     "api-service",
			suffix:       "a0b1c",
			expectedName: "api-service-a0b1c",
		},
		"all same character": {
			baseName:     "test",
			suffix:       "aaaaa",
			expectedName: "test-aaaaa",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mock := &mockStringGenerator{value: tt.suffix}
			generator := &NameGenerator{generator: mock}

			result := generator.GenerateName(tt.baseName)

			assert.Equal(t, tt.expectedName, result)
		})
	}
}

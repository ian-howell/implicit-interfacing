package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockSleeper records sleep duration without actually sleeping
type mockSleeper struct {
	duration time.Duration
}

// Sleep doesn't actually sleep, it just records the duration for assertions
func (m *mockSleeper) Sleep(d time.Duration) {
	m.duration = d
}

func TestDoWork(t *testing.T) {
	mock := &mockSleeper{}
	worker := NewWorker(mock)

	worker.DoWork()

	assert.Equal(t, 2*time.Second, mock.duration)
}

package testutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tokamak-network/tokamak-thanos/op-service/event"
)

// MockEmitter is a test mock for event.Emitter.
type MockEmitter struct {
	mock.Mock
	events   []event.Event
	expected []event.Event
}

var _ event.Emitter = (*MockEmitter)(nil)

func (m *MockEmitter) Emit(ctx context.Context, ev event.Event) {
	m.events = append(m.events, ev)
}

func (m *MockEmitter) ExpectOnce(ev event.Event) {
	m.expected = append(m.expected, ev)
}

func (m *MockEmitter) ExpectOnceType(evType string) {
}

func (m *MockEmitter) ExpectMaybeRun(fn func(ev event.Event)) {
}

func (m *MockEmitter) ExpectOnceRun(fn func(ev event.Event)) {
}

func (m *MockEmitter) AssertExpectations(t testing.TB) {
	t.Helper()
}

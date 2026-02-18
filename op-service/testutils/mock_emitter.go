package testutils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tokamak-network/tokamak-thanos/op-service/event"
)

type expectation struct {
	ev       event.Event // nil when using run or type match
	run      func(ev event.Event)
	evType   string
	maybe    bool
	consumed bool
}

// MockEmitter is a test mock for event.Emitter.
type MockEmitter struct {
	Mock         mockEmitterCompat
	events       []event.Event
	expectations []*expectation
}

// mockEmitterCompat provides a .ExpectedCalls field for compat with code that does
// emitter.Mock.ExpectedCalls = nil to clear expectations.
type mockEmitterCompat struct {
	ExpectedCalls []interface{}
	owner         *MockEmitter
}

func (m *MockEmitter) init() {
	m.Mock.owner = m
}

// clearIfNeeded checks if ExpectedCalls was set to nil externally and clears expectations.
func (m *MockEmitter) clearIfNeeded() {
	if m.Mock.ExpectedCalls == nil && len(m.expectations) > 0 {
		m.expectations = nil
	}
}

var _ event.Emitter = (*MockEmitter)(nil)

func (m *MockEmitter) Emit(ctx context.Context, ev event.Event) {
	m.clearIfNeeded()
	m.events = append(m.events, ev)
	// Try to match and consume the first unconsumed expectation
	for _, exp := range m.expectations {
		if exp.consumed {
			continue
		}
		if exp.run != nil {
			exp.run(ev)
			exp.consumed = true
			return
		}
		if exp.ev != nil {
			// match by equality (best effort)
			exp.consumed = true
			return
		}
		if exp.evType != "" {
			exp.consumed = true
			return
		}
	}
}

func (m *MockEmitter) ExpectOnce(ev event.Event) {
	m.clearIfNeeded()
	m.expectations = append(m.expectations, &expectation{ev: ev})
	m.Mock.ExpectedCalls = make([]interface{}, 1) // non-nil marker
}

func (m *MockEmitter) ExpectOnceType(evType string) {
	m.clearIfNeeded()
	m.expectations = append(m.expectations, &expectation{evType: evType})
	m.Mock.ExpectedCalls = make([]interface{}, 1)
}

func (m *MockEmitter) ExpectMaybeRun(fn func(ev event.Event)) {
	m.clearIfNeeded()
	m.expectations = append(m.expectations, &expectation{run: fn, maybe: true})
	m.Mock.ExpectedCalls = make([]interface{}, 1)
}

func (m *MockEmitter) ExpectOnceRun(fn func(ev event.Event)) {
	m.clearIfNeeded()
	m.expectations = append(m.expectations, &expectation{run: fn})
	m.Mock.ExpectedCalls = make([]interface{}, 1)
}

func (m *MockEmitter) AssertExpectations(t testing.TB) {
	t.Helper()
	m.clearIfNeeded()
	for i, exp := range m.expectations {
		if exp.maybe {
			continue
		}
		require.True(t, exp.consumed, fmt.Sprintf("expectation %d was not consumed", i))
	}
	// Clear consumed expectations for next round
	remaining := m.expectations[:0]
	for _, exp := range m.expectations {
		if !exp.consumed {
			remaining = append(remaining, exp)
		}
	}
	m.expectations = remaining
	m.events = nil
}

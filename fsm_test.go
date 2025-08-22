package fsm

import (
	"sync"
	"sync/atomic"
	"testing"
)

// Test basic transitions
func TestFSM_BasicTransitions(t *testing.T) {
	transitions := []Transition{
		{"S0", "go", "S1", nil, nil},
		{"S1", "back", "S0", nil, nil},
	}

	f := NewFSM("S0", transitions)

	// initial state
	if f.Current() != "S0" {
		t.Errorf("expected initial state S0, got %s", f.Current())
	}

	// valid step
	if err := f.Step("go"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if f.Current() != "S1" {
		t.Errorf("expected state S1, got %s", f.Current())
	}

	// another valid step
	if err := f.Step("back"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if f.Current() != "S0" {
		t.Errorf("expected state S0, got %s", f.Current())
	}
}

// Test guard blocking
func TestFSM_GuardBlocking(t *testing.T) {
	guardCalled := false
	transitions := []Transition{
		{"S0", "go", "S1", func() bool {
			guardCalled = true
			return false
		}, nil},
	}
	f := NewFSM("S0", transitions)

	err := f.Step("go")
	if err == nil {
		t.Errorf("expected error due to guard, got nil")
	}
	if !guardCalled {
		t.Errorf("guard was not called")
	}
	if f.Current() != "S0" {
		t.Errorf("expected state S0 (unchanged), got %s", f.Current())
	}
}

// Test action execution
func TestFSM_ActionExecution(t *testing.T) {
	var actionCount int32
	transitions := []Transition{
		{"S0", "go", "S1", nil, func() { atomic.AddInt32(&actionCount, 1) }},
	}
	f := NewFSM("S0", transitions)

	_ = f.Step("go")
	if f.Current() != "S1" {
		t.Errorf("expected state S1, got %s", f.Current())
	}
	if actionCount != 1 {
		t.Errorf("expected actionCount=1, got %d", actionCount)
	}
}

// Test reset
func TestFSM_Reset(t *testing.T) {
	transitions := []Transition{
		{"S0", "go", "S1", nil, nil},
	}
	f := NewFSM("S0", transitions)

	_ = f.Step("go")
	if f.Current() != "S1" {
		t.Errorf("expected S1, got %s", f.Current())
	}

	f.Reset()
	if f.Current() != "S0" {
		t.Errorf("expected state reset to S0, got %s", f.Current())
	}
}

// Test invalid transition
func TestFSM_InvalidTransition(t *testing.T) {
	transitions := []Transition{
		{"S0", "go", "S1", nil, nil},
	}
	f := NewFSM("S0", transitions)

	err := f.Step("invalid")
	if err == nil {
		t.Errorf("expected error for invalid transition, got nil")
	}
	if f.Current() != "S0" {
		t.Errorf("state should remain S0, got %s", f.Current())
	}
}

// Concurrency test
func TestFSM_Concurrent(t *testing.T) {
	transitions := []Transition{
		{"S0", "go", "S1", nil, nil},
		{"S1", "back", "S0", nil, nil},
	}
	f := NewFSM("S0", transitions)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				_ = f.Step("go")
			} else {
				_ = f.Step("back")
			}
		}(i)
	}
	wg.Wait()

	// final state must be either S0 or S1, not corrupted
	state := f.Current()
	if state != "S0" && state != "S1" {
		t.Errorf("unexpected state after concurrency test: %s", state)
	}
}

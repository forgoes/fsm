package fsm

import (
	"errors"
	"sync"
)

type Event = string

type Transition struct {
	From   string
	Event  Event
	To     string
	Guard  func() bool // check before transition
	Action func()      // run after transition
}

type FSM struct {
	state       string
	initial     string
	transitions []Transition
	mu          sync.RWMutex
}

func NewFSM(initial string, transitions []Transition) *FSM {
	return &FSM{
		state:       initial,
		initial:     initial,
		transitions: transitions,
	}
}

func (f *FSM) Current() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.state
}

func (f *FSM) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.state = f.initial
}

func (f *FSM) Step(event Event) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, t := range f.transitions {
		if t.From == f.state && t.Event == event {
			// guard check
			if t.Guard != nil && !t.Guard() {
				return errors.New("transition blocked by guard")
			}
			// do transition
			f.state = t.To
			// post action
			if t.Action != nil {
				t.Action()
			}
			return nil
		}
	}

	return errors.New("invalid transition")
}

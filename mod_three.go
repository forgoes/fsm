package fsm

func NewModThreeFSM() *FSM {
	transitions := []Transition{
		// From remainder 0
		{"S0", "0", "S0", nil, nil},
		{"S0", "1", "S1", nil, nil},

		// From remainder 1
		{"S1", "0", "S2", nil, nil},
		{"S1", "1", "S0", nil, nil},

		// From remainder 2
		{"S2", "0", "S1", nil, nil},
		{"S2", "1", "S2", nil, nil},
	}
	return NewFSM("S0", transitions)
}

// Convert binary string to remainder
func modThreeFSM(input string) int {
	fsm := NewModThreeFSM()
	for _, c := range input {
		_ = fsm.Step(string(c))
	}
	switch fsm.Current() {
	case "S0":
		return 0
	case "S1":
		return 1
	case "S2":
		return 2
	default:
		panic("invalid state")
	}
}

package scribble

import "fmt"

type StateMachine struct {
	states         map[string][]ParseStateTransition
	terminalStates []string
}

func NewStateMachine() *StateMachine {
	return &StateMachine{
		states:         make(map[string][]ParseStateTransition),
		terminalStates: make([]string, 0),
	}
}

func (sm *StateMachine) HasState(state string) bool {
	_, ok := sm.states[state]
	return ok
}

func (sm *StateMachine) IsTerminal(state string) bool {
	for _, v := range sm.terminalStates {
		if v == state {
			return true
		}
	}
	return false
}

func (sm *StateMachine) Terminal(state string) {
	sm.terminalStates = append(sm.terminalStates, state)
}

func (sm *StateMachine) State(state, trigger, newState string) {
	if _, ok := sm.states[state]; !ok {
		sm.states[state] = make([]ParseStateTransition, 0)
	}
	sm.states[state] = append(sm.states[state], ParseStateTransition{
		CurrentState: state,
		Token:        trigger,
		NextState:    newState,
	})
}

func (sm *StateMachine) Parse(initialState string, toks []Token) ([]ParseState, error) {
	states := make([]ParseState, 0)

	i := 0
	newState := initialState
	newState, isTerminal, err := sm.Next(initialState, toks[i].Name)
	if err != nil {
		return states, err
	}
	states = append(states, ParseState{
		State:   newState,
		Token:   toks[i].Name,
		Literal: toks[i].Literal,
	})
	for isTerminal == false && err == nil {
		i++
		if i >= len(toks) {
			break
		}
		newState, isTerminal, err = sm.Next(newState, toks[i].Name)
		states = append(states, ParseState{
			State:   newState,
			Token:   toks[i].Name,
			Literal: toks[i].Literal,
		})
	}

	return states, err
}

func (sm *StateMachine) Next(currentState, trigger string) (newState string, isTerminal bool, err error) {
	if sm.IsTerminal(currentState) {
		return currentState, true, nil
	}
	if _, ok := sm.states[currentState]; !ok {
		return "", false, fmt.Errorf("no such state: %q", currentState)
	}
	for _, v := range sm.states[currentState] {
		if v.Token == trigger {
			return v.NextState, sm.IsTerminal(v.NextState), nil
		}
	}
	return "", false, fmt.Errorf("illegal state transition from %q with token %q", currentState, trigger)
}

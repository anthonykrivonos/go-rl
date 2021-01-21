package mdp

import "errors"

var defaultStatesCapacity = 128

// A generic Markov Decision Process.
type MDP interface {
	R(state string) float32
	RByIndex(stateIndex int) float32
	T(state string, action string) Transition
	TByIndex(stateIndex int, action string) Transition
	SetState(state string, index int, terminal bool, reward float32, transitions map[string]Transition) error
	SetInitialState(state string, reward float32, transitions map[string]Transition) error
	AddState(state string, terminal bool, reward float32, transitions map[string]Transition) error
	RemoveStateByIndex(index int) error
	RemoveStateByName(state string) error
	AddAction(action string) error
	RemoveAction(action string) error
	SetDiscountRate(discountRate float32) error
	SetTransition(startState, endState string, action string, probability float32)
	RemoveTransition(startState, endState string)
	RemoveTransitionByAction(startState, action string)
}

type mdp struct {
	initialState State

	states []State
	statesCapacity int
	statesSize int

	actions map[string]Action

	rewards RewardsTable
	transitions TransitionTable
	discountRate float32

	stateMap map[string]State
	stateIndexMap map[int]State
}

// NewMDP constructs a new Markov Decision process.
// `initialState` is the name of the starting state.
// `states` is a list of string states in the MDP. This does not necessarily have to include the `initialState`.
// `terminals` is a subset of `states` that indicates terminal states.
// `actions` is a list of string actions in the MDP.
// `rewards` is a 1:1 mapping of string states to the rewards associated with being in each state.
// `transitions` is a mapping of states -> actions -> Transitions, where a Transition contains a probability and a next state.
// `discountRate` is the discount rate for learning, ɣ (gamma).
// Returns an MDP and a nil error on success or returns a nil MDP and a non-nil error on failure.
func NewMDP(initialState string, states, terminals, actions []string, rewards map[string]float32, transitions map[string]map[string]Transition, discountRate float32) (MDP, error) {
	if discountRate <= 0 || discountRate > 1.0 {
		return nil, errors.New("discount rate must be in (0, 1.0]")
	}

	m := &mdp{}

	m.states = make([]State, defaultStatesCapacity)
	m.stateMap = *new(map[string]State)
	m.stateIndexMap = *new(map[int]State)
	m.actions = *new(map[string]Action)
	m.rewards = NewRewards(nil)
	m.transitions = NewTransitionTable(nil)

	// Create initial state
	if initialState != "" {
		m.initialState = NewState(initialState, 0, false)
		m.states = append(m.states, m.initialState)
		m.stateMap[initialState] = m.initialState
		m.stateIndexMap[0] = m.initialState
	} else {
		m.initialState = nil
	}

	// Create the map of actions
	for _, action := range actions {
		a := NewAction(action)
		m.actions[action] = a
	}

	// Hashify the list of terminal sets
	terminalMap := make(map[string]bool)
	for _, terminal := range terminals {
		terminalMap[terminal] = true
	}

	// Construct all other states and transitions from these states
	i := 1
	for _, state := range states {
		// Determine if the state is terminal
		isTerminal := false
		if _, ok := terminalMap[state]; ok {
			isTerminal = true
		}

		s := NewState(state, i, isTerminal)
		if initialState != state {
			m.states = append(m.states, s)
			m.stateMap[state] = s
			m.stateIndexMap[i] = s
			m.rewards.Set(s, rewards[state])
			entry := NewTransitionTableEntry(nil)
			for action, _ := range transitions[state] {
				if _, ok := terminalMap[state]; !ok {
					return nil, errors.New("action with name " + action + " not in MDP")
				}
				entry.Set(m.actions[action], transitions[state][action].Probability(), transitions[state][action].NextState())
			}
			m.transitions.Set(s, entry)
			i++
		}
	}

	// Set discount rate (gamma)
	m.discountRate = discountRate

	return m, nil
}

// NewDefaultMDP creates an empty MDP with no initial state, no states, no terminals, no actions, no rewards, no transitions, and a gamma of 1.0.
func NewDefaultMDP() (MDP, error) {
	return NewMDP("", make([]string, 0), make([]string, 0), make([]string, 0), *new(map[string]float32), *new(map[string]map[string]Transition), 1)
}

// appendStateToList adds a State object to the MDP's `states` list. Returns nil on success, or an error on failure.
func (m *mdp) appendStateToList(state State) error {
	index := state.Index()
	if index < 0 {
		return errors.New("index must be non-negative")
	}
	if index >= m.statesCapacity {
		m.statesCapacity = index * 2
		newStates := make([]State, m.statesCapacity)
		for i, state := range m.states {
			newStates[i] = state
		}
		m.states = newStates
	}
	m.states[index] = state
	return nil
}

// getStateByName returns the State object using a given state name. Returns nil if not found.
func (m *mdp) getStateByName(state string) State {
	if s, ok := m.stateMap[state]; ok {
		return s
	}
	return nil
}

// getStateByName returns the State object using the state's index. Use 0 for initial state. Returns nil if not found.
func (m *mdp) getStateByIndex(index int) State {
	if s, ok := m.stateIndexMap[index]; ok {
		return s
	}
	return nil
}

// getAction returns an Action with the given name. Returns nil if not found.
func (m *mdp) getAction(action string) Action {
	if a, ok := m.actions[action]; ok {
		return a
	}
	return nil
}

// R returns the reward value for being in the state with the provided name.
func (m *mdp) R(state string) float32 {
	return m.rewards.Get(m.getStateByName(state))
}

// RByIndex returns the reward value for being in the state with the provided index.
func (m *mdp) RByIndex(stateIndex int) float32 {
	return m.rewards.Get(m.getStateByIndex(stateIndex))
}

// T returns the Transition object (probability and next state) given a state name and action name.
func (m *mdp) T(state string, action string) Transition {
	return m.transitions.Get(m.getStateByName(state)).Get(m.getAction(action))
}

// TByIndex returns the Transition object (probability and next state) given a state index and action name.
func (m *mdp) TByIndex(stateIndex int, action string) Transition {
	return m.transitions.Get(m.getStateByIndex(stateIndex)).Get(m.getAction(action))
}

// SetState sets a state with provided properties in the MDP. Overwrites if necessary.
func (m *mdp) SetState(state string, index int, terminal bool, reward float32, transitions map[string]Transition) error {
	if index < 0 {
		return errors.New("index must be non-negative")
	}

	s := NewState(state, index, terminal)
	if m.getStateByIndex(index) != nil {
		// Delete the old state at the given index
		sOld := m.getStateByIndex(index)
		m.rewards.Remove(sOld)
		m.transitions.Remove(sOld)
	}
	m.appendStateToList(s)
	m.stateMap[state] = s
	m.stateIndexMap[index] = s
	m.rewards.Set(s, reward)

	// Create a new set of transitions from the given state
	entry := NewTransitionTableEntry(nil)
	for action, _ := range transitions {
		entry.Set(m.getAction(action), transitions[action].Probability(), transitions[action].NextState())
	}
	m.transitions.Set(s, entry)

	// Update the initial state if the index is 0
	if index == 0 {
		m.initialState = s
	}

	return nil
}

// SetInitialState sets a new initial state (same as SetState on index 0).
func (m *mdp) SetInitialState(state string, reward float32, transitions map[string]Transition) error {
	return m.SetState(state, 0, false, reward, transitions)
}

// AddState creates a new State object without overwriting any states.
func (m *mdp) AddState(state string, terminal bool, reward float32, transitions map[string]Transition) error {
	// Find first empty index
	index := 0
	for _, state := range m.states {
		if state == nil {
			break
		}
		index++
	}

	return m.SetState(state, index, terminal, reward, transitions)
}

// RemoveStateByIndex removes a State object at the provided `index`.
func (m *mdp) RemoveStateByIndex(index int) error {
	// Delete the old state at the given index
	sOld := m.getStateByIndex(index)
	if sOld == nil {
		return errors.New("state at index " + string(index) + " doesn't exist")
	}
	m.states[index] = nil
	m.rewards.Remove(sOld)
	m.transitions.Remove(sOld)
	return nil
}

// RemoveStateByName removes a State object with the provided name.
func (m *mdp) RemoveStateByName(state string) error {
	if m.getStateByName(state) == nil {
		return errors.New("state with name " + state + " doesn't exist")
	}
	return m.RemoveStateByIndex(m.getStateByName(state).Index())
}

// AddAction adds an Action object with the provided name.
func (m *mdp) AddAction(action string) error {
	if m.getAction(action) != nil {
		return errors.New("action with name " + action + " already exists")
	}
	a := NewAction(action)
	m.actions[action] = a
	return nil
}

// AddAction removes an Action object with the provided name.
func (m *mdp) RemoveAction(action string) error {
	a := m.getAction(action)
	if a == nil {
		return errors.New("action with name " + action + " doesn't exist")
	}
	delete(m.actions, action)
	for _, state := range m.states {
		m.transitions.Get(state).Remove(a)
	}
	return nil
}

// SetDiscountRate updates the MDP's discount rate (gamma).
func (m *mdp) SetDiscountRate(discountRate float32) error {
	if discountRate <= 0 || discountRate > 1.0 {
		return errors.New("discount rate must be in (0, 1.0]")
	}
	m.discountRate = discountRate
	return nil
}

// SetTransition updates the transition between two states with provided names, given the action and the probability of
// pursuing this action.
func (m *mdp) SetTransition(startState, endState string, action string, probability float32) {
	m.transitions.Update(m.getStateByName(startState), m.getAction(action), probability, m.getStateByName(endState))
}

// RemoveTransition removes the transition between the two states with given names.
func (m *mdp) RemoveTransition(startState, endState string) {
	m.transitions.Get(m.getStateByName(startState)).RemoveTransition(m.getStateByName(endState))
}

// RemoveTransition removes the transition taken via the provided action name from the given `startState`
func (m *mdp) RemoveTransitionByAction(startState, action string) {
	m.transitions.Get(m.getStateByName(startState)).Remove(m.getAction(action))
}

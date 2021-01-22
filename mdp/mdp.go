package mdp

import (
	"errors"
	"fmt"
)

// Starting number of available states
var defaultStatesCapacity = 128

// A generic Markov Decision Process.
type MDP interface {
	R(state string) float32
	RByIndex(stateIndex int) float32
	T(state string, action string) Transition
	TByIndex(stateIndex int, action string) Transition
	SetState(state string, index int, terminal bool, reward float32, transitions map[string]Transition) error
	SetStateObject(state State, reward float32, transitions map[Action]Transition) error
	SetInitialState(state string, reward float32, transitions map[string]Transition) error
	SetInitialStateObject(state State, reward float32, transitions map[Action]Transition) error
	AddState(state string, terminal bool, reward float32, transitions map[string]Transition) error
	AddStateObject(state State, reward float32, transitions map[Action]Transition) error
	RemoveStateByIndex(index int) error
	RemoveStateByName(state string) error
	RemoveStateByObject(state State) error
	AddAction(action string) error
	AddActionObject(action Action) error
	RemoveAction(action string) error
	RemoveActionObject(action Action) error
	SetDiscountRate(discountRate float32) error
	SetTransition(startState, endState string, action string, probability float32)
	RemoveTransition(startState, endState string)
	RemoveTransitionByAction(startState, action string)
	String() string
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

	m.statesCapacity = defaultStatesCapacity
	m.states = make([]State, m.statesCapacity)
	m.stateMap = make(map[string]State)
	m.stateIndexMap = make(map[int]State)
	m.actions = make(map[string]Action)
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
	return NewMDP("", make([]string, 0), make([]string, 0), make([]string, 0), make(map[string]float32), make(map[string]map[string]Transition), 1)
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

// SetState creates and sets a state with provided properties in the MDP. Overwrites if necessary.
func (m *mdp) SetState(state string, index int, terminal bool, reward float32, transitions map[string]Transition) error {
	if index < 0 {
		return errors.New("index must be non-negative")
	} else if state == "" {
		return errors.New("name must be provided")
	}

	s := NewState(state, index, terminal)
	if m.getStateByIndex(index) != nil {
		// Delete the old state at the given index
		sOld := m.getStateByIndex(index)
		m.rewards.Remove(sOld)
		m.transitions.Remove(sOld)
	}
	err := m.appendStateToList(s)
	if err != nil {
		return err
	}
	m.stateMap[state] = s
	m.stateIndexMap[index] = s
	m.rewards.Set(s, reward)

	// Create a new set of transitions from the given state
	entry := NewTransitionTableEntry(nil)
	for action, _ := range transitions {
		a := m.getAction(action)
		if a == nil {
			a = NewAction(action)
			err = m.AddActionObject(a)
			if err != nil {
				return err
			}
		}
		entry.Set(a, transitions[action].Probability(), transitions[action].NextState())
	}
	m.transitions.Set(s, entry)

	// Update the initial state if the index is 0
	if index == 0 {
		m.initialState = s
	}

	return nil
}

// SetStateObject sets a state with provided properties in the MDP. Overwrites if necessary.
func (m *mdp) SetStateObject(state State, reward float32, transitions map[Action]Transition) error {
	if state.Index() < 0 {
		return errors.New("index must be non-negative")
	} else if state.Name() == "" {
		return errors.New("name must be provided")
	}

	if m.getStateByIndex(state.Index()) != nil {
		// Delete the old state at the given index
		sOld := m.getStateByIndex(state.Index())
		m.rewards.Remove(sOld)
		m.transitions.Remove(sOld)
	}
	err := m.appendStateToList(state)
	if err != nil {
		return err
	}
	m.stateMap[state.Name()] = state
	m.stateIndexMap[state.Index()] = state
	m.rewards.Set(state, reward)

	// Create a new set of transitions from the given state
	entry := NewTransitionTableEntry(nil)
	for action, _ := range transitions {
		if m.getAction(action.Name()) == nil {
			err = m.AddActionObject(action)
			if err != nil {
				return err
			}
		}
		entry.Set(action, transitions[action].Probability(), transitions[action].NextState())
	}
	m.transitions.Set(state, entry)

	// Update the initial state if the index is 0
	if state.Index() == 0 {
		m.initialState = state
	}

	return nil
}

// SetInitialState sets a new initial state (same as SetState on index 0).
func (m *mdp) SetInitialState(state string, reward float32, transitions map[string]Transition) error {
	return m.SetState(state, 0, false, reward, transitions)
}

// SetInitialState sets a new initial State object (same as SetStateObject on index 0).
func (m *mdp) SetInitialStateObject(state State, reward float32, transitions map[Action]Transition) error {
	return m.SetStateObject(state, reward, transitions)
}

// AddState creates and adds a new State object without overwriting any states.
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

// AddState adds a State object without overwriting any states.
func (m *mdp) AddStateObject(state State, reward float32, transitions map[Action]Transition) error {
	// Find first empty index
	index := 0
	for _, state := range m.states {
		if state == nil {
			break
		}
		index++
	}

	return m.SetStateObject(state, reward, transitions)
}

// RemoveStateByIndex removes a State object at the provided `index`.
func (m *mdp) RemoveStateByIndex(index int) error {
	// Delete the old state at the given index
	sOld := m.getStateByIndex(index)
	if sOld == nil {
		return errors.New("state at index " + fmt.Sprint(index) + " doesn't exist")
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

// RemoveStateByName removes a State object.
func (m *mdp) RemoveStateByObject(state State) error {
	if m.getStateByName(state.Name()) == nil {
		return errors.New("state " + state.String() + " doesn't exist")
	}
	return m.RemoveStateByIndex(m.getStateByName(state.Name()).Index())
}

// AddAction creates and adds an Action object with the provided name.
func (m *mdp) AddAction(action string) error {
	if m.getAction(action) != nil {
		return errors.New("action with name " + action + " already exists")
	}
	a := NewAction(action)
	m.actions[action] = a
	return nil
}

// AddActionObject adds an Action object.
func (m *mdp) AddActionObject(action Action) error {
	if m.getAction(action.Name()) != nil {
		return errors.New("action " + action.String() + " already exists")
	}
	m.actions[action.Name()] = action
	return nil
}

// RemoveAction removes an Action object with the provided name.
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

// RemoveActionObject removes an Action object.
func (m *mdp) RemoveActionObject(action Action) error {
	return m.RemoveAction(action.Name())
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

func (m *mdp) String() string {
	// Construct list of states
	states := ""
	for _, state := range m.states {
		if state != nil {
			states += state.String() + ", "
		}
	}
	states = states[:len(states) - 2]

	// Construct list of actions
	actions := ""
	for _, action := range m.actions {
		actions += action.String()
	}
	actions = actions[:len(actions) - 2]

	return fmt.Sprintf(
		"M := (\n" +
			"	S = %s\n" +
			"	A = %s\n" +
			"	R = %s\n" +
			"	T = %s\n" +
			"	ɣ = %s\n)\n",
		states,
		actions,
		m.rewards.String("	"),
		m.transitions.String("	"),
		fmt.Sprintf("%.4f", m.discountRate),
	)
}

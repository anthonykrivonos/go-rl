package mdp

type MDP interface {
}

type mdp struct {
	// Public
	initialState State
	states []State
	terminals []State
	actions []Action
	rewards RewardsTable
	transitions TransitionTable
	discountRate float32

	// Private
	stateMap map[string]State
	actionMap map[string]Action
}

func NewMDP(initialState string, states, terminals, actions []string, rewards map[string]float32, transitions map[string]map[string]Transition, discountRate float32) {
	m := &mdp{}
	m.initialState = NewState(initialState, 0)
	m.states = make([]State, 0)
	m.stateMap = *new(map[string]State)
	m.actionMap = *new(map[string]Action)
	i := 0
	m.states = append(m.states, m.initialState)
	m.rewards = NewRewards(nil)
	m.transitions = NewTransitionTable(nil)
	for _, state := range states {
		s := NewState(state, i)
		m.stateMap[state] = s
		if m.initialState.Name() != state {
			m.states = append(m.states, s)
			i++
		}
		m.rewards.Set(s, rewards[state])
		entry := NewTransitionTableEntry(nil)
		for action, _ := range transitions[state] {
			entry.Set(NewAction(action), transitions[state][action].Probability(), transitions[state][action].NextState())
		}
		m.transitions.Set(s, entry)
	}
	m.terminals = NewStates(terminals)
	m.actions = make([]Action, 0)
	for _, action := range actions {
		a := NewAction(action)
		m.actionMap[action] = a
		m.actions = append(m.actions, a)
	}
	m.discountRate = discountRate
}

func (m *mdp) getState(state string) State {
	return m.stateMap[state]
}

func (m *mdp) getAction(action string) Action {
	return m.actionMap[action]
}

func (m *mdp) R(state string) float32 {
	return m.rewards.Get(m.getState(state))
}

func (m *mdp) T(state string, action string) Transition {
	return m.transitions.Get(m.getState(state)).Get(m.getAction(action))
}

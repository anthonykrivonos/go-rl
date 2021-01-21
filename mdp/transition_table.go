package mdp

type TransitionTableEntry interface {
	Get(Action) Transition
	Set(Action, float32, State)
	Remove(Action)
	RemoveTransition(nextState State)
}

type transitionTableEntry struct {
	entry map[Action]Transition
}

func (t *transitionTableEntry) Get(action Action) Transition {
	if t, ok := t.entry[action]; ok {
		return t
	}
	return nil
}

func (t *transitionTableEntry) Set(action Action, probability float32, nextState State) {
	t.entry[action] = NewTransition(probability, nextState)
}

func (t *transitionTableEntry) Remove(action Action) {
	delete(t.entry, action)
}

func (t *transitionTableEntry) RemoveTransition(nextState State) {
	for action, transition := range t.entry {
		if transition.NextState().Equals(nextState) {
			t.Remove(action)
		}
	}
}

func NewTransitionTableEntry(entry *map[Action]Transition) TransitionTableEntry {
	t := &transitionTableEntry{}
	if entry == nil {
		t.entry = *new(map[Action]Transition)
	} else {
		t.entry = *entry
	}
	return t
}


type TransitionTable interface {
	Get(State) TransitionTableEntry
	Set(State, TransitionTableEntry)
	Remove(State)
	Update(state State, action Action, probability float32, nextState State)
}

type transitionTable struct {
	table map[State]TransitionTableEntry
}

func (t *transitionTable) Get(state State) TransitionTableEntry {
	if e, ok := t.table[state]; ok {
		return e
	}
	return nil
}

func (t *transitionTable) Set(state State, entry TransitionTableEntry) {
	t.table[state] = entry
}

func (t *transitionTable) Update(state State, action Action, probability float32, nextState State) {
	t.table[state].Set(action, probability, nextState)
}

func (t *transitionTable) Remove(state State) {
	delete(t.table, state)
}

func NewTransitionTable(table *map[State]map[Action]Transition) TransitionTable {
	t := &transitionTable{}
	t.table = *new(map[State]TransitionTableEntry)
	if table != nil {
		for state, entry := range *table {
			t.table[state] = NewTransitionTableEntry(&entry)
		}
	}
	return t
}

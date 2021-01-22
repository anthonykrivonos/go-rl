package mdp

type TransitionTableEntry interface {
	Get(Action) Transition
	Set(Action, float32, State)
	Remove(Action)
	RemoveTransition(nextState State)
	String(prefix string) string
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

func (t * transitionTableEntry) String(prefix string) string {
	res := "{\n"
	for action, transition := range t.entry {
		res += prefix + "	" + action.String() + ": " + transition.String() + ",\n"
	}
	res = res[:len(res) - 2]
	res += "\n" + prefix + "}"
	return res
}

func NewTransitionTableEntry(entry *map[Action]Transition) TransitionTableEntry {
	t := &transitionTableEntry{}
	if entry == nil {
		t.entry = make(map[Action]Transition)
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
	String(prefix string) string
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

func (t * transitionTable) String(prefix string) string {
	res := "{\n"
	for state, entry := range t.table {
		res += prefix + "	" + state.String() + ": " + entry.String("		") + ",\n"
	}
	res = res[:len(res) - 2]
	res += "\n" + prefix + "}"
	return res
}

func NewTransitionTable(table *map[State]map[Action]Transition) TransitionTable {
	t := &transitionTable{}
	t.table = make(map[State]TransitionTableEntry)
	if table != nil {
		for state, entry := range *table {
			t.table[state] = NewTransitionTableEntry(&entry)
		}
	}
	return t
}

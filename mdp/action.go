package mdp

type Action interface {
	Name() string
	Equals(Action) bool
	String() string
}

type action struct {
	name string
}

func (a *action) Name() string {
	return a.name
}

func (a *action) Equals(other Action) bool {
	return a.name == other.Name()
}

func (a *action) String() string {
	return a.name
}

func NewAction(name string) Action {
	a := &action{}
	a.name = name
	return a
}

func NewActions(names []string) []Action {
	var actions []Action
	for i := 0; i < len(names); i++ {
		a := &action{}
		a.name = names[i]
		actions = append(actions, a)
	}
	return actions
}

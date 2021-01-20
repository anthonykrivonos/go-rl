package mdp

type State interface {
	Name() string
	Index() int
	Equals(State) bool
	String() string
}

type state struct {
	name string
	index int
}

func (s *state) Name() string {
	return s.name
}

func (s *state) Index() int {
	return s.index
}

func (s *state) Equals(other State) bool {
	return s.index == other.Index()
}

func (s *state) String() string {
	return "S_" + string(s.index) + ": " + s.name
}

func NewState(name string, index int) State {
	s := &state{}
	s.name = name
	s.index = index
	return s
}

func NewStates(names []string) []State {
	var states []State
	for i := 0; i < len(names); i++ {
		s := &state{}
		s.name = names[i]
		s.index = i
		states = append(states, s)
	}
	return states
}

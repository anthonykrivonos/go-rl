package mdp

type State interface {
	Name() string
	Index() int
	Terminal() bool
	Equals(State) bool
	String() string
}

type state struct {
	name string
	index int
	terminal bool
}

func (s *state) Name() string {
	return s.name
}

func (s *state) Index() int {
	return s.index
}

func (s *state) Terminal() bool {
	return s.terminal
}

func (s *state) Equals(other State) bool {
	return s.index == other.Index()
}

func (s *state) String() string {
	return "S_" + string(s.index) + ": " + s.name
}

func NewState(name string, index int, terminal bool) State {
	s := &state{}
	s.name = name
	s.index = index
	s.terminal = terminal
	return s
}

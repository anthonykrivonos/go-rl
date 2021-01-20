package mdp

type RewardsTable interface {
	Get(state State) float32
	Set(state State, reward float32)
}

type rewardsTable struct {
	table map[int]float32
}

func (r *rewardsTable) Get(state State) float32 {
	return r.table[state.Index()]
}

func (r *rewardsTable) Set(state State, reward float32) {
	r.table[state.Index()] = reward
}

func NewRewards(table *map[int]float32) RewardsTable {
	t := &rewardsTable{}
	if table == nil {
		t.table = *new(map[int]float32)
	} else {
		t.table = *table
	}
	return t
}

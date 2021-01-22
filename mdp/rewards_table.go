package mdp

import (
	"fmt"
)

type RewardsTable interface {
	Get(state State) float32
	Set(state State, reward float32)
	Remove(state State)
	String(prefix string) string
}

type rewardsTable struct {
	table map[int]float32
	stateMap map[int]State
}

func (r *rewardsTable) Get(state State) float32 {
	return r.table[state.Index()]
}

func (r *rewardsTable) Set(state State, reward float32) {
	r.table[state.Index()] = reward
	r.stateMap[state.Index()] = state
}

func (r *rewardsTable) Remove(state State) {
	delete(r.table, state.Index())
	delete(r.stateMap, state.Index())
}

func (r * rewardsTable) String(prefix string) string {
	res := "{\n"
	for index, reward := range r.table {
		res += prefix + "	" + r.stateMap[index].String() + ": " + fmt.Sprint(reward) + ",\n"
	}
	res = res[:len(res) - 2]
	res += "\n" + prefix + "}"
	return res
}

func NewRewards(table *map[State]float32) RewardsTable {
	t := &rewardsTable{}
	t.table = make(map[int]float32)
	t.stateMap = make(map[int]State)
	if table != nil {
		for state, reward := range *table {
			t.table[state.Index()] = reward
			t.stateMap[state.Index()] = state
		}
	}
	return t
}

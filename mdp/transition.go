package mdp

import (
	"fmt"
)

type Transition interface {
	Probability() float32
	NextState() State
	String() string
}

type transition struct {
	probability float32
	nextState State
}

func NewTransition(probability float32, nextState State) Transition {
	t := &transition{}
	t.probability = probability
	t.nextState = nextState
	return t
}

func (t *transition) Probability() float32 {
	return t.probability
}

func (t *transition) NextState() State {
	return t.nextState
}

func (t *transition) String() string {
	return "(" + fmt.Sprintf("%.4f", t.probability) + ", " + t.nextState.Name() + ")"
}

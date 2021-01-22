package mdp

import (
	"fmt"
	"github.com/anthonykrivonos/go-rl/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetup(t *testing.T) {
}

func TestDefaultMDP(t *testing.T) {
	mdp, err := NewDefaultMDP()

	// Ensure the default MDP is created successfully
	assert.NoError(t, err)

	// Action constants
	goUp	:= NewAction("U")
	goRight	:= NewAction("R")
	goLeft	:= NewAction("L")
	goDown	:= NewAction("D")

	// Grid constants
	topLeft 		:= NewState("TL", 0, false)
	topCenter 		:= NewState("TC", 1, false)
	topRight 		:= NewState("TR", 2, false)
	middleLeft		:= NewState("ML", 3, false)
	middleCenter	:= NewState("MC", 4, false)
	middleRight		:= NewState("MR", 5, false)
	bottomLeft		:= NewState("BL", 6, false)
	bottomMiddle	:= NewState("BM", 7, false)
	bottomRight		:= NewState("BR", 8, true)

	// Create a grid board
	uniformMove := func(up, right, down, left State) map[Action]Transition {
		dirs := make(map[Action]Transition)
		support := utils.BoolToInt(up != nil) + utils.BoolToInt(down != nil) + utils.BoolToInt(left != nil) + utils.BoolToInt(right != nil)
		probability := utils.Uniform(support)
		if up != nil {
			dirs[goUp] = NewTransition(probability, up)
		}
		if right != nil {
			dirs[goRight] = NewTransition(probability, right)
		}
		if down != nil {
			dirs[goDown] = NewTransition(probability, down)
		}
		if left != nil {
			dirs[goLeft] = NewTransition(probability, left)
		}
		return dirs
	}

	// Top left
	err = mdp.AddStateObject(topLeft, 0, uniformMove(nil, topCenter, middleLeft, nil))
	assert.NoError(t, err)

	// Top center (HOLE)
	err = mdp.AddStateObject(topCenter, -2, uniformMove(nil, topRight, middleCenter, topLeft))
	assert.NoError(t, err)

	// Top right
	err = mdp.AddStateObject(topRight, 0, uniformMove(nil, nil, middleRight, topCenter))
	assert.NoError(t, err)

	// Middle left
	err = mdp.AddStateObject(middleLeft, 0, uniformMove(topLeft, middleCenter, bottomLeft, nil))
	assert.NoError(t, err)

	// Middle center (HOLE)
	err = mdp.AddStateObject(middleCenter, -2, uniformMove(topCenter, middleRight, bottomMiddle, middleLeft))
	assert.NoError(t, err)

	// Middle right
	err = mdp.AddStateObject(middleRight, 0, uniformMove(topRight, nil, bottomRight, middleCenter))
	assert.NoError(t, err)

	// Bottom left
	err = mdp.AddStateObject(bottomLeft, 0, uniformMove(middleLeft, bottomMiddle, nil, nil))
	assert.NoError(t, err)

	// Bottom center
	err = mdp.AddStateObject(bottomMiddle, 0, uniformMove(middleCenter, middleRight, nil, middleLeft))
	assert.NoError(t, err)

	// Bottom center (terminal)
	err = mdp.AddStateObject(bottomRight, 10, uniformMove(middleRight, nil, nil, middleCenter))
	assert.NoError(t, err)

	fmt.Print(mdp.String())
}
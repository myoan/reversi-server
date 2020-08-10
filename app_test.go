package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func  TestApp_Broadcast(t *testing.T) {
	testcases := []struct {
		input    *RequestRPC
		expected *BroadcastRPC
	} {
		{
			input:    &RequestRPC{Color: 1, Cell: &Cell{X: 0, Y: 0, State: 1}},
			expected: &BroadcastRPC{Status: 100, Color: 1, Board: [][]*Cell{}},
		},
	}

	app := NewApp()
	for _, tc := range testcases {
		input, _ := json.Marshal(tc.input)
		actual := app.Broadcast(input)
		expected, _ := json.Marshal(tc.expected)
		if bytes.Compare(actual, expected) != 0 {
			t.Errorf("got: %v\nwant: %v", string(actual), string(expected))
		}
	}
}

func TestApp_SetStone_whenNotMyTurn(t *testing.T) {
	app := NewApp()
	app.updateGameState(WhiteTurn)
	err := app.SetStone(1, &Position{X: 0, Y: 0})
	if err == nil {
		t.Errorf("Black can SetStone instead of white turn")
	}
}

func TestApp_SetStone_whenMyTurn(t *testing.T) {
	app := NewApp()
	app.updateGameState(BlackTurn)
	app.SetStone(1, &Position{X: 0, Y: 0})
	c := app.board.Cell(0, 0)
	if c.State != 1 {
		t.Errorf("got: %d\nwant: %d", c.State, 1)
	}
}

func TestApp_SetStone_whenOccupied(t *testing.T) {
	app := NewApp()
	app.updateGameState(BlackTurn)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			c := app.board.Cell(i, j)
			c.Update(1)
		}
	}
	app.SetStone(1, &Position{X: 0, Y: 0})
	if app.GameState != Finish {
		t.Errorf("board occupied but state not finish, %v", app.GameState)
	}
}

func TestApp_SetStone_whenOpponentCanAllocate(t *testing.T) {
	app := NewApp()
	app.updateGameState(BlackTurn)
	app.SetStone(1, &Position{X: 2, Y: 4})
	if app.GameState != WhiteTurn {
		t.Errorf("white can allocate but state is %v", app.GameState)
	}
}

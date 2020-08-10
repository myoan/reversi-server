package main

import (
	"encoding/json"
	"fmt"

	"github.com/myoan/go-reversi"
)

type CmdSetStoneBody struct {
	Color int           `json:"color"`
	Cell  *reversi.Cell `json:"cell"`
}

func CmdSetStone(game *reversi.Game, msg json.RawMessage) {
	var body *CmdSetStoneBody
	fmt.Printf("request: %s\n", msg)
	json.Unmarshal(msg, &body)

	switch game.GameState {
	case reversi.Prepare:
		game.GameState = reversi.BlackTurn
	case reversi.BlackTurn:
		fmt.Println("Black turn")
		if body.Color != int(reversi.BlackTurn) {
			break
		}
		pos := &reversi.Position{
			X: body.Cell.X,
			Y: body.Cell.Y,
		}
		game.SetStone(1, pos)
	case reversi.WhiteTurn:
		fmt.Println("White turn")
		if body.Color != int(reversi.WhiteTurn) {
			break
		}
		pos := &reversi.Position{
			X: body.Cell.X,
			Y: body.Cell.Y,
		}
		game.SetStone(2, pos)
	case reversi.Finish:
		fmt.Println("Finish")
		fmt.Printf("%d win!\n", game.Winner())
		return
	}
	game.Show()
}

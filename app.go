package main

import (
	"encoding/json"
	"fmt"

	"github.com/myoan/go-reversi"
)

type Request struct {
	Cmd  string `json:"cmd"`
	Body string `json:"body"`
}

type RequestRPC struct {
	Color int           `json:"color"`
	Cell  *reversi.Cell `json:"cell"`
}

type BroadcastRPC struct {
	Status int               `json:"status"`
	Color  int               `json:"color"`
	Board  [][]*reversi.Cell `json:"board"`
}

type GameEvent interface {
	OnRegister(client *Client)
	OnUnregister(client *Client)
	OnFromClient(msg []byte)
}

type App struct {
	game     *reversi.Game
	black    *Client
	white    *Client
	handlers map[string]func([]byte)
}

func NewApp() *App {
	game := reversi.NewGame()
	game.GameState = reversi.BlackTurn
	return &App{
		game:     game,
		handlers: make(map[string]func([]byte)),
	}
}

func (app *App) OnRegister(client *Client) {
	if app.black == nil {
		fmt.Println("register black")
		app.black = client
		message := app.Broadcast()
		app.black.send <- message
	} else if app.white == nil {
		fmt.Println("register white")
		app.white = client
		message := app.Broadcast()
		app.white.send <- message
	} else {
		fmt.Println("no register")
	}
}

func (app *App) OnUnregister(client *Client) {
}

func (app *App) OnFromClient(msg []byte) {
	app.Read(msg)
	message := app.Broadcast()
	app.black.send <- message
	app.white.send <- message
}

func (app *App) Broadcast() []byte {
	res := *&BroadcastRPC{
		Status: 0,
		Color:  int(app.game.GameState),
		Board:  app.game.GetBoard(),
	}
	ret, _ := json.Marshal(res)
	fmt.Printf("response: %s\n", string(ret))

	return ret
}

func (app *App) Read(msg []byte) {
	var req *RequestRPC
	fmt.Printf("request: %s\n", string(msg))
	json.Unmarshal(msg, &req)

	switch app.game.GameState {
	case reversi.Prepare:
		app.game.GameState = reversi.BlackTurn
	case reversi.BlackTurn:
		fmt.Println("Black turn")
		if req.Color != int(reversi.BlackTurn) {
			break
		}
		pos := &reversi.Position{
			X: req.Cell.X,
			Y: req.Cell.Y,
		}
		app.game.SetStone(1, pos)
	case reversi.WhiteTurn:
		fmt.Println("White turn")
		if req.Color != int(reversi.WhiteTurn) {
			break
		}
		pos := &reversi.Position{
			X: req.Cell.X,
			Y: req.Cell.Y,
		}
		app.game.SetStone(2, pos)
	case reversi.Finish:
		fmt.Println("Finish")
		fmt.Printf("%d win!\n", app.game.Winner())
		return
	}
	app.game.Show()
}

/*
func (app *App) HandlerFunc(key string, f func()) {
	app.handlers[key] = f
}

func (app *App) Call(key string, input []byte) {
	app.handlers[key](input)
}
*/

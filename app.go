package main

import (
	"encoding/json"
	"fmt"

	"github.com/myoan/go-reversi"
)

type Request struct {
	Cmd  string          `json:"cmd"`
	Body json.RawMessage `json:"body"`
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
	handlers map[string]func(*reversi.Game, json.RawMessage)
}

func NewApp() *App {
	game := reversi.NewGame()
	game.GameState = reversi.BlackTurn
	return &App{
		game:     game,
		handlers: make(map[string]func(*reversi.Game, json.RawMessage)),
	}
}

func (app *App) OnRegister(client *Client) {
	if app.black == nil {
		fmt.Println("register black")
		app.black = client
		message := app.Broadcast(1)
		app.black.send <- message
	} else if app.white == nil {
		fmt.Println("register white")
		app.white = client
		message := app.Broadcast(2)
		app.white.send <- message
	} else {
		fmt.Println("no register")
	}
}

func (app *App) OnUnregister(client *Client) {
}

func (app *App) OnFromClient(msg []byte) {
	var req *Request
	fmt.Printf("request: %s\n", string(msg))
	json.Unmarshal(msg, &req)
	app.Call(req.Cmd, req.Body)

	message := app.Broadcast(int(app.game.GameState))
	app.black.send <- message
	app.white.send <- message
}

func (app *App) Broadcast(color int) []byte {
	res := *&BroadcastRPC{
		Status: 0,
		Color:  color,
		Board:  app.game.GetBoard(),
	}
	ret, _ := json.Marshal(res)
	fmt.Printf("response: %s\n", string(ret))

	return ret
}

func (app *App) HandlerFunc(cmd string, f func(*reversi.Game, json.RawMessage)) {
	fmt.Printf("regist handler: %s\n", cmd)
	app.handlers[cmd] = f
}

func (app *App) Call(cmd string, input json.RawMessage) {
	fmt.Printf("call handler: %s\n", cmd)
	app.handlers[cmd](app.game, input)
}

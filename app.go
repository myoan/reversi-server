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
	OnFromClient(client *Client, msg []byte)
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
		message := app.SendColor(1)
		fmt.Printf("send color to black: %s\n", string(message))
		app.black.send <- message
	} else if app.white == nil {
		fmt.Println("register white")
		app.white = client
		message := app.SendColor(2)
		fmt.Printf("send color to white: %s\n", string(message))
		app.white.send <- message
		app.SendTurnToBlack()
	} else {
		fmt.Println("no register")
	}
}

func (app *App) OnUnregister(client *Client) {
}

func (app *App) OnFromClient(c *Client, msg []byte) {
	var req *Request
	fmt.Printf("request: %s\n", string(msg))
	json.Unmarshal(msg, &req)
	app.Call(req.Cmd, req.Body)

	// if app.game.board.IsOcupied() {
	// app.black.send <- app.SendResult(1)
	// app.white.send <- app.SendResult(0)
	// }
	if app.IsBlack(c) {
		app.black.send <- app.SendBoard()
		app.white.send <- app.SendBoard()
		app.white.send <- app.SendTurn()
	} else if app.IsWhite(c) {
		app.black.send <- app.SendBoard()
		app.white.send <- app.SendBoard()
		app.black.send <- app.SendTurn()
	}
}

type ColorRPC struct {
	Status int `json:"status"`
	Color  int `json:"color"`
}

func (app *App) SendColor(color int) []byte {
	res := *&ColorRPC{
		Status: 0,
		Color:  color,
	}
	ret, _ := json.Marshal(res)
	return ret
}

func (app *App) IsBlack(c *Client) bool {
	return c == app.black
}

func (app *App) IsWhite(c *Client) bool {
	return c == app.white
}

type TurnRPC struct {
	Status int `json:"status"`
}

func (app *App) SendTurnToBlack() {
	message := app.SendTurn()
	fmt.Printf("send turn to black: %s\n", string(message))
	app.black.send <- message
}

func (app *App) SendTurnToWhite() {
	message := app.SendTurn()
	fmt.Printf("send turn to white: %s\n", string(message))
	app.white.send <- message
}

func (app *App) SendTurn() []byte {
	res := *&TurnRPC{
		Status: 1,
	}
	ret, _ := json.Marshal(res)
	return ret
}

type BoardRPC struct {
	Status int               `json:"status"`
	Board  [][]*reversi.Cell `json:"board"`
}

func (app *App) SendBoard() []byte {
	res := *&BoardRPC{
		Status: 2,
		Board:  app.game.GetBoard(),
	}
	ret, _ := json.Marshal(res)
	fmt.Printf("send board: %s\n", string(ret))

	return ret
}

type ResultRPC struct {
	Status int `json:"status"`
	Result int `json:"result"`
}

func (app *App) SendResult(result int) []byte {
	res := *&ResultRPC{
		Status: 3,
		Result: result,
	}
	ret, _ := json.Marshal(res)
	return ret
}

func (app *App) Broadcast(color int) []byte {
	res := *&BroadcastRPC{
		Status: 0,
		Color:  color,
		Board:  app.game.GetBoard(),
	}
	ret, _ := json.Marshal(res)
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

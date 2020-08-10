package main

import (
	"flag"
	"log"
	"net/http"
)

/*
type Command int

const Command = {
	Register = iota
	SetStone
}
*/

func main() {
	var (
		addr = flag.String("addr", ":8080", "http service address")
	)

	flag.Parse()
	app := NewApp()
	app.HandlerFunc("set_stone", CmdSetStone)
	/*
		app.HandlerFunc("register", func() {
			fmt.Println("register")
		})
		app.HandlerFunc("setstone", func() {
			fmt.Println("set stone")
		})
	*/
	hub := newHub(app)
	go hub.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

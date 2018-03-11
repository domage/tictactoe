package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/domage/tictactoe/game"
	"github.com/domage/tictactoe/sockets"
)

var g *game.Game

func localGame() {
	g = game.NewGame()
	// The players take turns.
	err := g.TakeTurn(0, 0, "0")
	if err != nil {
		fmt.Println(err)
	}
	g.TakeTurn(2, 2, "X")
	g.TakeTurn(1, 0, "0")
	g.TakeTurn(0, 1, "X")
	g.TakeTurn(0, 2, "0")
	g.TakeTurn(2, 0, "X")
	g.TakeTurn(2, 1, "0")
	g.TakeTurn(1, 1, "X")
	g.TakeTurn(1, 2, "0")

	fmt.Print(g)
	status, winner := game.BoardStatus(g)

	fmt.Printf("%v %s", status, winner)
}

func main() {
	connect := flag.String("connect", "", "IP address of process to join. If empty, go into listen mode.")
	flag.Parse()

	// If the connect flag is set, go into client mode.
	if *connect != "" {
		err := sockets.StartClient()
		if err != nil {
			log.Println("Error:" + fmt.Sprint(err))
		}
		log.Println("Client done.")
		return
	}

	// Else go into server mode.
	err := sockets.StartServer()
	if err != nil {
		log.Println(err)
	}

	log.Println("Server done.")
}

package sockets

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/domage/tictactoe/game"
)

var g *game.Game
var players int
var status string

const (
	StopCommand  = "Exit"
	BoardCommand = "Board"
	InitCommand  = "Init"
	TurnRequest  = "Turn"
	InitPhase    = "Init"
	GamePhase    = "Game"
	ExitPhase    = "Exit"
	ErrorMessage = "can't parce the message"
)

func encode(x int, y int, turn string) string {
	return strconv.Itoa(x) + strconv.Itoa(y) + turn
}

func decode(message string) (x int, y int, turn string, err error) {
	err = nil

	if len(message) != 3 {
		err = fmt.Errorf(ErrorMessage)
	}

	x = int(message[0]) - '0'
	y = int(message[1]) - '0'
	turn = string(message[2])

	if x < 0 || x > 2 || y < 0 || y > 2 || (turn != "X" && turn != "0") {
		err = fmt.Errorf(ErrorMessage)
	}

	return
}

func socketServer(port int) error {
	listen, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	defer listen.Close()
	if err != nil {
		log.Fatalf("Socket listen port %d failed,%s", port, err)
		return err
	}
	log.Printf("Begin listen port: %d", port)
	for {
		conn, err := listen.Accept()
		switch {
		case err != nil:
			log.Fatalln(err)
			continue
		case players == 0:
			status = InitPhase
			g = game.NewGame()
			players++
			go handler(conn)
			continue
		case players == 1 && status == InitPhase:
			status = GamePhase
			players++
			go handler(conn)
			continue
		case players > 2:
			conn.Close()
			log.Printf("Two players maximum")
		case status == GamePhase || status == ExitPhase:
			conn.Close()
			log.Printf("The game did not finished yet")
		}
	}
}

func handler(conn net.Conn) {
	defer conn.Close()
	var (
		buf = make([]byte, 1024)
		r   = bufio.NewReader(conn)
		w   = bufio.NewWriter(conn)
	)

ILOOP:
	for {
		n, err := r.Read(buf)
		data := string(buf[:n])
		switch err {
		case io.EOF:
			players--
			status = ExitPhase
			break ILOOP

		case nil:

			log.Println("Receive:", data)
			x, y, turn, parsingError := decode(data)
			_, winner := game.BoardStatus(g)

			if winner != "" || status == ExitPhase {
				message := fmt.Sprint(g)
				if winner != "" {
					message = message + "The winner is: " + winner + "\n" + StopCommand
				}
				if status == ExitPhase {
					message = message + "The second player has left the game \n" + StopCommand
				}
				w.Write([]byte(message))
				w.Flush()
				log.Printf("Sent: %s", message)
				players--
				status = ExitPhase
				break ILOOP
			}

			switch {

			case strings.HasPrefix(data, InitCommand):
				var turn string
				if strings.HasSuffix(data, "X") || strings.HasSuffix(data, "0") {
					s := "You are already initiated"
					w.Write([]byte(s))
					w.Flush()
					log.Printf("Sent: %s", s)
					break
				}
				if players == 1 {
					turn = "X"
				} else {
					turn = "0"
				}
				w.Write([]byte(turn))
				w.Flush()
				log.Printf("Sent: %s", turn)

			case strings.HasPrefix(data, TurnRequest):
				w.Write([]byte(game.WhoseTurn(g)))
				w.Flush()
				log.Printf("Sent: %s", game.WhoseTurn(g))

			case strings.HasPrefix(data, StopCommand):
				players--
				status = ExitPhase
				break ILOOP

			case strings.HasPrefix(data, BoardCommand):
				board := fmt.Sprint(g)
				w.Write([]byte(board))
				w.Flush()
				log.Printf("Sent: %s", board)

			case parsingError == nil:
				res := g.TakeTurn(x, y, turn)
				_, winner := game.BoardStatus(g)
				board := fmt.Sprint(g)
				var message string

				if winner != "" {
					message = board + "The winner is: " + winner + "\n" + StopCommand
					w.Write([]byte(message))
					w.Flush()
					log.Printf("Sent: %s", message)
					players--
					break ILOOP
				}

				if res != nil {
					message = board + res.Error()
				} else {
					message = board
				}
				w.Write([]byte(message))
				w.Flush()
				log.Printf("Sent: %s", message)

			default:
				w.Write([]byte(ErrorMessage))
				w.Flush()
				log.Printf("Sent: %s", ErrorMessage)
			}

		default:
			log.Fatalf("Receive data failed:%s", err)
			return
		}
	}
}

func StartServer() error {
	port := 3333
	players = 0
	status = InitPhase
	err := socketServer(port)
	return err
}

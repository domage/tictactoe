package game

import "fmt"

type Game struct {
	board [][]string
	turn  string
}

func NewGame() *Game {
	g := new(Game)
	g.board = make([][]string, 3)
	g.board[0] = []string{"_", "_", "_"}
	g.board[1] = []string{"_", "_", "_"}
	g.board[2] = []string{"_", "_", "_"}
	g.turn = "X"
	return g
}

func rowStatus(a string, b string, c string) (finished bool, winner string) {
	row := a + b + c
	finished = false
	winner = ""
	if row == "XXX" {
		finished = true
		winner = "X"
		return
	}

	if row == "000" {
		finished = true
		winner = "0"
		return
	}

	return
}

func BoardStatus(g *Game) (finished bool, winner string) {
	finished = false
	winner = ""
	rowFinished := false
	rowWinner := ""
	draw := true
	b := g.board

	for i := 0; i < 3; i++ {
		rowFinished, rowWinner = rowStatus(b[i][0], b[i][1], b[i][2])
		if rowFinished == true {
			finished = true
			winner = rowWinner
			return
		}

		rowFinished, rowWinner = rowStatus(b[0][i], b[1][i], b[2][i])
		if rowFinished == true {
			finished = true
			winner = rowWinner
			return
		}

		for j := 0; j < 3; j++ {
			if b[i][j] == "_" {
				draw = false
			}
		}

	}

	if draw {
		finished = true
		winner = "draw"
		return
	}

	rowFinished, rowWinner = rowStatus(b[0][0], b[1][1], b[2][2])
	if rowFinished == true {
		finished = true
		winner = rowWinner
		return
	}

	rowFinished, rowWinner = rowStatus(b[2][0], b[1][1], b[0][2])
	if rowFinished == true {
		finished = true
		winner = rowWinner
		return
	}

	return
}

func (g Game) String() string {
	s := ""
	b := g.board
	for i := 0; i < len(b); i++ {
		for j := 0; j < len(b[i]); j++ {
			s += fmt.Sprintf("%s ", b[i][j])
		}
		s += fmt.Sprintf("\n")
	}
	return s
}

func (g *Game) TakeTurn(x int, y int, mark string) error {
	b := g.board
	if mark != g.turn {
		return fmt.Errorf("it's %s-s turn", g.turn)
	}
	_, winner := BoardStatus(g)
	if winner != "" {
		return fmt.Errorf("the game is finished")
	}
	if b[x][y] != "_" {
		return fmt.Errorf("place your %s on an empty space", g.turn)
	}
	b[x][y] = mark
	if g.turn == "X" {
		g.turn = "0"
	} else {
		g.turn = "X"
	}
	return nil
}

func WhoseTurn(g *Game) string {
	return g.turn
}

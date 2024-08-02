package main

import (
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

type Game struct {
	bat *Bat
	status string
	screen tcell.Screen
}

type Bat struct {
	pos [2]int
	dir Direction
}

type Direction int

const (
	up Direction = iota
	down
	static
)

func (d Direction) string() string {
	switch d {
	case up:
		return "up"
	case down:
		return "down"
	case static:
		return "static"
	default:
		return	"no direction registered"
	}

}

func main() {
	g := newGame()

	defer g.quit()

	for g.status != "off" {
		g.screen.Show()
		switch g.bat.dir {
		case up:
			g.bat.pos[1]--
			g.bat.dir = static
		case down:
			g.bat.pos[1]++
			g.bat.dir = static
		case static:
		}
		g.draw(tcell.StyleDefault)
	}

}


func newGame() *Game {
	b := newBat()	

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Init(); err != nil {
		log.Fatal(err)
	}

	game := &Game{
		bat: b,
		status: "on",
		screen: s,
	}

	go game.playerInput()

	return game
}

func newBat() *Bat {
	maxX, maxY := getSize()
	b := &Bat{
		pos: [2]int{maxX / 2, maxY / 2},
		dir: static,
	}

	return b
}

func (g *Game) draw(style tcell.Style) {
	g.screen.Clear()	
	
	var quitText = "hit ctrl + c to quit"
	for i, r := range []rune(quitText) {
		g.screen.SetContent(i, 0, r, nil, style)
	}

	direction := g.bat.dir 
	for i, r := range []rune(direction.string()) {
		g.screen.SetContent(i, 2, r, nil, style)
	}

	g.screen.SetContent(g.bat.pos[0], g.bat.pos[1] - 1, '|', nil, style)
	g.screen.SetContent(g.bat.pos[0], g.bat.pos[1] + 0, '|', nil, style)
	g.screen.SetContent(g.bat.pos[0], g.bat.pos[1] + 1, '|', nil, style)
	g.screen.Show()

	time.Sleep(time.Millisecond * 20)
}

func (g *Game) playerInput() {

	defer g.quit()

	for {
		ev := g.screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			g.screen.Sync()
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyUp:
				g.bat.dir = up
			case tcell.KeyDown:
				g.bat.dir = down
			case tcell.KeyCtrlC:
				g.status = "off"
			}
		}
	}
}

func (g *Game) quit() {
	g.screen.Fini()
	os.Exit(0)
}

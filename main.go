package main

import (
	"log"
	"os"
	"strconv"
	"time"
	"github.com/gdamore/tcell"
)

type Game struct {
	players [2]*player
	status  string
	screen  tcell.Screen
}

type player struct {
	name  string
	score int
	bat   *Bat
}

type Bat struct {
	pos  [2]int
	dir  Direction
	size int
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
		return "no direction registered"
	}

}

func newGame() *Game {
	p1 := newPlayer(1)
	p2 := newPlayer(2)

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Init(); err != nil {
		log.Fatal(err)
	}

	game := &Game{
		players: [2]*player{p1, p2},
		status:  "on",
		screen:  s,
	}
	go game.playerInput()
	return game
}

func newPlayer(playerNum int) *player {
	lengthFromBorder := 20
	xPos, _ := getSize()
	if playerNum == 1 {
		xPos = lengthFromBorder
	} else {
		xPos -= lengthFromBorder
	}
	p := &player{
		name:  "Player " + strconv.Itoa(playerNum),
		score: 0,
		bat:   newBat(xPos),
	}
	return p
}

func newBat(xPos int) *Bat {
	_, maxY := getSize()
	b := &Bat{
		pos:  [2]int{xPos, maxY / 2},
		dir:  static,
		size: 5,
	}
	return b
}

func (b *Bat) move() {
  _, maxY := getSize()
	switch b.dir {
	case up:
		if b.pos[1] > (b.size / 2) { 
			b.pos[1]--
		}
		b.dir = static
	case down:
		if b.pos[1] < maxY - ((b.size / 2) + 1) {
			b.pos[1]++
		}
		b.dir = static
	case static:
	}
}

func (g *Game) draw(style tcell.Style) {
	g.screen.Clear()
  maxX, _ := getSize()

  g.drawText("hit 'ctrl + c' or 'q' to quit", 1, 0, style)
  g.drawText("s = ↓   d = ↑", 5, 1, style)
  g.drawText("j = ↓   k = ↑", (maxX - 20), 1, style)
	g.drawBats(style)

	time.Sleep(time.Millisecond * 20)
}

func (g *Game) drawText(text string, x int, y int, style tcell.Style) {
	for i, r := range []rune(text) {
		g.screen.SetContent(i + x, y, r, nil, style)
	}
}

func (g *Game) drawBats(style tcell.Style) {
  var yOffset int
  for index := range g.players {
		for i := 0; i < g.players[index].bat.size; i++ { 
      if i == 0 {
        yOffset = 0
      } else if i%2 != 0 {
        yOffset = -((i + 1) / 2)
      } else {
        yOffset = (i / 2)
      }
			g.screen.SetContent(g.players[index].bat.pos[0], g.players[index].bat.pos[1]+yOffset, '|', nil, style)
		}
	}
}

func (g *Game) playerInput() {

	defer g.quit()

	for {
		ev := g.screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			g.screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyUp || ev.Rune() == 'k' {
				g.players[1].bat.dir = up
			} else if ev.Key() == tcell.KeyDown || ev.Rune() == 'j' {
				g.players[1].bat.dir = down
			} else if ev.Rune() == 'd'{
        g.players[0].bat.dir = up       
      } else if ev.Rune() == 's' {
				g.players[0].bat.dir = down
      } else if ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' { //ev.Rune() for chars
				g.status = "off"
			}
		}
	}
}

func (g *Game) quit() {
	g.screen.Fini()
	os.Exit(0)
}

func main() {
	g := newGame()
	defer g.quit()
	for g.status != "off" {
		g.screen.Show()
		g.players[0].bat.move()
		g.players[1].bat.move()
		g.draw(tcell.StyleDefault)
	}
}

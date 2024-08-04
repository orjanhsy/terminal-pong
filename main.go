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
  name    string
  score   int
	bat     *Bat
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
		players: [2]*player{p1,p2},
		status: "on",
		screen: s,
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
    name: "Player " + strconv.Itoa(playerNum),
    score: 0,
    bat: newBat(xPos),
  }    
  return p
}

func newBat(xPos int) *Bat {
	_, maxY := getSize()
	b := &Bat{
		pos: [2]int{xPos, maxY / 2},
		dir: static,
	}
	return b
}

func (b *Bat) moveBat() {
  _, maxY := getSize()
  switch b.dir {
  case up:
    if b.pos[1] > 1 { // account for bat-size here
      b.pos[1]--
    }
    b.dir = static
  case down:
    if b.pos[1] < maxY - 2 {
      b.pos[1]++
    }
    b.dir = static
  case static:
  }
}

func (g *Game) draw(style tcell.Style) {
	g.screen.Clear()	
	
	var quitText = "hit 'ctrl + c' or 'q' to quit"
	for i, r := range []rune(quitText) {
		g.screen.SetContent(i, 0, r, nil, style)
	}

	direction := g.players[0].bat.dir
	for i, r := range []rune(direction.string()) {
		g.screen.SetContent(i, 2, r, nil, style)
	}
  
  // player 1
	g.screen.SetContent(g.players[0].bat.pos[0], g.players[0].bat.pos[1]-1, '|', nil, style)
	g.screen.SetContent(g.players[0].bat.pos[0], g.players[0].bat.pos[1]+0, '|', nil, style)
	g.screen.SetContent(g.players[0].bat.pos[0], g.players[0].bat.pos[1]+1, '|', nil, style)
  // player 2
	g.screen.SetContent(g.players[1].bat.pos[0], g.players[1].bat.pos[1]-1, '|', nil, style)
	g.screen.SetContent(g.players[1].bat.pos[0], g.players[1].bat.pos[1]+0, '|', nil, style)
	g.screen.SetContent(g.players[1].bat.pos[0], g.players[1].bat.pos[1]+1, '|', nil, style)
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
      if ev.Key() == tcell.KeyUp {
        g.players[0].bat.dir = up
      } else if ev.Key() == tcell.KeyDown {
        g.players[0].bat.dir = down
      } else if ev.Key() == tcell.KeyCtrlC || ev.Rune() =='q' {  //ev.Rune() for chars
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
    g.players[0].bat.moveBat()
    //g.players[1].bat.moveBat()
		g.draw(tcell.StyleDefault)
	}
}

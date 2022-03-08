package foight

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
	"math/rand"
)

var img *ebiten.Image

func init() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("gopher.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Player struct {
	name string

	x, y   float32
	dx, dy float32

	speed float32

	draw_options *ebiten.DrawImageOptions

	messages chan UpdateMessage
}

func (p *Player) Update(dt float32) {
	select {
	case um := <-p.messages:
		p.applyUpdateMessage(&um)
	default:
	}

	p.x += p.dx * dt * p.speed
	p.y += p.dy * dt * p.speed

	p.draw_options.GeoM.Reset()
	p.draw_options.GeoM.Translate(float64(p.x), float64(p.y))
}

func (p *Player) Render(screen *ebiten.Image) {
	screen.DrawImage(img, p.draw_options)

}

func (p *Player) applyUpdateMessage(um *UpdateMessage) {
	p.dx = float32(um.dx)
	p.dy = float32(um.dy)

}

func (g *Game) AddPlayer(name string) *Player {
	player := &Player{
		name: name,

		x: 100 + rand.Float32()*100.,
		y: 100 + rand.Float32()*100.,

		speed: 100,

		draw_options: &ebiten.DrawImageOptions{},

		messages: make(chan UpdateMessage, 1024),
	}

	g.players = append(g.players, player)

	return player
}

func (g *Game) indexOfPlayer(p *Player) int {
	for k, v := range g.players {
		if p == v {
			return k
		}
	}
	return -1 //not found.
}

func (g *Game) RemovePlayer(p *Player) {
	idx := g.indexOfPlayer(p)
	log.Printf("remove player idx", idx)
	if idx > -1 {
		g.players = append(g.players[:idx], g.players[idx+1:]...)
	}
}

package foight

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
	"math"
	"math/rand"
)

var img *ebiten.Image
var img_w, img_h float64

func init() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("ploier.png")
	if err != nil {
		log.Fatal(err)
	}

	img_w = float64(img.Bounds().Dx())
	img_h = float64(img.Bounds().Dy())
}

type Player struct {
	name  string
	color uint32

	x, y   float32
	dx, dy float32
	angle  float64

	speed float32

	draw_options *ebiten.DrawImageOptions

	messages chan UpdateMessage
}

func (p *Player) Update(dt float32) {
	p.readMessages()

	p.x += p.dx * dt * p.speed
	p.y += p.dy * dt * p.speed

	p.draw_options.GeoM.Reset()
	p.draw_options.GeoM.Translate(img_h/-2, img_h/-2)

	if p.dx != 0 && p.dy != 0 {
		p.angle = math.Atan2(float64(p.dy), float64(p.dx)) + math.Pi/2
	}
	p.draw_options.GeoM.Rotate(p.angle)
	//log.Println(angle)

	p.draw_options.GeoM.Translate(float64(p.x), float64(p.y))
}

func (p *Player) Render(screen *ebiten.Image) {
	screen.DrawImage(img, p.draw_options)

}

func (p *Player) SetColor(color uint32) {

	p.draw_options.ColorM.Scale(0, 0, 0, 1)

	r := float64((color&0xFF0000)>>4) / 0xff
	g := float64((color&0x00FF00)>>2) / 0xff
	b := float64((color&0x0000FF)>>0) / 0xff

	p.draw_options.ColorM.Translate(r, g, b, 0)
}

func (p *Player) readMessages() {
	for {
		select {
		case um := <-p.messages:
			p.applyUpdateMessage(&um)
		default:
			return
		}
	}
}

func (p *Player) applyUpdateMessage(um *UpdateMessage) {
	p.dx = float32(um.dx) / 50
	p.dy = float32(um.dy) / 50

}

func (g *Game) AddPlayer(name string, color uint32) *Player {
	player := &Player{
		name:  name,
		color: color,

		x: 100 + rand.Float32()*300,
		y: 100 + rand.Float32()*300,

		speed: 100,

		draw_options: &ebiten.DrawImageOptions{},

		messages: make(chan UpdateMessage, 1024),
	}

	player.SetColor(color)

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
	log.Printf("remove player [%s] idx %i", p.name, idx)

	if idx > -1 {
		g.players = append(g.players[:idx], g.players[idx+1:]...)
	}
}

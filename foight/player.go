package foight

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jakecoffman/cp"
	_ "image/png"
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

	x, y   float64
	dx, dy float64
	angle  float64

	speed float64

	body  *cp.Body
	shape *cp.Shape

	draw_options *ebiten.DrawImageOptions
	messages     chan UpdateMessage
}

func (p *Player) UpdateInputs(dt float64) {
	p.readMessages()

	tx := p.dx * dt * p.speed
	ty := p.dy * dt * p.speed

	p.body.SetVelocity(tx, ty)
}

func (p *Player) Update(dt float64) {
	position := p.body.Position()
	p.x = position.X
	p.y = position.Y

	p.draw_options.GeoM.Reset()
	p.draw_options.GeoM.Translate(img_w/-2, img_h/-2)

	if p.dx != 0 && p.dy != 0 {
		p.angle = math.Atan2(float64(p.dy), float64(p.dx)) + math.Pi/2
	}
	p.draw_options.GeoM.Rotate(p.angle)

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
	p.dx = float64(um.dx) / 50
	p.dy = float64(um.dy) / 50
}

func (g *Game) AddPlayer(name string, color uint32) *Player {

	player := &Player{
		name:  name,
		color: color,

		x: 100 + rand.Float64()*300,
		y: 100 + rand.Float64()*300,

		speed: 10000,

		draw_options: &ebiten.DrawImageOptions{},

		messages: make(chan UpdateMessage, 1024),
	}

	player.SetColor(color)

	g.players = append(g.players, player)

	body := g.space.AddBody(cp.NewBody(1, cp.INFINITY))
	body.SetPosition(cp.Vector{player.x, player.y})

	shape := g.space.AddShape(cp.NewCircle(body, img_w/2, cp.Vector{img_w / 2, img_h / 2}))
	shape.SetElasticity(0)
	shape.SetFriction(0)
	shape.SetCollisionType(1)

	player.body = body
	player.shape = shape

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

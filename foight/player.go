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

  x, y float32
  dx, dy float32

  speed float32

  draw_options *ebiten.DrawImageOptions

  messages chan UpdateMessage;
}

func (p *Player) Update (dt float32) {

  p.x += p.dx * dt
  p.y += p.dy * dt

  p.draw_options.GeoM.Reset()
  p.draw_options.GeoM.Translate(float64(p.x), float64(p.y))
}

func (p *Player) Render (screen *ebiten.Image) {
  screen.DrawImage(img, p.draw_options)

}

func (g *Game) AddPlayer(name string) *Player {
  player := &Player{
    name: name,

    x: 100 + rand.Float32() * 100.,
    y: 100 + rand.Float32() * 100.,

    speed: 100,

    draw_options: &ebiten.DrawImageOptions{},

    messages: make(chan UpdateMessage, 1024),
  }

  g.players = append(g.players, player)

  return player
}



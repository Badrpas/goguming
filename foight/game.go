package foight

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"

	"image/color"
)

type Game struct {
	players      []*Player
	last_message string

	space *cp.Space
}

func NewGame() *Game {
	game := &Game{
		space: cp.NewSpace(),
	}

	return game
}

func (g *Game) Layout(outWidth, outHeight int) (width, height int) {
	return 800, 600
}

func (g *Game) Update() error {
	var dt = 1. / 60. // Really disliking that

	for _, player := range g.players {
		player.UpdateInputs(dt)
	}

	g.space.Step(dt)

	for _, player := range g.players {
		player.Update(dt)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for _, player := range g.players {
		player.Render(screen)
	}

}

package foight

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	_ "image/png"
)


type Game struct {
	players []*Player;
	last_message string

}

func(g *Game) Layout (outWidth, outHeight int) (width, height int) {
	return 800, 600
}

func (g *Game) Update () error {
	var dt float32 = 1. / 60. // Really disliking that

	g.readMessages()

	for _, player := range g.players {
		player.Update(dt);
	}

	return nil
}


func (g *Game) Draw (screen *ebiten.Image) {
	screen.Fill(color.Black)

	for _, player := range g.players {
		player.Render(screen);
	}

	if g.last_message != "" {
		ebitenutil.DebugPrintAt(screen, g.last_message, 100, 100)
	}
}
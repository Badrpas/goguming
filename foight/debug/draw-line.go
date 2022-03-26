package debug

import (
	"game/foight/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"image/color"
)

func DrawLine(image *ebiten.Image, p1, p2 cp.Vector, c color.Color) {
	line := util.Makeline(int(p1.X), int(p1.Y), int(p2.X), int(p2.Y))
	for _, point := range line {
		image.Set(point.X, point.Y, c)
	}
}

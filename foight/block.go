package foight

import (
	imagestore "game/img"
	"github.com/jakecoffman/cp"
)

type Block struct {
	*Entity
}

func NewBlock(x, y float64) *Block {
	block := &Block{
		Entity: NewEntity(
			x, y,
			nil,
			nil,
			imagestore.Images["block.png"],
		),
	}

	return block
}

func (b *Block) Init(game *Game) int32 {
	space := game.Space

	b.Body = space.AddBody(cp.NewBody(1, cp.INFINITY))
	b.Body.SetPosition(cp.Vector{b.X, b.Y})
	b.Body.SetType(cp.BODY_STATIC)

	b.Shape = space.AddShape(cp.NewBox(b.Body, float64(b.Img.Bounds().Dx()), float64(b.Img.Bounds().Dy()), 0))
	b.Shape.SetElasticity(1)
	b.Shape.SetFriction(1)

	return game.AddEntity(b.Entity)
}

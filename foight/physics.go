package foight

import "github.com/jakecoffman/cp"

func AddCirclePhysicsToEntity(g *Game, e *Entity) (*cp.Body, *cp.Shape) {
	body := g.Space.AddBody(cp.NewBody(1, cp.INFINITY))
	body.SetPosition(cp.Vector{e.X, e.Y})
	body.UserData = e

	radius := float64(e.Img.Bounds().Dx() / 2)
	shape := g.Space.AddShape(cp.NewCircle(body, radius, cp.Vector{}))
	shape.SetElasticity(0.3)
	shape.SetFriction(0)
	shape.SetCollisionType(1)

	if e.Team != 0 {
		shape.Filter.Group = uint(e.ID)
	}

	e.Body = body
	e.Shape = shape

	return body, shape
}

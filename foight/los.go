package foight

import (
	"game/foight/pathfind"
	"game/foight/util"
)

func isInLos(e1, e2 *Entity, nav *pathfind.Nav) bool {
	size := nav.GetTileSize()
	pos1 := e1.GetPosition().Mult(1 / size)
	pos2 := e2.GetPosition().Mult(1 / size)
	line := util.Makeline(int(pos1.X), int(pos1.Y), int(pos2.X), int(pos2.Y))

	for _, point := range line {
		if tile := nav.GetTileAt(point.X, point.Y); tile != nil && tile.Type == pathfind.NavTileWall {
			return false
		}
	}

	return true
}

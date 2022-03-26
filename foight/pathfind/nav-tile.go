package pathfind

import (
	"github.com/beefsack/go-astar"
	"log"
	"math"
)

type NavTilePos struct {
	X, Y int
}

type NavTileType uint8

const (
	NavTileEmpty NavTileType = 1 << iota
	NavTileWall
	NavTileFilledGap
)

func init() {
	log.Println(NavTileEmpty, NavTileWall, NavTileFilledGap)
}

type NavTile struct {
	NavTilePos
	Type NavTileType
	Nav  *Nav
}

func (t *NavTile) PathNeighborCost(to astar.Pather) float64 {
	return math.Pow(float64(to.(*NavTile).X-t.X), 2) + math.Pow(float64(to.(*NavTile).X-t.X), 2)
}

func (t *NavTile) PathEstimatedCost(to astar.Pather) float64 {
	return t.PathNeighborCost(to)
}

var NEIGHBOR_DELTAS = []NavTilePos{
	{-1, +0},
	{+1, +0},
	{+0, -1},
	{+0, +1},

	{-1, -1},
	{-1, +1},
	{+1, -1},
	{+1, +1},
}

func (t *NavTile) PathNeighbors() []astar.Pather {
	n := make([]astar.Pather, len(NEIGHBOR_DELTAS))
	idx := 0

	for _, d := range NEIGHBOR_DELTAS {
		x, y := t.X+d.X, t.Y+d.Y
		if tile := t.Nav.GetTileAtSafe(x, y); tile != nil && tile.Type == NavTileEmpty {
			n[idx] = tile
			idx++
		}
	}

	return n[:idx]
}

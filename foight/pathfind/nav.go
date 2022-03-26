package pathfind

import (
	"game/foight/util"
)

type Nav struct {
	Width, Height int

	tile_size float64
	dirty     bool
	findFn    func(sx, sy, tx, ty int) []NavPathNode

	tiles     map[int]map[int]*NavTile
	gap_tiles []*NavTile
}

func NewNav(w, h int) *Nav {
	return &Nav{
		dirty: true,
		tiles: map[int]map[int]*NavTile{},
	}
}

func (nav *Nav) SetSize(w, h int) {
	nav.Width, nav.Height = w, h
	nav.dirty = true
}
func (nav *Nav) SetTileSize(x float64) {
	nav.tile_size = x
	nav.dirty = true
}
func (nav *Nav) GetTileSize() float64 {
	return nav.tile_size
}

func (nav *Nav) GetTileAt(x, y int) *NavTile {
	if nav.tiles[x] == nil {
		return nav.SetTileAt(x, y, NavTileEmpty)
	}

	tile := nav.tiles[x][y]
	if tile == nil {
		return nav.SetTileAt(x, y, NavTileEmpty)
	}

	return tile
}

func (nav *Nav) GetTileAtSafe(x, y int) *NavTile {
	if nav.tiles[x] == nil {
		return nil
	}

	tile := nav.tiles[x][y]
	if tile == nil {
		return nil
	}

	return tile
}

func (nav *Nav) SetTileAt(x, y int, t NavTileType) *NavTile {
	if nav.tiles[x] == nil {
		nav.tiles[x] = map[int]*NavTile{}
	}

	tile := nav.tiles[x][y]

	if nav.tiles[x][y] == nil {
		tile = &NavTile{
			NavTilePos{x, y},
			t,
			nav,
		}
		nav.tiles[x][y] = tile
	} else {
		tile.Type = t
	}

	return tile
}

func (nav *Nav) FixHolesWithActorSize(size int) {
	for x, row := range nav.tiles {
		for y, node := range row {
			for i := 1; i <= size; i++ {
				neighbors := nav.findNeighborsInRadius(node, i)
				for _, neighbor := range neighbors {
					for _, p := range util.Makeline(x, y, neighbor.X, neighbor.Y) {
						if nav.GetTileAt(p.X, p.Y).Type == NavTileEmpty {
							nav.SetGapWall(p.X, p.Y)
						}
					}
				}
			}
		}
	}
}

func (nav *Nav) findNeighborsInRadius(node *NavTile, radius int) []*NavTile {
	check := func(x, y int) (bool, *NavTile) {
		if nav.tiles[x] == nil {
			return false, nil
		}
		if nav.tiles[x][y] != nil {
			return true, nav.tiles[x][y]
		}
		return false, nil
	}

	x, y := node.X, node.Y
	nodes := []*NavTile{}

	for i := 0; i <= radius; i++ {

		if ok, node := check(x+radius, y+i); ok {
			nodes = append(nodes, node)
		}
		if ok, node := check(x+radius, y-i); ok {
			nodes = append(nodes, node)
		}
		if ok, node := check(x-radius, y+i); ok {
			nodes = append(nodes, node)
		}
		if ok, node := check(x-radius, y-i); ok {
			nodes = append(nodes, node)
		}

		if ok, node := check(x+i, y+radius); ok {
			nodes = append(nodes, node)
		}
		if ok, node := check(x-i, y+radius); ok {
			nodes = append(nodes, node)
		}
		if ok, node := check(x+i, y-radius); ok {
			nodes = append(nodes, node)
		}
		if ok, node := check(x-i, y-radius); ok {
			nodes = append(nodes, node)
		}

	}

	return nodes
}

func (nav *Nav) SetWall(x, y int) {
	nav.SetTileAt(x, y, NavTileWall)
	nav.dirty = true
}

func (nav *Nav) SetGapWall(x, y int) {
	nav.SetTileAt(x, y, NavTileFilledGap)
	nav.dirty = true
}

func (nav *Nav) IsWallAt(x, y int) bool {
	return nav.GetTileAt(x, y).Type != NavTileEmpty
}

var _IterationPrealloc = make([]*NavTile, 100_000)

func (nav *Nav) IterateTiles(t NavTileType) []*NavTile {
	idx := 0

	for _, col := range nav.tiles {
		for _, tile := range col {
			if tile.Type&t != 0 {
				_IterationPrealloc[idx] = tile
				idx++
			}
		}
	}

	return _IterationPrealloc[:idx]
}

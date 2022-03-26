package pathfind

import (
	"game/foight/util"
	"github.com/beefsack/go-astar"
	"github.com/jakecoffman/cp"
	"log"
)

type NavPathNode struct {
	cp.Vector
}

func (nav *Nav) Init() error {
	for x := 0; x < nav.Width; x++ {
		for y := 0; y < nav.Height; y++ {
			nav.GetTileAt(x, y)
		}
	}

	nav.findFn = func(sx, sy, tx, ty int) []NavPathNode {
		start_time := util.TimeNow()
		defer func() {
			total := util.TimeNow() - start_time
			log.Println(total)
		}()

		if nav.IsWallAt(tx, ty) {
			return nil
		}
		p, _, ok := astar.Path(nav.GetTileAt(sx, sy), nav.GetTileAt(tx, ty))
		if !ok {
			return nil
		}

		l := len(p)
		path := make([]NavPathNode, l)
		idx := 0

		for i := 1; i < l-1; i++ {
			prev := p[i-1].(*NavTile)
			c := p[i].(*NavTile)
			next := p[i+1].(*NavTile)

			path[idx].X = float64(c.X)*nav.tile_size - nav.tile_size/2
			path[idx].Y = float64(c.Y)*nav.tile_size - nav.tile_size/2
			idx++
			continue

			if next.X != prev.X && next.Y != prev.Y {
				path[idx].X = float64(c.X)*nav.tile_size - nav.tile_size/2
				path[idx].Y = float64(c.Y)*nav.tile_size - nav.tile_size/2
				idx++
			}

			// last
			if i == l-2 {
				path[idx].X = float64(next.X)*nav.tile_size - nav.tile_size/2
				path[idx].Y = float64(next.Y)*nav.tile_size - nav.tile_size/2
				idx++
			}
		}

		path = path[:idx]
		for i := 0; i < len(path)/2; i++ {
			path[i], path[len(path)-1-i] = path[len(path)-1-i], path[i]
		}

		return path
	}

	nav.dirty = false

	return nil
}

func (nav *Nav) Find(sx, sy, tx, ty int) []NavPathNode {

	if nav.dirty {
		err := nav.Init()
		if err != nil {
			return nil
		}
	}

	return nav.findFn(sx, sy, tx, ty)
}

func (nav *Nav) FindVec(s cp.Vector, t cp.Vector) []NavPathNode {
	s = s.Mult(1 / nav.tile_size)
	t = t.Mult(1 / nav.tile_size)
	return nav.Find(int(s.X), int(s.Y), int(t.X), int(t.Y))
}

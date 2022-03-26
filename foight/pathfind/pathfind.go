package pathfind

import (
  "game/foight/util"
  "github.com/jakecoffman/cp"
  "github.com/jpierer/astar"
  "log"
  "math"
)

type NavPathNode struct {
  cp.Vector
}

type Nav struct {
  astar.Config
  tile_size  float64
  dirty      bool
  findFn     func(sx, sy, tx, ty int) []NavPathNode
  FixedHoles []astar.Node
  cache      map[int]map[int]*astar.Node
}

func NewNav(w, h int) *Nav {
  return &Nav{
    Config: astar.Config{
      GridWidth: w,
      GridHeight: h,
    },
    dirty: true,
    FixedHoles: nil,
    cache: map[int]map[int]*astar.Node{},
  }
}

func (nav *Nav) SetSize(w, h int) {
  nav.GridWidth, nav.GridHeight = w, h
  nav.dirty = true
}
func (nav *Nav) SetTileSize(x float64) {
  nav.tile_size = x
  nav.dirty = true
}
func (nav *Nav) GetTileSize() float64 {
  return nav.tile_size
}

func (nav *Nav) FixHolesWithActorSize(size int) {
  for x, row := range nav.cache {
    for y, node := range row {
      for i := 1; i <= size; i++ {
        neighbors := nav.findNeighborsInRadius(node, i)
        for _, neighbor := range neighbors {
          for _, p := range makeline(x, y, neighbor.X, neighbor.Y) {
            nav.SetGapWall(p.x, p.y)
          }
        }
      }
    }
  }
}
type point struct{ x, y int }

func makeline(x1, y1, x2, y2 int) []point {
  if x1 > x2 {
    x1, x2 = x2, x1
    y1, y2 = y2, y1
  }

  dx := x2 - x1
  dy := y2 - y1
  adx, ady := math.Abs(float64(dx)), math.Abs(float64(dy))

  i := 0
  if adx > ady {
    if x1 > x2 {
      x1, x2 = x2, x1
      y1, y2 = y2, y1
    }

    dx = x2 - x1
    dy = y2 - y1

    max := int(math.Max(float64(dx), float64(dy)))
    points := make([]point, max+1)

    for x := x1; x <= x2; x++ {
      points[i] = point{ x: x, y: y1 + dy * (x - x1) / dx }
      i++
    }

    return points
  } else {
    if y1 > y2 {
      x1, x2 = x2, x1
      y1, y2 = y2, y1
    }

    dx = x2 - x1
    dy = y2 - y1

    max := int(math.Max(float64(dx), float64(dy)))
    points := make([]point, max+1)

    for y := y1; y <= y2; y++ {
      points[i] = point{ x: x1 + dx * (y - y1) / dy, y: y }
      i++
    }

    return points
  }
}

func (nav *Nav) findNeighborsInRadius(node *astar.Node, radius int) []*astar.Node {
  check := func(x,y int) (bool, *astar.Node) {
    if nav.cache[x] == nil {
      return false, nil
    }
    if nav.cache[x][y] != nil {
      return true, nav.cache[x][y]
    }
    return false, nil
  }

  x, y := node.X, node.Y
  nodes := []*astar.Node{}

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

func sign(x int) int {
  if x < 0 {
    return -1
  }
  return 1
}

func (nav *Nav) SetWall(x, y int) {
  if nav.cache[x] == nil {
    nav.cache[x] = make(map[int]*astar.Node)
  }

  if nav.cache[x][y] != nil {
    return
  }

  node := astar.Node{
    X: x,
    Y: y,
  }
  nav.InvalidNodes = append(nav.InvalidNodes, node)
  nav.cache[x][y] = &node
  nav.dirty = true
}

func (nav *Nav) SetGapWall(x, y int) {
  if nav.cache[x] == nil {
    nav.cache[x] = make(map[int]*astar.Node)
  }

  if nav.cache[x][y] != nil {
    return
  }

  for _, node := range nav.FixedHoles {
    if node.X == x && node.Y == y {
      return
    }
  }

  node := astar.Node{
    X: x,
    Y: y,
  }
  nav.FixedHoles = append(nav.FixedHoles, node)
  nav.dirty = true
}

func (nav *Nav) IsWallAt(x, y int) bool {
  for _, hole := range nav.FixedHoles {
    if hole.X == x && hole.Y == y {
      return true
    }
  }
  return nav.cache[x] != nil && nav.cache[x][y] != nil
}


func (nav *Nav) Init() error {
  config := nav.Config
  config.InvalidNodes = append(config.InvalidNodes, nav.FixedHoles...)
  algo, err := astar.New(config)
  if err != nil {
    return err
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
    p, err := algo.FindPath(astar.Node{X: sx, Y: sy}, astar.Node{X: tx, Y: ty}, 50)
    if err != nil {
      return nil
    }

    l := len(p)
    path := make([]NavPathNode, l)
    idx := 0

    for i := 1; i < l - 1; i++ {
      prev := p[i-1]
      c := p[i]
      next := p[i+1]

      if next.X != prev.X && next.Y != prev.Y {
        path[idx].X = float64(c.X) * nav.tile_size - nav.tile_size / 2
        path[idx].Y = float64(c.Y) * nav.tile_size - nav.tile_size / 2
        idx++
      }

      // last
      if i == l-2 {
        path[idx].X = float64(next.X) * nav.tile_size - nav.tile_size / 2
        path[idx].Y = float64(next.Y) * nav.tile_size - nav.tile_size / 2
        idx++
      }
    }

    path = path[:idx]
    for i := 0; i < len(path)/2; i++ {
      path[i], path[len(path) - 1 - i] = path[len(path) - 1 - i], path[i]
    }

    return path
  }

  nav.dirty = false

  return nil
}

func (nav *Nav) Find(sx, sy, tx, ty int) []NavPathNode {
  if sx < 0 || sx > nav.GridWidth || sy < 0 || sy > nav.GridHeight {
    return nil
  }

  if nav.dirty {
    err := nav.Init()
    if err != nil {
      return nil
    }
  }

  return nav.findFn(sx, sy, tx, ty)
}

func (nav *Nav) FindVec(s cp.Vector, t cp.Vector) []NavPathNode {
  s = s.Mult(1/nav.tile_size)
  t = t.Mult(1/nav.tile_size)
  return nav.Find(int(s.X), int(s.Y), int(t.X), int(t.Y))
}

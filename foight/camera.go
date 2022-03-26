package foight

import (
  "github.com/hajimehoshi/ebiten/v2"
  "github.com/jakecoffman/cp"
  camera "github.com/melonfunction/ebiten-camera"
  "image"
  "math"
)

const (
  ZOOM_MAX = 2.0
  ZOOM_MIN = 0.7
)

func UpdateCamera(game *Game, dt float64) {
  return
  var players []*Entity
  for _, entity := range game.Entities {
    _, ok := entity.Holder.(*Player)
    if ok {
      players = append(players, entity)
    }
  }

  player_count := len(players)
  if player_count == 0 {
    return
  }

  f := players[0]
  left, right, top, bot := f.X, f.X, f.Y, f.Y

  avg := cp.Vector{f.X, f.Y}
  for i := 1; i < player_count; i++ {
    player := players[i]
    avg = avg.Add(player.GetPosition())

    if player.X < left {
      left = player.X
    } else if player.X > right {
      right = player.X
    }
    if player.Y < top {
      top = player.Y
    } else if player.Y > bot {
      bot = player.Y
    }
  }

  target := avg.Mult(1. / float64(player_count))
  current := cp.Vector{game.Camera.X, game.Camera.Y}
  diff := target.Sub(current)
  length := diff.Length()
  speed := math.Max(700, length * 0.01)
  delta := speed * dt

  if delta >= length {
    game.Camera.SetPosition(target.X, target.Y)
  } else {
    step := current.Add(diff.Normalize().Mult(delta))
    game.Camera.SetPosition(step.X, step.Y)
  }

  {
    w := right - left
    h := bot - top
    const (
      PADDING = 200
    )

    inv_w := w * 1.2 / float64(game.Camera.Width-PADDING)
    inv_h := h * 1.2 / float64(game.Camera.Height-PADDING)

    inv := math.Max(inv_w, inv_h)
    inv = cp.Clamp(inv, ZOOM_MIN, ZOOM_MAX)

    zoom := 1 / inv

    const ZOOM_SPEED = 0.4
    current_zoom := game.Camera.Scale
    diff := zoom - current_zoom
    delta := ZOOM_SPEED * dt * sign(diff)
    if math.Abs(delta) < math.Abs(diff) {
      zoom = current_zoom + delta
    }

    SetZoom(game.Camera, zoom)
  }
}


func SetZoom (c *camera.Camera, zoom float64) {
  c.Scale = zoom
  if c.Scale <= 0.01 {
    c.Scale = 0.01
  }
  Resize(c, c.Width, c.Height)
}

var _init_surf = false
func Resize(c *camera.Camera, w, h int) {
  c.Width = w
  c.Height = h
  var (
    DEF_WIDTH = ((w) * ZOOM_MAX) //* 2.0
    DEF_HEIGHT = ((h) * ZOOM_MAX) //* 2.0
  )
  if !_init_surf {
    _init_surf = true
    c.Surface.Dispose()
    c.Surface = ebiten.NewImage(DEF_WIDTH, DEF_HEIGHT)
  }
}


func Blit(c *camera.Camera, screen *ebiten.Image) {
  op := &ebiten.DrawImageOptions{}
  w, h := float64(c.Width) / c.Scale, float64(c.Height) / c.Scale
  cx := float64(w) / 2.0
  cy := float64(h) / 2.0

  op.GeoM.Translate(-cx, -cy)
  op.GeoM.Scale(c.Scale, c.Scale)
  op.GeoM.Translate(cx*c.Scale, cy*c.Scale)

  r := image.Rectangle{
    Min: image.Point{0, 0},
    Max: image.Point{int(w)+1, int(h)+1},
  }
  subImage := c.Surface.SubImage(r)
  img := subImage.(*ebiten.Image)
  screen.DrawImage(img, op)
}



func sign(x float64) float64 {
  if x < 0 {
    return -1
  }
  return 1
}
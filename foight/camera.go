package foight

import (
  "github.com/jakecoffman/cp"
  "math"
)

func UpdateCamera(game *Game, dt float64) {
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

    inv = cp.Clamp(inv, 0.7, 3)

    zoom := 1 / inv

    const ZOOM_SPEED = 0.4
    current_zoom := game.Camera.Scale
    diff := zoom - current_zoom
    delta := ZOOM_SPEED * dt * sign(diff)
    if math.Abs(delta) < math.Abs(diff) {
      zoom = current_zoom + delta
    }

    game.Camera.SetZoom(zoom)
    //log.Println(zoom, inv, inv_w, inv_h)
  }
}

func sign(x float64) float64 {
  if x < 0 {
    return -1
  }
  return 1
}
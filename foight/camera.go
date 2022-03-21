package foight

import "github.com/jakecoffman/cp"

func UpdateCamera(game *Game) {
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

  //left, right, top, bot := 0,0,0,0

  avg := cp.Vector{players[0].X, players[0].Y}
  for i := 1; i < player_count; i++ {
    player := players[i]
    avg = avg.Add(player.GetPosition())
  }

  avg = avg.Mult(1. / float64(player_count))

  game.Camera.SetPosition(avg.X, avg.Y)
}

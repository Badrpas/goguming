package levelmap

import (
	"game/foight"
	"game/foight/pathfind"
	imagestore "game/img"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/lafriks/go-tiled"
	"image/color"
	"log"
	"math/rand"
	"strings"
)

func LoadToGameTiled(path string, game *foight.Game) error {
	file, err := tiled.LoadFile(path)
	if err != nil {
		log.Println(err)
		return err
	}

	layer := file.Layers[0]
	tiles := layer.Tiles

	game.Nav.SetSize(file.Width, file.Height)
	game.Nav.SetTileSize(float64(file.TileWidth))

	cell_w := float64(file.TileWidth)
	cell_h := float64(file.TileHeight)

	for idx, tile := range tiles {
		if tile.Nil {
			continue
		}

		idx_x := (idx % (file.Width))
		idx_y:= (idx / (file.Width))
		x := cell_w * float64(idx_x)
		y := cell_h * float64(idx_y)

		img := GetImage(tile)
		b := foight.NewBlock(x, y, img)
		b.Init(game)

		c := uint8(155 + rand.Int()%100)
		b.SetColor(color.RGBA{c, c, c, 255})

		SetWallAroundPoint(game.Nav, idx_x, idx_y, 1)
	}

	game.Nav.FixHolesWithActorSize(2)
	game.Nav.Init()

	for _, objectGroup := range file.ObjectGroups {
		var points = make([]cp.Vector, len(objectGroup.Objects))
		switch objectGroup.Name {
		case "player_spawn_points":
			game.PlayerSpawnPoints = points
		case "item_spawn_points":
			game.ItemSpawnPoints = points
		default:
			log.Println("Unknown object group name", objectGroup.Name)
			continue
		}

		for i, info := range objectGroup.Objects {
			points[i] = cp.Vector{info.X, info.Y}
		}
	}

	return nil
}

func SetWallAroundPoint(nav *pathfind.Nav, x, y, radius int) {
	//nav.SetWall(x, y)

	for i := 0; i <= radius; i++ {
		nav.SetWall(x+radius, y+i)
		nav.SetWall(x+radius, y-i)
		nav.SetWall(x-radius, y+i)
		nav.SetWall(x-radius, y-i)

		nav.SetWall(x+i, y+radius)
		nav.SetWall(x-i, y+radius)
		nav.SetWall(x+i, y-radius)
		nav.SetWall(x-i, y-radius)
	}

}

func GetImage(t *tiled.LayerTile) *ebiten.Image {
	for _, prototile := range t.Tileset.Tiles {
		if prototile.ID == t.ID {
			img_name := prototile.Image.Source
			img_name = strings.Replace(img_name, "../img/", "", 1)
			img, ok := imagestore.Images[img_name]
			if !ok {
				print("Can't find image", img_name)
			}
			return img
		}
	}
	return nil
}

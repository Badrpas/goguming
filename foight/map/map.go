package levelmap

import (
	"game/foight"
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
	gameMap, err := tiled.LoadFile(path)
	if err != nil {
		log.Println(err)
		return err
	}

	layer := gameMap.Layers[0]
	tiles := layer.Tiles

	cell_w := float64(gameMap.TileWidth)
	cell_h := float64(gameMap.TileHeight)

	for idx, tile := range tiles {
		if tile.Nil {
			continue
		}

		x := cell_w * float64(idx%(gameMap.Width))
		y := cell_h * float64(idx/(gameMap.Width))

		img := GetImage(tile)
		b := foight.NewBlock(x, y, img)
		b.Init(game)

		c := uint8(155 + rand.Int()%100)
		b.SetColor(color.RGBA{c, c, c, 255})
	}

	for _, objectGroup := range gameMap.ObjectGroups {
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

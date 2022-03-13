package imagestore

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/png"
	"io/fs"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

var Images = make(map[string]*ebiten.Image)

func init() {
	_ = filepath.WalkDir("img", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		is_png, _ := regexp.MatchString("\\.png$", path)
		if !is_png {
			return nil
		}

		{
			img, _, err := ebitenutil.NewImageFromFile(path)
			if err != nil {
				return err
			}

			key := strings.ReplaceAll(path, "\\", "/")
			key = strings.Replace(key, "img/", "", 1)
			Images[key] = img

			log.Println("Loaded image", key)
		}

		return nil
	})
}

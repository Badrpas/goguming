package util

import (
  "github.com/hajimehoshi/ebiten/v2"
  "image/color"
)

func SetDrawOptsColor(options *ebiten.DrawImageOptions, color color.Color) {
  if color == nil {
    return
  }

  options.ColorM.Scale(0, 0, 0, 1)

  rb, gb, bb, _ := color.RGBA()
  r := float64(rb) / 0xFFFF
  g := float64(gb) / 0xFFFF
  b := float64(bb) / 0xFFFF

  options.ColorM.Translate(r, g, b, 0)
}

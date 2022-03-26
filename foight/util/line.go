package util

import "math"

type point struct {
	X, Y int
}

var _line_buf []point

func init() {
	_line_buf = make([]point, 1024*4)
}

func Makeline(x1, y1, x2, y2 int) []point {
	if x1 > x2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}

	dx := x2 - x1
	dy := y2 - y1
	adx, ady := math.Abs(float64(dx)), math.Abs(float64(dy))
	if dx == 0 && dy == 0 {
		return nil
	}

	i := 0
	if adx > ady {
		if x1 > x2 {
			x1, x2 = x2, x1
			y1, y2 = y2, y1
		}

		dx = x2 - x1
		dy = y2 - y1

		//max := int(math.Max(float64(dx), float64(dy)))
		points := _line_buf

		for x := x1; x <= x2; x++ {
			points[i].X, points[i].Y = x, y1+dy*(x-x1)/dx
			i++
		}

		return points[:i]
	} else {
		if y1 > y2 {
			x1, x2 = x2, x1
			y1, y2 = y2, y1
		}

		dx = x2 - x1
		dy = y2 - y1

		//max := int(math.Max(float64(dx), float64(dy)))
		points := _line_buf

		for y := y1; y <= y2; y++ {
			points[i].X, points[i].Y = x1+dx*(y-y1)/dy, y
			i++
		}

		return points[:i]
	}
}

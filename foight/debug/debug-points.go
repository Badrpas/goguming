package debug

import (
	"github.com/jakecoffman/cp"
	"time"
)

type DebugPoint struct {
	cp.Vector

	ttl int64
}

var Points []*DebugPoint

func AddDebugPoint(pos cp.Vector, ms int64) {
	Points = append(Points, &DebugPoint{
		pos,
		time.Now().UnixMilli() + ms,
	})
}

func Update() {
	now := time.Now().UnixMilli()
	for idx, point := range Points {
		if now > point.ttl {
			Points[idx] = nil
		}
	}

	idx := 0
	for _, point := range Points {
		if point != nil {
			Points[idx] = point
			idx++
		}
	}

	Points = Points[:idx]
}

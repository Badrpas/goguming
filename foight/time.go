package foight

import "time"

func TimeNow() int64 {
	return time.Now().UnixMilli()
}

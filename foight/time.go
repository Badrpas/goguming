package foight

import "time"

func TimeNow() int64 {
	return time.Now().UnixMilli()
}

type Timeout struct {
	cb      func()
	emit_at int64
}

type TimeHolder struct {
	timers []Timeout
}

func (t *TimeHolder) SetTimeout(cb func(), ms int64) {
	t.timers = append(t.timers, Timeout{cb, TimeNow() + ms})
}

func (t *TimeHolder) Update() {
	now := TimeNow()
	idx := 0

	for _, to := range t.timers {
		if now >= to.emit_at {
			to.cb()
		} else {
			t.timers[idx] = to
			idx++
		}
	}

	t.timers = t.timers[:idx]
}

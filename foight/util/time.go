package util

import (
	"time"
)

func TimeNow() int64 {
	return time.Now().UnixMilli()
}

var _id_counter uint64

type Timeout struct {
	cb      func()
	emit_at int64
	id      uint64
}

type Interval struct {
	fn         func()
	id         uint64
	timeout_id uint64
}

type TimeHolder struct {
	timeouts     []*Timeout
	new_timeouts []*Timeout
	intervals    []*Interval
}

func (t *TimeHolder) SetTimeout(cb func(), ms int64) uint64 {
	_id_counter++
	id := _id_counter
	t.new_timeouts = append(t.new_timeouts, &Timeout{cb, TimeNow() + ms, id})
	return id
}

func (t *TimeHolder) ClearTimeout(id uint64) {
	for idx, to := range t.new_timeouts {
		if to.id == id {
			t.new_timeouts = append(t.new_timeouts[:idx], t.new_timeouts[idx+1:]...)
			break
		}
	}
	for idx, to := range t.timeouts {
		if to.id == id {
			t.timeouts[idx] = nil
			t.timeouts = append(t.timeouts[:idx], t.timeouts[idx+1:]...)
			break
		}
	}
}

func (t *TimeHolder) SetInterval(cb func(), ms int64) uint64 {
	_id_counter++
	id := _id_counter
	interval := &Interval{nil, id, 0}
	interval.fn = func() {
		cb()
		interval.timeout_id = t.SetTimeout(interval.fn, ms)
	}

	interval.timeout_id = t.SetTimeout(interval.fn, ms)

	t.intervals = append(t.intervals, interval)

	return id
}

func (t *TimeHolder) ClearInterval(id uint64) {
	for idx, interval := range t.intervals {
		if interval.id == id {
			t.intervals = append(t.intervals[:idx], t.intervals[idx+1:]...)
			t.ClearTimeout(interval.timeout_id)
			break
		}
	}
}

func (t *TimeHolder) Update() {
	now := TimeNow()
	idx := 0
	t.timeouts = append(t.timeouts, t.new_timeouts...)
	t.new_timeouts = nil

	for _, to := range t.timeouts {
		if to == nil {
			continue
		}
		if now >= to.emit_at {
			to.cb()
		} else {
			t.timeouts[idx] = to
			idx++
		}
	}

	t.timeouts = t.timeouts[:idx]
}

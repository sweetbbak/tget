package main

import (
	"time"

	"golang.org/x/time/rate"
)

func setRate(maxFriendsPerSec int) *rate.Limiter {
	return rate.NewLimiter(per(maxFriendsPerSec, time.Second), maxFriendsPerSec)
}

// for setting rate limit - taken from https://github.com/fanpei91/torsniff MIT licensed
func per(events int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(events))
}

// https://gitlab.com/axet/libtorrent/-/blob/master/libtorrent.go
func limit(kbps int) *rate.Limiter {
	var l = rate.NewLimiter(rate.Inf, 0)

	if kbps > 0 {
		b := kbps
		if b < 16*1024 {
			b = 16 * 1024
		}
		l = rate.NewLimiter(rate.Limit(kbps), b)
	}

	return l
}

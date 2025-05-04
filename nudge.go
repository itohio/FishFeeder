package main

import (
	"time"
)

type Nudge struct {
	timestamp time.Time
	delay     time.Duration
	nudge     func(*Nudge) bool
}

func (n *Nudge) Elapsed() time.Duration {
	return time.Since(n.timestamp)
}
func (n *Nudge) ETA() time.Duration {
	return n.delay - n.Elapsed()
}
func (n *Nudge) ETAPercent() float64 {
	return n.Elapsed().Seconds() / n.delay.Seconds()
}

func (n *Nudge) Check() {
	if n.ETA() < 0 {
		if n.nudge(n) {
			n.Reset()
		}
	}
}

func (n *Nudge) Reset() {
	n.timestamp = time.Now()
}

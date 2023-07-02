package lives

import (
	"runtime"
	"time"
)

type BoardLoop struct {
	onTick   func(float64)
	tickRate time.Duration
	quit     chan bool
}

func NewBoardLoop(onTick func(float64), tickRate time.Duration) *BoardLoop {
	return &BoardLoop{
		onTick:   onTick,
		tickRate: tickRate,
		quit:     make(chan bool),
	}
}

func (loop *BoardLoop) Start() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	startedAt := time.Now().UnixNano()
	ticker := time.NewTicker(time.Second / loop.tickRate)

	for {
		select {
		case <-ticker.C:
			now := time.Now().UnixNano()
			delta := float64(now-startedAt) / float64(time.Second)
			startedAt = now
			loop.onTick(delta)
		case <-loop.quit:
			ticker.Stop()
		}
	}
}

func (loop *BoardLoop) Stop() {
	loop.quit <- true
}

func (loop *BoardLoop) Restart() {
	loop.Stop()
	loop.Start()
}

func (loop *BoardLoop) TickRate() time.Duration {
	return loop.tickRate
}

func (loop *BoardLoop) SetTickRate(tickRate time.Duration) {
	loop.tickRate = tickRate
	loop.Restart()
}

func (loop *BoardLoop) SetOnTick(onTick func(float64)) {
	loop.onTick = onTick
}

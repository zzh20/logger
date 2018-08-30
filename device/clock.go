package device

import (
	"fmt"
	"sync/atomic"
	"time"
)

var clock = &Clock{}

type Clock struct {
	lastDateTime string
	lastTime     uint32
	lastDate     uint32
}

func TickOfTheClock() {
	t := time.Now()
	dt := uint32(t.Year()%100*10000 + int(t.Month())*100 + t.Day())
	tm := uint32(t.Hour()*10000 + t.Minute()*100 + t.Second())

	atomic.StoreUint32(&clock.lastDate, dt)
	atomic.StoreUint32(&clock.lastTime, tm)

	clock.lastDateTime = fmt.Sprintf(" %04d %06d", dt%10000, tm)
}

func GetLastDate() uint32 {
	return atomic.LoadUint32(&clock.lastDate)
}

func GetLastDateTime() string {
	return clock.lastDateTime
}

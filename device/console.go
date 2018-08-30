package device

import (
	"os"
	"sync"
)

// Console 控制台设备
type Console struct {
	sync.Mutex
}

func newConsoleDevice() Device {
	return &Console{}
}

func (console *Console) Write(b []byte) {
	console.Lock()
	defer console.Unlock()
	os.Stdout.Write(b)
}

func (console *Console) Flush() {
}

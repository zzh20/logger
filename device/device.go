package device

const (
	ConsoleDev = "console"
	File       = "file"
)

// Device 日志输出设备
type Device interface {
	Write(b []byte)
	Flush()
}

// NewDevice 创建一个新的日志输出设备
func NewDevice(deviceName, args string) Device {
	switch deviceName {
	case ConsoleDev:
		return newConsoleDevice()
	case File:
		return newFileDevice(args)
	}

	return nil
}

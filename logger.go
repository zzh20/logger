package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/zzh20/logger/device"
)

// Logger 日志对象
type Logger struct {
	minLevel uint8
	format   Formatter
	writers  []Writer

	ctx        context.Context
	cancelFunc context.CancelFunc
}

// Writer 日志输出对象
type Writer struct {
	level  uint8
	device device.Device
}

var worker = &Logger{
	minLevel: 99,
	format:   &DefaultFormatter{},
}

// Init 日志库初始化
func Init(filename string) {
	config, err := loadConfig(filename)
	if err != nil || len(config) == 0 {
		config = append(config, ConfigItem{LogLevelDebug, device.ConsoleDev, ""})
	}

	for _, writer := range config {
		worker.writers = append(worker.writers, NewWriter(writer.Level, writer.Device, writer.Args))
	}
	worker.UpdateLevel()

	worker.ctx, worker.cancelFunc = context.WithCancel(context.Background())

	go work()
}

func work() {
	device.TickOfTheClock()
	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-worker.ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			device.TickOfTheClock()
			worker.Flush()
		}
	}
}

func Stop() {
	worker.cancelFunc()
}

// NewWriter 创建新的日志输出对象
func NewWriter(level uint8, deviceName, args string) Writer {
	return Writer{
		level:  level,
		device: device.NewDevice(deviceName, args),
	}
}

// UpdateLevel 更新日志对象的最小输出级别
func (log *Logger) UpdateLevel() {
	for _, writer := range log.writers {
		if writer.level < log.minLevel {
			log.minLevel = writer.level
		}
	}
}

// Flush 刷新日志
func (log *Logger) Flush() {
	for _, writer := range log.writers {
		writer.device.Flush()
	}
}

// Write 输出日志
func (log *Logger) Write(level uint8, format string, a ...interface{}) {
	if level < log.minLevel {
		return
	}
	var msg string
	if len(a) == 0 {
		msg = format
	} else {
		msg = fmt.Sprintf(format, a...)
	}
	buff := log.format.Format(level, msg)
	b := buff.Bytes()
	for _, writer := range log.writers {
		if level >= writer.level {
			writer.device.Write(b)
		}
	}
	buffs.put(buff)
}

// Debug 输出DEBUG级别日志
func Debug(format string, a ...interface{}) {
	worker.Write(LogLevelDebug, format, a...)
}

// Info 输出INFO级别日志
func Info(format string, a ...interface{}) {
	worker.Write(LogLevelInfo, format, a...)
}

// Warn 输出WARN级别日志
func Warn(format string, a ...interface{}) {
	worker.Write(LogLevelWarn, format, a...)
}

// Error 输出ERROR级别日志
func Error(format string, a ...interface{}) {
	worker.Write(LogLevelError, format, a...)
}

// Fatal 输出FATAL级别日志
func Fatal(format string, a ...interface{}) {
	worker.Write(LogLevelFatal, format, a...)
	worker.Flush()
	os.Exit(1)
}

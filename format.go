package logger

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"

	"time"
)

// Formatter 日志格式化接口
type Formatter interface {
	Format(level uint8, msg string) *bytes.Buffer
}

type buffPool struct {
	pool sync.Pool
}

var buffs = &buffPool{
	pool: sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0))
		},
	},
}

func (b *buffPool) get() *bytes.Buffer {
	return b.pool.Get().(*bytes.Buffer)
}

func (b *buffPool) put(buf *bytes.Buffer) {
	buf.Reset()
	b.pool.Put(buf)
}

// DefaultFormatter 默认格式化
type DefaultFormatter struct {
	format string
}

func getLevelStr(level uint8) string {
	switch level {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	default:
		fmt.Printf("ERROR: logger level unknown: %v\n", level)
		return "INFO"
	}
}

// Format 格式化
func (format *DefaultFormatter) Format(level uint8, msg string) *bytes.Buffer {
	buff := buffs.get()
	fmt.Fprintf(buff, "%s ", time.Now().Format("2006-01-02 15:04:05"))
	_, file, line, ok := runtime.Caller(4)
	if ok {
		var i = len(file) - 2
		for ; i >= 0; i-- {
			if file[i] == '/' {
				i++
				break
			}
		}
		buff.WriteString(file[i:])
		buff.WriteByte(':')
		buff.WriteString(strconv.FormatInt(int64(line), 10))
	}
	fmt.Fprintf(buff, " %s - ", getLevelStr(level))
	buff.WriteString(msg)
	buff.WriteByte('\n')
	return buff
}

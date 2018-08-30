package device

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

// File 文件输出设备
type FileDevice struct {
	file     *os.File
	writer   *bufio.Writer
	prefix   string
	lastDate uint32
	sync.Mutex
}

func newFileDevice(args string) Device {
	return &FileDevice{
		prefix: args,
	}
}

func (fileDevice *FileDevice) Write(b []byte) {
	var err error

	date := GetLastDate()
	fileDevice.Lock()
	defer fileDevice.Unlock()

	if fileDevice.lastDate != date {
		if fileDevice.file != nil {
			fileDevice.writer.Flush()
			err = fileDevice.file.Close()
			if err != nil {
				fmt.Printf("ERROR: logger cannot close file: %v\n", err.Error())
			}
			fileDevice.file = nil
		}
	}
	if fileDevice.file == nil {
		filename := fmt.Sprintf("%s-%v.log", fileDevice.prefix, date)
		filePath := fileDevice.prefix[0:strings.LastIndex(fileDevice.prefix, "/")]
		err := os.MkdirAll(filePath, 0711)
		fileDevice.file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Printf("ERROR: logger cannot open file: %v\n", err.Error())
			return
		}
		fileDevice.writer = bufio.NewWriter(fileDevice.file)
		fileDevice.lastDate = date
	}
	_, err = fileDevice.writer.Write(b)

	if err != nil {
		fmt.Printf("ERROR: logger cannot write file: %v\n", err.Error())
	}
	return
}

// Flush 刷新到设备
func (fileDevice *FileDevice) Flush() {
	fileDevice.Lock()
	defer fileDevice.Unlock()

	if fileDevice.writer != nil {
		fileDevice.writer.Flush()
	}
}

package logger

import (
	"encoding/json"
	"io/ioutil"
)

// ConfigItem 日志配置项
type ConfigItem struct {
	Level  uint8
	Device string
	Args   string
}

type Config struct {
	Writers []ConfigItem
}

func loadConfig(filename string) (config []ConfigItem, err error) {
	var data []byte
	data, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return
	}

	return cfg.Writers, nil
}

package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	cfgPath = flag.String("c", "config.yaml", "path to config file")
	dump    = flag.Bool("d", false, "dump config to file")
)

type FileType int

const (
	FileTypeYAML = iota
	FileTypeJSON
)

func LoadConfig(filename string, cfg interface{}, typ FileType) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("ReadFile: %v", err)
	}
	switch typ {
	case FileTypeYAML:
		return yaml.Unmarshal(data, cfg)
	case FileTypeJSON:
		return json.Unmarshal(data, cfg)
	}
	return errors.New("unsupported file type")
}

func DumpConfig(filename string, cfg interface{}, typ FileType) error {
	var (
		data []byte
		err  error
	)
	switch typ {
	case FileTypeYAML:
		data, err = yaml.Marshal(cfg)
	case FileTypeJSON:
		data, err = json.Marshal(cfg)
	}
	if err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	_ = f.Close()
	return err
}

func init() {
	// 默认从环境变量加载配置
	_ = flag.Set("c", os.Getenv("CONFIG_PATH"))
	_ = flag.Set("d", os.Getenv("DUMP_CONFIG"))
}

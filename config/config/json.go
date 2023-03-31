package config

import (
	"flag"
	"os"
)

func LoadJSONConfig(filename string, cfg interface{}) error {
	return LoadConfig(filename, cfg, FileTypeJSON)
}

func LoadJSONConfigFromFlag(cfg interface{}) error {
	if !flag.Parsed() {
		flag.Parse()
	}
	return LoadJSONConfig(*cfgPath, cfg)
}

func DumpJSONConfig(filename string, cfg interface{}) error {
	return DumpConfig(filename, cfg, FileTypeJSON)
}

func DumpJSONConfigFromFlag(cfg interface{}) error {
	if !flag.Parsed() {
		flag.Parse()
	}
	return DumpJSONConfig(*cfgPath, cfg)
}

func LoadOrDumpJSONConfigFromFlag(cfg interface{}) error {
	if *dump {
		if err := DumpJSONConfigFromFlag(cfg); err != nil {
			return err
		}
		os.Exit(0)
	}
	return LoadJSONConfigFromFlag(cfg)
}

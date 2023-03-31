package config

import (
	"flag"
	"os"
)

func LoadYAMLConfig(filename string, cfg interface{}) error {
	return LoadConfig(filename, cfg, FileTypeYAML)
}

func LoadYAMLConfigFromFlag(cfg interface{}) error {
	if !flag.Parsed() {
		flag.Parse()
	}
	return LoadYAMLConfig(*cfgPath, cfg)
}

func DumpYAMLConfig(filename string, cfg interface{}) error {
	return DumpConfig(filename, cfg, FileTypeYAML)
}

func DumpYAMLConfigFromFlag(cfg interface{}) error {
	if !flag.Parsed() {
		flag.Parse()
	}
	return DumpYAMLConfig(*cfgPath, cfg)
}

func LoadOrDumpYAMLConfigFromFlag(cfg interface{}) error {
	if *dump {
		if err := DumpYAMLConfigFromFlag(cfg); err != nil {
			return err
		}
		os.Exit(0)
	}
	return LoadYAMLConfigFromFlag(cfg)
}

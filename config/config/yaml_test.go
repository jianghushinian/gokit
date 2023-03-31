package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Config struct {
	Username string
	Password string
	Server   struct {
		Endpoint string
	}
}

var expCfg = Config{
	Username: "user",
	Password: "pass",
	Server: struct {
		Endpoint string
	}{
		Endpoint: "https://jianghushinian.cn/",
	},
}

func TestLoadYAMLConfig(t *testing.T) {
	var cfg Config
	err := LoadYAMLConfig("testdata/config.yaml", &cfg)
	assert.NoError(t, err)
	assert.Equal(t, expCfg, cfg)
}

func TestLoadYAMLConfigFromFlag(t *testing.T) {
	_ = flag.Set("c", "testdata/config.yaml")

	var cfg Config
	err := LoadYAMLConfigFromFlag(&cfg)
	assert.NoError(t, err)
	assert.Equal(t, expCfg, cfg)
}

func TestDumpYAMLConfig(t *testing.T) {
	f, _ := os.CreateTemp("", "TEST_DUMP")
	defer os.Remove(f.Name())

	err := DumpYAMLConfig(f.Name(), &expCfg)
	assert.NoError(t, err)

	var cfg Config
	err = LoadYAMLConfig(f.Name(), &cfg)
	assert.NoError(t, err)
	assert.Equal(t, expCfg, cfg)
}

func TestDumpYAMLConfigFromFlag(t *testing.T) {
	f, _ := os.CreateTemp("", "TEST_DUMP")
	defer os.Remove(f.Name())

	_ = flag.Set("c", f.Name())

	err := DumpYAMLConfigFromFlag(&expCfg)
	assert.NoError(t, err)

	var cfg Config
	err = LoadYAMLConfig(f.Name(), &cfg)
	assert.NoError(t, err)
	assert.Equal(t, expCfg, cfg)
}

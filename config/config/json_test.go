package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadJSONConfig(t *testing.T) {
	var cfg Config
	err := LoadJSONConfig("testdata/config.json", &cfg)
	assert.NoError(t, err)
	assert.Equal(t, expCfg, cfg)
}

func TestLoadJSONConfigFromFlag(t *testing.T) {
	_ = flag.Set("c", "testdata/config.json")

	var cfg Config
	err := LoadJSONConfigFromFlag(&cfg)
	assert.NoError(t, err)
	assert.Equal(t, expCfg, cfg)
}

func TestDumpJSONConfig(t *testing.T) {
	f, _ := os.CreateTemp("", "TEST_DUMP")
	defer os.Remove(f.Name())

	err := DumpJSONConfig(f.Name(), &expCfg)
	assert.NoError(t, err)

	var cfg Config
	err = LoadJSONConfig(f.Name(), &cfg)
	assert.NoError(t, err)
	assert.Equal(t, expCfg, cfg)
}

func TestDumpJSONConfigFromFlag(t *testing.T) {
	f, _ := os.CreateTemp("", "TEST_DUMP")
	defer os.Remove(f.Name())

	_ = flag.Set("c", f.Name())

	err := DumpJSONConfigFromFlag(&expCfg)
	assert.NoError(t, err)

	var cfg Config
	err = LoadJSONConfig(f.Name(), &cfg)
	assert.NoError(t, err)
	assert.Equal(t, expCfg, cfg)
}

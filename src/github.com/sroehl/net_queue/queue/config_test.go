package queue

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_read_yaml_config(t *testing.T) {
	assert := assert.New(t)

	f, err := os.CreateTemp("", "config")
	if err != nil {
		assert.Fail("Could not create temp file: %v", err)
	}
	yaml_string := "server:\n  listenhost: \"localhost\"\n  port: 4545\n"
	_, err = f.Write([]byte(yaml_string))
	if err != nil {
		assert.Fail("Could not write file: %v", err)
	}
	f.Close()

	cfg, err := Read_config(f.Name())
	if err != nil {
		assert.Fail("Failed to read file: %v", err)
	}
	assert.Equal(4545, cfg.Server.Port)
	assert.Equal("localhost", cfg.Server.ListenHost)

	os.Remove(f.Name())
}

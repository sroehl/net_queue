package queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_new_entry(t *testing.T) {
	msg := "This is a test message"
	e := new_entry(msg)
	assert.Equal(t, e.msg, msg)
	assert.Equal(t, e.read, false)
	assert.Greater(t, e.time, time.Now().Unix()-1)
}

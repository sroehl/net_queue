package net_queue

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Create_Queue_API(t *testing.T) {
	assert := assert.New(t)
	s := new_server("localhost", 4545)
	s.connect()
	resp, err := s.create_queue("testQueue")
	if err != nil {
		assert.Fail(fmt.Sprintf("Create queue should have response, got error: %v", err))
		return
	}
	assert.Equal(SUCCESS, resp.Status)
	assert.Equal("", resp.Msg)
}

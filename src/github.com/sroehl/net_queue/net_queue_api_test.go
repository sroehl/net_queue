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

func Test_Create_Queue_Already_Error_API(t *testing.T) {
	assert := assert.New(t)
	s := new_server("localhost", 4545)
	s.connect()
	resp, err := s.create_queue("testAlreadyCreated")
	if err != nil {
		assert.Fail(fmt.Sprintf("Create queue should have response, got error: %v", err))
		return
	}
	assert.Equal(SUCCESS, resp.Status)
	assert.Equal("", resp.Msg)

	resp2, err2 := s.create_queue("testAlreadyCreated")
	if err2 != nil {
		assert.Fail(fmt.Sprintf("Create queue should have response, got error: %v", err))
		return
	}
	assert.Equal(ERROR, resp2.Status)
	assert.Equal("Queue already exists", resp2.Msg)
}

func Test_Delete_Queue_API(t *testing.T) {
	assert := assert.New(t)
	s := new_server("localhost", 4545)
	s.connect()
	resp, err := s.create_queue("deleteQueueTest")
	if err != nil {
		assert.Fail(fmt.Sprintf("Create queue should have response, got error: %v", err))
		return
	}
	assert.Equal(SUCCESS, resp.Status)
	assert.Equal("", resp.Msg)

	resp2, err2 := s.delete_queue("deleteQueueTest")
	if err2 != nil {
		assert.Fail(fmt.Sprintf("Delete queue should have response, got error: %v", err))
		return
	}
	assert.Equal(SUCCESS, resp2.Status)
}

func Test_Delete_Nonexistant_Queue_API(t *testing.T) {
	assert := assert.New(t)
	s := new_server("localhost", 4545)
	s.connect()

	resp, err := s.delete_queue("deleteNonexistantQueueTest")
	if err != nil {
		assert.Fail(fmt.Sprintf("Delete queue should have response, got error: %v", err))
		return
	}
	assert.Equal(ERROR, resp.Status)
}

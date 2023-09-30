package api

import (
	"fmt"
	"net_queue/queue"
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
	assert.Equal(queue.SUCCESS, resp.Status)
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
	assert.Equal(queue.SUCCESS, resp.Status)
	assert.Equal("", resp.Msg)

	resp2, err2 := s.create_queue("testAlreadyCreated")
	if err2 != nil {
		assert.Fail(fmt.Sprintf("Create queue should have response, got error: %v", err))
		return
	}
	assert.Equal(queue.ERROR, resp2.Status)
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
	assert.Equal(queue.SUCCESS, resp.Status)
	assert.Equal("", resp.Msg)

	resp2, err2 := s.delete_queue("deleteQueueTest")
	if err2 != nil {
		assert.Fail(fmt.Sprintf("Delete queue should have response, got error: %v", err))
		return
	}
	assert.Equal(queue.SUCCESS, resp2.Status)
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
	assert.Equal(queue.ERROR, resp.Status)
}

func Test_List_Queue(t *testing.T) {
	assert := assert.New(t)
	s := new_server("localhost", 4545)
	s.connect()
	create_resp1, create_err1 := s.create_queue("listQueue1")
	if create_err1 != nil {
		assert.Fail(fmt.Sprintf("Failed to create queue, got error: %v", create_err1))
	}
	assert.Equal(queue.SUCCESS, create_resp1.Status)
	create_resp2, create_err2 := s.create_queue("listQueue2")
	if create_err2 != nil {
		assert.Fail(fmt.Sprintf("Failed to create queue, got error: %v", create_err1))
	}
	assert.Equal(queue.SUCCESS, create_resp2.Status)

	resp, err := s.list_queues()
	if err != nil {
		assert.Fail(fmt.Sprintf("Failed to list queues, got error: %v", create_err1))
	}
	assert.Equal(queue.SUCCESS, resp.Status)
	assert.Contains(resp.Msg, "listQueue1")
	assert.Contains(resp.Msg, "listQueue2")
}

func Test_Purge_Queue(t *testing.T) {
	assert := assert.New(t)
	s := new_server("localhost", 4545)
	s.connect()
	create_resp, err := s.create_queue("purgeQueue")
	if err != nil {
		assert.Fail(fmt.Sprintf("Failed to create queue, got error: %v", err))
	}
	assert.Equal(queue.SUCCESS, create_resp.Status)

	list_resp, err := s.list_queues()
	if err != nil {
		assert.Fail(fmt.Sprintf("Failed to list queues, got error: %v", err))
	}
	assert.Contains(list_resp.Msg, "purgeQueue")

	write_resp, err := s.write_msg("purgeQueue", "This is a test")
	if err != nil {
		assert.Fail(fmt.Sprintf("Failed to write to queue, got error: %v", err))
	}
	assert.Equal(queue.SUCCESS, write_resp.Status)
	assert.Equal("", write_resp.Msg)

	purge_resp, err := s.purge_queue("purgeQueue")
	if err != nil {
		assert.Fail(fmt.Sprintf("Failed to purge queue, got error: %v", err))
	}
	assert.Equal(queue.SUCCESS, purge_resp.Status)

	read_resp, err := s.read_msg("purgeQueue", 0)
	if err != nil {
		assert.Fail(fmt.Sprintf("Failed to read from queue, got error: %v", err))
	}
	assert.Equal(queue.NO_MSG, read_resp.Status)
	assert.Equal("", read_resp.Msg)
}

func Test_write_read_api(t *testing.T) {
	assert := assert.New(t)
	s := new_server("localhost", 4545)
	s.connect()
	create_resp, err := s.create_queue("writeQueue")
	if err != nil {
		assert.Fail(fmt.Sprintf("Failed to create queue, got error: %v", err))
	}
	assert.Equal(queue.SUCCESS, create_resp.Status)

	list_resp, err := s.list_queues()
	if err != nil {
		assert.Fail(fmt.Sprintf("Failed to list queues, got error: %v", err))
	}
	assert.Contains(list_resp.Msg, "writeQueue")

	write_resp, err := s.write_msg("writeQueue", "This is a test")
	if err != nil {
		assert.Fail(fmt.Sprintf("Failed to write to queue, got error: %v", err))
	}
	assert.Equal(queue.SUCCESS, write_resp.Status)
	assert.Equal("", write_resp.Msg)

	read_resp, err := s.read_msg("writeQueue", 0)
	if err != nil {
		assert.Fail(fmt.Sprintf("Failed to read from queue, got error: %v", err))
	}
	assert.Equal(queue.SUCCESS, read_resp.Status)
	assert.Equal("This is a test", read_resp.Msg)
}

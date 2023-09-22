package net_queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_write_entry(t *testing.T) {
	assert := assert.New(t)
	queue_name := "testQueue"
	queues := make(map[string]*Queue)
	queues[queue_name] = new_queue(queue_name)

	net_entry := NetMessageEntry{
		queue: queue_name,
		msg:   "Test message",
	}

	resp1 := net_entry.write_entry(queues)
	assert.Equal(SUCCESS, resp1.Status)
	assert.Equal(1, queues[queue_name].size)

	net_entry2 := NetMessageEntry{
		queue: "badQueue",
		msg:   "Test message",
	}
	resp2 := net_entry2.write_entry(queues)
	assert.Equal(ERROR, resp2.Status)
	assert.Equal("Queue does not exist", resp2.Msg)
}

func Test_read_entry(t *testing.T) {
	assert := assert.New(t)
	queue_name := "testQueue"
	queues := make(map[string]*Queue)
	queues[queue_name] = new_queue(queue_name)
	queues[queue_name].add_msg("Test message")
	assert.Equal(1, queues[queue_name].size)

	result, err := queues[queue_name].read(false, false, 0)
	if err != nil {
		assert.FailNowf("Should not have error: %v", err.Error())
	}
	assert.Equal(1, queues[queue_name].size)
	assert.Equal("Test message", result.entry.msg)

	net_entry := NetMessageEntry{
		queue: queue_name,
	}

	resp := net_entry.read_entry(queues)
	assert.Equal(SUCCESS, resp.Status)
	assert.Equal("Test message", resp.Msg)
}

func Test_read_many(t *testing.T) {
	assert := assert.New(t)
	queue_name := "testQueue"
	queues := make(map[string]*Queue)
	queues[queue_name] = new_queue(queue_name)
	queues[queue_name].add_msg("Test message1")
	queues[queue_name].add_msg("Test message2")
	queues[queue_name].add_msg("Test message3")
	queues[queue_name].add_msg("Test message4")
	assert.Equal(4, queues[queue_name].size)

	net_entry := NetMessageEntry{
		queue: queue_name,
	}
	resp := net_entry.read_entry(queues)
	found := 1
	has_more := resp.Status == HAS_MORE
	idx := resp.Index
	for has_more {
		net_entry2 := NetMessageEntry{
			queue: queue_name,
			index: idx + 1,
		}
		resp2 := net_entry2.read_entry(queues)
		if idx+1 < 3 {
			assert.Equal(HAS_MORE, resp2.Status)
			assert.Equal(idx+1, resp2.Index)
		}
		if idx+1 == 4 {
			assert.Equal(SUCCESS, resp2.Status)
			assert.Equal(4, resp2.Index)
		}
		if resp2.Status == SUCCESS || resp2.Status == HAS_MORE {
			found++
		}
		has_more = resp2.Status == HAS_MORE
		idx = resp2.Index
	}
	assert.Equal(4, found)
}

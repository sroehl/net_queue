package queue

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
		Queue: queue_name,
		Msg:   "Test message",
	}

	resp1 := net_entry.Write_entry(queues)
	assert.Equal(SUCCESS, resp1.Status)
	assert.Equal(1, queues[queue_name].size)

	net_entry2 := NetMessageEntry{
		Queue: "badQueue",
		Msg:   "Test message",
	}
	resp2 := net_entry2.Write_entry(queues)
	assert.Equal(ERROR, resp2.Status)
	assert.Equal("Queue does not exist", resp2.Msg)
}

func Test_Read_entry(t *testing.T) {
	assert := assert.New(t)
	queue_name := "testQueue"
	queues := make(map[string]*Queue)
	queues[queue_name] = new_queue(queue_name)
	queues[queue_name].add_msg("Test message")
	assert.Equal(1, queues[queue_name].size)

	ch := make(chan queueReadResult)
	done_ch := make(chan bool, 1)
	opts := ReadOptions{
		Index:      0,
		Unread:     false,
		Delete:     false,
		Continuous: false,
	}

	go queues[queue_name].read(opts, ch, done_ch)
	result := <-ch
	if result.err != nil {
		assert.FailNowf("Should not have error: %v", result.err.Error())
	}
	assert.Equal(1, queues[queue_name].size)
	assert.Equal("Test message", result.entryResult.entry.msg)
	done_ch <- true

	net_entry := NetMessageEntry{
		Queue: queue_name,
		Opt:   opts,
	}

	net_response_ch := make(chan NetResponse)
	done2 := make(chan bool)
	go net_entry.Read_entry(queues, net_response_ch, done2)
	resp := <-net_response_ch
	assert.Equal(SUCCESS, resp.Status)
	assert.Equal("Test message", resp.Msg)
	done2 <- true
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
		Queue: queue_name,
	}
	net_response_ch := make(chan NetResponse)
	done := make(chan bool)
	go net_entry.Read_entry(queues, net_response_ch, done)
	resp := <-net_response_ch
	found := 1
	has_more := resp.Status == HAS_MORE
	idx := resp.Index
	for has_more {
		net_entry2 := NetMessageEntry{
			Queue: queue_name,
			Opt: ReadOptions{
				Index: idx + 1,
			},
		}
		net_response_ch = make(chan NetResponse)
		done = make(chan bool)
		go net_entry2.Read_entry(queues, net_response_ch, done)
		resp2 := <-net_response_ch
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

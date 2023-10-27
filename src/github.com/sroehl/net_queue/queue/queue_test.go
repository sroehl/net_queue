package queue

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func create_queue(num_msgs int) *Queue {
	queue := new_queue("TestQueue1")
	for i := 1; i < num_msgs+1; i++ {
		msg := fmt.Sprintf("Test msg num:%v", i)
		queue.add_msg(msg)
	}
	return queue
}

func Test_queue_creation(t *testing.T) {
	assert := assert.New(t)
	queue_name := "TestQueue1"
	queue := new_queue(queue_name)
	assert.Equal(queue.name, queue_name)
	assert.Equal(queue.size, 0)
}

func Test_queue_add(t *testing.T) {
	assert := assert.New(t)
	q := create_queue(3)
	assert.Equal(3, q.size)
	assert.Equal("Test msg num:1", q.entries[0].msg)
	assert.Equal(false, q.entries[0].read)
	assert.Equal("Test msg num:2", q.entries[1].msg)
	assert.Equal(false, q.entries[1].read)
}

func Test_queue_read_unread_no_remove(t *testing.T) {
	assert := assert.New(t)
	q := create_queue(5)
	opts := ReadOptions{
		Index:      0,
		Delete:     false,
		Unread:     false,
		Continuous: false,
	}
	ch := make(chan queueReadResult)
	done := make(chan bool)
	go q.read(opts, ch, done)
	result := <-ch
	if result.err != nil {
		assert.Fail("Should have result")
	}
	assert.Equal("Test msg num:1", result.entryResult.entry.msg)
	assert.Equal(true, result.entryResult.has_more)
	assert.Equal(5, q.size)
	ch = make(chan queueReadResult)
	done = make(chan bool)
	go q.read(opts, ch, done)
	result2 := <-ch
	if result2.err != nil {
		assert.Fail("Should have result")
	}

	assert.Equal("Test msg num:1", result2.entryResult.entry.msg)
	assert.Equal(true, result2.entryResult.has_more)
	assert.Equal(5, q.size)
}

func Test_queue_read_unread_remove(t *testing.T) {
	assert := assert.New(t)
	q := create_queue(5)
	opts := ReadOptions{
		Index:      0,
		Delete:     true,
		Unread:     false,
		Continuous: false,
	}
	ch := make(chan queueReadResult)
	done := make(chan bool)
	go q.read(opts, ch, done)
	result := <-ch
	if result.err != nil {
		assert.Fail("Should have result")
	}
	assert.Equal("Test msg num:1", result.entryResult.entry.msg)
	assert.Equal(true, result.entryResult.has_more)
	assert.Equal(4, q.size)

	ch = make(chan queueReadResult)
	done = make(chan bool)
	go q.read(opts, ch, done)
	result2 := <-ch
	if result2.err != nil {
		assert.Fail("Should have result")
	}
	assert.Equal("Test msg num:2", result2.entryResult.entry.msg)
	assert.Equal(true, result2.entryResult.has_more)
	assert.Equal(3, q.size)
}

func Test_queue_read_already_read_no_remove(t *testing.T) {
	assert := assert.New(t)
	q := create_queue(5)
	opts := ReadOptions{
		Index:      0,
		Delete:     false,
		Unread:     true,
		Continuous: false,
	}
	ch := make(chan queueReadResult)
	done := make(chan bool)
	go q.read(opts, ch, done)
	result := <-ch
	if result.err != nil {
		assert.Fail("Should have result")
	}
	assert.Equal("Test msg num:1", result.entryResult.entry.msg)
	assert.Equal(true, result.entryResult.has_more)
	assert.Equal(5, q.size)

	ch = make(chan queueReadResult)
	done = make(chan bool)
	go q.read(opts, ch, done)
	result2 := <-ch
	if result2.err != nil {
		assert.Fail("Should have result")
	}
	assert.Equal("Test msg num:2", result2.entryResult.entry.msg)
	assert.Equal(true, result2.entryResult.has_more)
	assert.Equal(5, q.size)
}

func Test_queue_read_already_read_remove(t *testing.T) {
	assert := assert.New(t)
	q := create_queue(5)
	opts := ReadOptions{
		Index:      0,
		Delete:     true,
		Unread:     true,
		Continuous: false,
	}
	ch := make(chan queueReadResult)
	done := make(chan bool)
	go q.read(opts, ch, done)
	result := <-ch
	if result.err != nil {
		assert.Fail("Should have result")
	}
	assert.Equal("Test msg num:1", result.entryResult.entry.msg)
	assert.Equal(true, result.entryResult.has_more)
	assert.Equal(4, q.size)

	ch = make(chan queueReadResult)
	done = make(chan bool)
	go q.read(opts, ch, done)
	result2 := <-ch
	if result2.err != nil {
		assert.Fail("Should have result")
	}
	assert.Equal("Test msg num:2", result2.entryResult.entry.msg)
	assert.Equal(true, result2.entryResult.has_more)
	assert.Equal(3, q.size)
}

func Test_queue_read_0_entries(t *testing.T) {
	assert := assert.New(t)
	q := create_queue(0)
	opts := ReadOptions{
		Index:      0,
		Delete:     true,
		Unread:     true,
		Continuous: false,
	}
	ch := make(chan queueReadResult)
	done := make(chan bool)
	go q.read(opts, ch, done)
	result := <-ch
	if result.err == nil {
		assert.Fail("Should have error")
	}
	assert.Equal("Queue is empty or index is too large", result.err.Error())
	assert.Equal(Entry{}, result.entryResult.entry)
	assert.Equal(EntryResult{}, result.entryResult)
}

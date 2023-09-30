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
	result, err := q.read(false, false, 0)
	if err != nil {
		assert.Fail("Should have result")
	}
	assert.Equal("Test msg num:1", result.entry.msg)
	assert.Equal(true, result.has_more)
	assert.Equal(5, q.size)
	result2, err2 := q.read(false, false, 0)
	if err2 != nil {
		assert.Fail("Should have result")
	}
	assert.Equal("Test msg num:1", result2.entry.msg)
	assert.Equal(true, result2.has_more)
	assert.Equal(5, q.size)
}

func Test_queue_read_unread_remove(t *testing.T) {
	assert := assert.New(t)
	q := create_queue(5)
	result, err := q.read(false, true, 0)
	if err != nil {
		assert.Fail("Should have result")
	}
	assert.Equal("Test msg num:1", result.entry.msg)
	assert.Equal(true, result.has_more)
	assert.Equal(4, q.size)

	result2, err2 := q.read(false, true, 0)
	if err2 != nil {
		assert.Fail("Should have result")
	}
	assert.Equal("Test msg num:2", result2.entry.msg)
	assert.Equal(true, result2.has_more)
	assert.Equal(3, q.size)
}

func Test_queue_read_already_read_no_remove(t *testing.T) {
	assert := assert.New(t)
	q := create_queue(5)
	result, err := q.read(true, false, 0)
	if err != nil {
		assert.Fail("Should have result")
	}
	assert.Equal("Test msg num:1", result.entry.msg)
	assert.Equal(true, result.has_more)
	assert.Equal(5, q.size)

	result2, err2 := q.read(true, false, 0)
	if err2 != nil {
		assert.Fail("Should have result")
	}
	assert.Equal("Test msg num:2", result2.entry.msg)
	assert.Equal(true, result2.has_more)
	assert.Equal(5, q.size)
}

func Test_queue_read_already_read_remove(t *testing.T) {
	assert := assert.New(t)
	q := create_queue(5)
	result, err := q.read(true, true, 0)
	if err != nil {
		assert.Fail("Should have result")
	}
	assert.Equal("Test msg num:1", result.entry.msg)
	assert.Equal(true, result.has_more)
	assert.Equal(4, q.size)

	result2, err2 := q.read(true, true, 0)
	if err2 != nil {
		assert.Fail("Should have result")
	}
	assert.Equal("Test msg num:2", result2.entry.msg)
	assert.Equal(true, result2.has_more)
	assert.Equal(3, q.size)
}

func Test_queue_read_0_entries(t *testing.T) {
	assert := assert.New(t)
	q := create_queue(0)
	result, err := q.read(false, false, 0)
	if err == nil {
		assert.Fail("Should have error")
	}
	assert.Equal("Queue is empty or index is too large", err.Error())
	assert.Equal(Entry{}, result.entry)
	assert.Equal(EntryResult{}, result)
}

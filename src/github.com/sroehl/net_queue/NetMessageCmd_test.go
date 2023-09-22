package net_queue

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	fmt.Println("Starting listener")
	go start_server(4545)
	m.Run()
	fmt.Println("Listener down")
}

func Test_get_queue(t *testing.T) {
	assert := assert.New(t)
	queues := make(map[string]*Queue)
	net_cmd := NetMessageCmd{
		Command:   CREATE_QUEUE,
		Arguments: "testQueue",
	}
	resp1 := net_cmd.handle_cmd(queues)
	assert.Equal(SUCCESS, resp1.Status)
	resp2 := net_cmd.handle_cmd(queues)
	assert.Equal(ERROR, resp2.Status)
	assert.Equal("Queue already exists", resp2.Msg)
	assert.Equal(1, len(queues))
}

func Test_delete_queue(t *testing.T) {
	assert := assert.New(t)
	queue_name := "testQueue"
	queues := make(map[string]*Queue)
	queues["testQueue0"] = new_queue("testQueue0")
	queues[queue_name] = new_queue(queue_name)
	queues["testQueue2"] = new_queue("testQueue2")
	net_cmd := NetMessageCmd{
		Command:   DELETE_QUEUE,
		Arguments: queue_name,
	}
	resp1 := net_cmd.handle_cmd(queues)
	assert.Equal(SUCCESS, resp1.Status)

	resp2 := net_cmd.handle_cmd(queues)
	assert.Equal(ERROR, resp2.Status)
	assert.Equal("Queue does not exist", resp2.Msg)
	assert.Equal(2, len(queues))
}

func Test_list_queues(t *testing.T) {
	assert := assert.New(t)
	queues := make(map[string]*Queue)
	queues["testQueue0"] = new_queue("testQueue0")
	queues["testQueue"] = new_queue("testQueue")
	queues["testQueue2"] = new_queue("testQueue2")
	net_cmd := NetMessageCmd{
		Command: LIST_QUEUES,
	}
	resp1 := net_cmd.handle_cmd(queues)
	assert.Equal(SUCCESS, resp1.Status)
	assert.Equal(3, len(strings.Split(resp1.Msg, ",")))

	queues["testQueue3"] = new_queue("testQueue3")
	resp2 := net_cmd.handle_cmd(queues)
	assert.Equal(SUCCESS, resp2.Status)
	assert.Equal(4, len(strings.Split(resp2.Msg, ",")))
}

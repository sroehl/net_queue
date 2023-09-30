package queue

import "fmt"

const (
	EXIT         = "exit"
	CREATE_QUEUE = "create_queue"
	DELETE_QUEUE = "delete_queue"
	PURGE_QUEUE  = "purge_queue"
	LIST_QUEUES  = "list_queues"
	NONE         = "none"
)

type NetMessageCmd struct {
	Command   string
	Arguments string
}

func (cmd NetMessageCmd) Handle_cmd(queues map[string]*Queue) NetResponse {
	fmt.Printf("Processing Command:%v\n", cmd.Command)
	switch cmd.Command {
	case EXIT:
		// TODO: close out server
	case CREATE_QUEUE:
		queue_name := cmd.Arguments
		_, ok := queues[queue_name]
		if ok {
			return new_netresponse(ERROR, "Queue already exists")
		} else {
			queues[queue_name] = new_queue(queue_name)
		}
	case DELETE_QUEUE:
		queue_name := cmd.Arguments
		_, ok := queues[queue_name]
		if ok {
			delete(queues, queue_name)
		} else {
			return new_netresponse(ERROR, "Queue does not exist")
		}
	case PURGE_QUEUE:
		queue_name := cmd.Arguments
		_, ok := queues[queue_name]
		if ok {
			queues[queue_name] = new_queue(queue_name)
		} else {
			return new_netresponse(ERROR, "Queue does not exist")
		}
	case LIST_QUEUES:
		queue_names := ""
		first := true
		for k := range queues {
			if !first {
				queue_names += ","
			} else {
				first = false
			}
			queue_names += k
		}
		return new_netresponse(SUCCESS, queue_names)
	default:
		return new_netresponse(ERROR, "Unknown Command")
	}
	return new_netresponse(SUCCESS, "")
}

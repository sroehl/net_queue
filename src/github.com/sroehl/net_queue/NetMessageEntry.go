package net_queue

type NetMessageEntry struct {
	queue string
	msg   string
	index int
}

func (entry NetMessageEntry) write_entry(queues map[string]*Queue) NetResponse {
	queue, ok := queues[entry.queue]
	if ok {
		queue.add_msg(entry.msg)
		return new_netresponse(SUCCESS, "")
	} else {
		return new_netresponse(ERROR, "Queue does not exist")
	}
}

func (entry NetMessageEntry) read_entry(queues map[string]*Queue) NetResponse {
	queue, ok := queues[entry.queue]
	if ok {
		result, err := queue.read(false, false, entry.index)
		if err != nil {
			return new_netresponse(NO_MSG, "")
		}
		status := SUCCESS
		if result.has_more {
			status = HAS_MORE
		}
		resp := new_netresponse(status, result.entry.msg)
		if result.index > -1 {
			resp.Index = result.index
		}
		return resp
	}
	return new_netresponse(ERROR, "Queue does not exist")
}

package queue

type NetMessageEntry struct {
	Queue string
	Msg   string
	Index int
}

func (entry NetMessageEntry) Write_entry(queues map[string]*Queue) NetResponse {
	queue, ok := queues[entry.Queue]
	if ok {
		queue.add_msg(entry.Msg)
		return new_netresponse(SUCCESS, "")
	} else {
		return new_netresponse(ERROR, "Queue does not exist")
	}
}

func (entry NetMessageEntry) Read_entry(queues map[string]*Queue) NetResponse {
	queue, ok := queues[entry.Queue]
	if ok {
		result, err := queue.read(false, false, entry.Index)
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

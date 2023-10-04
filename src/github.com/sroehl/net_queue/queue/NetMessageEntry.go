package queue

type NetMessageEntryOptions struct {
	Index  int
	Unread bool
	Delete bool
}

type NetMessageEntry struct {
	Queue string
	Msg   string
	Opt   NetMessageEntryOptions
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
		result, err := queue.read(entry.Opt.Unread, entry.Opt.Delete, entry.Opt.Index)
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

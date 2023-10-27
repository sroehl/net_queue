package queue

type ReadOptions struct {
	Index      int
	Unread     bool
	Delete     bool
	Continuous bool
}

type NetMessageEntry struct {
	Queue string
	Msg   string
	Opt   ReadOptions
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

func (entry NetMessageEntry) Read_entry(queues map[string]*Queue, ch chan NetResponse, done chan bool) {
	queue, ok := queues[entry.Queue]
	if ok {
		close_reader := make(chan bool)
		read_queue_ch := make(chan queueReadResult)
		go queue.read(entry.Opt, read_queue_ch, close_reader)
		for {
			result := <-read_queue_ch
			if result.err != nil {
				ch <- new_netresponse(NO_MSG, "")
			} else {
				status := SUCCESS
				if result.entryResult.has_more {
					status = HAS_MORE
				}
				resp := new_netresponse(status, result.entryResult.entry.msg)
				if result.entryResult.index > -1 {
					resp.Index = result.entryResult.index
				}
				ch <- resp
			}
			closed := <-done
			if closed {
				close_reader <- true
				break
			}
		}
	} else {
		ch <- new_netresponse(ERROR, "Queue does not exist")
	}
	close(ch)
}

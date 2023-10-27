package queue

import (
	"errors"
	"fmt"
	"time"
)

type Queue struct {
	name    string
	entries []Entry
	size    int
}

func new_queue(name string) *Queue {
	return &Queue{
		name: name,
	}
}

func (q *Queue) add_msg(msg string) {
	e := new_entry(msg)
	q.entries = append(q.entries, e)
	q.size++
}

type EntryResult struct {
	entry    Entry
	has_more bool
	index    int
}

func new_entryresult(entry Entry, has_more bool, index int) EntryResult {
	if index > -1 {
		return EntryResult{
			entry:    entry,
			has_more: has_more,
			index:    index,
		}
	} else {
		return EntryResult{
			entry:    entry,
			has_more: has_more,
		}
	}
}

type queueReadResult struct {
	entryResult EntryResult
	err         error
}

func (q *Queue) read(opt ReadOptions, ch chan queueReadResult, done chan bool) {
	starting_idx := opt.Index
	closed := false
	for !closed {
		select {
		case is_done := <-done:
			fmt.Println("closing down read func")
			closed = is_done
		default:
		}
		if q.size == 0 || starting_idx >= q.size+1 {
			if opt.Continuous {
				time.Sleep(1000 * time.Millisecond)
			} else {
				ch <- queueReadResult{
					entryResult: EntryResult{},
					err:         errors.New("Queue is empty or index is too large"),
				}
			}
		} else {
			has_more := false
			e := Entry{}
			idx := -1
			if !opt.Unread {
				q.entries[starting_idx].read = true
				e = q.entries[starting_idx]
				idx = starting_idx
				if len(q.entries) > starting_idx+1 {
					has_more = true
				}
			} else {
				for i := starting_idx; i < q.size; i++ {
					if !q.entries[i].read {
						if idx == -1 {
							e = q.entries[i]
							e.read = true
							q.entries[i] = e
							idx = i
						} else {
							// Look for additional entries to set has_more flag
							has_more = true
							break
						}
					}
				}
			}
			if idx == -1 {
				ch <- queueReadResult{
					entryResult: EntryResult{},
					err:         errors.New("no entry found"),
				}
				continue
			}
			if opt.Delete {
				ret := []Entry{}
				ret = append(ret, q.entries[:idx]...)
				ret = append(ret, q.entries[idx+1:]...)
				q.entries = ret
				q.size--
			}
			ch <- queueReadResult{
				entryResult: new_entryresult(e, has_more, idx),
				err:         nil,
			}
		}
		if !opt.Continuous {
			break
		}
	}
	close(ch)
}

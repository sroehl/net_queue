package net_queue

import "time"

type Entry struct {
	time int64
	msg  string
	read bool
}

func new_entry(msg string) Entry {
	return Entry{
		msg:  msg,
		read: false,
		time: time.Now().Unix(),
	}
}

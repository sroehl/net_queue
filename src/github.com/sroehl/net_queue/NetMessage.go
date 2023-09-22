package net_queue

const (
	CMD         = 1
	WRITE_ENTRY = 2
	READ_ENTRY  = 3
)

type NetMessage struct {
	NetMessageCmd
	NetMessageEntry
	Msg_type int
}

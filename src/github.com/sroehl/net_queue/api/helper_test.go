package api

import (
	"fmt"
	"net_queue/queue"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("Starting listener")
	go queue.Start_server(4545)
	m.Run()
	fmt.Println("Listener down")
}

package api

import (
	"fmt"
	"net_queue/queue"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("Starting listener")
	cfg := queue.New_config("localhost", 4545)
	go queue.Start_server(cfg)
	m.Run()
	fmt.Println("Listener down")
}

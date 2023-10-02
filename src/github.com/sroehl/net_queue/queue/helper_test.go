package queue

import (
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("Starting listener")
	cfg := New_config("localhost", 4545)
	go Start_server(cfg)
	m.Run()
	fmt.Println("Listener down")
}

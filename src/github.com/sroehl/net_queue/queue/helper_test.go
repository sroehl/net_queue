package queue

import (
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("Starting listener")
	go Start_server(4545)
	m.Run()
	fmt.Println("Listener down")
}

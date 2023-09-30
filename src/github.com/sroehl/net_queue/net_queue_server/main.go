package main

import (
	"fmt"
	"net_queue/queue"
)

func main() {
	port := 4545
	fmt.Printf("Starting server on %v", port)
	queue.Start_server(port)
}

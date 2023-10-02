package main

import (
	"flag"
	"fmt"
	"net_queue/queue"
	"os"
)

func main() {
	config_name := flag.String("config", ".config", "")
	flag.Parse()
	cfg, err := queue.Read_config(*config_name)
	if err != nil {
		fmt.Printf("Failed to read config file: %v", err)
		os.Exit(-1)
	}
	queue.Start_server(cfg)
}

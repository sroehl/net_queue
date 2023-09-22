package net_queue

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

var queues map[string]*Queue

type Net_Queue_Server struct {
	//running bool
}

/*func new_net_queue_server() Net_Queue_Server {
	return Net_Queue_Server{
		running: true,
	}
}*/

func start_server(port int) {
	queues = make(map[string]*Queue)
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handle_client(conn)
	}
}

func handle_client(c net.Conn) {
	defer c.Close()
	fmt.Printf("Client connected on %v\n", c.RemoteAddr())

	br := bufio.NewReader(c)
	for {
		buf, err := br.ReadBytes('\n')
		if err != nil {
			fmt.Printf("Failed reading bytes: %v\n", err)
			return
		}
		var net_message = &NetMessage{}
		err = json.Unmarshal(buf, net_message)
		if err != nil {
			fmt.Printf("Client parse failed: %v", err)
			return
		}
		var resp NetResponse
		if net_message.Msg_type == CMD {
			resp = net_message.handle_cmd(queues)
		} else if net_message.Msg_type == WRITE_ENTRY {
			resp = net_message.write_entry(queues)
		} else if net_message.Msg_type == READ_ENTRY {
			resp = net_message.read_entry(queues)
		}
		bytes, err := json.Marshal(resp)
		if err != nil {
			fmt.Printf("Client resp bad: %v", err)
			return
		}
		c.Write(bytes)
		c.Write([]byte{'\n'})
	}
}

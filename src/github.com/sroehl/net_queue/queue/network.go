package queue

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

func Start_server(config Config) {
	queues = make(map[string]*Queue)

	listen_addr := fmt.Sprintf("%v:%v", config.Server.ListenHost, config.Server.Port)
	fmt.Printf("Listening on %v\n", listen_addr)
	listener, err := net.Listen("tcp", listen_addr)
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
			resp = net_message.Handle_cmd(queues)
		} else if net_message.Msg_type == WRITE_ENTRY {
			resp = net_message.Write_entry(queues)
		} else if net_message.Msg_type == READ_ENTRY {
			resp = net_message.Read_entry(queues)
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

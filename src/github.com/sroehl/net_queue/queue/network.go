package queue

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"time"
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
			return
		}
		var net_message = &NetMessage{}
		err = json.Unmarshal(buf, net_message)
		if err != nil {
			fmt.Printf("Client parse failed: %v", err)
			return
		}
		var resp NetResponse
		if net_message.Msg_type == READ_ENTRY {
			ch := make(chan NetResponse)
			done := make(chan bool)
			go net_message.Read_entry(queues, ch, done)
			for {
				closed := check_closed(c)
				if !closed {
					resp, more := <-ch
					bytes, err := json.Marshal(resp)
					if err != nil {
						fmt.Printf("Client resp bad: %v", err)
						return
					}
					_, write_err := c.Write(append(bytes, '\n'))
					if write_err != nil {
						closed = true
					}
					if !more {
						break
					}
				}
				if closed {
					done <- true
					break
				} else {
					done <- false
				}
			}
		} else {
			if net_message.Msg_type == CMD {
				resp = net_message.Handle_cmd(queues)
			} else if net_message.Msg_type == WRITE_ENTRY {
				resp = net_message.Write_entry(queues)
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
}

func check_closed(c net.Conn) bool {
	one := make([]byte, 1)
	c.SetReadDeadline(time.Now())
	if _, err := c.Read(one); err == io.EOF {
		c.Close()
		return true
	} else {
		c.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
		return false
	}
}

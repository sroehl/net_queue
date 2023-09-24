package net_queue

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type Server struct {
	host      string
	port      int
	connected bool
	conn      net.Conn
}

func new_server(host string, port int) Server {
	return Server{
		host: host,
		port: port,
	}
}

func (s Server) create_queue(queue_name string) (NetResponse, error) {
	net_cmd := NetMessageCmd{
		Command:   CREATE_QUEUE,
		Arguments: queue_name,
	}
	return s.send_cmd(net_cmd)
}

func (s Server) delete_queue(queue_name string) (NetResponse, error) {
	net_cmd := NetMessageCmd{
		Command:   DELETE_QUEUE,
		Arguments: queue_name,
	}
	return s.send_cmd(net_cmd)
}

func (s Server) send_cmd(net_cmd NetMessageCmd) (NetResponse, error) {
	s.connect()
	net_msg := NetMessage{
		Msg_type:      CMD,
		NetMessageCmd: net_cmd,
	}
	encoded, err := json.Marshal(net_msg)
	if err != nil {
		fmt.Printf("Error marshalling in send_cmd")
		return NetResponse{}, err
	}
	br := bufio.NewReader(s.conn)
	s.conn.Write(encoded)
	s.conn.Write([]byte{'\n'})
	buffer, err := br.ReadBytes('\n')
	if err != nil {
		fmt.Printf("Error reading bytes %v\n", err)
		return NetResponse{}, err
	}
	var net_response = &NetResponse{}
	err = json.Unmarshal(buffer, net_response)
	if err != nil {
		fmt.Printf("Error unmarshalling in send_cmd")
		return NetResponse{}, err
	}
	return *net_response, nil
}

func (s *Server) connect() {
	if !s.connected {
		var err error
		tcpAddr, _ := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%v:%v", s.host, s.port))
		s.conn, err = net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			return
		}
	}
}

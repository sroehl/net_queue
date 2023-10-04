package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net_queue/queue"
)

type Server struct {
	host      string
	port      int
	connected bool
	conn      net.Conn
}

func New_server(host string, port int) Server {
	return Server{
		host: host,
		port: port,
	}
}

func (s Server) Create_queue(queue_name string) (queue.NetResponse, error) {
	net_cmd := queue.NetMessageCmd{
		Command:   queue.CREATE_QUEUE,
		Arguments: queue_name,
	}
	return s.send_cmd(net_cmd)
}

func (s Server) Delete_queue(queue_name string) (queue.NetResponse, error) {
	net_cmd := queue.NetMessageCmd{
		Command:   queue.DELETE_QUEUE,
		Arguments: queue_name,
	}
	return s.send_cmd(net_cmd)
}

func (s Server) List_queues() (queue.NetResponse, error) {
	net_cmd := queue.NetMessageCmd{
		Command: queue.LIST_QUEUES,
	}
	return s.send_cmd(net_cmd)
}

func (s Server) Purge_queue(queue_name string) (queue.NetResponse, error) {
	net_cmd := queue.NetMessageCmd{
		Command:   queue.PURGE_QUEUE,
		Arguments: queue_name,
	}
	return s.send_cmd(net_cmd)
}

func (s Server) send_cmd(net_cmd queue.NetMessageCmd) (queue.NetResponse, error) {
	s.Connect()
	net_msg := queue.NetMessage{
		Msg_type:      queue.CMD,
		NetMessageCmd: net_cmd,
	}
	encoded, err := json.Marshal(net_msg)
	if err != nil {
		fmt.Printf("Error marshalling in send_cmd")
		return queue.NetResponse{}, err
	}
	br := bufio.NewReader(s.conn)
	s.conn.Write(encoded)
	s.conn.Write([]byte{'\n'})
	buffer, err := br.ReadBytes('\n')
	if err != nil {
		fmt.Printf("Error reading bytes %v\n", err)
		return queue.NetResponse{}, err
	}
	var net_response = &queue.NetResponse{}
	err = json.Unmarshal(buffer, net_response)
	if err != nil {
		fmt.Printf("Error unmarshalling in send_cmd")
		return queue.NetResponse{}, err
	}
	return *net_response, nil
}

func (s Server) Write_msg(queue_name string, msg string) (queue.NetResponse, error) {
	net_entry := queue.NetMessageEntry{
		Queue: queue_name,
		Msg:   msg,
	}
	return s.send_msg(net_entry, queue.WRITE_ENTRY)
}

func (s Server) Read_msg(queue_name string, index int, unread bool, delete bool) (queue.NetResponse, error) {
	opt := queue.NetMessageEntryOptions{
		Index:  index,
		Unread: unread,
		Delete: delete,
	}
	net_entry := queue.NetMessageEntry{
		Queue: queue_name,
		Opt:   opt,
	}
	return s.send_msg(net_entry, queue.READ_ENTRY)
}

func (s Server) send_msg(net_entry queue.NetMessageEntry, msg_type int) (queue.NetResponse, error) {
	s.Connect()
	net_msg := queue.NetMessage{
		Msg_type:        msg_type,
		NetMessageEntry: net_entry,
	}
	encoded, err := json.Marshal(net_msg)
	if err != nil {
		fmt.Printf("Error marshalling in send_msg")
		return queue.NetResponse{}, err
	}
	br := bufio.NewReader(s.conn)
	s.conn.Write(encoded)
	s.conn.Write([]byte{'\n'})
	buffer, err := br.ReadBytes('\n')
	if err != nil {
		fmt.Printf("Error reading bytes %v\n", err)
		return queue.NetResponse{}, err
	}
	var net_response = &queue.NetResponse{}
	err = json.Unmarshal(buffer, net_response)
	if err != nil {
		fmt.Printf("Error unmarshalling in send_msg")
		return queue.NetResponse{}, err
	}
	return *net_response, nil
}

func (s *Server) Connect() {
	if !s.connected {
		var err error
		tcpAddr, _ := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%v:%v", s.host, s.port))
		s.conn, err = net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			return
		}
	}
}

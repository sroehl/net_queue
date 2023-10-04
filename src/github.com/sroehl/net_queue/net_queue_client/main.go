package main

import (
	"fmt"
	"net_queue/api"
	"net_queue/queue"

	"github.com/jessevdk/go-flags"
)

type options struct {
	Host           string `short:"s" long:"server" description:"Hostname or IP address of net_queue server" required:"true"`
	Port           int    `short:"p" long:"port" description:"Host of net_queue server" required:"true"`
	Queue          string `short:"q" long:"queue" description:"Name of the queue" required:"true"`
	Message        string `short:"m" long:"message" description:"Message to send"`
	Read           bool   `short:"r" long:"read" description:"Set to read a message from queue"`
	Stay_Connected bool   `long:"continous" description:"Continously check for new messages until stopped"`
	Create         bool   `long:"create" description:"Create queue"`
	Purge          bool   `long:"purge" description:"Purge queue"`
	Delete         bool   `long:"delete" description:"Delete queue"`
}

func main() {
	var opts options
	_, err := flags.Parse(&opts)
	if err != nil {
		fmt.Printf("Failed to parse args: %v\n", err)
	}

	fmt.Printf("Connecting to %v:%v\n", opts.Host, opts.Port)
	s := api.New_server(opts.Host, opts.Port)
	s.Connect()
	if opts.Read {
		read_queue(s, opts.Queue, opts.Stay_Connected)
	} else if len(opts.Message) > 0 {
		write_queue(s, opts.Queue, opts.Message)
	} else if opts.Create {
		create_queue(s, opts.Queue)
	} else if opts.Purge {
		purge_queue(s, opts.Queue)
	} else if opts.Delete {
		delete_queue(s, opts.Queue)
	} else {
		fmt.Printf("Invalid usage: must read/write message or send command (create, purge, delete)\n")
	}
}

func read_queue(s api.Server, queue_name string, stay_connected bool) {
	for {
		resp, err := s.Read_msg(queue_name, 0, true, true)
		if err != nil {
			fmt.Printf("Failed: %v", err)
			return
		}
		if resp.Status == queue.SUCCESS {
			fmt.Printf("%v\n", resp.Msg)
		} else if resp.Status == queue.NO_MSG && !stay_connected {
			fmt.Printf("No messages!\n")
		} else if resp.Status == queue.ERROR {
			fmt.Printf("Failed: %v\n", resp.Msg)
		}
		if !stay_connected {
			break
		}
	}
}

func write_queue(s api.Server, queue_name string, msg string) {
	resp, err := s.Write_msg(queue_name, msg)
	if !handle_error(err, resp) {
		return
	}
	fmt.Printf("Message sent\n")
}

func create_queue(s api.Server, queue_name string) {
	resp, err := s.Create_queue(queue_name)
	if !handle_error(err, resp) {
		return
	}
	fmt.Printf("'%v' created successfully\n", queue_name)
}

func purge_queue(s api.Server, queue_name string) {
	resp, err := s.Purge_queue(queue_name)
	if !handle_error(err, resp) {
		return
	}
	fmt.Printf("'%v' purged successfully\n", queue_name)
}

func delete_queue(s api.Server, queue_name string) {
	resp, err := s.Delete_queue(queue_name)
	if !handle_error(err, resp) {
		return
	}
	fmt.Printf("'%v' deleted successfully\n", queue_name)
}

func handle_error(err error, resp queue.NetResponse) bool {
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		return false
	}
	if resp.Status != queue.SUCCESS {
		fmt.Printf("Failed: %v\n", resp.Msg)
		return false
	}
	return true
}

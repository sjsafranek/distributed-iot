package main

import (
	"bufio"
	"fmt"
	"net"
)

type StreamClient struct {
	conn net.Conn
}

func (self *StreamClient) ReadLine() (string, error) {
	return bufio.NewReader(self.conn).ReadString('\n')
}

func (self *StreamClient) ReadAll(quit chan bool, clbk func(string, error)) {
	for {
		select {
		case <-quit:
			break
		default:
			clbk(self.ReadLine())
		}
	}
}

func OpenStream(host string, port int) (*StreamClient, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", host, port))
	return &StreamClient{conn: conn}, err
}

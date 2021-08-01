package main

/*

import (
	"bufio"
	"fmt"

	"github.com/tarm/serial"
)

type SerialClient struct {
	conn *serial.Port
}

func (self *SerialClient) ReadLine() (string, error) {
	return bufio.NewReader(self.conn).ReadString('\n')
}

func (self *SerialClient) ReadAll(quit chan bool, clbk func(string, error)) {
	scanner := bufio.NewScanner(self.conn)
	scanner.Split(bufio.ScanLines)

	for {
		select {
		case <-quit:
			break
		default:
			for scanner.Scan() {
				fmt.Println("Scanning")
				line := scanner.Text() // Println will add back the final '\n'
				fmt.Println(line)
				clbk(line, nil)
			}
			// err := scanner.Err()
			// line, err := reader.ReadString('\n')
			// fmt.Println(line, err)
			// clbk(line, err)
		}
	}

}

func OpenSerial(port string) (*SerialClient, error) {
	config := &serial.Config{
		Name:        port,
		Baud:        9600,
		ReadTimeout: 5,
	}
	conn, err := serial.OpenPort(config)
	return &SerialClient{conn: conn}, err
}

*/

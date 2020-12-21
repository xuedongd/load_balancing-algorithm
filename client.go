package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "49.232.103.154:8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	//read msg
	go func(conn net.Conn) error {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				return err
			}
			fmt.Printf(string(buf[:n]))
		}
	}(conn)
	//write msg
	buf := make([]byte, 1024)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		n, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func failOnErr(err error, str string) {
	if err != nil {
		fmt.Println(str, err.Error())
		os.Exit(1)
	}
}
func handleConn(conn net.Conn) {
	buff := make([]byte, 1024)
	for {
		_, err := conn.Read(buff)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading sent data", err.Error())
				os.Exit(1)
			}
			break
		}
		res := "+" + "PONG\r\n"
		writeBuffer := []byte(res)
		_, err = conn.Write(writeBuffer)

		if err != nil {
			fmt.Println("Error writing data", err.Error())
			os.Exit(1)
		}
	}
	conn.Close()
}

func main() {
	fmt.Println("Logs from your program will appear here!")
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	failOnErr(err, "Fail to listen")
	for err == nil {
		conn, err := l.Accept()
		failOnErr(err, "Fail to Accept")
		fmt.Println("Accepted a connection")
		go handleConn(conn)

	}
}

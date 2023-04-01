package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func failOnErr(err error, str string) {
	if err != nil {
		fmt.Println(str, err.Error())
		os.Exit(1)
	}
}
func handleConn(conn net.Conn) {

	buff := make([]byte, 1024)
	n, err := conn.Read(buff)

	strArr := strings.Split(string(buff[:n]), "\r\n")

	// fmt.Print(len(strArr), strArr)
	if err != nil {
		fmt.Println("Error reading sent data", err.Error())
		os.Exit(1)
	}

	res := "$" + strconv.Itoa(len(strArr)-1) + "\r\n" + "PONG\r\n"
	writeBuffer := []byte(res)
	_, err = conn.Write(writeBuffer)

	if err != nil {
		fmt.Println("Error writing data", err.Error())
		os.Exit(1)
	}

	// n, err = conn.Read(buff[:])
	// fmt.Printf("%s", strings.Split(string(buff[:n]), "\r\n"))

	conn.Close()
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
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

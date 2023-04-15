package main

import (
	"fmt"
	"io"
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

func parseCommand(command string) string {

	list := strings.Split(command, "\r\n")
	// fmt.Println(list)
	_, err := strconv.Atoi(strings.TrimPrefix(list[0], "*"))
	failOnErr(err, "Failed to parse")

	var strOutput []string
	switch strings.ToLower(list[2]) {
	case "echo":

		echoConcat := ""
		for i := 4; i < len(list); i += 2 {
			echoConcat += list[i]
		}
		strOutput = []string{"$" + strconv.Itoa(len(echoConcat)), "\r\n", echoConcat, "\r\n"}

	case "ping":
		strOutput = []string{"+PONG\r\n"}
	default:
		fmt.Println("Something wrong happened")
	}

	//fmt.Println(strOutput)

	return strings.Join(strOutput, "")
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

		response := parseCommand(string(buff))

		writeBuffer := []byte(response)
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

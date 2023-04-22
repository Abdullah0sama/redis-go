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

func decode(command string) []string {

	out := []string{}
	commandParams := strings.Split(command, "\r\n")
	switch command[0] {
	case '*':
		for i := 2; i < len(commandParams); i += 2 {
			out = append(out, commandParams[i])
		}
	case '+':
		out = append(out, commandParams[1])
	}

	return out
}

func parseCommand(command string) string {

	commandDecoded := decode(command)
	var strOutput []string
	// fmt.Println(commandDecoded)
	switch strings.ToLower(commandDecoded[0]) {
	case "echo":
		echoConcat := strings.Join(commandDecoded[1:], " ")
		strOutput = []string{"$" + strconv.Itoa(len(echoConcat)), "\r\n", echoConcat, "\r\n"}
	case "ping":
		strOutput = []string{"+PONG\r\n"}
	case "set":
		key, value := commandDecoded[1], commandDecoded[2]
		var expiryTime int64 = -1
		if len(commandDecoded) > 3 && strings.ToLower(commandDecoded[3]) == "px" {
			parsedTime, _ := strconv.Atoi(commandDecoded[4])
			expiryTime = int64(parsedTime)
		}

		storage.set(key, value, expiryTime)
		strOutput = []string{"+OK\r\n"}
	case "get":
		key := commandDecoded[1]
		value, ok := storage.get(key)
		if ok {
			strOutput = []string{"+", value, "\r\n"}
		} else {
			strOutput = []string{"$-1\r\n"}
		}
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

var storage *Storage

func main() {
	storage = NewStorage()

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

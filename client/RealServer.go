package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Enter IP: ")
	// IP, _ := reader.ReadString('\n')
	// stripIP := strings.TrimSpace(IP)

	for {

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		strEcho, _ := reader.ReadString('\n')
		strip := strings.TrimSpace(strEcho)
		buf := []byte(strEcho)

		if strip == "quit" {
			fmt.Println("bye")
			return
		}

		// servAddr := stripIP + ":3333"
		// tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
		// if err != nil {
		// 	println("ResolveTCPAddr failed:", err.Error())
		// 	os.Exit(1)
		// }

		// conn, err := net.DialTCP("tcp", nil, tcpAddr)
		// if err != nil {
		// 	println("Dial failed:", err.Error())
		// 	os.Exit(1)
		// }

		if buf[0] == 'c' && buf[1] == 'd' {
			fmt.Print("Enter full directory: ")
			directory, _ := reader.ReadString('\n')
			strEcho = "cd " + directory
		}
		_, err := conn.Write([]byte(strEcho))
		if err != nil {
			println("Write to server failed:", err.Error())
		}

		println("write to server =", strEcho)

		reply := make([]byte, 8192)

		_, err = conn.Read(reply)
		if err != nil {
			println("no result output", err.Error())
		}

		println("reply from server=", string(reply))

		// conn.Close()
	}

}

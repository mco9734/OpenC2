package main

import (
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
	fmt.Println("Welcome to OpenC2!\nListening on " + CONN_HOST + ":" + CONN_PORT)
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

		path, patherr := os.Getwd()
		if patherr != nil {
			fmt.Println(patherr)
		}

		// reader := bufio.NewReader(os.Stdin)
		// fmt.Print("Enter text: ")
		// strEcho, _ := reader.ReadString('\n')
		strEcho := "cd .."
		strip := strings.TrimSpace(strEcho)
		// buf := []byte(strEcho)
		stripSplit := strings.Split(strip, " ")

		if strip == "quit" {
			fmt.Println("bye")
			return
		}

		if stripSplit[0] == "cd" {
			// fmt.Print("Enter full directory: ")
			// directory, _ := reader.ReadString('\n')
			// strEcho = "cd " + directory
			if stripSplit[1] == ".." {
				splitPath := strings.Split(path, "\\")
				newDirectory := ""
				for i := 0; i < len(splitPath)-1; i++ {
					newDirectory += splitPath[i] + "\\"
				}
				strEcho = "cd" + newDirectory
			} else {
				strEcho = "cd" + path + "\\" + stripSplit[1]
			}
		}

		communicate(conn, strEcho)
		// conn.Close()
	}

}

func communicate(conn net.Conn, strEcho string) {
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
}

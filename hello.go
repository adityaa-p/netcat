package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)

		if err != nil {
			fmt.Println("Error reading:", err)
			break
		}
		fmt.Print(string(buffer[:n]))

		response := "Thanks for connecting!\n"
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing:", err)
			break
		}
	}
}

func runTCPServer() {
	listener := startTcpServer("tcp")
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn) // Launch a goroutine to handle each connection
	}
}

func runUDPServer() {
	conn := startUdpServer()
	if conn == nil {
		fmt.Println("Error starting udp server")
		return
	}

	defer conn.Close()
	handleUdpConnection(conn)
}

func runLocalhostMode(ports string) {
	if strings.Contains(ports, "-") {
		start, _ := strconv.Atoi(strings.Split(ports, "-")[0])
		end, _ := strconv.Atoi(strings.Split(ports, "-")[1])

		for i := start; i <= end; i++ {
			port := ":" + strconv.Itoa(i)
			_, err := net.Dial("tcp", port)
			if err != nil {
				fmt.Println("Error connecting to the server. Port: ", strconv.Itoa(i))
				continue
			}

			fmt.Println("Connection successfull")
			break
		}
	} else {
		_, err := net.Dial("tcp", ":"+ports)
		if err != nil {
			fmt.Println("Error connecting to the server. Port: ", ports)
			return
		}

		fmt.Println("Connection succesfull")
	}
}

func runCommandMode() {
	listener := startTcpServer("tcp")

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		response := "Thanks for connecting!\n"
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing:", err)
			break
		}

		cmd := exec.Command("/bin/bash")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Println("Failed to get stdin:", err)
			return
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Println("Failed to get stdout:", err)
			return
		}

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			// Reading from stdout and sending to WebSocket
			buf := make([]byte, 1024)
			for {
				n, err := stdout.Read(buf)
				if err != nil {
					log.Println("Error reading stdout:", err)
					break
				}
				fmt.Println(string(buf[:n]))
				if err != nil {
					log.Println("Error sending message:", err)
					break
				}
			}
			wg.Done()
		}()

		go func() {
			// Reading messages from the WebSocket and writing to stdin
			for {
				buffer := make([]byte, 1024)
				input, err := conn.Read(buffer)
				if err != nil {
					log.Println("Error reading message from WebSocket:", err)
					break
				}
				fmt.Println(string(buffer[:input]))
				_, err = stdin.Write(buffer[:input])
				if err != nil {
					log.Println("Error writing to stdin:", err)
					break
				}
			}
			wg.Done()
		}()

		err = cmd.Start()
		if err != nil {
			log.Println("Failed to start command:", err)
			return
		}
		wg.Wait()

		// Ensure the command has finished before returning
		_ = cmd.Wait()
	}
}

func main() {
	mode := os.Args[1]

	switch mode {
	case "tcp":
		runTCPServer()

	case "udp":
		runUDPServer()

	case "localhost":
		runLocalhostMode(os.Args[2])

	case "-e":
		runCommandMode()

	default:
		fmt.Println("Unknown mode")
	}
}

func handleUdpConnection(conn *net.UDPConn) {
	buf := make([]byte, 1024)

	fmt.Println("UDP server listening on port", ":8080")

	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		fmt.Printf("Received %d bytes from %s: %s\n", n, addr, string(buf[:n]))

		// You can optionally send a response here
		// _, err = conn.WriteToUDP([]byte("Hello from server!"), addr)
		// if err != nil {
		//     fmt.Println("Error sending response:", err)
		// }
	}
}

func startUdpServer() *net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return nil
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return nil
	}
	return conn
}

func startTcpServer(mode string) net.Listener {
	listener, err := net.Listen(mode, ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return listener
	} else {
		fmt.Println("Server listening on port 8080")
		return listener
	}
}

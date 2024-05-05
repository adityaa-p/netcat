package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func sendResponseToClient(conn net.Conn) {
	defer conn.Close()
	for {
		fmt.Println("Waiting to send data to client")
		time.Sleep(8 * time.Second)
		fmt.Println(os.Args)

		// args := os.Args[1]
		// // response := "Thanks for connecting!\n"

		_, err := conn.Write([]byte("Hello Aditya"))
		if err != nil {
			fmt.Println("Error writing:", err)
			break
		}
	}
}

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

		//Simulate processing the request
		response := "Thanks for connecting!\n"
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing:", err)
			break
		}
	}
}

func main() {

	args := os.Args
	mode := args[1]

	var listener net.Listener
	if mode == "tcp" {
		listener = startServer(mode)
		defer listener.Close()

		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				continue
			}
			go handleConnection(conn) // Launch a goroutine to handle each connection
			// go sendResponseToClient(conn)
		}
	} else {
		addr, err := net.ResolveUDPAddr("udp", ":8080")
		if err != nil {
			fmt.Println("Error resolving address:", err)
			return
		}

		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			fmt.Println("Error listening:", err)
			return
		}
		defer conn.Close()

		buf := make([]byte, 1024)

		fmt.Println("UDP server listening on port", ":8080")

		// Infinite loop to receive messages
		for {
			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("Error reading:", err)
				continue
			}

			// Process the received data
			fmt.Printf("Received %d bytes from %s: %s\n", n, addr, string(buf[:n]))

			// You can optionally send a response here
			// _, err = conn.WriteToUDP([]byte("Hello from server!"), addr)
			// if err != nil {
			//     fmt.Println("Error sending response:", err)
			// }
		}
	}

}

func startServer(mode string) net.Listener {
	listener, err := net.Listen(mode, ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return listener
	} else {
		fmt.Println("Server listening on port 8080")
		return listener
	}
}

package handlers

import (
	"ascii/utils"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func Setup(port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server listening on %v\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {

			if err == io.EOF || strings.Contains(err.Error(), "forcibly closed") || strings.Contains(err.Error(), "use of closed") {
				log.Println("Connection closed !!!")
				return
			}

			log.Println("Error reading:", err)
			return
		}

		log.Printf("READ bytes: %v\n", n)
		packet, err := utils.Deserialize(buf[:n])

		if err != nil {
			log.Println("NOT ABLE TO DESERIALIZE: ", err)
			continue
		}

		HandlePacketType(&packet, conn)
	}
}

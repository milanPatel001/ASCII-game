package main

import (
	"ascii/utils"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {

	url := "localhost:3000"

	conn, err := net.Dial("tcp", url)
	if err != nil {
		log.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	isAuthenticated := authenticateClient(conn)

	if !isAuthenticated {
		conn.Close()
		return
	}

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)

			if err == io.EOF || strings.Contains(err.Error(), "forcibly closed") {
				return
			}

			continue
		}

		log.Printf("READ: %v bytes\n", n)
	}

}

func authenticateClient(conn net.Conn) bool {

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Enter Auth token: ")

	if !scanner.Scan() {
		return false
	}

	input := scanner.Text()
	ip, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		log.Println(err)
		return false
	}

	fmt.Println(ip)
	if ip == "::1" {
		ip = "127.0.0.1"
	}

	packet := utils.NewPacket(net.ParseIP(ip), utils.AUTH, []byte(input))

	pb, err := packet.Serialize()
	if err != nil {
		return false
	}

	_, err = conn.Write(pb)
	if err != nil {
		log.Println(err)
		return false
	}

	buf := make([]byte, 240)

	n, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
		return false
	}

	log.Println(string(buf[:n]))

	return true

}

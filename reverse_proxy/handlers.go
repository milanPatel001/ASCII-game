package main

import (
	"ascii/utils"
	"log"
	"net"
)

func HandlePacketType(packet *utils.Packet, conn net.Conn) {
	switch packet.MessageType {
	case utils.AUTH:
		if !authenticate(string(packet.Payload)) {
			conn.Write([]byte("Not authenticated !!!"))
			conn.Close()
			return
		}
		conn.Write([]byte("Authenticated !!!"))
		break
	default:
		log.Println("Unkown packet type !!!")
		conn.Close()
	}
}

func authenticate(token string) bool {
	tokens := []string{"ABCosp", "OPOOO", "JKASSS"}

	for _, t := range tokens {
		if t == token {
			return true
		}
	}

	return false
}

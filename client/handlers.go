package main

import (
	"ascii/utils"
	"log"
	"net"
)

func HandleServerPacket(packet *utils.Packet, conn net.Conn) {
	switch packet.MessageType {
	case utils.CREATE_GROUP:

	default:
		log.Println("Unkown packet type !!!")
		conn.Close()
	}
}

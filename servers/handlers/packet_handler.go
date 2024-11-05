package handlers

import (
	"ascii/client/game"
	"ascii/utils"
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strings"
)

var PlayerIdMap = make(map[string]net.Conn)
var RoomMap = make(map[string]*game.Room) // room code -> info

func HandlePacketType(packet *utils.Packet, conn net.Conn) {
	switch packet.MessageType {
	case utils.AUTH:
		handleAuthPacket(packet, conn)

	case utils.CREATE_GROUP:
		handleGroupCreatePacket(packet)

	case utils.JOIN_GROUP:
		handleGroupJoinPacket(packet, conn)

	case utils.DESTROY_ROOM:
		roomCode := string(packet.Payload)
		pkt, _ := utils.CreatePacketAndSerialize("127.0.0.1", utils.DESTROY_ROOM, nil)

		room := RoomMap[roomCode].PlayersJoined

		for i, id := range room {
			if i == 0 {
				continue
			}

			PlayerIdMap[id].Write(pkt)
		}

		delete(RoomMap, roomCode)

	case utils.PLAYER_LEFT:
		s := strings.Split(string(packet.Payload), " | ")
		roomCode := s[0]
		playerLeftId := s[1]

		room := RoomMap[roomCode]

		pkt, _ := utils.CreatePacketAndSerialize("127.0.0.1", utils.PLAYER_LEFT, []byte(playerLeftId))

		i := 0
		for ind, id := range room.PlayersJoined {
			if id == playerLeftId {
				i = ind
			} else {
				PlayerIdMap[id].Write(pkt)
			}
		}

		if i < len(room.PlayersJoined)-1 {
			room.PlayersJoined = append(room.PlayersJoined[:i], room.PlayersJoined[i+1:]...)
		} else {
			room.PlayersJoined = room.PlayersJoined[:i]
		}

	default:
		log.Println("Unkown packet type !!!")
		conn.Close()
	}
}

func handleGroupJoinPacket(packet *utils.Packet, conn net.Conn) {
	log.Println("JOIN GROUP ENTERED ...")
	s := strings.Split(string(packet.Payload), " | ")
	code := s[0]
	playerId := s[1]

	if room, exists := RoomMap[code]; exists {
		room.PlayersJoined = append(room.PlayersJoined, playerId)

		b, err := utils.ConvComplexPayloadToBytes(room)
		if err != nil {
			log.Fatal(err)
		}

		packet, _ := utils.CreatePacketAndSerialize("127.0.0.1", utils.JOIN_GROUP, b)

		fmt.Printf("Struct after adding: %+v\n", room)
		conn.Write(packet)

		// Broadcast to other player waiting ...
		pkt, err := utils.CreatePacketAndSerialize("127.0.0.1", utils.BROADCAST, []byte(playerId))
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Players joined", len(room.PlayersJoined))

		for _, id := range room.PlayersJoined {
			if id == playerId {
				continue
			}

			_, err = PlayerIdMap[id].Write(pkt)
			if err != nil {
				log.Println(err)
			}
		}

	} else {
		emptyPkt, err := utils.CreatePacketAndSerialize("127.0.0.1", utils.NOT_FOUND, nil)
		if err != nil {
			log.Fatal(err)
		}
		conn.Write(emptyPkt)
	}
}

func handleGroupCreatePacket(packet *utils.Packet) {
	log.Println("CREATE GROUP ENTERED ...")
	buf := bytes.NewBuffer(packet.Payload)

	var room game.Room
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&room)
	if err != nil {
		log.Fatal("gob.Decode failed:", err)
	}

	// Print the restored struct
	fmt.Printf("Restored struct: %+v\n", room)
	RoomMap[room.Code] = &room
}

func handleAuthPacket(packet *utils.Packet, conn net.Conn) {
	if !verifySessionToken(string(packet.Payload), utils.ConvBytesToIpv4(packet.SrcIP)) {
		conn.Write([]byte("-1"))
		conn.Close()
		return
	}

	// add id to player conn map
	playerId := utils.GeneratePlayerId()
	PlayerIdMap[playerId] = conn

	conn.Write([]byte(playerId))
}

func verifySessionToken(token string, ip net.IP) bool {

	rdb := GetRedisInstance()

	resToken, err := rdb.Get(context.Background(), ip.String()).Result()

	if err != nil {
		log.Println(err)
		return false
	}

	return token == resToken

}

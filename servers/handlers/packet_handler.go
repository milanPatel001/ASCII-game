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
	"strconv"
	"strings"
)

var PlayerIdMap = make(map[string]net.Conn) // player id -> tcp conn
var RoomMap = make(map[string]*game.Room)   // room code -> info
//var

func HandlePacketType(packet *utils.Packet, conn net.Conn) {
	switch packet.MessageType {
	case utils.AUTH:
		handleAuthPacket(packet, conn)

	case utils.CREATE_GROUP:
		handleGroupCreatePacket(packet)

	case utils.JOIN_GROUP:
		handleGroupJoinPacket(packet, conn)

	case utils.DESTROY_ROOM:
		handleRoomDestroyPacket(packet)

	case utils.PLAYER_LEFT:
		handlePlayerLeftPacket(packet)

	case utils.START_GAME:
		handleStartGamePacket(packet)

	case utils.PLAYER_MOVE:
		handlePlayerMovement(packet)

	case utils.TERRAIN_EXIT:
		handleTerrainExit(packet)

	default:
		log.Println("Unkown packet type !!!")
		conn.Close()
	}
}

func handlePlayerMovement(packet *utils.Packet) {
	// curr pos, last pos, room code, player id
	var playerMovementPayload game.PlayerMovementPayload

	err := utils.GetComplexPayloadFromBytes(packet.Payload, &playerMovementPayload)
	if err != nil {
		log.Println(err)
		return
	}

	players := RoomMap[playerMovementPayload.RoomCode].GameState.Players
	pkt, _ := packet.Serialize()

	// update player coords in server and send payload to other players
	for i := range players {
		if players[i].Id == playerMovementPayload.PlayerId {
			players[i].Pos = playerMovementPayload.CurrPos
			continue
		}

		PlayerIdMap[players[i].Id].Write(pkt)

	}

}

// *********************************************** GAME MENU STUFF ***********************************************

func handleTerrainExit(packet *utils.Packet) {
	payload := string(packet.Payload)

	str := strings.Split(payload, " | ")

	roomCode := str[0]
	playerId := str[1]
	exitSeed, err := strconv.Atoi(str[2])

	if err != nil {
		log.Fatal(err)
	}

	room := RoomMap[roomCode]

	pkt, err := utils.CreatePacketAndSerialize("127.0.0.1", utils.TERRAIN_EXIT, []byte(fmt.Sprintf("%v | %v", playerId, exitSeed)))

	for i, pl := range room.GameState.Players {
		if pl.Id == playerId {
			room.GameState.Players[i].CurrSeed = exitSeed

			exitCoords := room.GameState.Terrains[game.GetTerrainIndexUsingSeed(room.GameState.Terrains, exitSeed)].ExitCoord

			seedIndex := 0
			for j, ex := range exitCoords {
				if ex.ExitSeed == exitSeed {
					seedIndex = j
					break
				}
			}

			room.GameState.Players[i].Pos = exitCoords[seedIndex].Pos
		} else {
			PlayerIdMap[room.PlayersJoined[i]].Write(pkt)
		}
	}

}

func handleStartGamePacket(packet *utils.Packet) {
	// TODO : CREATE new packet type and CREATE NEW GAME STATE HERE, and then send all the seeds to clients

	roomCode := string(packet.Payload)
	room := RoomMap[roomCode]

	playerSeeds, middleGroundSeeds := game.RandomSeedAssigner(len(room.PlayersJoined))
	room.GameState = game.NewGameState(86, 26, room, playerSeeds, middleGroundSeeds)

	payload, err := utils.ConvComplexPayloadToBytes(game.GameStartPayload{PlayerSeeds: playerSeeds, MiddleGroundSeeds: middleGroundSeeds})
	if err != nil {
		log.Fatal(err)
	}

	pkt, _ := utils.CreatePacketAndSerialize("127.0.0.1", utils.START_GAME, payload)

	for i := 0; i < len(room.PlayersJoined); i++ {
		PlayerIdMap[room.PlayersJoined[i]].Write(pkt)
	}
}

func handleRoomDestroyPacket(packet *utils.Packet) {
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

}

func handlePlayerLeftPacket(packet *utils.Packet) {
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

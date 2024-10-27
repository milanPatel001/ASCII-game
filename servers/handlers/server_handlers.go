package handlers

import (
	"ascii/utils"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func GetRedisInstance() *redis.Client {

	addr := os.Getenv("REDIS_ADDR")
	psk := os.Getenv("REDIS_PSWD")

	if redisClient == nil {
		return redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: psk,
		})
	}

	return redisClient
}

func DisconnectRedis() error {
	if redisClient == nil {
		return nil
	}

	err := redisClient.Close()
	redisClient = nil

	return err
}

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

func HandlePacketType(packet *utils.Packet, conn net.Conn) {
	switch packet.MessageType {
	case utils.AUTH:
		if !verifySessionToken(string(packet.Payload), utils.ConvBytesToIpv4(packet.SrcIP)) {
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

func verifySessionToken(token string, ip net.IP) bool {

	rdb := GetRedisInstance()

	resToken, err := rdb.Get(context.Background(), ip.String()).Result()

	if err != nil {
		log.Println(err)
		return false
	}

	return token == resToken

}

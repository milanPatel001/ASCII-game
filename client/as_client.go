package main

import (
	"ascii/client/game"
	"ascii/utils"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

func authServerFlow() {

	url := "http://127.0.0.1:3000/auth"

	resp, err := sendRestAuthToken(url)

	if err != nil {
		log.Fatal(err)
	}
	s := strings.Split(string(resp), " | ")

	token := s[0]
	addr := s[1]

	connectToGameServer(token, addr)

}

func connectToGameServer(token, addr string) {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	ip, err := utils.GetIpFromRemoteAddr(conn.RemoteAddr().String())
	if err != nil {
		log.Fatal(err)
	}

	playerId, err := sendSessionToken(conn, ip, token)

	if err != nil {
		log.Fatal(err)
	}

	// Here the game starts
	game.InitializeGame(conn, string(playerId))

}

// authentication by auth server (REST)
func sendRestAuthToken(url string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "TOKENn")

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		log.Fatal("Not authenticated !!!")
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

// Authentication by game server (SOCKET)
func sendSessionToken(conn net.Conn, ip string, token string) ([]byte, error) {
	pb, err := utils.CreatePacketAndSerialize(ip, utils.AUTH, []byte(token))
	if err != nil {
		return nil, err
	}

	_, err = conn.Write(pb)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 240)

	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	if string(buf[:n]) == "-1" {
		return nil, utils.AUTH_ERROR
	}

	return buf[:n], nil
}

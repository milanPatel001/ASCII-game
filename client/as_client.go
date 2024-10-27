package main

import (
	"ascii/utils"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

func authServerFlow() {

	url := "http://127.0.0.1:3000/auth"

	resp, err := sendAuthToken(url)

	if err != nil {
		log.Fatal(err)
	}

	connectToGameServer(string(resp))

}

func connectToGameServer(resp string) error {

	s := strings.Split(resp, " | ")

	token := s[0]
	addr := s[1]

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	defer conn.Close()

	ip, err := utils.GetIpFromRemoteAddr(conn.RemoteAddr().String())
	if err != nil {
		return err
	}

	pb, err := utils.CreatePacketAndSerialize(ip, utils.AUTH, token)
	if err != nil {
		return err
	}

	_, err = conn.Write(pb)
	if err != nil {
		return err
	}

	buf := make([]byte, 240)

	n, err := conn.Read(buf)
	if err != nil {
		return err
	}

	log.Println(string(buf[:n]))

	return nil
}

func sendAuthToken(url string) ([]byte, error) {
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

package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	AUTH_SERVER_FLOW = iota
	REVERSE_PROXY_FLOW
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}

	flow, err := strconv.Atoi(os.Getenv("FLOW"))
	if err != nil {
		log.Fatal(err)
	}

	if flow == REVERSE_PROXY_FLOW {
		reverseProxyFlow()
	} else {
		authServerFlow()
	}

}

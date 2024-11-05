package main

import (
	"ascii/servers/handlers"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("GAME_SERVER_PORT")
	if err != nil {
		log.Fatal(err)
	}

	handlers.Setup(port)
}

package main

import (
	"ascii/servers/handlers"
)

func main() {
	handlers.Setup(":3001")
}

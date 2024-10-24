package reverse_proxy

import (
	"ascii/servers/handlers"
)

func Setup() {
	handlers.Setup(":3000")
}

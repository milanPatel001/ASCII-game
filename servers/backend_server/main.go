package main

import (
	"ascii/servers/backend_server/auth_server"
	"ascii/servers/backend_server/reverse_proxy"
)

const (
	AUTH_SERVER_FLOW = 0
	REVERSE_PROXY_FLOW
)

func main() {
	flow := AUTH_SERVER_FLOW // TODO use cli args or env vars

	if flow == REVERSE_PROXY_FLOW {
		reverse_proxy.Setup()
	} else {
		auth_server.Setup()
	}

}

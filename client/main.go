package main

const (
	AUTH_SERVER_FLOW = 0
	REVERSE_PROXY_FLOW
)

func main() {
	flow := AUTH_SERVER_FLOW // TODO get using env vars or cli arg

	if flow == REVERSE_PROXY_FLOW {
		reverseProxyFlow()
	} else {
		authServerFlow()
	}

}

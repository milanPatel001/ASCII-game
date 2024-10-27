package auth_server

import (
	"ascii/servers/handlers"
	"ascii/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func Setup() {

	port := ":3000"
	gameServerAddr := ":3001"

	rdb := handlers.GetRedisInstance()

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		arr := r.Header["Authorization"]

		if len(arr) == 0 {
			http.Error(w, "Auth failed : No token provided in authorization header", http.StatusUnauthorized)
		}

		token := arr[0]

		if !authenticateClient(token) {
			http.Error(w, "Auth failed : Invalid token provided", http.StatusUnauthorized)
		}

		// STORE in CACHE
		ip, err := utils.GetIpFromRemoteAddr(r.RemoteAddr)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		err = rdb.Set(context.Background(), ip, "my_session_token", time.Minute*60).Err()

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		// TODO generate a new session token and then store token + src ip in global cache

		sessionToken := fmt.Sprintf("my_session_token | %v", gameServerAddr)
		w.Write([]byte(sessionToken))
	})

	log.Printf("Starting server on port %v\n", port)
	err := http.ListenAndServe(port, nil)

	if err != nil && err != http.ErrServerClosed {
		log.Printf("ListenAndServe() error: %v", err)
	}

	log.Println("Server stopped")
}

func authenticateClient(token string) bool {

	return true
}

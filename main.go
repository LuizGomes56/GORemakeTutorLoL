package main

import (
	"fmt"
	"golang/functions"
	"golang/routes"
	"log"
	"net/http"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	db, err := functions.ConnectDB()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/game/last", func(w http.ResponseWriter, r *http.Request) {
		routes.LastByCode(w, r, db)
	})

	handler := corsMiddleware(mux)

	fmt.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal("Erro ao inicializar o servidor: ", err)
	}
}

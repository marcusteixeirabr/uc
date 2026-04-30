package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8180"
	}

	fmt.Printf("Servidor rodando na porta %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

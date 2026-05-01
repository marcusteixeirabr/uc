package main

import (
	"log"
	"net/http"
	"os"

	"github.com/marcusteixeirabr/uc/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8180"
	}

	http.HandleFunc("/ucs", handlers.UCList)
	http.HandleFunc("/comunicacoes", handlers.ComunicacaoList)
	http.HandleFunc("/comunicacoes/nova", handlers.ComunicacaoForm)

	log.Printf("Servidor rodando em http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

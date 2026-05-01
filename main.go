package main

import (
	"fmt"
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

	fmt.Printf("Servidor rodando na porta %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	http.HandleFunc("/ucs", handlers.UCList)
	http.HandleFunc("/comunicacoes", handlers.ComunicacaoList)
	http.HandleFunc("/comunicacoes/nova", handlers.ComunicacaoForm)

}

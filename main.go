package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/marcusteixeirabr/uc/internal/db"
	"github.com/marcusteixeirabr/uc/internal/handlers"
	"github.com/marcusteixeirabr/uc/internal/middleware"
)

func main() {
	godotenv.Overload() // carrega .env sobrescrevendo o ambiente; garante .env como fonte de verdade no dev

	client, err := db.Connect()
	if err != nil {
		log.Fatalf("falha ao conectar ao Supabase: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8180"
	}

	h := &handlers.Handler{DB: client}

	mux := http.NewServeMux()
	mux.HandleFunc("/ucs", h.UCList)
	mux.HandleFunc("/comunicacoes", h.ComunicacaoList)
	mux.HandleFunc("/comunicacoes/nova", h.ComunicacaoForm)

	log.Printf("Servidor rodando em http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, middleware.Logger(mux)))
}

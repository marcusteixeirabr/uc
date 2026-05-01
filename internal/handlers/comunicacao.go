package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type Comunicacao struct {
	Titulo   string
	Email    string
	DataHora string
	Status   int
	NomeUC   string
}

func ComunicacaoForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/comunicacao_form.html",
	)
	if err != nil {
		log.Printf("erro ao carregar template: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base", nil)
}

func ComunicacaoList(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/comunicacao_list.html",
	)
	if err != nil {
		log.Printf("erro ao carregar template: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
	comunicacoes := []Comunicacao{
		{Titulo: "Lixo na trilha", Email: "user1@email.com", DataHora: "2026-04-20 10:30", Status: 0, NomeUC: "Parque do Rio Vermelho"},
		{Titulo: "Placa quebrada", Email: "user2@email.com", DataHora: "2026-04-21 14:00", Status: 1, NomeUC: "Parque Raimundo Malta"},
	}
	tmpl.ExecuteTemplate(w, "base", comunicacoes)
}

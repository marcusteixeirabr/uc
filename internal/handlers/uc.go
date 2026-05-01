package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type UC struct {
	Nome        string
	DataCriacao string
	Instituicao string
	Descricao   string
}

func UCList(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/uc_list.html",
	)
	if err != nil {
		log.Printf("erro ao carregar template: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
	ucs := []UC{
		{Nome: "Parque do Rio Vermelho", DataCriacao: "2007-03-24", Instituicao: "IMA SC", Descricao: "O Parque do Rio Vermelho é uma unidade de conservação localizada em Florianópolis, Santa Catarina. Criado em 2007, o parque tem como objetivo proteger a biodiversidade local e oferecer um espaço de lazer para a população."},
		{Nome: "Parque Estadual da Serra do Tabuleiro", DataCriacao: "1975-06-10", Instituicao: "IMA SC", Descricao: "O Parque Estadual da Serra do Tabuleiro é uma unidade de conservação situada em Santa Catarina. Criado em 1975, o parque é conhecido por sua rica biodiversidade e por abrigar importantes nascentes de rios."},
		{Nome: "Parque Nacional de São Joaquim", DataCriacao: "1961-09-24", Instituicao: "ICMBio", Descricao: "O Parque Nacional de São Joaquim é uma unidade de conservação localizada em Santa Catarina. Criado em 1961, o parque é famoso por suas paisagens montanhosas, cachoeiras e pela presença do Morro da Igreja, o ponto mais alto do estado."},
	}
	tmpl.ExecuteTemplate(w, "base", ucs)
}

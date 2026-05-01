package handlers

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	postgrest "github.com/supabase-community/postgrest-go"
	supabase "github.com/supabase-community/supabase-go"
)

type Handler struct {
	DB        *supabase.Client
	Templates embed.FS
}

type UC struct {
	Nome        string
	DataCriacao string
	Descricao   string
	Instituicao string
}

func (h *Handler) UCList(w http.ResponseWriter, r *http.Request) {
	// PostgREST: "instituicao(nome)" faz o JOIN via foreign key
	type ucRow struct {
		Nome        string `json:"nome"`
		DataCriacao string `json:"data_criacao"`
		Descricao   string `json:"descricao"`
		Instituicao struct {
			Nome string `json:"nome"`
		} `json:"instituicao"`
	}

	var rows []ucRow
	_, err := h.DB.From("unidade_conservacao").
		Select("nome, data_criacao, descricao, instituicao(nome)", "", false).
		Order("nome", &postgrest.OrderOpts{Ascending: true}).
		ExecuteTo(&rows)
	if err != nil {
		log.Printf("erro na query: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	ucs := make([]UC, len(rows))
	for i, row := range rows {
		ucs[i] = UC{
			Nome:        row.Nome,
			DataCriacao: row.DataCriacao,
			Descricao:   row.Descricao,
			Instituicao: row.Instituicao.Nome,
		}
	}

	tmpl, err := template.ParseFS(h.Templates, "templates/base.html", "templates/uc_list.html")
	if err != nil {
		log.Printf("erro no template: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base", ucs)

}

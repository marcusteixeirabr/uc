package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	postgrest "github.com/supabase-community/postgrest-go"
)

type OpcaoUC struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

type Comunicacao struct {
	Titulo   string
	Email    string
	DataHora string
	Status   int
	NomeUC   string
}

func (h *Handler) ComunicacaoForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.comunicacaoCreate(w, r)
		return
	}

	var opcoes []OpcaoUC
	_, err := h.DB.From("unidade_conservacao").
		Select("id, nome", "", false).
		Order("nome", &postgrest.OrderOpts{Ascending: true}).
		ExecuteTo(&opcoes)
	if err != nil {
		log.Printf("erro na query: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFS(h.Templates, "templates/base.html", "templates/comunicacao_form.html")
	if err != nil {
		log.Printf("erro no template: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base", opcoes)
}

func (h *Handler) comunicacaoCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	titulo := r.FormValue("titulo")
	descricao := r.FormValue("descricao")
	email := r.FormValue("email")
	unidadeID, err := strconv.Atoi(r.FormValue("unidade_id"))
	if err != nil || titulo == "" || descricao == "" || email == "" {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	type insertRow struct {
		Titulo    string `json:"titulo"`
		Descricao string `json:"descricao"`
		DataHora  string `json:"data_hora"`
		Email     string `json:"email"`
		Status    int    `json:"status"`
		UnidadeID int    `json:"unidade_id"`
	}

	_, _, err = h.DB.From("comunicacao").
		Insert(insertRow{
			Titulo:    titulo,
			Descricao: descricao,
			DataHora:  time.Now().Format("2006-01-02T15:04:05"),
			Email:     email,
			Status:    0,
			UnidadeID: unidadeID,
		}, false, "", "minimal", "").
		Execute()
	if err != nil {
		log.Printf("erro ao inserir: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	// dispara notificação em paralelo — não bloqueia o redirect
	// argumentos passados por valor: garante cópias próprias para a goroutine
	go func(titulo, email string) {
		log.Printf("[notif] nova comunicação '%s' registrada por %s", titulo, email)
		// aqui entraria: enviarEmail(email, titulo)
	}(titulo, email)

	http.Redirect(w, r, "/comunicacoes", http.StatusSeeOther)
}

func (h *Handler) ComunicacaoList(w http.ResponseWriter, r *http.Request) {
	type comunicacaoRow struct {
		Titulo             string `json:"titulo"`
		Email              string `json:"email"`
		DataHora           string `json:"data_hora"`
		Status             int    `json:"status"`
		UnidadeConservacao struct {
			Nome string `json:"nome"`
		} `json:"unidade_conservacao"`
	}

	var rows []comunicacaoRow
	_, err := h.DB.From("comunicacao").
		Select("titulo, email, data_hora, status, unidade_conservacao(nome)", "", false).
		Order("data_hora", &postgrest.OrderOpts{Ascending: false}).
		ExecuteTo(&rows)
	if err != nil {
		log.Printf("erro na query: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	comunicacoes := make([]Comunicacao, len(rows))
	for i, row := range rows {
		dataHora := row.DataHora
		if t, err := time.Parse("2006-01-02T15:04:05", dataHora); err == nil {
			dataHora = t.Format("02/01/2006 15:04")
		}
		comunicacoes[i] = Comunicacao{
			Titulo:   row.Titulo,
			Email:    row.Email,
			DataHora: dataHora,
			Status:   row.Status,
			NomeUC:   row.UnidadeConservacao.Nome,
		}
	}

	tmpl, err := template.ParseFS(h.Templates, "templates/base.html", "templates/comunicacao_list.html")
	if err != nil {
		log.Printf("erro no template: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base", comunicacoes)
}

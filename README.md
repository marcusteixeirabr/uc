# HOW VI — Sistema de Unidades de Conservação

Trabalho prático da disciplina **Oficinas de Integração VI** — ADS UNIVALI.

Aplicação web desenvolvida em **Go puro** (sem framework) com banco de dados **Supabase (PostgreSQL na nuvem)**, rodando em ambiente **DevContainer** (Docker + VS Code).

---

## Funcionalidades

| Rota | Descrição |
|------|-----------|
| `GET /ucs` | Lista todas as Unidades de Conservação com a instituição responsável |
| `GET /comunicacoes` | Lista todas as comunicações registradas com status (Pendente/Atendida) |
| `GET /comunicacoes/nova` | Formulário para registrar nova comunicação |
| `POST /comunicacoes/nova` | Grava a comunicação no banco e redireciona para a lista |

---

## Stack

```
┌────────────────────────────────────────────────┐
│  Browser (cliente)                              │
│     ↕  HTTP (GET / POST)                       │
├────────────────────────────────────────────────┤
│  Go net/http  →  middleware  →  handlers        │
│  html/template embutido via //go:embed          │
├────────────────────────────────────────────────┤
│  supabase-go  →  PostgREST API (HTTPS/443)     │
├────────────────────────────────────────────────┤
│  Supabase (PostgreSQL na nuvem)                │
└────────────────────────────────────────────────┘
```

**Por que Go sem framework?** Em Spring Boot (`@GetMapping`, `@Controller`) o framework faz roteamento, injeção de dependências e serialização de forma invisível. Em Go, escrevemos tudo explicitamente — você entende exatamente o que cada linha faz e consegue adaptar ou construir um framework quando necessário.

**Por que Supabase em vez de PostgreSQL local?** O `supabase-go` usa a API **PostgREST** via HTTPS — sem abrir porta TCP 5432, sem problemas de firewall em redes corporativas, sem container de banco no `docker-compose.yml`.

**Por que DevContainer?** Garante que todos os colaboradores têm exatamente as mesmas ferramentas (`air`, `goose`, `sqlc`, `go 1.26`) independente do sistema operacional.

---

## Estrutura do Projeto

```
uc/
├── .devcontainer/
│   ├── Dockerfile              # Go 1.26 Alpine + air + goose + sqlc + openssh
│   ├── devcontainer.json       # VS Code no container, porta 8180
│   └── setup-ssh.sh            # copia chaves SSH do host → container
├── .air.toml                   # hot reload: monitora .go e .html
├── .env                        # credenciais (NÃO vai pro Git)
├── .env.example                # template para novos colaboradores
├── .gitignore                  # .env, tmp/, *.exe
├── docker-compose.yml          # só o container app (banco no Supabase)
├── go.mod / go.sum             # módulo e checksums de dependências
├── main.go                     # embed templates, ServeMux, middleware, servidor
├── db/
│   └── migrations/
│       ├── 00001_schema_inicial.sql    # DDL completo (up/down)
│       └── 00002_dados_iniciais.sql    # dados de exemplo (up/down)
├── internal/
│   ├── db/db.go                # Connect() → cliente Supabase
│   ├── handlers/
│   │   ├── uc.go               # Handler struct, UCList
│   │   └── comunicacao.go      # ComunicacaoForm/List/Create + goroutine
│   └── middleware/
│       └── logger.go           # Logger: método, rota, status, duração
└── templates/                  # embutidos no binário via //go:embed
    ├── base.html
    ├── uc_list.html
    ├── comunicacao_form.html
    └── comunicacao_list.html
```

---

## Diagrama de Classes (structs Go)

```mermaid
classDiagram
    class Handler {
        +DB *supabase.Client
        +Templates embed.FS
        +UCList(w ResponseWriter, r Request)
        +ComunicacaoList(w ResponseWriter, r Request)
        +ComunicacaoForm(w ResponseWriter, r Request)
        -comunicacaoCreate(w ResponseWriter, r Request)
    }

    class UC {
        +Nome string
        +DataCriacao string
        +Descricao string
        +Instituicao string
    }

    class Comunicacao {
        +Titulo string
        +Email string
        +DataHora string
        +Status int
        +NomeUC string
    }

    class OpcaoUC {
        +ID int
        +Nome string
    }

    Handler ..> UC : produz []UC para o template
    Handler ..> Comunicacao : produz []Comunicacao para o template
    Handler ..> OpcaoUC : produz []OpcaoUC para o select
```

> Em Go não há classes nem herança — apenas **structs** e **composição**. Métodos são definidos com *receivers*: `func (h *Handler) UCList(...)`.

---

## DER — Diagrama Entidade-Relacionamento

```mermaid
erDiagram
    INSTITUICAO {
        serial id PK
        varchar nome
        varchar email
    }
    UNIDADE_CONSERVACAO {
        serial id PK
        varchar nome
        date data_criacao
        text descricao
        varchar imagem_url
        int instituicao_id FK
    }
    MUNICIPIO {
        serial id PK
        varchar nome
        char estado
    }
    UNIDADE_MUNICIPIO {
        int unidade_id FK
        int municipio_id FK
    }
    COMUNICACAO {
        serial id PK
        varchar titulo
        text descricao
        timestamp data_hora
        varchar email
        int status
        int unidade_id FK
    }

    INSTITUICAO ||--o{ UNIDADE_CONSERVACAO : "gerencia"
    UNIDADE_CONSERVACAO ||--o{ UNIDADE_MUNICIPIO : "abrange"
    MUNICIPIO       ||--o{ UNIDADE_MUNICIPIO : "pertence a"
    UNIDADE_CONSERVACAO ||--o{ COMUNICACAO : "recebe"
```

> `UNIDADE_MUNICIPIO` é uma **tabela associativa** (N:N). Uma UC pode abranger vários municípios; um município pode ter várias UCs.

---

## DDL — Criação das Tabelas (PostgreSQL / Supabase)

> O schema está versionado em `db/migrations/00001_schema_inicial.sql`. Para criar via goose, veja a seção [Migrations](#migrations). Para criar manualmente, execute no **SQL Editor do Supabase** na ordem abaixo.

```sql
CREATE TABLE instituicao (
  id    SERIAL PRIMARY KEY,
  nome  VARCHAR(150) NOT NULL,
  email VARCHAR(150) NOT NULL UNIQUE
);

CREATE TABLE unidade_conservacao (
  id             SERIAL PRIMARY KEY,
  nome           VARCHAR(150) NOT NULL,
  data_criacao   DATE         NOT NULL,
  descricao      TEXT,
  imagem_url     VARCHAR(255),
  instituicao_id INT NOT NULL REFERENCES instituicao(id)
);

CREATE TABLE municipio (
  id     SERIAL PRIMARY KEY,
  nome   VARCHAR(150) NOT NULL,
  estado CHAR(2)      NOT NULL
);

CREATE TABLE unidade_municipio (
  unidade_id   INT NOT NULL REFERENCES unidade_conservacao(id),
  municipio_id INT NOT NULL REFERENCES municipio(id),
  PRIMARY KEY (unidade_id, municipio_id)
);

CREATE TABLE comunicacao (
  id         SERIAL PRIMARY KEY,
  titulo     VARCHAR(150) NOT NULL,
  descricao  TEXT         NOT NULL,
  data_hora  TIMESTAMP    NOT NULL,
  email      VARCHAR(150) NOT NULL,
  status     INT          NOT NULL DEFAULT 0,
  unidade_id INT          NOT NULL REFERENCES unidade_conservacao(id)
);
```

> O arquivo `material/UnidadeConservacao.sql` usa sintaxe **MySQL** (`AUTO_INCREMENT`, `DATETIME`). O DDL acima é adaptado para **PostgreSQL** (`SERIAL`, `TIMESTAMP`).

---

## DML — Dados de Exemplo

> Também versionado em `db/migrations/00002_dados_iniciais.sql`.

```sql
INSERT INTO instituicao (nome, email) VALUES
('ICMBio', 'icmbio@org.br'), ('IMA SC', 'ima@sc.gov.br');

INSERT INTO municipio (nome, estado) VALUES
('Florianopolis','SC'),('Balneario Camboriu','SC'),('Itajai','SC'),('Bombinhas','SC');

INSERT INTO unidade_conservacao (nome, data_criacao, descricao, imagem_url, instituicao_id) VALUES
('Parque do Rio Vermelho','2007-03-24','Area de preservacao ambiental','img1.jpg',2),
('Parque Raimundo Malta', '1993-07-15','Parque com trilhas e vegetacao nativa','img2.jpg',2),
('Reserva do Arvoredo',   '1990-04-12','Reserva marinha protegida','img3.jpg',1),
('Parque da Ressacada',   '2008-09-10','Area voltada para educacao ambiental','img4.jpg',2);

INSERT INTO unidade_municipio (unidade_id, municipio_id) VALUES
(1,1),(2,2),(3,1),(3,4),(4,3);

INSERT INTO comunicacao (titulo, descricao, data_hora, email, status, unidade_id) VALUES
('Lixo na trilha',  'Tem lixo acumulado em alguns pontos',  '2026-04-20 10:30:00','user1@email.com',0,1),
('Placa quebrada',  'Uma placa informativa esta danificada','2026-04-21 14:00:00','user2@email.com',1,2),
('Pesca irregular', 'Possivel atividade ilegal na area',   '2026-04-22 09:15:00','user3@email.com',0,3);
```

> `status = 0` → Pendente | `status = 1` → Atendida

---

## Como Rodar o Projeto

### Pré-requisitos

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [VS Code](https://code.visualstudio.com/) com extensão **Dev Containers**
- Conta no [Supabase](https://supabase.com) com o banco criado

### 1. Configure o `.env`

Copie `.env.example` para `.env` e preencha:

```bash
# Supabase → Settings → API → Project URL
DATABASE_URL=https://xxxxxxxxxxxx.supabase.co

# Supabase → Settings → API → anon/public key
DATABASE_KEY=eyJhbGci...

# Supabase → Settings → Database → Connection string → URI  (para goose)
GOOSE_DBSTRING=postgresql://postgres.xxxx:SENHA@host.pooler.supabase.com:6543/postgres

PORT=8180
```

> `DATABASE_URL` deve ser **apenas o domínio**, sem `/rest/v1`.

### 2. Abra no DevContainer

`Ctrl+Shift+P` → **Dev Containers: Reopen in Container** → aguarde ~2 min na primeira vez.

### 3. Crie o schema com goose (primeira vez)

```bash
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="$(grep GOOSE_DBSTRING .env | cut -d= -f2-)"
goose -dir db/migrations up
```

### 4. Execute o servidor com hot reload

```bash
air
```

Acesse: `http://localhost:8180/ucs`

> Salve qualquer `.go` ou `.html` e o servidor reinicia automaticamente.

---

## Distribuição para Windows (sem instalar nada)

Para enviar o executável para colegas rodarem no Windows:

```bash
# dentro do DevContainer, na pasta uc/
GOOS=windows GOARCH=amd64 go build -o uc.exe .
```

O colega recebe `uc.exe` (~13MB, templates embutidos) + cria um `.env` com as credenciais. Executa:

```powershell
cd pasta-do-exe
.\uc.exe
# acessa http://localhost:8180/ucs
```

Nenhuma instalação necessária — Go compila binários estáticos e autossuficientes.

---

## Migrations com goose

```bash
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="$(grep GOOSE_DBSTRING .env | cut -d= -f2-)"

goose -dir db/migrations status   # ver o que está aplicado
goose -dir db/migrations up       # aplicar migrations pendentes
goose -dir db/migrations down     # reverter a última migration
```

Para evoluir o schema, crie um novo arquivo — nunca edite os existentes:

```
db/migrations/00003_minha_mudanca.sql
```

```sql
-- +goose Up
ALTER TABLE comunicacao ADD COLUMN telefone VARCHAR(20);

-- +goose Down
ALTER TABLE comunicacao DROP COLUMN telefone;
```

---

## Conceitos Implementados

| # | Conceito | Onde no código |
|---|----------|----------------|
| 1 | DevContainer + Docker Compose sem banco local | `.devcontainer/`, `docker-compose.yml` |
| 2 | Servidor HTTP sem framework, `http.ServeMux` | `main.go` |
| 3 | Templates HTML com herança (`base.html`) | `templates/` |
| 4 | Dependency Injection via struct | `handlers.Handler{DB, Templates}` |
| 5 | PostgREST: JOIN via `"tabela(coluna)"` | `uc.go`, `comunicacao.go` |
| 6 | Padrão PRG (Post-Redirect-Get) | `comunicacaoCreate` |
| 7 | Hot reload com `air` | `.air.toml` |
| 8 | Middleware como function wrapping | `middleware/logger.go` |
| 9 | Struct embedding para estender interface | `responseWriter` em `logger.go` |
| 10 | Templates embutidos no binário (`//go:embed`) | `main.go` |
| 11 | Goroutine fire-and-forget | `comunicacaoCreate` |
| 12 | Tratamento de erros idiomático (`%w`) | `db.go`, todos os handlers |
| 13 | Migrations versionadas com goose | `db/migrations/` |
| 14 | Cross-compile para Windows | `GOOS=windows go build` |

---

## Armadilhas Resolvidas

| Problema | Causa | Solução |
|----------|-------|---------|
| `git push` falha no container | `openssh-client` ausente na imagem Alpine | Adicionar no `apk add` do Dockerfile |
| SSH rejeita chaves no container | Proprietário das chaves (uid 1000) ≠ root (uid 0) | `setup-ssh.sh` copia e corrige permissões |
| `pgx` não conecta ao Supabase | Supabase bloqueia TCP 5432 externamente | Usar `supabase-go` (HTTPS/PostgREST) |
| Erro PGRST125 "Invalid path" | `DATABASE_URL` com `/rest/v1` — a lib duplica o sufixo | `DATABASE_URL` = só o domínio |
| `.env` ignorado após container subir | `godotenv.Load()` não sobrescreve vars já injetadas | Usar `godotenv.Overload()` |
| Datas chegam como string | PostgREST retorna JSON, não protocolo binário | Parsear com `time.Parse` |
| `.exe` sem templates quebra | `template.ParseFiles` lê do disco em runtime | Usar `//go:embed` + `template.ParseFS` |

---

## Dependências

| Pacote | Versão | Função |
|--------|--------|--------|
| `github.com/joho/godotenv` | v1.5.1 | Carrega o `.env` |
| `github.com/supabase-community/supabase-go` | v0.0.4 | Cliente Supabase (PostgREST) |
| `github.com/supabase-community/postgrest-go` | v0.0.11 | Query builder PostgREST |

Ferramentas no container (não são dependências do módulo Go):

| Ferramenta | Versão | Função |
|-----------|--------|--------|
| `air` | v1.65.1 | Hot reload |
| `goose` | v3.27.1 | Migrations de banco |
| `sqlc` | v1.31.1 | Geração de código a partir de SQL |

---

## Próximos Passos

### Para quem quer rodar e entender o projeto agora
- Seguir as seções "Como Rodar" e "Migrations" acima
- Ler os comentários no código — cada decisão não-óbvia está explicada

### Para quem quer continuar o aprendizado Go

**Fundamentos da linguagem** (base para tudo abaixo)
- Tipos, zero values, slices, maps
- Interfaces — o padrão que você já viu em `http.Handler`
- Ponteiros — quando usar `*T` vs `T`

**Próximas funcionalidades no projeto**
- [ ] **Testes** — `testing` + `httptest`: testar handlers sem subir servidor
- [ ] **sqlc** — gerar structs e funções Go tipadas a partir das queries SQL; elimina os mapeamentos manuais
- [ ] **Repositório pattern** — desacoplar handlers do banco com interfaces; facilita testes e troca de banco
- [ ] **Validação com feedback visual** — erros de formulário mostrados ao usuário

**Infraestrutura e deploy**
- [ ] **Docker multi-stage build** — imagem de produção ~10MB (sem toolchain Go)
- [ ] **GitHub Actions** — CI/CD: compilar e testar automaticamente a cada push; gerar binários para Windows/Linux/Mac
- [ ] **Deploy** — Fly.io ou Railway (Go tem deploy direto de binário, sem JVM)
- [ ] **Autenticação** — Supabase Auth + validação de JWT no middleware

**Horizonte (Web3 + infra)**
- [ ] Go tem o ecossistema mais forte para Web3 — `go-ethereum` é Go puro
- [ ] gRPC, Protocol Buffers — comunicação entre serviços
- [ ] Kubernetes operators — escritos em Go (controller-runtime)

---

## Material do Projeto

```
../material/
├── DiagramaDER.png
├── DiagramaDeClasses.png
├── SQL_DDL.png
├── SQL_DML.png
├── UnidadeConservacao.sql   # DDL original (sintaxe MySQL — ver adaptação acima)
├── InsereDados.sql
└── Enunciado do Trabalho.pdf
```

-- +goose Up
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

-- +goose Down
DROP TABLE IF EXISTS comunicacao;
DROP TABLE IF EXISTS unidade_municipio;
DROP TABLE IF EXISTS municipio;
DROP TABLE IF EXISTS unidade_conservacao;
DROP TABLE IF EXISTS instituicao;

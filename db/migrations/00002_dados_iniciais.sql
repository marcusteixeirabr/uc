-- +goose Up
INSERT INTO instituicao (nome, email) VALUES
('ICMBio', 'icmbio@org.br'),
('IMA SC', 'ima@sc.gov.br');

INSERT INTO municipio (nome, estado) VALUES
('Florianopolis', 'SC'),
('Balneario Camboriu', 'SC'),
('Itajai', 'SC'),
('Bombinhas', 'SC');

INSERT INTO unidade_conservacao (nome, data_criacao, descricao, imagem_url, instituicao_id) VALUES
('Parque do Rio Vermelho', '2007-03-24', 'Area de preservacao ambiental',       'img1.jpg', 2),
('Parque Raimundo Malta',  '1993-07-15', 'Parque com trilhas e vegetacao nativa','img2.jpg', 2),
('Reserva do Arvoredo',    '1990-04-12', 'Reserva marinha protegida',            'img3.jpg', 1),
('Parque da Ressacada',    '2008-09-10', 'Area voltada para educacao ambiental', 'img4.jpg', 2);

INSERT INTO unidade_municipio (unidade_id, municipio_id) VALUES
(1, 1), (2, 2), (3, 1), (3, 4), (4, 3);

INSERT INTO comunicacao (titulo, descricao, data_hora, email, status, unidade_id) VALUES
('Lixo na trilha',  'Tem lixo acumulado em alguns pontos',  '2026-04-20 10:30:00', 'user1@email.com', 0, 1),
('Placa quebrada',  'Uma placa informativa esta danificada', '2026-04-21 14:00:00', 'user2@email.com', 1, 2),
('Pesca irregular', 'Possivel atividade ilegal na area',    '2026-04-22 09:15:00', 'user3@email.com', 0, 3);

-- +goose Down
DELETE FROM comunicacao;
DELETE FROM unidade_municipio;
DELETE FROM unidade_conservacao;
DELETE FROM municipio;
DELETE FROM instituicao;

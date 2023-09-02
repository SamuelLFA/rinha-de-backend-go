-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pessoas (
    id UUID NOT NULL PRIMARY KEY,
    apelido VARCHAR(32) NOT NULL UNIQUE,
    nome VARCHAR(100) NOT NULL,
    nascimento VARCHAR(10) NOT NULL,
    stack VARCHAR(32)[]
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pessoas
-- +goose StatementEnd

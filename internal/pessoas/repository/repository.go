package repository

import (
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/SamuelLFA/rinha-de-backend-go/internal/pessoas/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PessoasRepository interface {
	ListPessoas(searchTerm string) ([]model.Pessoa, error)
	CountPessoas() (int64, error)
	CreatePessoa(pessoa model.Pessoa) (*model.Pessoa, error)
	GetPessoa(id uuid.UUID) (model.Pessoa, error)
}

type pessoasRepository struct {
	db *sql.DB
}

func New(db *sql.DB) PessoasRepository {
	return &pessoasRepository{
		db: db,
	}
}

func (h *pessoasRepository) ListPessoas(searchTerm string) ([]model.Pessoa, error) {
	pessoas := []model.Pessoa{}

	query := sq.Select("id", "apelido", "nome", "nascimento", "stack").From("pessoas")
	query = query.Where(
		sq.Or{
			sq.Like{"apelido": "%" + searchTerm + "%"},
			sq.Like{"nome": "%" + searchTerm + "%"},
			sq.Like{"nascimento": "%" + searchTerm + "%"},
			sq.Expr("stack IN (?)", "{"+searchTerm+"}"),
		},
	).PlaceholderFormat(sq.Dollar)
	rows, err := query.RunWith(h.db).Query()

	if err != nil {
		return pessoas, err
	}
	defer rows.Close()

	next := rows.Next()
	for next {
		var pessoa model.Pessoa
		rows.Scan(&pessoa.ID, &pessoa.Apelido, &pessoa.Nome, &pessoa.Nascimento, pq.Array(&pessoa.Stack))
		pessoas = append(pessoas, pessoa)
		next = rows.Next()
	}
	return pessoas, nil
}

func (h *pessoasRepository) GetPessoa(id uuid.UUID) (model.Pessoa, error) {
	pessoa := model.Pessoa{}

	query := sq.Select("id", "apelido", "nome", "nascimento", "stack").From("pessoas")
	query = query.Where(
		sq.Eq{"id": id},
	).PlaceholderFormat(sq.Dollar)
	rows, err := query.RunWith(h.db).Query()

	if err != nil {
		return pessoa, err
	}
	defer rows.Close()

	next := rows.Next()
	if !next {
		return pessoa, fmt.Errorf("pessoa not found")
	}
	rows.Scan(&pessoa.ID, &pessoa.Apelido, &pessoa.Nome, &pessoa.Nascimento, pq.Array(&pessoa.Stack))
	return pessoa, nil
}

func (h *pessoasRepository) CountPessoas() (int64, error) {
	query := sq.Select("COUNT(*)").From("pessoas")
	rows, err := query.RunWith(h.db).Query()

	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var count int64
	rows.Next()
	rows.Scan(&count)
	return count, nil
}

func (h *pessoasRepository) CreatePessoa(pessoa model.Pessoa) (*model.Pessoa, error) {
	pessoa.ID = uuid.New()

	query := sq.Insert("pessoas").Columns("id", "apelido", "nome", "nascimento", "stack")
	query = query.Values(pessoa.ID, pessoa.Apelido, pessoa.Nome, pessoa.Nascimento, "{"+strings.Join(pessoa.Stack, ",")+"}").PlaceholderFormat(sq.Dollar)
	_, err := query.RunWith(h.db).Exec()

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "duplicate key value violates unique constraint") {
			return nil, fmt.Errorf("apelido %s já existe", pessoa.Apelido)
		}
		return nil, err
	}

	return &pessoa, nil
}

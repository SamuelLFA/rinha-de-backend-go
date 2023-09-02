package mocks

import (
	"github.com/SamuelLFA/rinha-de-backend-go/internal/pessoas/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockPessoasRepository is a mock implementation of PessoasRepository for testing.
type MockPessoasRepository struct {
	mock.Mock
}

// ListPessoas is a mocked implementation of the ListPessoas method.
func (m *MockPessoasRepository) ListPessoas(searchTerm string) ([]model.Pessoa, error) {
	args := m.Called(searchTerm)
	return args.Get(0).([]model.Pessoa), args.Error(1)
}

// GetPessoa is a mocked implementation of the GetPessoa method.
func (m *MockPessoasRepository) GetPessoa(id uuid.UUID) (model.Pessoa, error) {
	args := m.Called(id)
	return args.Get(0).(model.Pessoa), args.Error(1)
}

// CountPessoas is a mocked implementation of the CountPessoas method.
func (m *MockPessoasRepository) CountPessoas() (int64, error) {
	args := m.Called()
	return int64(args.Int(0)), args.Error(1)
}

// CreatePessoa is a mocked implementation of the CreatePessoa method.
func (m *MockPessoasRepository) CreatePessoa(pessoa model.Pessoa) (*model.Pessoa, error) {
	args := m.Called(pessoa)
	return args.Get(0).(*model.Pessoa), args.Error(1)
}

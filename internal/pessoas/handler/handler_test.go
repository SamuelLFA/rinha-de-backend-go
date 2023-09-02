package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SamuelLFA/rinha-de-backend-go/internal/pessoas/model"
	"github.com/SamuelLFA/rinha-de-backend-go/internal/pessoas/repository/mocks"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestListPessoas(t *testing.T) {
	// given
	repoMock := new(mocks.MockPessoasRepository)
	handler := New(repoMock)
	expectedPessoas := []model.Pessoa{
		{
			ID:         uuid.New(),
			Nome:       "Alice",
			Apelido:    "Ali",
			Nascimento: "1990-01-02",
			Stack:      nil,
		},
		{
			ID:         uuid.New(),
			Nome:       "Robert",
			Apelido:    "Rob",
			Nascimento: "1995-04-16",
			Stack:      []string{"Node", "Go"},
		},
	}
	repoMock.On("ListPessoas", "search").Return(expectedPessoas, nil)

	// when
	req := httptest.NewRequest(http.MethodGet, "/pessoas?t=search", nil)
	recorder := httptest.NewRecorder()
	handler.ListPessoas(recorder, req, nil)

	// then
	assert.Equal(t, http.StatusOK, recorder.Code)
	repoMock.AssertExpectations(t)
}

func TestGetPessoas(t *testing.T) {
	// given
	repoMock := new(mocks.MockPessoasRepository)
	handler := New(repoMock)
	id := uuid.New()
	expectedPessoa := model.Pessoa{
		ID:         id,
		Nome:       "Robert",
		Apelido:    "Rob",
		Nascimento: "1995-04-16",
		Stack:      []string{"Node", "Go"},
	}
	repoMock.On("GetPessoa", id).Return(expectedPessoa, nil)

	// when
	params := httprouter.Params{
		{
			Key:   "id",
			Value: id.String(),
		},
	}
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/pessoas/%s", id.String()), nil)
	recorder := httptest.NewRecorder()
	handler.GetPessoa(recorder, req, params)

	// then
	assert.Equal(t, http.StatusOK, recorder.Code)
	repoMock.AssertExpectations(t)
}

func TestCreatePessoa(t *testing.T) {
	// given
	repoMock := new(mocks.MockPessoasRepository)
	handler := New(repoMock)
	id := uuid.New()
	request := PessoaRequest{
		Nome:       "Robert",
		Apelido:    "Rob",
		Nascimento: "1995-04-16",
		Stack:      []string{"Node", "Go"},
	}
	pessoa := model.Pessoa{
		Nome:       request.Nome,
		Apelido:    request.Apelido,
		Nascimento: request.Nascimento,
		Stack:      request.Stack,
	}
	repoMock.On("CreatePessoa", pessoa).Return(
		&model.Pessoa{
			ID:         id,
			Nome:       pessoa.Nome,
			Apelido:    pessoa.Apelido,
			Nascimento: pessoa.Nascimento,
			Stack:      pessoa.Stack,
		}, nil)

	// when
	data, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/pessoas", bytes.NewReader(data))
	recorder := httptest.NewRecorder()
	handler.CreatePessoa(recorder, req, nil)

	// then
	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Equal(t, fmt.Sprintf("/pessoas/%s", id), recorder.Header().Get("Location"))
	repoMock.AssertExpectations(t)
}

func TestCountPessoas(t *testing.T) {
	// given
	repoMock := new(mocks.MockPessoasRepository)
	handler := New(repoMock)
	repoMock.On("CountPessoas").Return(10, nil)

	// when
	req := httptest.NewRequest(http.MethodGet, "/pessoas/contagem-pessoas", nil)
	recorder := httptest.NewRecorder()
	handler.CountPessoas(recorder, req, nil)

	// then
	assert.Equal(t, http.StatusOK, recorder.Code)
	repoMock.AssertExpectations(t)
}

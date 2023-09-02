package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/SamuelLFA/rinha-de-backend-go/internal/pessoas/model"
	"github.com/SamuelLFA/rinha-de-backend-go/internal/pessoas/repository"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

var regex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

type PessoasHandler interface {
	ListPessoas(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	CountPessoas(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	CreatePessoa(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	GetPessoa(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
}

type pessoasHandler struct {
	repo repository.PessoasRepository
}

func New(repo repository.PessoasRepository) PessoasHandler {
	return &pessoasHandler{
		repo: repo,
	}
}

func (h *pessoasHandler) ListPessoas(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	searchTerm := r.URL.Query().Get("t")
	if searchTerm == "" {
		http.Error(w, "t is required", http.StatusBadRequest)
		return
	}

	pessoas, err := h.repo.ListPessoas(searchTerm)
	if err != nil {
		log.Default().Printf("error listing pessoas: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pessoas); err != nil {
		log.Default().Printf("error encoding pessoas: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *pessoasHandler) GetPessoa(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	paramId := params.ByName("id")
	id, err := uuid.Parse(paramId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pessoa, err := h.repo.GetPessoa(id)
	if err != nil {
		if strings.Compare(fmt.Sprint(err), "pessoa not found") == 0 {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Default().Printf("error getting pessoa: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pessoa); err != nil {
		log.Default().Printf("error encoding pessoa: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *pessoasHandler) CountPessoas(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	pessoas, err := h.repo.CountPessoas()
	if err != nil {
		log.Default().Printf("error getting number of pessoas: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%d", pessoas)
}

func (h *pessoasHandler) CreatePessoa(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var pessoa PessoaRequest

	if err := json.NewDecoder(r.Body).Decode(&pessoa); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(pessoa.Apelido) <= 0 || len(pessoa.Apelido) > 32 {
		http.Error(w, "Invalid \"apelido\": required, min = 0, max = 32", http.StatusUnprocessableEntity)
		return
	}

	if len(pessoa.Nome) <= 0 || len(pessoa.Nome) > 100 {
		http.Error(w, "Invalid \"nome\": required, min = 0, max = 100", http.StatusUnprocessableEntity)
		return
	}

	if !regex.MatchString(pessoa.Nascimento) {
		http.Error(w, "Invalid \"nascimento\": required, format YYYY-MM-DD", http.StatusUnprocessableEntity)
		return
	}

	for _, stack := range pessoa.Stack {
		if len(stack) <= 0 || len(stack) > 32 {
			http.Error(w, "Invalid \"stack\": required, min = 0, max = 32", http.StatusUnprocessableEntity)
			return
		}
	}

	model := model.Pessoa{
		Apelido:    pessoa.Apelido,
		Nome:       pessoa.Nome,
		Nascimento: pessoa.Nascimento,
		Stack:      pessoa.Stack,
	}

	newPessoa, err := h.repo.CreatePessoa(model)
	if err != nil {
		if strings.Compare(fmt.Sprint(err), fmt.Sprintf("apelido %s já existe", model.Apelido)) == 0 {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		log.Default().Printf("error creating new pessoa :%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/pessoas/%s", newPessoa.ID))
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newPessoa); err != nil {
		log.Default().Printf("error encoding response :%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

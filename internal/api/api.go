package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/SamuelLFA/rinha-de-backend-go/internal/pessoas/handler"
	"github.com/SamuelLFA/rinha-de-backend-go/internal/pessoas/repository"
	httprouter "github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

type HealthCheckResponse struct {
	Status string `json:"status"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthCheckResponse{Status: "OK"})
}

func Init(db *sql.DB) {
	pessoasRepository := repository.New(db)
	pessoasHandler := handler.New(pessoasRepository)

	router := httprouter.New()
	router.GET("/", HealthCheck)
	router.GET("/pessoas", pessoasHandler.ListPessoas)
	router.GET("/pessoas/:id", pessoasHandler.GetPessoa)
	router.POST("/pessoas", pessoasHandler.CreatePessoa)
	router.GET("/contagem-pessoas", pessoasHandler.CountPessoas)

	handler := cors.Default().Handler(router)
	http.ListenAndServe(":80", handler)
}

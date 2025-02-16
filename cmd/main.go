package main

import (
	"log"
	"net/http"

	"github.com/EBayego/scrapad-backend/internal/repository"
	"github.com/EBayego/scrapad-backend/internal/rest"
	"github.com/EBayego/scrapad-backend/internal/service"
	"github.com/gorilla/mux"
)

func main() {
	// Inicializa repositorio (SQLite)
	db, err := repository.NewSQLiteConnection("C:/Users/Edu/Desktop/Projects/Backend-Test-Medior/scrapad_backend/database/scrapad_assigment.db")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Crea instancias de repositorio
	repo := repository.NewSQLiteRepository(db)

	// Crea los servicios de negocio
	financeService := service.NewFinanceService(repo)
	offerService := service.NewOfferService(repo, financeService)

	// Crea el router y registra handlers
	r := mux.NewRouter()
	rest.RegisterHandlers(r, offerService)

	// Inicia servidor
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

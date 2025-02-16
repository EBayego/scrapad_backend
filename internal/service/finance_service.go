package service

import (
	"fmt"

	"github.com/EBayego/scrapad-backend/internal/repository"
)

type FinanceService interface {
	RequestFinancing(partnerSlug string, total int) (int, error)
}

type financeService struct {
	repo *repository.SQLiteRepository
}

func NewFinanceService(repo *repository.SQLiteRepository) FinanceService {
	return &financeService{repo: repo}
}

// RequestFinancing simula la llamada a un partner financiero
func (f *financeService) RequestFinancing(partnerSlug string, total int) (int, error) {
	// Obtenemos el provider para saber el porcentaje
	fp, err := f.repo.GetFinancingProviderBySlug(partnerSlug)
	if err != nil {
		return 0, err
	}
	// Calculamos la cantidad neta a percibir
	fee := total * fp.FinancingPercentage / 100
	net := total - fee

	// Simulación: “llamada exitosa al partner”
	// Realmente podríamos hacer una request HTTP, etc.
	fmt.Printf("Requesting financing from %s. Original: %d, fee: %d, net: %d\n",
		partnerSlug, total, fee, net)

	return net, nil
}

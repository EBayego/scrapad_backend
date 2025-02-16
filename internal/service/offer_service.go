package service

import (
	"errors"
	"time"

	"github.com/EBayego/scrapad-backend/internal/domain"
	"github.com/EBayego/scrapad-backend/internal/repository"
	"github.com/google/uuid"
)

const (
	FinancingBankSlug    = "financing_bank"
	FinancingFintechSlug = "financing_fintech"
)

type OfferService interface {
	CreateOffer(req domain.CreateOfferRequest) (*domain.Offer, error)
	GetPendingOffersByOrg(orgID string) ([]domain.Offer, error)
	RequestFinancing(offerID string, req domain.FinancingRequest) (int, error)
	AcceptOffer(offerID string, req domain.AcceptOfferRequest) error
}

type offerService struct {
	repo           *repository.SQLiteRepository
	financeService FinanceService
}

func NewOfferService(r *repository.SQLiteRepository, f FinanceService) OfferService {
	return &offerService{
		repo:           r,
		financeService: f,
	}
}

// CreateOffer implementa las reglas de negocio
func (s *offerService) CreateOffer(req domain.CreateOfferRequest) (*domain.Offer, error) {
	// Validar Ad
	ad, err := s.repo.GetAdByID(req.Ad)
	if err != nil {
		return nil, errors.New("ad not found")
	}

	// Obtener Org
	org, err := s.repo.GetOrganizationByID(ad.OrgID)
	if err != nil {
		return nil, errors.New("organization not found")
	}

	// Revisar reglas de financing
	financingProviderID, err := s.decideFinancingProvider(*org)
	if err != nil {
		// no hay financiamiento
		financingProviderID = 0
	}

	// Crear la Offer
	offer := domain.Offer{
		ID:                uuid.New().String(),
		PaymentMethod:     req.PaymentMethod,
		FinancingProvider: financingProviderID,
		Amount:            req.Amount,
		Accepted:          0,
		Price:             req.Price,
		AdId:              req.Ad,
	}

	createdOffer, err := s.repo.CreateOffer(offer)
	if err != nil {
		return nil, err
	}

	return createdOffer, nil
}

func (s *offerService) decideFinancingProvider(org domain.Organization) (int, error) {
	sumAds, err := s.repo.GetSumAdsPublishedByOrg(org.ID)
	if err != nil {
		return 0, err
	}

	oneYearAgo := time.Now().AddDate(-1, 0, 0).UTC()
	if sumAds > 10000 && oneYearAgo.Before(org.CreatedDate) { // si cumple las condiciones necesarias para todos los casos
		if org.Country == "SPAIN" || org.Country == "FRANCE" {
			fp, err := s.repo.GetFinancingProviderBySlug(FinancingBankSlug)
			if err != nil {
				return 0, err
			}
			return fp.ID, nil
		} else {
			fp, err := s.repo.GetFinancingProviderBySlug(FinancingFintechSlug)
			if err != nil {
				return 0, err
			}
			return fp.ID, nil
		}
	}

	// si no cumple
	return 0, errors.New("no financing applicable")
}

// Ofertas pendientes
func (s *offerService) GetPendingOffersByOrg(orgID string) ([]domain.Offer, error) {
	offers, err := s.repo.GetOffersByOrgID(orgID)
	if err != nil {
		return nil, err
	}

	// Filtrar por accepted=0
	var pending []domain.Offer
	for _, offer := range offers {
		if offer.Accepted == 0 {
			pending = append(pending, offer)
		}
	}
	return pending, nil
}

// Simula la llamada al partner seleccionado en la offer
func (s *offerService) RequestFinancing(offerID string, req domain.FinancingRequest) (int, error) {
	offer, err := s.repo.GetOfferByID(offerID)
	if err != nil {
		return 0, errors.New("offer not found")
	}

	// Validar que la financingProvider coincida con la slug
	partner, err := s.repo.GetFinancingProviderBySlug(req.FinancingPartner)
	if err != nil {
		return 0, errors.New("financing partner not found")
	}
	if offer.FinancingProvider != partner.ID {
		return 0, errors.New("offer is not set for this financing partner")
	}

	net, err := s.financeService.RequestFinancing(partner.Slug, req.TotalToPerceive)
	if err != nil {
		return 0, err
	}

	return net, nil // net => la cantidad final que percibe el seller
}

// Marca la oferta como aceptada, y si se pasa un financingPartner, se cambia de tipo de financiacion
func (s *offerService) AcceptOffer(offerID string, req domain.AcceptOfferRequest) error {
	offer, err := s.repo.GetOfferByID(offerID)
	if err != nil {
		return errors.New("offer not found")
	}

	if req.FinancingPartner != "" {
		fp, err := s.repo.GetFinancingProviderBySlug(req.FinancingPartner)
		if err != nil {
			return errors.New("financing provider not found")
		}
		offer.FinancingProvider = fp.ID
	}

	offer.Accepted = 1
	return s.repo.UpdateOffer(*offer)
}

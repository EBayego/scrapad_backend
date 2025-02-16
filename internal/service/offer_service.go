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
	// 1. Validar Ad
	ad, err := s.repo.GetAdByID(req.Ad)
	if err != nil {
		return nil, errors.New("ad not found")
	}

	// 2. Obtener Org
	org, err := s.repo.GetOrganizationByID(ad.OrgID)
	if err != nil {
		return nil, errors.New("organization not found")
	}

	// 3. Revisar reglas de financing
	financingProviderID, err := s.decideFinancingProvider(*org)
	if err != nil {
		// no hay financiamiento
		financingProviderID = 0
	}

	// 4. Crear Offer
	offer := domain.Offer{
		ID:                uuid.New().String(),
		PaymentMethod:     req.PaymentMethod,
		FinancingProvider: financingProviderID,
		Amount:            req.Amount,
		Accepted:          0, // 0 => false
		Price:             req.Price,
		AdId:              req.Ad,
	}

	createdOffer, err := s.repo.CreateOffer(offer)
	if err != nil {
		return nil, err
	}

	return createdOffer, nil
}

// Aplica las reglas definidas para decidir si va con "financing_by_bank" o "financing_by_fintech"
// Devuelve el ID del provider en la base, o un error si no aplica
func (s *offerService) decideFinancingProvider(org domain.Organization) (int, error) {
	// Reglas:
	//   - org_country in ('SPAIN','FRANCE')
	//     and sum_ads_published > 10000
	//     and org_created_date < (now() - 1 año)
	//     => financing_bank
	//
	//   - org_country not in ('SPAIN','FRANCE')
	//     and sum_ads_published > 10000
	//     and org_created_date < (now() - 1 año)
	//     => financing_fintech

	// sumAds se asume la multiplicación de amount*price de la tabla ads
	sumAds, err := s.repo.GetSumAdsPublishedByOrg(org.ID)
	if err != nil {
		return 0, err
	}

	oneYearAgo := time.Now().AddDate(-1, 0, 0).UTC()
	if sumAds > 10000 && oneYearAgo.Before(org.CreatedDate) {
		if org.Country == "SPAIN" || org.Country == "FRANCE" {
			// financing_by_bank
			fp, err := s.repo.GetFinancingProviderBySlug(FinancingBankSlug)
			if err != nil {
				return 0, err
			}
			return fp.ID, nil
		} else {
			// financing_by_fintech
			fp, err := s.repo.GetFinancingProviderBySlug(FinancingFintechSlug)
			if err != nil {
				return 0, err
			}
			return fp.ID, nil
		}
	}

	// Si no cumple, error
	return 0, errors.New("no financing applicable")
}

// GetPendingOffersByOrg retorna ofertas “pendientes” (no aceptadas)
func (s *offerService) GetPendingOffersByOrg(orgID string) ([]domain.Offer, error) {
	offers, err := s.repo.GetOffersByOrgID(orgID)
	if err != nil {
		return nil, err
	}

	// Filtramos por accepted=0
	var pending []domain.Offer
	for _, offer := range offers {
		if offer.Accepted == 0 {
			pending = append(pending, offer)
		}
	}
	return pending, nil
}

// RequestFinancing simula la llamada al partner seleccionado en la offer
func (s *offerService) RequestFinancing(offerID string, req domain.FinancingRequest) (int, error) {
	// 1. Obtenemos la oferta
	offer, err := s.repo.GetOfferByID(offerID)
	if err != nil {
		return 0, errors.New("offer not found")
	}

	// 2. Validamos que la financingProvider coincida con la slug
	partner, err := s.repo.GetFinancingProviderBySlug(req.FinancingPartner)
	if err != nil {
		return 0, errors.New("financing partner not found")
	}
	if offer.FinancingProvider != partner.ID {
		return 0, errors.New("offer is not set for this financing partner")
	}

	// 3. Llamamos al FinanceService
	net, err := s.financeService.RequestFinancing(partner.Slug, req.TotalToPerceive)
	if err != nil {
		return 0, err
	}
	// net => la cantidad final que percibe el seller

	return net, nil
}

// AcceptOffer marca la oferta como aceptada.
// Opcionalmente, si se pasa un financingPartner, se podría forzar un partner
func (s *offerService) AcceptOffer(offerID string, req domain.AcceptOfferRequest) error {
	offer, err := s.repo.GetOfferByID(offerID)
	if err != nil {
		return errors.New("offer not found")
	}

	// Si se pasa un slug nuevo, forzamos
	if req.FinancingPartner != "" {
		fp, err := s.repo.GetFinancingProviderBySlug(req.FinancingPartner)
		if err != nil {
			return errors.New("financing provider not found")
		}
		offer.FinancingProvider = fp.ID
	}

	// Aceptamos
	offer.Accepted = 1
	return s.repo.UpdateOffer(*offer)
}

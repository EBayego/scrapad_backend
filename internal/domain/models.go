package domain

import (
	"time"
)

// Organization representa la tabla organizations
type Organization struct {
	ID          string    `json:"id"`
	Country     string    `json:"country"`
	CreatedDate time.Time `json:"created_date"`
}

// Ad representa la tabla ads
type Ad struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"`
	Price  int    `json:"price"`
	OrgID  string `json:"org_id"`
}

// FinancingProvider representa la tabla financing_providers
type FinancingProvider struct {
	ID                  int    `json:"id"`
	Slug                string `json:"slug"`
	PaymentMethod       string `json:"payment_method"`
	FinancingPercentage int    `json:"financing_percentage"`
}

// Offer representa la tabla offers
type Offer struct {
	ID                string `json:"id"`
	PaymentMethod     string `json:"payment_method"`
	FinancingProvider int    `json:"financing_privder"` // referencia a FinancingProvider.ID
	Amount            int    `json:"amount"`
	Accepted          int    `json:"accepted"` // 0 o 1
	Price             int    `json:"price"`
	// otras propiedades según sea necesario
}

// Solicitudes (DTOs) para los endpoints
type CreateOfferRequest struct {
	Ad            string `json:"ad"`             // uuid
	Amount        int    `json:"amount"`         // en centavos
	Price         int    `json:"price"`          // en centavos
	PaymentMethod string `json:"payment_method"` // ejemplo: "100_in_unload"
}

type FinancingRequest struct {
	FinancingPartner string `json:"financingPartner"` // ej: "financing_by_bank"
	TotalToPerceive  int    `json:"totalToPerceive"`  // valor en centavos
}

type AcceptOfferRequest struct {
	FinancingPartner string `json:"financingPartner,omitempty"` // opcional
}

package domain

import (
	"time"
)

type Organization struct {
	ID          string    `json:"id"`
	Country     string    `json:"country"`
	CreatedDate time.Time `json:"created_date"`
}

type Ad struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"`
	Price  int    `json:"price"`
	OrgID  string `json:"org_id"`
}

type FinancingProvider struct {
	ID                  int    `json:"id"`
	Slug                string `json:"slug"`
	PaymentMethod       string `json:"payment_method"`
	FinancingPercentage int    `json:"financing_percentage"`
}

type Offer struct {
	ID                string `json:"id"`
	PaymentMethod     string `json:"payment_method"`
	FinancingProvider int    `json:"financing_privder"`
	Amount            int    `json:"amount"`
	Accepted          int    `json:"accepted"`
	Price             int    `json:"price"`
	AdId              string `json:"ad_id"`
}

// Solicitudes o DTOs para los endpoints
type CreateOfferRequest struct {
	Ad            string `json:"ad"`
	Amount        int    `json:"amount"`
	Price         int    `json:"price"`
	PaymentMethod string `json:"payment_method"`
}

type FinancingRequest struct {
	FinancingPartner string `json:"financingPartner"`
	TotalToPerceive  int    `json:"totalToPerceive"`
}

type AcceptOfferRequest struct {
	FinancingPartner string `json:"financingPartner,omitempty"` // opcional
}

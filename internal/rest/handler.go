package rest

import (
	"encoding/json"
	"net/http"

	"github.com/EBayego/scrapad-backend/internal/domain"
	"github.com/EBayego/scrapad-backend/internal/service"
	"github.com/gorilla/mux"
)

type offerHandler struct {
	offerService service.OfferService
}

func RegisterHandlers(r *mux.Router, offerSvc service.OfferService) {
	h := &offerHandler{
		offerService: offerSvc,
	}

	r.HandleFunc("/offers", h.CreateOffer).Methods("POST")
	r.HandleFunc("/orgs/{orgID}/offers/pending", h.GetPendingOffers).Methods("GET")
	r.HandleFunc("/offers/{offerID}/financing", h.RequestFinancing).Methods("POST")
	r.HandleFunc("/offers/{offerID}/accept", h.AcceptOffer).Methods("POST")
}

func (h *offerHandler) CreateOffer(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateOfferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offer, err := h.offerService.CreateOffer(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(offer)
}

func (h *offerHandler) GetPendingOffers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgID := vars["orgID"]

	offers, err := h.offerService.GetPendingOffersByOrg(orgID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(offers)
}

func (h *offerHandler) RequestFinancing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	offerID := vars["offerID"]

	var req domain.FinancingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	net, err := h.offerService.RequestFinancing(offerID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{
		"net_amount": net,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *offerHandler) AcceptOffer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	offerID := vars["offerID"]

	var req domain.AcceptOfferRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Evita campos desconocidos en el JSON
	if err := decoder.Decode(&req); err != nil && err.Error() != "EOF" {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := h.offerService.AcceptOffer(offerID, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"offer accepted"}`))
}

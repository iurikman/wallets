package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/iurikman/wallets/internal/models"
	log "github.com/sirupsen/logrus"
)

type service interface {
	CreateWallet(context context.Context) (*models.Wallet, error)
	GetWallet(ctx context.Context, id uuid.UUID) (*models.Wallet, error)
	Deposit(ctx context.Context, transaction models.Transaction) error
	Withdraw(ctx context.Context, transaction models.Transaction) error
}

type HTTPResponse struct {
	Data  any    `json:"data"`
	Error string `json:"error"`
}

func (s *Server) createWallet(w http.ResponseWriter, r *http.Request) {
	createdWallet, err := s.service.CreateWallet(r.Context())
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "internal server error")
		log.Warnf("failed to create new wallet: %v", err)

		return
	}

	writeOkResponse(w, http.StatusCreated, createdWallet)
}

func (s *Server) getWallet(w http.ResponseWriter, r *http.Request) {
	walletID := chi.URLParam(r, "id")
	if walletID == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Missing 'id' parameter")

		return
	}

	walletIDParsed, err := uuid.Parse(walletID)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())

		return
	}

	wallet, err := s.service.GetWallet(r.Context(), walletIDParsed)

	switch {
	case errors.Is(err, models.ErrWalletNotFound):
		writeErrorResponse(w, http.StatusNotFound, err.Error())

		return
	case err != nil:
		writeErrorResponse(w, http.StatusInternalServerError, "internal server error")

		return
	}

	writeOkResponse(w, http.StatusOK, wallet)
}

func (s *Server) deposit(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())

		return
	}

	if err := transaction.Validate(); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())

		return
	}

	err := s.service.Deposit(r.Context(), transaction)

	switch {
	case errors.Is(err, models.ErrWalletNotFound):
		writeErrorResponse(w, http.StatusNotFound, err.Error())

		return
	case err != nil:
		writeErrorResponse(w, http.StatusInternalServerError, "internal server error")
		log.Warnf("failed to deposit transaction: %v", err)

		return
	}

	writeOkResponse(w, http.StatusOK, nil)
}

func (s *Server) withdraw(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())

		return
	}

	if err := transaction.Validate(); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())

		return
	}

	err := s.service.Withdraw(r.Context(), transaction)

	switch {
	case errors.Is(err, models.ErrWalletNotFound):
		writeErrorResponse(w, http.StatusNotFound, err.Error())

		return
	case errors.Is(err, models.ErrBalanceBelowZero):
		writeErrorResponse(w, http.StatusBadRequest, err.Error())

		return
	case err != nil:
		writeErrorResponse(w, http.StatusInternalServerError, "internal server error")
		log.Warnf("failed to withdraw transaction: %v", err)

		return
	}

	writeOkResponse(w, http.StatusOK, nil)
}

func writeOkResponse(w http.ResponseWriter, statusCode int, respData any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(HTTPResponse{Data: respData}); err != nil {
		log.Warnf("json.NewEncoder(w).Encode(HTTPResponse{Data: respData}) err: %v", err)
	}
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, description string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(HTTPResponse{Error: description}); err != nil {
		log.Warnf("json.NewEncoder(w).Encode(HTTPResponse{Error: description}) err: %s", err)
	}
}

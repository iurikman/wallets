package tests

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/iurikman/wallets/internal/models"
	"github.com/iurikman/wallets/internal/rest"
)

func (s *IntegrationTestSuite) TestWallets() {
	testWallet1 := models.Wallet{}
	testWalletID := uuid.New()

	s.Run("POST", func() {
		s.Run("201/statusCreated", func() {
			createdWallet := new(models.Wallet)

			resp := s.sendRequest(
				context.Background(),
				http.MethodPost,
				"/",
				testWallet1,
				&rest.HTTPResponse{Data: &createdWallet},
			)
			s.Require().Equal(http.StatusCreated, resp.StatusCode)
			s.Require().Equal(testWallet1.Balance, createdWallet.Balance)
			testWalletID = createdWallet.ID
		})
	})

	s.Run("GET", func() {
		s.Run("200/statusOK", func() {
			wallet := new(models.Wallet)

			resp := s.sendRequest(
				context.Background(),
				http.MethodGet,
				"/"+testWalletID.String(),
				nil,
				&rest.HTTPResponse{Data: &wallet},
			)
			s.Require().Equal(http.StatusOK, resp.StatusCode)
			s.Require().Equal(testWalletID, wallet.ID)
			s.Require().Equal(0.0, wallet.Balance)
			s.Require().Equal(false, wallet.Deleted)
		})
	})

	s.Run("PUT", func() {
		s.Run("/deposit", func() {
			s.Run("200/statusOK", func() {
				executedTransaction := new(models.Transaction)

				testDepositOperation := models.Transaction{
					TransactionID: uuid.New(),
					WalletID:      testWalletID,
					Amount:        500,
					OperationType: "DEPOSIT",
				}

				resp := s.sendRequest(
					context.Background(),
					http.MethodPut,
					"/deposit",
					testDepositOperation,
					&rest.HTTPResponse{Data: &executedTransaction},
				)
				s.Require().Equal(http.StatusOK, resp.StatusCode)
			})

			s.Run("404/StatusNotFound(random wallet id)", func() {
				testDepositOperation := models.Transaction{
					TransactionID: uuid.New(),
					WalletID:      uuid.New(),
					Amount:        500,
					OperationType: "DEPOSIT",
				}

				resp := s.sendRequest(
					context.Background(),
					http.MethodPut,
					"/deposit",
					testDepositOperation,
					nil,
				)
				s.Require().Equal(http.StatusNotFound, resp.StatusCode)
			})

			s.Run("400/StatusBadRequest(bad request)", func() {
				resp := s.sendRequest(
					context.Background(),
					http.MethodPut,
					"/deposit",
					"bad request",
					nil,
				)
				s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
			})
		})
		s.Run("/withdraw", func() {
			s.Run("200/statusOK", func() {
				executedTransaction := new(models.Transaction)

				testDepositOperation := models.Transaction{
					TransactionID: uuid.New(),
					WalletID:      testWalletID,
					Amount:        250,
					OperationType: "WITHDRAW",
				}

				resp := s.sendRequest(
					context.Background(),
					http.MethodPut,
					"/withdraw",
					testDepositOperation,
					&rest.HTTPResponse{Data: &executedTransaction},
				)
				s.Require().Equal(http.StatusOK, resp.StatusCode)
			})

			s.Run("404/StatusBadRequest(operation type not allowed", func() {
				executedTransaction := new(models.Transaction)

				testDepositOperation := models.Transaction{
					TransactionID: uuid.New(),
					WalletID:      testWalletID,
					Amount:        250,
					OperationType: "bad operation",
				}

				resp := s.sendRequest(
					context.Background(),
					http.MethodPut,
					"/withdraw",
					testDepositOperation,
					&rest.HTTPResponse{Data: &executedTransaction},
				)
				s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
			})

			s.Run("404/StatusNotFound(random wallet id)", func() {
				testDepositOperation := models.Transaction{
					TransactionID: uuid.New(),
					WalletID:      uuid.New(),
					Amount:        500,
					OperationType: "DEPOSIT",
				}

				resp := s.sendRequest(
					context.Background(),
					http.MethodPut,
					"/withdraw",
					testDepositOperation,
					nil,
				)
				s.Require().Equal(http.StatusNotFound, resp.StatusCode)
			})

			s.Run("400/StatusBadRequest(balance below zero)", func() {
				testDepositOperation := models.Transaction{
					TransactionID: uuid.New(),
					WalletID:      testWalletID,
					Amount:        5000,
					OperationType: "DEPOSIT",
				}

				resp := s.sendRequest(
					context.Background(),
					http.MethodPut,
					"/withdraw",
					testDepositOperation,
					nil,
				)
				s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
			})

			s.Run("400/StatusBadRequest(bad request)", func() {
				resp := s.sendRequest(
					context.Background(),
					http.MethodPut,
					"/withdraw",
					"bad request",
					nil,
				)
				s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
			})
		})
	})
}

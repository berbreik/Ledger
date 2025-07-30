package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"ledger/internal/domain"
	"ledger/internal/handler"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// mock service
type mockTransactionService struct {
	ProcessFunc func(ctx context.Context, tx *domain.Transaction) error
	HistoryFunc func(ctx context.Context, accountID int64) ([]*domain.Transaction, error)
}

func (m *mockTransactionService) ProcessTransaction(ctx context.Context, tx *domain.Transaction) error {
	return m.ProcessFunc(ctx, tx)
}

func (m *mockTransactionService) GetTransactionHistory(ctx context.Context, accountID int64) ([]*domain.Transaction, error) {
	return m.HistoryFunc(ctx, accountID)
}

func TestProcessTransaction_Success(t *testing.T) {
	mockService := &mockTransactionService{
		ProcessFunc: func(ctx context.Context, tx *domain.Transaction) error {
			tx.ID = "txn123"
			return nil
		},
	}

	h := handler.NewTransactionHandler(mockService)

	tx := domain.Transaction{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        100.0,
		Currency:      "USD",
	}

	body, _ := json.Marshal(tx)
	req := httptest.NewRequest(http.MethodPost, "/transaction", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.ProcessTransaction(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	var result map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&result)

	if result["status"] != "success" {
		t.Errorf("expected success, got %s", result["status"])
	}

	if result["transaction_id"] != "txn123" {
		t.Errorf("expected txn123, got %s", result["transaction_id"])
	}
}

func TestProcessTransaction_BadRequest(t *testing.T) {
	mockService := &mockTransactionService{}
	h := handler.NewTransactionHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/transaction", bytes.NewBuffer([]byte(`bad json`)))
	w := httptest.NewRecorder()

	h.ProcessTransaction(w, req)
	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Result().StatusCode)
	}
}

func TestProcessTransaction_MissingFields(t *testing.T) {
	mockService := &mockTransactionService{}
	h := handler.NewTransactionHandler(mockService)

	tx := domain.Transaction{} // empty
	body, _ := json.Marshal(tx)
	req := httptest.NewRequest(http.MethodPost, "/transaction", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.ProcessTransaction(w, req)
	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Result().StatusCode)
	}
}

func TestProcessTransaction_ServiceError(t *testing.T) {
	mockService := &mockTransactionService{
		ProcessFunc: func(ctx context.Context, tx *domain.Transaction) error {
			return errors.New("internal error")
		},
	}

	h := handler.NewTransactionHandler(mockService)

	tx := domain.Transaction{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        100.0,
		Currency:      "USD",
	}

	body, _ := json.Marshal(tx)
	req := httptest.NewRequest(http.MethodPost, "/transaction", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.ProcessTransaction(w, req)
	if w.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Result().StatusCode)
	}
}

func TestGetTransactionHistory_Success(t *testing.T) {
	mockService := &mockTransactionService{
		HistoryFunc: func(ctx context.Context, accountID int64) ([]*domain.Transaction, error) {
			return []*domain.Transaction{
				{
					ID:            "txn1",
					FromAccountID: accountID,
					ToAccountID:   2,
					Amount:        100.0,
					Currency:      "USD",
					CreatedAt:     time.Now().String(),
				},
			}, nil
		},
	}

	h := handler.NewTransactionHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/transactions?account_id=1", nil)
	w := httptest.NewRecorder()

	h.GetTransactionHistory(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Result().StatusCode)
	}

	var txs []*domain.Transaction
	_ = json.NewDecoder(w.Body).Decode(&txs)

	if len(txs) != 1 {
		t.Errorf("expected 1 transaction, got %d", len(txs))
	}

	if txs[0].FromAccountID != 1 {
		t.Errorf("expected accountID 1, got %d", txs[0].FromAccountID)
	}
}

func TestGetTransactionHistory_InvalidAccountID(t *testing.T) {
	h := handler.NewTransactionHandler(&mockTransactionService{})

	req := httptest.NewRequest(http.MethodGet, "/transactions?account_id=abc", nil)
	w := httptest.NewRecorder()

	h.GetTransactionHistory(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Result().StatusCode)
	}
}

func TestGetTransactionHistory_MissingParam(t *testing.T) {
	h := handler.NewTransactionHandler(&mockTransactionService{})

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	w := httptest.NewRecorder()

	h.GetTransactionHistory(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Result().StatusCode)
	}
}

func TestGetTransactionHistory_ErrorFromService(t *testing.T) {
	h := handler.NewTransactionHandler(&mockTransactionService{
		HistoryFunc: func(ctx context.Context, accountID int64) ([]*domain.Transaction, error) {
			return nil, errors.New("db error")
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/transactions?account_id=1", nil)
	w := httptest.NewRecorder()

	h.GetTransactionHistory(w, req)

	if w.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Result().StatusCode)
	}
}

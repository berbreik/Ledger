package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"ledger/internal/domain"
	"ledger/internal/handler"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockAccountService struct {
	CreateAccountFn        func(ctx context.Context, owner string, balance float64) error
	GetAccountFn           func(ctx context.Context, id string) (*domain.Account, error)
	UpdateAccountBalanceFn func(ctx context.Context, id string, balance float64) error
	GetAllAccountsFn       func(ctx context.Context) ([]*domain.Account, error)
	DeleteAccountFn        func(ctx context.Context, id string) error
}

func (m *mockAccountService) CreateAccount(ctx context.Context, owner string, balance float64) error {
	return m.CreateAccountFn(ctx, owner, balance)
}
func (m *mockAccountService) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	return m.GetAccountFn(ctx, id)
}
func (m *mockAccountService) UpdateAccountBalance(ctx context.Context, id string, balance float64) error {
	return m.UpdateAccountBalanceFn(ctx, id, balance)
}
func (m *mockAccountService) GetAllAccounts(ctx context.Context) ([]*domain.Account, error) {
	return m.GetAllAccountsFn(ctx)
}
func (m *mockAccountService) DeleteAccount(ctx context.Context, id string) error {
	return m.DeleteAccountFn(ctx, id)
}

func TestCreateAccount_Success(t *testing.T) {
	h := handler.NewAccountHandler(&mockAccountService{
		CreateAccountFn: func(ctx context.Context, owner string, balance float64) error {
			return nil
		},
	})

	body := `{"owner_name": "Alice", "initial_balance": 100}`
	r := httptest.NewRequest("POST", "/accounts", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.CreateAccount(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Account created successfully")) {
		t.Error("expected success message")
	}
}

func TestCreateAccount_InvalidInput(t *testing.T) {
	h := handler.NewAccountHandler(nil)

	r := httptest.NewRequest("POST", "/accounts", bytes.NewBufferString(`invalid-json`))
	w := httptest.NewRecorder()

	h.CreateAccount(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetAccount_Success(t *testing.T) {
	h := handler.NewAccountHandler(&mockAccountService{
		GetAccountFn: func(ctx context.Context, id string) (*domain.Account, error) {
			return &domain.Account{ID: "123", OwnerName: "Alice", Balance: 100}, nil
		},
	})
	r := httptest.NewRequest("GET", "/accounts?id=123", nil)
	w := httptest.NewRecorder()

	h.GetAccount(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp domain.Account
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp.ID != "123" {
		t.Errorf("expected account ID 123, got %s", resp.ID)
	}
}

func TestDeleteAccount_Success(t *testing.T) {
	h := handler.NewAccountHandler(&mockAccountService{
		DeleteAccountFn: func(ctx context.Context, id string) error {
			return nil
		},
	})

	r := httptest.NewRequest("DELETE", "/accounts?id=123", nil)
	w := httptest.NewRecorder()

	h.DeleteAccount(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestUpdateBalance_InvalidInput(t *testing.T) {
	h := handler.NewAccountHandler(nil)

	r := httptest.NewRequest("PUT", "/accounts/balance", bytes.NewBufferString(`bad-json`))
	w := httptest.NewRecorder()

	h.UpdateBalance(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetAllAccounts(t *testing.T) {
	h := handler.NewAccountHandler(&mockAccountService{
		GetAllAccountsFn: func(ctx context.Context) ([]*domain.Account, error) {
			return []*domain.Account{{ID: "1", OwnerName: "Alice", Balance: 100}}, nil
		},
	})

	r := httptest.NewRequest("GET", "/accounts", nil)
	w := httptest.NewRecorder()

	h.GetAllAccounts(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp []domain.Account
	body, _ := io.ReadAll(w.Body)
	_ = json.Unmarshal(body, &resp)
	if len(resp) != 1 || resp[0].OwnerName != "Alice" {
		t.Error("unexpected account list")
	}
}

func TestUpdateBalance_Success(t *testing.T) {
	h := handler.NewAccountHandler(&mockAccountService{
		UpdateAccountBalanceFn: func(ctx context.Context, id string, balance float64) error {
			if id != "123" || balance != 200 {
				return errors.New("invalid input")
			}
			return nil
		},
	})

	body := `{"id": "123", "balance": 200}`
	r := httptest.NewRequest("PUT", "/accounts/balance", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.UpdateBalance(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Balance updated successfully")) {
		t.Error("expected success message")
	}
}

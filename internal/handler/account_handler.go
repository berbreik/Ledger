package handler

import (
	"encoding/json"
	"net/http"
	_ "strconv"

	"ledger/internal/domain"
)

// AccountHandler handles HTTP requests related to accounts.
type AccountHandler struct {
	AccountService domain.AccountService
}

// NewAccountHandler creates a new AccountHandler instance.
func NewAccountHandler(service domain.AccountService) *AccountHandler {
	return &AccountHandler{
		AccountService: service,
	}
}

// CreateAccount handles POST /accounts
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {

	var req struct {
		OwnerName      string  `json:"owner_name"`
		InitialBalance float64 `json:"initial_balance"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if req.OwnerName == "" || req.InitialBalance < 0 {
		http.Error(w, "Invalid account data", http.StatusBadRequest)
		return
	}

	err := h.AccountService.CreateAccount(r.Context(), req.OwnerName, req.InitialBalance)
	if err != nil {
		http.Error(w, "Account creation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Account created successfully"))
}

// GetAccount handles GET /accounts/{id}
func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	accountID := r.URL.Query().Get("id")
	if accountID == "" {
		http.Error(w, "Missing account ID", http.StatusBadRequest)
		return
	}

	account, err := h.AccountService.GetAccount(r.Context(), accountID)
	if err != nil {
		http.Error(w, "Account not found: "+err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, account)
}

// UpdateBalance handles PUT /accounts/balance
func (h *AccountHandler) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID      string  `json:"id"`
		Balance float64 `json:"balance"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if req.ID == "" || req.Balance < 0 {
		http.Error(w, "Invalid account data", http.StatusBadRequest)
		return
	}

	err := h.AccountService.UpdateAccountBalance(r.Context(), req.ID, req.Balance)
	if err != nil {
		http.Error(w, "Failed to update balance: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Balance updated successfully"))
}

// GetAllAccounts handles GET /accounts
func (h *AccountHandler) GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.AccountService.GetAllAccounts(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch accounts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, accounts)
}

// DeleteAccount handles DELETE /accounts/{id}
func (h *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {

	accountID := r.URL.Query().Get("id")
	if accountID == "" {
		http.Error(w, "Missing account ID", http.StatusBadRequest)
		return
	}

	err := h.AccountService.DeleteAccount(r.Context(), accountID)
	if err != nil {
		http.Error(w, "Failed to delete account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// Helper: write JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

package handler

import (
	"encoding/json"
	"ledger/internal/domain"
	"net/http"
	"strconv"
)

type TransactionHandler struct {
	TransactionService domain.TransactionService
}

// NewTransactionHandler creates a new TransactionHandler instance.
func NewTransactionHandler(service domain.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		TransactionService: service,
	}
}

func (t *TransactionHandler) ProcessTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var tx domain.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "Invalid JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if tx.FromAccountID == 0 || tx.ToAccountID == 0 || tx.Amount <= 0 || tx.Currency == "" {
		http.Error(w, "Missing required transaction fields", http.StatusBadRequest)
		return
	}

	err := t.TransactionService.ProcessTransaction(ctx, &tx)
	if err != nil {
		http.Error(w, "Failed to process transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{
		"status":         "success",
		"transaction_id": tx.ID,
	})
	if err != nil {
		return
	}
}

func (t *TransactionHandler) GetTransactionHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	accountIDStr := r.URL.Query().Get("account_id")
	if accountIDStr == "" {
		http.Error(w, "Missing query param: account_id", http.StatusBadRequest)
		return
	}

	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid account_id", http.StatusBadRequest)
		return
	}

	txList, err := t.TransactionService.GetTransactionHistory(ctx, accountID)
	if err != nil {
		http.Error(w, "Error retrieving transactions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(txList)
	if err != nil {
		return
	}
}

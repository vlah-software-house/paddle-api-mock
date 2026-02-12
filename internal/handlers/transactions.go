package handlers

import (
	"net/http"
	"strings"

	"github.com/vlah-software-house/paddle-api-mock/internal/models"
	"github.com/vlah-software-house/paddle-api-mock/internal/store"
)

type TransactionsHandler struct {
	Store *store.Store
}

func (h *TransactionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/transactions")
	path = strings.TrimPrefix(path, "/")

	if r.Method != http.MethodGet {
		respondError(w, r, http.StatusMethodNotAllowed, "request_error", "method_not_allowed", "Method not allowed")
		return
	}

	if path == "" {
		h.list(w, r)
		return
	}

	h.get(w, r, path)
}

func (h *TransactionsHandler) list(w http.ResponseWriter, r *http.Request) {
	txns := h.Store.ListTransactions()

	// Filter by subscription_id
	if subID := r.URL.Query().Get("subscription_id"); subID != "" {
		filtered := make([]*models.Transaction, 0)
		for _, t := range txns {
			if t.SubscriptionID != nil && *t.SubscriptionID == subID {
				filtered = append(filtered, t)
			}
		}
		txns = filtered
	}

	// Filter by customer_id
	if cid := r.URL.Query().Get("customer_id"); cid != "" {
		filtered := make([]*models.Transaction, 0)
		for _, t := range txns {
			if t.CustomerID == cid {
				filtered = append(filtered, t)
			}
		}
		txns = filtered
	}

	respondList(w, r, txns, len(txns))
}

func (h *TransactionsHandler) get(w http.ResponseWriter, r *http.Request, id string) {
	txn, ok := h.Store.GetTransaction(id)
	if !ok {
		respondError(w, r, http.StatusNotFound, "request_error", "not_found", "Transaction not found")
		return
	}
	respond(w, r, http.StatusOK, txn)
}

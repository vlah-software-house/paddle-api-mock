package handlers

import (
	"net/http"
	"strings"

	"github.com/vlah-software-house/paddle-api-mock/internal/store"
)

type PricesHandler struct {
	Store *store.Store
}

func (h *PricesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/prices")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		if r.Method == http.MethodGet {
			h.list(w, r)
			return
		}
		respondError(w, r, http.StatusMethodNotAllowed, "request_error", "method_not_allowed", "Method not allowed")
		return
	}

	if r.Method == http.MethodGet {
		h.get(w, r, path)
		return
	}
	respondError(w, r, http.StatusMethodNotAllowed, "request_error", "method_not_allowed", "Method not allowed")
}

func (h *PricesHandler) list(w http.ResponseWriter, r *http.Request) {
	prices := h.Store.ListPrices()

	// Optionally include product in response
	if r.URL.Query().Get("include") == "product" {
		for _, p := range prices {
			if prod, ok := h.Store.GetProduct(p.ProductID); ok {
				p.Product = prod
			}
		}
	}

	respondList(w, r, prices, len(prices))
}

func (h *PricesHandler) get(w http.ResponseWriter, r *http.Request, id string) {
	price, ok := h.Store.GetPrice(id)
	if !ok {
		respondError(w, r, http.StatusNotFound, "request_error", "not_found", "Price not found")
		return
	}

	if r.URL.Query().Get("include") == "product" {
		if prod, ok := h.Store.GetProduct(price.ProductID); ok {
			price.Product = prod
		}
	}

	respond(w, r, http.StatusOK, price)
}

package handlers

import (
	"net/http"
	"strings"

	"github.com/vlah-software-house/paddle-api-mock/internal/store"
)

type ProductsHandler struct {
	Store *store.Store
}

func (h *ProductsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// /v1/products or /v1/products/{id}
	path := strings.TrimPrefix(r.URL.Path, "/v1/products")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		if r.Method == http.MethodGet {
			h.list(w, r)
			return
		}
		respondError(w, r, http.StatusMethodNotAllowed, "request_error", "method_not_allowed", "Method not allowed")
		return
	}

	// GET /v1/products/{id}
	if r.Method == http.MethodGet {
		h.get(w, r, path)
		return
	}
	respondError(w, r, http.StatusMethodNotAllowed, "request_error", "method_not_allowed", "Method not allowed")
}

func (h *ProductsHandler) list(w http.ResponseWriter, r *http.Request) {
	products := h.Store.ListProducts()
	respondList(w, r, products, len(products))
}

func (h *ProductsHandler) get(w http.ResponseWriter, r *http.Request, id string) {
	product, ok := h.Store.GetProduct(id)
	if !ok {
		respondError(w, r, http.StatusNotFound, "request_error", "not_found", "Product not found")
		return
	}
	respond(w, r, http.StatusOK, product)
}

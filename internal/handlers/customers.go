package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/vlah-software-house/paddle-api-mock/internal/models"
	"github.com/vlah-software-house/paddle-api-mock/internal/store"
)

type CustomersHandler struct {
	Store *store.Store
}

func (h *CustomersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/customers")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.list(w, r)
		case http.MethodPost:
			h.create(w, r)
		default:
			respondError(w, r, http.StatusMethodNotAllowed, "request_error", "method_not_allowed", "Method not allowed")
		}
		return
	}

	// /v1/customers/{id}
	id := path
	switch r.Method {
	case http.MethodGet:
		h.get(w, r, id)
	case http.MethodPatch:
		h.update(w, r, id)
	default:
		respondError(w, r, http.StatusMethodNotAllowed, "request_error", "method_not_allowed", "Method not allowed")
	}
}

func (h *CustomersHandler) list(w http.ResponseWriter, r *http.Request) {
	customers := h.Store.ListCustomers()
	respondList(w, r, customers, len(customers))
}

func (h *CustomersHandler) get(w http.ResponseWriter, r *http.Request, id string) {
	customer, ok := h.Store.GetCustomer(id)
	if !ok {
		respondError(w, r, http.StatusNotFound, "request_error", "not_found", "Customer not found")
		return
	}
	respond(w, r, http.StatusOK, customer)
}

func (h *CustomersHandler) create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCustomerRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, r, http.StatusBadRequest, "request_error", "bad_request", "Invalid JSON body")
		return
	}
	if req.Email == "" {
		respondError(w, r, http.StatusBadRequest, "request_error", "validation_error", "email is required")
		return
	}

	now := time.Now().UTC()
	locale := req.Locale
	if locale == "" {
		locale = "en"
	}
	customer := &models.Customer{
		ID:         store.NextID("ctm"),
		Name:       req.Name,
		Email:      req.Email,
		Locale:     locale,
		Status:     "active",
		CustomData: req.CustomData,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if customer.CustomData == nil {
		customer.CustomData = map[string]string{}
	}
	h.Store.SetCustomer(customer)
	respond(w, r, http.StatusCreated, customer)
}

func (h *CustomersHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	customer, ok := h.Store.GetCustomer(id)
	if !ok {
		respondError(w, r, http.StatusNotFound, "request_error", "not_found", "Customer not found")
		return
	}

	var req models.UpdateCustomerRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, r, http.StatusBadRequest, "request_error", "bad_request", "Invalid JSON body")
		return
	}

	if req.Name != nil {
		customer.Name = req.Name
	}
	if req.Email != "" {
		customer.Email = req.Email
	}
	if req.Locale != "" {
		customer.Locale = req.Locale
	}
	if req.Status != "" {
		customer.Status = req.Status
	}
	if req.CustomData != nil {
		customer.CustomData = req.CustomData
	}
	customer.UpdatedAt = time.Now().UTC()

	h.Store.SetCustomer(customer)
	respond(w, r, http.StatusOK, customer)
}

package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/vlah-software-house/paddle-api-mock/internal/models"
	"github.com/vlah-software-house/paddle-api-mock/internal/store"
)

type EventsHandler struct {
	Store *store.Store
}

func (h *EventsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, r, http.StatusMethodNotAllowed, "request_error", "method_not_allowed", "Method not allowed")
		return
	}
	events := h.Store.ListEvents()
	respondList(w, r, events, len(events))
}

type NotificationSettingsHandler struct {
	Store *store.Store
}

func (h *NotificationSettingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/notification-settings")
	path = strings.TrimPrefix(path, "/")

	switch r.Method {
	case http.MethodGet:
		h.list(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		respondError(w, r, http.StatusMethodNotAllowed, "request_error", "method_not_allowed", "Method not allowed")
	}
}

func (h *NotificationSettingsHandler) list(w http.ResponseWriter, r *http.Request) {
	settings := h.Store.ListNotificationSettings()
	respondList(w, r, settings, len(settings))
}

func (h *NotificationSettingsHandler) create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateNotificationSettingRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, r, http.StatusBadRequest, "request_error", "bad_request", "Invalid JSON body")
		return
	}
	if req.Destination == "" {
		respondError(w, r, http.StatusBadRequest, "request_error", "validation_error", "destination is required")
		return
	}

	now := time.Now().UTC()
	active := true
	if req.Active != nil {
		active = *req.Active
	}
	apiVersion := 1
	if req.APIVersion != nil {
		apiVersion = *req.APIVersion
	}
	nsType := "url"
	if req.Type != "" {
		nsType = req.Type
	}
	subscribedEvents := req.SubscribedEvents
	if subscribedEvents == nil {
		subscribedEvents = []string{
			"subscription.created",
			"subscription.updated",
			"subscription.activated",
			"subscription.canceled",
			"subscription.past_due",
			"transaction.completed",
			"transaction.payment_failed",
		}
	}

	ns := &models.NotificationSetting{
		ID:                     store.NextID("ntfset"),
		Description:            req.Description,
		Destination:            req.Destination,
		Active:                 active,
		APIVersion:             apiVersion,
		IncludeSensitiveFields: req.IncludeSensitiveFields != nil && *req.IncludeSensitiveFields,
		SubscribedEvents:       subscribedEvents,
		Type:                   nsType,
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	h.Store.SetNotificationSetting(ns)
	respond(w, r, http.StatusCreated, ns)
}

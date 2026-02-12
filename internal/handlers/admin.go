package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/vlah-software-house/paddle-api-mock/internal/models"
	"github.com/vlah-software-house/paddle-api-mock/internal/seed"
	"github.com/vlah-software-house/paddle-api-mock/internal/store"
	"github.com/vlah-software-house/paddle-api-mock/internal/webhook"
)

type AdminHandler struct {
	Store             *store.Store
	Webhook           *webhook.Notifier
	SeedEnabled       bool
	DefaultWebhookURL string
}

func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/admin/")

	switch {
	case path == "reset" && r.Method == http.MethodPost:
		h.reset(w, r)
	case strings.HasPrefix(path, "advance-time/") && r.Method == http.MethodPost:
		subID := strings.TrimPrefix(path, "advance-time/")
		h.advanceTime(w, r, subID)
	case strings.HasPrefix(path, "trigger-webhook/") && r.Method == http.MethodPost:
		eventType := strings.TrimPrefix(path, "trigger-webhook/")
		h.triggerWebhook(w, r, eventType)
	default:
		respondError(w, r, http.StatusNotFound, "request_error", "not_found", "Admin route not found")
	}
}

func (h *AdminHandler) reset(w http.ResponseWriter, r *http.Request) {
	h.Store.Reset()
	if h.SeedEnabled {
		seed.Load(h.Store)
	}
	// Re-register the default webhook URL so webhooks continue to fire after reset.
	if h.DefaultWebhookURL != "" {
		now := time.Now().UTC()
		h.Store.SetNotificationSetting(&models.NotificationSetting{
			ID:          store.NextID("ntfset"),
			Description: "Default webhook (from CLI)",
			Destination: h.DefaultWebhookURL,
			Active:      true,
			APIVersion:  1,
			SubscribedEvents: []string{
				"subscription.created",
				"subscription.updated",
				"subscription.activated",
				"subscription.canceled",
				"subscription.past_due",
				"transaction.completed",
				"transaction.payment_failed",
			},
			Type:      "url",
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	respond(w, r, http.StatusOK, map[string]string{"status": "reset"})
}

func (h *AdminHandler) advanceTime(w http.ResponseWriter, r *http.Request, subID string) {
	sub, ok := h.Store.GetSubscription(subID)
	if !ok {
		respondError(w, r, http.StatusNotFound, "request_error", "not_found", "Subscription not found")
		return
	}

	now := time.Now().UTC()

	switch sub.Status {
	case "trialing":
		// Trial → active: simulate trial ending
		sub.Status = "active"
		sub.FirstBilledAt = &now
		sub.UpdatedAt = now

		nextBill := addPeriod(now, sub.BillingCycle.Interval, sub.BillingCycle.Frequency)
		sub.NextBilledAt = &nextBill
		sub.CurrentBillingPeriod = &models.BillingPeriodDates{
			StartsAt: now,
			EndsAt:   nextBill,
		}

		for i := range sub.Items {
			sub.Items[i].Status = "active"
			sub.Items[i].TrialDates = nil
			sub.Items[i].PreviouslyBilledAt = &now
			sub.Items[i].NextBilledAt = &nextBill
			sub.Items[i].UpdatedAt = now
		}

		h.Store.SetSubscription(sub)
		h.createTransaction(sub, "subscription_recurring")
		h.Webhook.Fire("subscription.activated", sub)
		h.Webhook.Fire("transaction.completed", h.Store.ListTransactions()[len(h.Store.ListTransactions())-1])

	case "active":
		// Active → simulate billing cycle. 50/50 chance of payment failure for testing,
		// but default to success. Use query param ?fail=true to force failure.
		if r.URL.Query().Get("fail") == "true" {
			sub.Status = "past_due"
			sub.UpdatedAt = now
			h.Store.SetSubscription(sub)

			// Create failed transaction
			txn := h.createFailedTransaction(sub)
			h.Webhook.Fire("subscription.past_due", sub)
			h.Webhook.Fire("transaction.payment_failed", txn)
		} else {
			// Successful billing — advance to next period
			prevEnd := now
			if sub.CurrentBillingPeriod != nil {
				prevEnd = sub.CurrentBillingPeriod.EndsAt
			}
			nextBill := addPeriod(prevEnd, sub.BillingCycle.Interval, sub.BillingCycle.Frequency)
			sub.NextBilledAt = &nextBill
			sub.CurrentBillingPeriod = &models.BillingPeriodDates{
				StartsAt: prevEnd,
				EndsAt:   nextBill,
			}
			sub.UpdatedAt = now

			for i := range sub.Items {
				sub.Items[i].PreviouslyBilledAt = &prevEnd
				sub.Items[i].NextBilledAt = &nextBill
				sub.Items[i].UpdatedAt = now
			}

			h.Store.SetSubscription(sub)
			h.createTransaction(sub, "subscription_recurring")
			h.Webhook.Fire("subscription.updated", sub)
			h.Webhook.Fire("transaction.completed", h.Store.ListTransactions()[len(h.Store.ListTransactions())-1])
		}

	case "past_due":
		// past_due → canceled
		sub.Status = "canceled"
		sub.CanceledAt = &now
		sub.UpdatedAt = now
		sub.ScheduledChange = nil
		h.Store.SetSubscription(sub)
		h.Webhook.Fire("subscription.canceled", sub)

	case "paused":
		// paused → active (resume)
		sub.Status = "active"
		sub.PausedAt = nil
		sub.UpdatedAt = now
		nextBill := addPeriod(now, sub.BillingCycle.Interval, sub.BillingCycle.Frequency)
		sub.NextBilledAt = &nextBill
		sub.CurrentBillingPeriod = &models.BillingPeriodDates{
			StartsAt: now,
			EndsAt:   nextBill,
		}
		h.Store.SetSubscription(sub)
		h.Webhook.Fire("subscription.updated", sub)

	default:
		respondError(w, r, http.StatusConflict, "request_error", "conflict", "Cannot advance subscription in status: "+sub.Status)
		return
	}

	respond(w, r, http.StatusOK, sub)
}

func (h *AdminHandler) createTransaction(sub *models.Subscription, origin string) {
	now := time.Now().UTC()
	txn := &models.Transaction{
		ID:             store.NextID("txn"),
		Status:         "completed",
		CustomerID:     sub.CustomerID,
		SubscriptionID: &sub.ID,
		CurrencyCode:   sub.CurrencyCode,
		CollectionMode: sub.CollectionMode,
		Origin:         origin,
		Items:          make([]models.TransactionItem, 0),
		BilledAt:       &now,
		CreatedAt:      now,
		UpdatedAt:      now,
		CustomData:     map[string]string{},
	}

	var totalAmount int
	for _, item := range sub.Items {
		txnItem := models.TransactionItem{
			PriceID:  item.Price.ID,
			Quantity: item.Quantity,
			Price:    item.Price,
			Product:  item.Product,
		}
		txn.Items = append(txn.Items, txnItem)
		totalAmount += parseAmount(item.Price.UnitPrice.Amount) * item.Quantity
	}

	total := formatAmount(totalAmount)
	txn.Details = models.TransactionDetails{
		Totals: models.TransactionTotals{
			Subtotal:     total,
			Tax:          "0",
			Total:        total,
			GrandTotal:   total,
			CurrencyCode: sub.CurrencyCode,
		},
	}

	h.Store.SetTransaction(txn)
}

func (h *AdminHandler) createFailedTransaction(sub *models.Subscription) *models.Transaction {
	now := time.Now().UTC()
	txn := &models.Transaction{
		ID:             store.NextID("txn"),
		Status:         "failed",
		CustomerID:     sub.CustomerID,
		SubscriptionID: &sub.ID,
		CurrencyCode:   sub.CurrencyCode,
		CollectionMode: sub.CollectionMode,
		Origin:         "subscription_recurring",
		Items:          make([]models.TransactionItem, 0),
		CreatedAt:      now,
		UpdatedAt:      now,
		CustomData:     map[string]string{},
	}

	var totalAmount int
	for _, item := range sub.Items {
		txnItem := models.TransactionItem{
			PriceID:  item.Price.ID,
			Quantity: item.Quantity,
			Price:    item.Price,
			Product:  item.Product,
		}
		txn.Items = append(txn.Items, txnItem)
		totalAmount += parseAmount(item.Price.UnitPrice.Amount) * item.Quantity
	}

	total := formatAmount(totalAmount)
	txn.Details = models.TransactionDetails{
		Totals: models.TransactionTotals{
			Subtotal:     total,
			Tax:          "0",
			Total:        total,
			GrandTotal:   total,
			CurrencyCode: sub.CurrencyCode,
		},
	}

	h.Store.SetTransaction(txn)
	return txn
}

func (h *AdminHandler) triggerWebhook(w http.ResponseWriter, r *http.Request, eventType string) {
	// Read optional JSON body as event data
	var data interface{}
	if r.Body != nil {
		_ = decodeJSON(r, &data)
	}
	if data == nil {
		data = map[string]string{"triggered": "manual"}
	}

	h.Webhook.Fire(eventType, data)
	respond(w, r, http.StatusOK, map[string]string{"status": "triggered", "event_type": eventType})
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vlah-software-house/paddle-api-mock/internal/handlers"
	"github.com/vlah-software-house/paddle-api-mock/internal/middleware"
	"github.com/vlah-software-house/paddle-api-mock/internal/models"
	"github.com/vlah-software-house/paddle-api-mock/internal/seed"
	"github.com/vlah-software-house/paddle-api-mock/internal/store"
	"github.com/vlah-software-house/paddle-api-mock/internal/webhook"
)

func main() {
	port := flag.Int("port", 8081, "Server port")
	noAuth := flag.Bool("no-auth", false, "Disable API key validation")
	noSeed := flag.Bool("no-seed", false, "Start with empty store")
	webhookURL := flag.String("webhook-url", "", "Default webhook URL to register on startup")
	signingSecret := flag.String("signing-secret", "pdl_test_signing_secret", "Webhook signing secret")
	apiKey := flag.String("api-key", "test_paddle_api_key", "API key for authentication")
	flag.Parse()

	s := store.New()
	if !*noSeed {
		seed.Load(s)
		log.Println("Seed data loaded")
	}

	notifier := webhook.New(s, *signingSecret)

	// Register default webhook URL if provided
	if *webhookURL != "" {
		now := time.Now().UTC()
		s.SetNotificationSetting(&models.NotificationSetting{
			ID:          store.NextID("ntfset"),
			Description: "Default webhook (from CLI)",
			Destination: *webhookURL,
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
		log.Printf("Registered default webhook URL: %s", *webhookURL)
	}

	// Set up handlers
	productsH := &handlers.ProductsHandler{Store: s}
	pricesH := &handlers.PricesHandler{Store: s}
	customersH := &handlers.CustomersHandler{Store: s}
	subscriptionsH := &handlers.SubscriptionsHandler{Store: s, Webhook: notifier}
	transactionsH := &handlers.TransactionsHandler{Store: s}
	eventsH := &handlers.EventsHandler{Store: s}
	notifSettingsH := &handlers.NotificationSettingsHandler{Store: s}
	adminH := &handlers.AdminHandler{Store: s, Webhook: notifier, SeedEnabled: !*noSeed}

	mux := http.NewServeMux()

	// Paddle API v1 routes
	mux.Handle("/v1/products", productsH)
	mux.Handle("/v1/products/", productsH)
	mux.Handle("/v1/prices", pricesH)
	mux.Handle("/v1/prices/", pricesH)
	mux.Handle("/v1/customers", customersH)
	mux.Handle("/v1/customers/", customersH)
	mux.Handle("/v1/subscriptions", subscriptionsH)
	mux.Handle("/v1/subscriptions/", subscriptionsH)
	mux.Handle("/v1/transactions", transactionsH)
	mux.Handle("/v1/transactions/", transactionsH)
	mux.Handle("/v1/events", eventsH)
	mux.Handle("/v1/notification-settings", notifSettingsH)
	mux.Handle("/v1/notification-settings/", notifSettingsH)

	// Admin routes
	mux.Handle("/admin/", adminH)

	// Health check
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Apply middleware
	var handler http.Handler = mux
	handler = middleware.Auth(*apiKey, *noAuth)(handler)
	handler = middleware.JSON(handler)
	handler = middleware.RequestID(handler)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Paddle Mock API starting on %s (auth=%v, seed=%v)", addr, !*noAuth, !*noSeed)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

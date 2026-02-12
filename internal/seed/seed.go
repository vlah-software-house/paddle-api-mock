package seed

import (
	"time"

	"github.com/vlah-software-house/paddle-api-mock/internal/models"
	"github.com/vlah-software-house/paddle-api-mock/internal/store"
)

// Load populates the store with the default seed data for Yieldly testing.
func Load(s *store.Store) {
	now := time.Now().UTC()

	// Product
	s.SetProduct(&models.Product{
		ID:          "prod_yieldly_base",
		Name:        "Yieldly Base Plan",
		Description: strPtr("Base subscription plan for Yieldly"),
		TaxCategory: "standard",
		Status:      "active",
		CustomData:  map[string]string{},
		CreatedAt:   now,
		UpdatedAt:   now,
	})

	// Price
	monthlyName := "Monthly"
	s.SetPrice(&models.Price{
		ID:        "pri_yieldly_monthly",
		ProductID: "prod_yieldly_base",
		Name:      &monthlyName,
		Description: "$5.00/month with 3-month trial",
		Type:      "standard",
		BillingCycle: &models.BillingCycle{
			Interval:  "month",
			Frequency: 1,
		},
		TrialPeriod: &models.TrialPeriod{
			Interval:  "month",
			Frequency: 3,
		},
		TaxMode: "account_setting",
		UnitPrice: models.Money{
			Amount:       "500",
			CurrencyCode: "USD",
		},
		UnitPriceOverrides: []interface{}{},
		Quantity: models.Quantity{
			Minimum: 1,
			Maximum: 100,
		},
		Status:     "active",
		CustomData: map[string]string{},
		CreatedAt:  now,
		UpdatedAt:  now,
	})

	// Customers
	aliceName := "Alice"
	s.SetCustomer(&models.Customer{
		ID:         "ctm_test_alice",
		Name:       &aliceName,
		Email:      "alice@test.com",
		Locale:     "en",
		Status:     "active",
		CustomData: map[string]string{},
		CreatedAt:  now,
		UpdatedAt:  now,
	})

	bobName := "Bob"
	s.SetCustomer(&models.Customer{
		ID:         "ctm_test_bob",
		Name:       &bobName,
		Email:      "bob@test.com",
		Locale:     "en",
		Status:     "active",
		CustomData: map[string]string{},
		CreatedAt:  now,
		UpdatedAt:  now,
	})

	// Subscription for Alice (trialing, trial ends in 90 days)
	trialEnd := now.Add(90 * 24 * time.Hour)
	product, _ := s.GetProduct("prod_yieldly_base")
	price, _ := s.GetPrice("pri_yieldly_monthly")

	s.SetSubscription(&models.Subscription{
		ID:             "sub_test_alice",
		Status:         "trialing",
		CustomerID:     "ctm_test_alice",
		CurrencyCode:   "USD",
		CreatedAt:      now,
		UpdatedAt:      now,
		StartedAt:      &now,
		CollectionMode: "automatic",
		BillingCycle: models.BillingCycle{
			Interval:  "month",
			Frequency: 1,
		},
		CurrentBillingPeriod: &models.BillingPeriodDates{
			StartsAt: now,
			EndsAt:   trialEnd,
		},
		NextBilledAt: &trialEnd,
		Items: []models.SubscriptionItem{
			{
				Status:    "trialing",
				Quantity:  1,
				Recurring: true,
				CreatedAt: now,
				UpdatedAt: now,
				NextBilledAt: &trialEnd,
				TrialDates: &models.BillingPeriodDates{
					StartsAt: now,
					EndsAt:   trialEnd,
				},
				Price:   *price,
				Product: product,
			},
		},
		CustomData: map[string]string{},
	})
}

func strPtr(s string) *string {
	return &s
}

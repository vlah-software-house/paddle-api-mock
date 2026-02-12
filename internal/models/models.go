package models

import "time"

// PaddleResponse is the standard envelope for all Paddle API responses.
type PaddleResponse struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

// PaddleListResponse is the envelope for list endpoints.
type PaddleListResponse struct {
	Data       interface{}     `json:"data"`
	Meta       Meta            `json:"meta"`
	Pagination *PaginationInfo `json:"pagination,omitempty"`
}

type Meta struct {
	RequestID string          `json:"request_id"`
	Pagination *PaginationInfo `json:"pagination,omitempty"`
}

type PaginationInfo struct {
	PerPage        int    `json:"per_page"`
	Next           string `json:"next,omitempty"`
	HasMore        bool   `json:"has_more"`
	EstimatedTotal int    `json:"estimated_total"`
}

// ErrorResponse is the Paddle error envelope.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
	Meta  Meta        `json:"meta"`
}

type ErrorDetail struct {
	Type   string `json:"type"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

// Product represents a Paddle product.
type Product struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description *string           `json:"description"`
	TaxCategory string            `json:"tax_category"`
	ImageURL    *string           `json:"image_url"`
	Status      string            `json:"status"`
	CustomData  map[string]string `json:"custom_data"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// Price represents a Paddle price.
type Price struct {
	ID              string            `json:"id"`
	ProductID       string            `json:"product_id"`
	Name            *string           `json:"name"`
	Description     string            `json:"description"`
	Type            string            `json:"type"`           // "standard" or "custom"
	BillingCycle    *BillingCycle     `json:"billing_cycle"`
	TrialPeriod     *TrialPeriod      `json:"trial_period"`
	TaxMode         string            `json:"tax_mode"`
	UnitPrice       Money             `json:"unit_price"`
	UnitPriceOverrides []interface{}  `json:"unit_price_overrides"`
	Quantity        Quantity          `json:"quantity"`
	Status          string            `json:"status"`
	CustomData      map[string]string `json:"custom_data"`
	Product         *Product          `json:"product,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

type BillingCycle struct {
	Interval  string `json:"interval"`  // "day", "week", "month", "year"
	Frequency int    `json:"frequency"`
}

type TrialPeriod struct {
	Interval  string `json:"interval"`
	Frequency int    `json:"frequency"`
}

type Money struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currency_code"`
}

type Quantity struct {
	Minimum int `json:"minimum"`
	Maximum int `json:"maximum"`
}

// Customer represents a Paddle customer.
type Customer struct {
	ID         string            `json:"id"`
	Name       *string           `json:"name"`
	Email      string            `json:"email"`
	Locale     string            `json:"locale"`
	Status     string            `json:"status"`
	CustomData map[string]string `json:"custom_data"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

type CreateCustomerRequest struct {
	Email      string            `json:"email"`
	Name       *string           `json:"name,omitempty"`
	Locale     string            `json:"locale,omitempty"`
	CustomData map[string]string `json:"custom_data,omitempty"`
}

type UpdateCustomerRequest struct {
	Name       *string           `json:"name,omitempty"`
	Email      string            `json:"email,omitempty"`
	Locale     string            `json:"locale,omitempty"`
	Status     string            `json:"status,omitempty"`
	CustomData map[string]string `json:"custom_data,omitempty"`
}

// Subscription represents a Paddle subscription.
type Subscription struct {
	ID                    string              `json:"id"`
	Status                string              `json:"status"` // trialing, active, past_due, paused, canceled
	CustomerID            string              `json:"customer_id"`
	AddressID             *string             `json:"address_id"`
	BusinessID            *string             `json:"business_id"`
	CurrencyCode          string              `json:"currency_code"`
	CreatedAt             time.Time           `json:"created_at"`
	UpdatedAt             time.Time           `json:"updated_at"`
	StartedAt             *time.Time          `json:"started_at"`
	FirstBilledAt         *time.Time          `json:"first_billed_at"`
	NextBilledAt          *time.Time          `json:"next_billed_at"`
	PausedAt              *time.Time          `json:"paused_at"`
	CanceledAt            *time.Time          `json:"canceled_at"`
	CollectionMode        string              `json:"collection_mode"`
	BillingDetails        *BillingDetails     `json:"billing_details"`
	CurrentBillingPeriod  *BillingPeriodDates `json:"current_billing_period"`
	BillingCycle          BillingCycle         `json:"billing_cycle"`
	ScheduledChange       *ScheduledChange    `json:"scheduled_change"`
	Items                 []SubscriptionItem  `json:"items"`
	CustomData            map[string]string   `json:"custom_data"`
	ManagementURLs        *ManagementURLs     `json:"management_urls"`
	Discount              *interface{}        `json:"discount"`
}

type BillingDetails struct {
	PaymentTerms     BillingCycle `json:"payment_terms"`
	EnableCheckout   bool         `json:"enable_checkout"`
	PurchaseOrderNum string       `json:"purchase_order_number"`
}

type BillingPeriodDates struct {
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
}

type ScheduledChange struct {
	Action      string    `json:"action"` // "cancel", "pause", "resume"
	EffectiveAt time.Time `json:"effective_at"`
	ResumeAt    *time.Time `json:"resume_at,omitempty"`
}

type SubscriptionItem struct {
	Status    string       `json:"status"`
	Quantity  int          `json:"quantity"`
	Recurring bool         `json:"recurring"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	PreviouslyBilledAt *time.Time `json:"previously_billed_at"`
	NextBilledAt       *time.Time `json:"next_billed_at"`
	TrialDates         *BillingPeriodDates `json:"trial_dates"`
	Price     Price        `json:"price"`
	Product   *Product     `json:"product,omitempty"`
}

type ManagementURLs struct {
	UpdatePaymentMethod string `json:"update_payment_method"`
	Cancel              string `json:"cancel"`
}

type CreateSubscriptionRequest struct {
	CustomerID     string              `json:"customer_id"`
	Items          []CreateSubItemReq  `json:"items"`
	CurrencyCode   string              `json:"currency_code,omitempty"`
	CollectionMode string              `json:"collection_mode,omitempty"`
	CustomData     map[string]string   `json:"custom_data,omitempty"`
}

type CreateSubItemReq struct {
	PriceID  string `json:"price_id"`
	Quantity int    `json:"quantity"`
}

type UpdateSubscriptionRequest struct {
	ScheduledChange *ScheduledChangeReq `json:"scheduled_change,omitempty"`
	Items           []CreateSubItemReq  `json:"items,omitempty"`
	ProrationBillingMode string         `json:"proration_billing_mode,omitempty"`
	CustomData      map[string]string   `json:"custom_data,omitempty"`
}

type ScheduledChangeReq struct {
	Action      string     `json:"action"` // "cancel", "pause", "resume"
	EffectiveAt *string    `json:"effective_at,omitempty"`
	ResumeAt    *string    `json:"resume_at,omitempty"`
}

type ChargeRequest struct {
	Items      []ChargeItem `json:"items"`
	EffectiveFrom string   `json:"effective_from,omitempty"` // "next_billing_period" or "immediately"
}

type ChargeItem struct {
	PriceID  string `json:"price_id"`
	Quantity int    `json:"quantity"`
}

// Transaction represents a Paddle transaction.
type Transaction struct {
	ID             string            `json:"id"`
	Status         string            `json:"status"` // "completed", "failed", "past_due"
	CustomerID     string            `json:"customer_id"`
	SubscriptionID *string           `json:"subscription_id"`
	CurrencyCode   string            `json:"currency_code"`
	CollectionMode string            `json:"collection_mode"`
	Origin         string            `json:"origin"` // "subscription_recurring", "subscription_charge", "api"
	Items          []TransactionItem `json:"items"`
	Details        TransactionDetails `json:"details"`
	BilledAt       *time.Time        `json:"billed_at"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	CustomData     map[string]string `json:"custom_data"`
}

type TransactionItem struct {
	PriceID  string `json:"price_id"`
	Quantity int    `json:"quantity"`
	Price    Price  `json:"price"`
	Product  *Product `json:"product,omitempty"`
}

type TransactionDetails struct {
	Totals TransactionTotals `json:"totals"`
}

type TransactionTotals struct {
	Subtotal    string `json:"subtotal"`
	Tax         string `json:"tax"`
	Total       string `json:"total"`
	GrandTotal  string `json:"grand_total"`
	CurrencyCode string `json:"currency_code"`
}

// Event represents a fired webhook event.
type Event struct {
	EventID    string      `json:"event_id"`
	EventType  string      `json:"event_type"`
	OccurredAt time.Time   `json:"occurred_at"`
	Data       interface{} `json:"data"`
}

// NotificationSetting represents a webhook endpoint configuration.
type NotificationSetting struct {
	ID              string    `json:"id"`
	Description     string    `json:"description"`
	Destination     string    `json:"destination"`
	Active          bool      `json:"active"`
	APIVersion      int       `json:"api_version"`
	IncludeSensitiveFields bool `json:"include_sensitive_fields"`
	SubscribedEvents []string `json:"subscribed_events"`
	Type            string    `json:"type"` // "url"
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateNotificationSettingRequest struct {
	Description     string   `json:"description"`
	Destination     string   `json:"destination"`
	SubscribedEvents []string `json:"subscribed_events,omitempty"`
	Active          *bool    `json:"active,omitempty"`
	APIVersion      *int     `json:"api_version,omitempty"`
	IncludeSensitiveFields *bool `json:"include_sensitive_fields,omitempty"`
	Type            string   `json:"type,omitempty"`
}

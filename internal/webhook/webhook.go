package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vlah-software-house/paddle-api-mock/internal/models"
	"github.com/vlah-software-house/paddle-api-mock/internal/store"
)

// Notifier handles firing webhooks to registered endpoints.
type Notifier struct {
	Store         *store.Store
	SigningSecret string
	Client        *http.Client
}

func New(s *store.Store, signingSecret string) *Notifier {
	return &Notifier{
		Store:         s,
		SigningSecret: signingSecret,
		Client:        &http.Client{Timeout: 10 * time.Second},
	}
}

// Fire creates an event and sends it to all registered webhook endpoints.
func (n *Notifier) Fire(eventType string, data interface{}) {
	event := &models.Event{
		EventID:    store.NextID("evt"),
		EventType:  eventType,
		OccurredAt: time.Now().UTC(),
		Data:       data,
	}
	n.Store.AddEvent(event)

	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("webhook: failed to marshal event: %v", err)
		return
	}

	settings := n.Store.ListNotificationSettings()
	for _, ns := range settings {
		if !ns.Active {
			continue
		}
		n.send(ns.Destination, payload)
	}
}

func (n *Notifier) send(url string, payload []byte) {
	ts := fmt.Sprintf("%d", time.Now().Unix())
	signedPayload := ts + ":" + string(payload)

	mac := hmac.New(sha256.New, []byte(n.SigningSecret))
	mac.Write([]byte(signedPayload))
	h1 := hex.EncodeToString(mac.Sum(nil))

	signature := fmt.Sprintf("ts=%s;h1=%s", ts, h1)

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		log.Printf("webhook: failed to create request for %s: %v", url, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Paddle-Signature", signature)

	resp, err := n.Client.Do(req)
	if err != nil {
		log.Printf("webhook: failed to POST to %s: %v", url, err)
		return
	}
	resp.Body.Close()
	log.Printf("webhook: POST %s → %d (%s)", url, resp.StatusCode, signature[:30]+"...")
}

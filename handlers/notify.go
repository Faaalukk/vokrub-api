package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Faaalukk/vokrub-api.git/models"
)

// notifyDiscordNewCustomer posts a "new customer registered" message to the
// Discord webhook in DISCORD_WEBHOOKS_LINK. No-op when the env var is empty.
// Fire-and-forget — never blocks or fails the request.
func notifyDiscordNewCustomer(customer *models.Customer, source string) {
	webhook := os.Getenv("DISCORD_WEBHOOKS_LINK")
	if webhook == "" {
		return
	}

	go func() {
		contact := "—"
		if customer.Email != nil {
			contact = *customer.Email
		} else if customer.Phone != nil {
			contact = *customer.Phone
		}

		payload := map[string]any{
			"embeds": []map[string]any{{
				"title": "🎉 New customer registered",
				"color": 5763719, // green
				"fields": []map[string]any{
					{"name": "Name", "value": customer.Name, "inline": true},
					{"name": "Contact", "value": contact, "inline": true},
					{"name": "Plan", "value": customer.Plan, "inline": true},
					{"name": "Source", "value": source, "inline": true},
				},
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			}},
		}

		body, err := json.Marshal(payload)
		if err != nil {
			log.Printf("discord notify: marshal failed: %v", err)
			return
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Post(webhook, "application/json", bytes.NewReader(body))
		if err != nil {
			log.Printf("discord notify: post failed: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			log.Printf("discord notify: webhook returned %d", resp.StatusCode)
		}
	}()
}

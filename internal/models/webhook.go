package models

type PipefyWebhookRequest struct {
	EventID      string `json:"event_id"`
	CardID       string `json:"card_id"`
	ClientEmail  string `json:"cliente_email"`
	Timestamp    string `json:"timestamp"`
}
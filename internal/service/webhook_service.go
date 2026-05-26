package service

import (
	"errors"
	"strconv"

	"client-core/internal/models"
	"client-core/internal/pipefy"
	"client-core/internal/repository"
)

type WebhookService struct {
	clientRepo  *repository.ClientRepository
	webhookRepo *repository.WebhookRepository
	pipefy      *pipefy.PipefyClient
}

func NewWebhookService(cRepo *repository.ClientRepository, wRepo *repository.WebhookRepository, p *pipefy.PipefyClient) *WebhookService {
	return &WebhookService{clientRepo: cRepo, webhookRepo: wRepo, pipefy: p}
}

func (s *WebhookService) ProcessWebhook(req models.PipefyWebhookRequest) (string, error) {
	if req.EventID == "" {
		return "", errors.New("event_id is required")
	}

	processed, err := s.webhookRepo.IsProcessed(req.EventID)
	if err != nil {
		return "", err
	}
	if processed {
		return "", nil
	}

	client, err := s.clientRepo.FindByEmail(req.ClientEmail)
	if err != nil {
		return "", err
	}

	priority := "prioridade_normal"
	if client.Assets >= 200000 {
		priority = "prioridade_alta"
	}

	client.Status = "Processado"
	client.Priority = priority

	if err := s.clientRepo.UpdateClient(*client); err != nil {
		return "", err
	}

	if err := s.webhookRepo.MarkProcessed(req.EventID, req.CardID); err != nil {
		return "", err
	}

	nodeID := int64(0)
	if parsed, perr := strconv.ParseInt(req.CardID, 10, 64); perr == nil {
		nodeID = parsed
	}

	values := map[string]string{
		"status":   client.Status,
		"priority": client.Priority,
	}

	mutation := s.pipefy.BuildUpdateFieldsValuesMutation(nodeID, values)

	return mutation, nil
}

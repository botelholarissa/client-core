package service

import (
	"database/sql"
	"errors"
	"strconv"

	"client-core/internal/models"
	"client-core/internal/repository"
)

type PipefyUpdater interface {
	BuildUpdateFieldsValuesMutation(int64, map[string]string) string
}

type WebhookService struct {
	clientRepo  *repository.ClientRepository
	webhookRepo *repository.WebhookRepository
	pipefy      PipefyUpdater
}

func NewWebhookService(cRepo *repository.ClientRepository, wRepo *repository.WebhookRepository, p PipefyUpdater) *WebhookService {
	return &WebhookService{clientRepo: cRepo, webhookRepo: wRepo, pipefy: p}
}

func (s *WebhookService) ProcessWebhook(req models.PipefyWebhookRequest) (string, error) {
	if req.EventID == "" {
		return "", errors.New("event_id é obrigatório")
	}

	processed, err := s.webhookRepo.IsProcessed(req.EventID)
	if err != nil {
		return "", errors.New("erro interno ao verificar evento")
	}
	if processed {
		return "", nil
	}

	client, err := s.clientRepo.FindByEmail(req.ClientEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("cliente não encontrado")
		}
		return "", errors.New("erro interno ao buscar cliente")
	}

	priority := "prioridade_normal"
	if client.Assets >= 200000 {
		priority = "prioridade_alta"
	}

	client.Status = "Processado"
	client.Priority = priority

	if err := s.clientRepo.UpdateClient(*client); err != nil {
		return "", errors.New("erro interno ao atualizar cliente")
	}

	if err := s.webhookRepo.MarkProcessed(req.EventID, req.CardID); err != nil {
		return "", errors.New("erro interno ao marcar evento")
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

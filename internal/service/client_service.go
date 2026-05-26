package service

import (
	"errors"
	"net/mail"

	"client-core/internal/models"
	"client-core/internal/pipefy"
	"client-core/internal/repository"
)

type ClientService struct {
	repository   *repository.ClientRepository
	pipefyClient *pipefy.PipefyClient
	pipeID       int64
}

func NewClientService(clientRepository *repository.ClientRepository, pipefyClient *pipefy.PipefyClient, pipeID int64) *ClientService {
	return &ClientService{repository: clientRepository, pipefyClient: pipefyClient, pipeID: pipeID}
}

func (s *ClientService) CreateClient(request models.CreateClientRequest) (string, error) {
	if request.ClientEmail == "" {
		return "", errors.New("client email is required")
	}

	if _, err := mail.ParseAddress(request.ClientEmail); err != nil {
		return "", errors.New("invalid email")
	}

	if request.ClientName == "" {
		return "", errors.New("client name is required")
	}

	client := models.Client{
		Name:     request.ClientName,
		Email:    request.ClientEmail,
		Assets:   request.AssetsValue,
		Status:   "Aguardando Análise",
		Priority: "",
	}

	if err := s.repository.CreateClient(client); err != nil {
		return "", err
	}

	mutation := s.pipefyClient.BuildCreateCardMutation(s.pipeID, request)
	return mutation, nil
}

package service

import (
	"database/sql"
	"errors"
	"net/mail"
	"strings"

	"client-core/internal/models"
	"client-core/internal/repository"
)

type PipefyCreateClient interface {
	BuildCreateCardMutation(int64, models.CreateClientRequest) string
}

type ClientService struct {
	repository        *repository.ClientRepository
	pipefyClient PipefyCreateClient
	pipeID       int64
}

func NewClientService(clientRepository *repository.ClientRepository, pipefyClient PipefyCreateClient, pipeID int64) *ClientService {
	return &ClientService{repository: clientRepository, pipefyClient: pipefyClient, pipeID: pipeID}
}

func (s *ClientService) CreateClient(request models.CreateClientRequest) (string, error) {
	if request.ClientEmail == "" {
		return "", errors.New("email do cliente é obrigatório")
	}

	if _, err := mail.ParseAddress(request.ClientEmail); err != nil {
		return "", errors.New("email inválido")
	}

	if request.ClientName == "" {
		return "", errors.New("nome do cliente é obrigatório")
	}

	client := models.Client{
		Name:     request.ClientName,
		Email:    request.ClientEmail,
		Assets:   request.AssetsValue,
		Status:   "Aguardando Análise",
		Priority: "",
	}

	if err := s.repository.CreateClient(client); err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return "", errors.New("email já cadastrado")
		}
		return "", errors.New("erro interno ao salvar cliente")
	}

	mutation := s.pipefyClient.BuildCreateCardMutation(s.pipeID, request)
	return mutation, nil
}

func (s *ClientService) GetClient(email string) (*models.Client, error) {
	client, err := s.repository.FindByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("cliente não encontrado")
		}
		return nil, errors.New("erro interno ao buscar cliente")
	}
	return client, nil
}

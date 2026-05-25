package service

import (
	"errors"
	"net/mail"
	"client-core/internal/models"
	"client-core/internal/repository"
)

type ClientService struct {
	repository *repository.ClientRepository
}

func NewClientService(clientRepository *repository.ClientRepository) *ClientService {
	return &ClientService{repository: clientRepository}
}

func (s *ClientService) CreateClient(request models.CreateClientRequest) error {
		if request.ClientEmail == "" {
		return errors.New("client email is required")
	}

	if _, err := mail.ParseAddress(request.ClientEmail); err != nil {
		return errors.New("invalid email")
	}

	if request.ClientName == "" {
		return errors.New("client name is required")
	}

	client := models.Client{
		Name:   request.ClientName,
		Email:  request.ClientEmail,
		Assets: request.AssetsValue,
		Status: "Aguardando Análise",
		Priority: "",
	}

	if err := s.repository.CreateClient(client); err != nil {
		return err
	}
	
	return nil
}


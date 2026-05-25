package pipefy

import "client-core/internal/models"

type PipefyClient struct{}

func NewPipefyClient() *PipefyClient {
	return &PipefyClient{}
}
// TODO: implement Pipefy GraphQL mutations based on doc

func (p *PipefyClient) BuildCreateCardMutation(request models.CreateClientRequest) string {
	return ""
}

func (p *PipefyClient) BuildUpdateCardMutation() string {
	return ""
}

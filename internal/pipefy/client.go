package pipefy

import (
	"fmt"
	"strings"

	"client-core/internal/models"
)

type PipefyClient struct {
	// In a real client we would hold auth, endpoint, etc.
}

func NewPipefyClient() *PipefyClient {
	return &PipefyClient{}
}

func (p *PipefyClient) BuildCreateCardMutation(pipeID int64, request models.CreateClientRequest) string {
	// map input fields to fields_attributes

	fields := []string{
		fmt.Sprintf(`{field_id: "cliente_nome", field_value: "%s"}`, escapeString(request.ClientName)),
		fmt.Sprintf(`{field_id: "cliente_email", field_value: "%s"}`, escapeString(request.ClientEmail)),
		fmt.Sprintf(`{field_id: "valor_patrimonio", field_value: "%v"}`, request.AssetsValue),
		fmt.Sprintf(`{field_id: "tipo_solicitacao", field_value: "%s"}`, escapeString(request.RequestType)),
	}

	mutation := fmt.Sprintf(`mutation { createCard(input: { pipe_id: %d, fields_attributes: [%s] }) { card { id } } }`, pipeID, strings.Join(fields, ", "))

	return mutation
}

func (p *PipefyClient) BuildUpdateFieldsValuesMutation(nodeID int64, values map[string]string) string {
	// values maps each Pipefy fieldId to its new value
	pairs := []string{}
	for fieldId, value := range values {
		pairs = append(pairs, fmt.Sprintf(`{fieldId: "%s", value: "%s"}`, fieldId, escapeString(value)))
	}

	mutation := fmt.Sprintf(`mutation { updateFieldsValues(input: { nodeId: %d, values: [%s] }) { success } }`, nodeID, strings.Join(pairs, ", "))

	return mutation
}

func (p *PipefyClient) BuildUpdateCardFieldMutation(cardID int64, fieldID string, newValue string) string {
	mutation := fmt.Sprintf(`mutation { updateCardField(input: { card_id: %d, field_id: "%s", new_value: "%s" }) { card { id } } }`, cardID, fieldID, escapeString(newValue))
	return mutation
}

func escapeString(s string) string {
	s = strings.ReplaceAll(s, `"`, `\\"`)
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}

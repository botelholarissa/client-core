package models

type Client struct {
	ID        int64
	Name      string
	Email     string
	Assets    float64
	Status    string
	Priority  string
}

type CreateClientRequest struct {
	ClientName     string  `json:"cliente_nome"`
	ClientEmail    string  `json:"cliente_email"`
	RequestType    string  `json:"tipo_solicitacao"`
	AssetsValue    float64 `json:"valor_patrimonio"`
}
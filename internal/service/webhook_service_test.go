package service

import (
	"database/sql"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"client-core/internal/models"
	"client-core/internal/repository"
)

type mockWebhookPipefy struct {
	mutation string
}

func (m *mockWebhookPipefy) BuildUpdateFieldsValuesMutation(nodeID int64, values map[string]string) string {
	return m.mutation
}

var _ = Describe("WebhookService", func() {
	var (
		db          *sql.DB
		clientRepo  *repository.ClientRepository
		webhookRepo *repository.WebhookRepository
		pipefy      *mockWebhookPipefy
		svc         *WebhookService
		req         models.PipefyWebhookRequest
	)

	BeforeEach(func() {
		var err error
		db, err = createInMemoryDB()
		Expect(err).ToNot(HaveOccurred())

		clientRepo = repository.NewClientRepository(db)
		webhookRepo = repository.NewWebhookRepository(db)
		pipefy = &mockWebhookPipefy{mutation: "update-mutation"}
		svc = NewWebhookService(clientRepo, webhookRepo, pipefy)
		req = models.PipefyWebhookRequest{
			EventID:     "evt_1",
			CardID:      "123",
			ClientEmail: "joao.silva@email.com",
			Timestamp:   "2026-05-18T12:00:00Z",
		}

		err = clientRepo.CreateClient(models.Client{
			Name:   "João Silva",
			Email:  req.ClientEmail,
			Assets: 250000,
			Status: "Aguardando Análise",
		})
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		db.Close()
	})

	Describe("ProcessWebhook", func() {
		It("updates the client and marks the event processed", func() {
			mutation, err := svc.ProcessWebhook(req)

			Expect(err).ToNot(HaveOccurred())
			Expect(mutation).To(Equal("update-mutation"))

			updated, err := clientRepo.FindByEmail(req.ClientEmail)
			Expect(err).ToNot(HaveOccurred())
			Expect(updated.Status).To(Equal("Processado"))
			Expect(updated.Priority).To(Equal("prioridade_alta"))

			processed, err := webhookRepo.IsProcessed(req.EventID)
			Expect(err).ToNot(HaveOccurred())
			Expect(processed).To(BeTrue())
		})

		It("sets prioridade_normal when assets are below 200000", func() {
			err := clientRepo.CreateClient(models.Client{
				Name:   "Maria Souza",
				Email:  "maria@email.com",
				Assets: 50000,
				Status: "Aguardando Análise",
			})
			Expect(err).ToNot(HaveOccurred())

			lowReq := models.PipefyWebhookRequest{
				EventID:     "evt_low",
				CardID:      "789",
				ClientEmail: "maria@email.com",
				Timestamp:   "2026-05-18T12:00:00Z",
			}

			mutation, err := svc.ProcessWebhook(lowReq)

			Expect(err).ToNot(HaveOccurred())
			Expect(mutation).To(Equal("update-mutation"))

			updated, err := clientRepo.FindByEmail("maria@email.com")
			Expect(err).ToNot(HaveOccurred())
			Expect(updated.Status).To(Equal("Processado"))
			Expect(updated.Priority).To(Equal("prioridade_normal"))
		})

		It("returns an error when event_id is empty", func() {
			req.EventID = ""

			_, err := svc.ProcessWebhook(req)

			Expect(err).To(MatchError("event_id é obrigatório"))
		})

		It("returns an error when the client email does not exist", func() {
			req.ClientEmail = "carlos@email.com"

			_, err := svc.ProcessWebhook(req)

			Expect(err).To(MatchError("cliente não encontrado"))
		})

		It("returns an empty mutation when the event was already processed", func() {
			err := webhookRepo.MarkProcessed(req.EventID, req.CardID)
			Expect(err).ToNot(HaveOccurred())

			mutation, err := svc.ProcessWebhook(req)

			Expect(err).ToNot(HaveOccurred())
			Expect(mutation).To(BeEmpty())
		})
	})
})

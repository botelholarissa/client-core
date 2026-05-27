package service

import (
	"database/sql"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"client-core/internal/models"
	"client-core/internal/repository"
)

type mockPipefyClient struct {
	mutation string
}

func (m *mockPipefyClient) BuildCreateCardMutation(pipeID int64, request models.CreateClientRequest) string {
	return m.mutation
}

var _ = Describe("ClientService", func() {
	var (
		db      *sql.DB
		repo    *repository.ClientRepository
		pipefy  *mockPipefyClient
		svc     *ClientService
		request models.CreateClientRequest
	)

	BeforeEach(func() {
		var err error
		db, err = createInMemoryDB()
		Expect(err).ToNot(HaveOccurred())

		repo = repository.NewClientRepository(db)
		pipefy = &mockPipefyClient{mutation: "mutation-placeholder"}
		svc = NewClientService(repo, pipefy, 34)
		request = models.CreateClientRequest{
			ClientName:  "João Silva",
			ClientEmail: "joao.silva@email.com",
			RequestType: "Atualização cadastral",
			AssetsValue: 250000,
		}
	})

	AfterEach(func() {
		db.Close()
	})

	Describe("CreateClient", func() {
		It("creates a client in the database and returns a Pipefy mutation", func() {
			mutation, err := svc.CreateClient(request)

			Expect(err).ToNot(HaveOccurred())
			Expect(mutation).To(Equal("mutation-placeholder"))

			saved, err := repo.FindByEmail(request.ClientEmail)
			Expect(err).ToNot(HaveOccurred())
			Expect(saved).ToNot(BeNil())
			Expect(saved.Name).To(Equal(request.ClientName))
			Expect(saved.Email).To(Equal(request.ClientEmail))
			Expect(saved.Status).To(Equal("Aguardando Análise"))
		})

		It("returns an error when the email is blank", func() {
			request.ClientEmail = ""

			_, err := svc.CreateClient(request)

			Expect(err).To(MatchError("email do cliente é obrigatório"))
		})

		It("returns an error when the email is invalid", func() {
			request.ClientEmail = "invalid-email"

			_, err := svc.CreateClient(request)

			Expect(err).To(MatchError("email inválido"))
		})

		It("returns an error when the name is blank", func() {
			request.ClientName = ""

			_, err := svc.CreateClient(request)

			Expect(err).To(MatchError("nome do cliente é obrigatório"))
		})

		It("returns a friendly error when saving a duplicate email", func() {
			err := repo.CreateClient(models.Client{
				Name:     "Duplicate",
				Email:    request.ClientEmail,
				Assets:   1000,
				Status:   "Aguardando Análise",
				Priority: "",
			})
			Expect(err).ToNot(HaveOccurred())

			_, err = svc.CreateClient(request)

			Expect(err).To(MatchError("email já cadastrado"))
		})
	})
})

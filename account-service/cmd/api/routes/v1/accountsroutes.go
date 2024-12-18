package controllersV1

import (
	"encoding/json"
	"github.com/matheus-oliveira-andrade/ledger/account-service/cmd/api/routes/utils"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/logger"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/usecases"
)

type AccountsController struct {
	logger               logger.LoggerInterface
	createAccountUseCase usecases.CreateAccountUseCaseInterface
}

func NewAccountsController(
	logger logger.LoggerInterface,
	createAccountUseCase usecases.CreateAccountUseCaseInterface) *AccountsController {
	return &AccountsController{
		logger:               logger,
		createAccountUseCase: createAccountUseCase,
	}
}

func (ar *AccountsController) RegisterRoutes(router chi.Router) {
	router.Route("/accounts", func(r chi.Router) {
		r.Post("/", ar.createAccountHandle)
	})
}

func (ar *AccountsController) createAccountHandle(w http.ResponseWriter, r *http.Request) {
	ar.logger.LogInformation("processing create account")

	reqBody := &CreateAccountRequest{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		ar.logger.LogError("invalid request", "err", err.Error())
		utils.RenderErrorJsonResponse(w, err)
		return
	}

	accId, err := ar.createAccountUseCase.Handle(reqBody.Document, reqBody.Name)
	if err != nil {
		utils.RenderErrorJsonResponse(w, err)
		return
	}

	response := struct {
		AccountId string `json:"accountId"`
	}{
		AccountId: accId,
	}

	utils.RenderJsonResponse(w, response, http.StatusOK)
}

type CreateAccountRequest struct {
	Document string `json:"document"`
	Name     string `json:"name"`
}

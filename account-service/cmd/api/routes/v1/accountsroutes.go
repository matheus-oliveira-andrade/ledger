package controllersV1

import (
	"encoding/json"
	"errors"
	"github.com/matheus-oliveira-andrade/ledger/account-service/cmd/api/routes/utils"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/slogger"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/usecases"
)

type AccountsController struct {
	logger               slogger.LoggerInterface
	createAccountUseCase usecases.CreateAccountUseCaseInterface
	getAccountUseCase    usecases.GetAccountUseCaseInterface
}

func NewAccountsController(
	logger slogger.LoggerInterface,
	createAccountUseCase usecases.CreateAccountUseCaseInterface,
	getAccountUseCase usecases.GetAccountUseCaseInterface,
) *AccountsController {
	return &AccountsController{
		logger:               logger,
		createAccountUseCase: createAccountUseCase,
		getAccountUseCase:    getAccountUseCase,
	}
}

func (ar *AccountsController) RegisterRoutes(router chi.Router) {
	router.Route("/accounts", func(r chi.Router) {
		r.Post("/", ar.createAccountHandle)
		r.Get("/{accountId}", ar.getAccountHandle)
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

func (ar *AccountsController) getAccountHandle(w http.ResponseWriter, r *http.Request) {

	accId := chi.URLParam(r, "accountId")
	if accId == "" {
		utils.RenderErrorJsonResponse(w, errors.New("required parameter accountId not found"))
		return
	}

	acc, err := ar.getAccountUseCase.Handle(accId)

	if err != nil {
		utils.RenderErrorJsonResponse(w, err)
		return
	}

	if acc == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	utils.RenderJsonResponse(w, GetAccountResponse{
		Id:        acc.Id,
		Name:      acc.Name,
		Document:  acc.Document,
		CreatedAt: acc.CreatedAt,
	}, http.StatusOK)
	return
}

type CreateAccountRequest struct {
	Document string `json:"document"`
	Name     string `json:"name"`
}

type GetAccountResponse struct {
	Id        string    `json:"id"`
	Document  string    `json:"document"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

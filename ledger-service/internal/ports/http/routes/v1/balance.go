package controllersV1

import (
	"errors"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/ports/http/routes/utils"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/utils/slogger"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/usecases"
)

type BalanceController struct {
	logger            slogger.LoggerInterface
	getBalanceUseCase usecases.GetBalanceUseCaseInterface
}

func NewBalanceController(
	logger slogger.LoggerInterface,
	getBalanceUseCase usecases.GetBalanceUseCaseInterface,
) *BalanceController {
	return &BalanceController{
		logger:            logger,
		getBalanceUseCase: getBalanceUseCase,
	}
}

func (c *BalanceController) RegisterRoutes(router chi.Router) {
	router.Route("/balance", func(r chi.Router) {
		r.Get("/{accountId}", c.getBalanceHandler)
	})
}

func (c *BalanceController) getBalanceHandler(w http.ResponseWriter, r *http.Request) {
	c.logger.LogInformation("processing get balance")

	accIdStr := chi.URLParam(r, "accountId")
	if accIdStr == "" {
		utils.RenderErrorJsonResponse(w, errors.New("required parameter accountId not found"))
		return
	}

	accId, _ := strconv.ParseInt(accIdStr, 10, 64)

	balance, err := c.getBalanceUseCase.Handle(accId)
	if err != nil {
		utils.RenderErrorJsonResponse(w, err)
		return
	}

	utils.RenderJsonResponse(w, &GetBalanceResponse{
		AccountId: accId,
		Balance:   balance,
	}, http.StatusOK)
}

type GetBalanceResponse struct {
	AccountId int64 `json:"accountId"`
	Balance   int64 `json:"balance"`
}

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

type StatementController struct {
	logger              slogger.LoggerInterface
	getStatementUseCase usecases.GetStatementUseCaseInterface
}

func NewStatementController(
	logger slogger.LoggerInterface,
	getBalanceUseCase usecases.GetStatementUseCaseInterface,
) *StatementController {
	return &StatementController{
		logger:              logger,
		getStatementUseCase: getBalanceUseCase,
	}
}

func (c *StatementController) RegisterRoutes(router chi.Router) {
	router.Route("/statement", func(r chi.Router) {
		r.Get("/{accountId}", c.getStatementHandler)
	})
}

func (c *StatementController) getStatementHandler(w http.ResponseWriter, r *http.Request) {
	c.logger.LogInformation("processing get balance")

	accIdStr := chi.URLParam(r, "accountId")
	if accIdStr == "" {
		utils.RenderErrorJsonResponse(w, errors.New("required parameter accountId not found"))
		return
	}
	accId, _ := strconv.ParseInt(accIdStr, 10, 64)

	transactionType, period, limit, page := c.readGetStatementQueryParams(r)

	result, err := c.getStatementUseCase.Handle(accId, transactionType, period, limit, page)

	if err != nil {
		utils.RenderErrorJsonResponse(w, err)
		return
	}

	utils.RenderJsonResponse(w, &result, http.StatusOK)
}

func (c *StatementController) readGetStatementQueryParams(r *http.Request) (string, int, int, int) {
	transactionType := r.URL.Query().Get("transactionType")
	if transactionType == "" {
		transactionType = "ALL"
	}

	periodStr := r.URL.Query().Get("period")
	if periodStr == "" {
		periodStr = "30"
	}
	period, _ := strconv.Atoi(periodStr)

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "50"
	}
	limit, _ := strconv.Atoi(limitStr)

	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, _ := strconv.Atoi(pageStr)

	return transactionType, period, limit, page
}

package controllersV1

import (
	"encoding/json"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/ports/http/routes/utils"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/utils/slogger"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/usecases"
)

type FundsTransferController struct {
	logger               slogger.LoggerInterface
	fundsTransferUseCase usecases.FundsTransferUseCaseInterface
}

func NewFundsTransferController(
	logger slogger.LoggerInterface,
	fundsTransferUseCase usecases.FundsTransferUseCaseInterface,
) *FundsTransferController {
	return &FundsTransferController{
		logger:               logger,
		fundsTransferUseCase: fundsTransferUseCase,
	}
}

func (c *FundsTransferController) RegisterRoutes(router chi.Router) {
	router.Route("/funds-transfer", func(r chi.Router) {
		r.Post("/", c.fundsTransferHandler)
	})
}

func (c *FundsTransferController) fundsTransferHandler(w http.ResponseWriter, r *http.Request) {
	c.logger.LogInformation("processing request create account")

	reqBody := &FundsTransferRequest{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		c.logger.LogError("invalid request", "err", err.Error())
		utils.RenderErrorJsonResponse(w, err)
		return
	}

	err = c.fundsTransferUseCase.Handle(r.Context(), reqBody.AccFrom, reqBody.AccTo, reqBody.Amount)
	if err != nil {
		utils.RenderErrorJsonResponse(w, err)
		return
	}

	utils.RenderJsonResponse(w, nil, http.StatusNoContent)
}

type FundsTransferRequest struct {
	AccFrom int64 `json:"accFrom"`
	AccTo   int64 `json:"accTo"`
	Amount  int64 `json:"Amount"`
}

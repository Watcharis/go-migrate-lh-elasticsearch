package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"watcharis/go-migrate-lotto-history-els/services"

	"github.com/labstack/echo/v4"
)

type researchElasticAndDatabase struct {
	lottoHistoryService services.LottoHistoryServices
}

func NewResearchElasticAndDatabase(lottoHistoryServices services.LottoHistoryServices) ResearchElasticAndDatabase {
	return &researchElasticAndDatabase{
		lottoHistoryService: lottoHistoryServices,
	}
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"string"`
}

func (h *researchElasticAndDatabase) ResearchElsWithDbHandler(c echo.Context) error {
	ctx := c.Request().Context()
	slog.InfoContext(ctx, "start research elastic with database handler")

	flow := c.QueryParam("flow")
	fmt.Println("flow :", flow)
	switch flow {
	case "1":
		if err := h.lottoHistoryService.CreateIndexIfNotExists(ctx); err != nil {
			slog.ErrorContext(ctx, "[ERROR] elastic lotto_history create index if not exists", slog.Any("error", err))
			return c.JSONPretty(http.StatusOK, Response{Code: 5000, Message: "failed flow 1"}, "")
		}
	case "2":
		if err := h.lottoHistoryService.MigrateDataLottoHistoryToElastic(ctx); err != nil {
			// log.Panic(err)
			slog.ErrorContext(ctx, "[ERROR] migrate lotto_history to elastic failed", slog.Any("error", err))
			return c.JSONPretty(http.StatusOK, Response{Code: 5000, Message: "failed flow 2"}, "")
		}
	case "3":
		if err := h.lottoHistoryService.GetLottoHistory(ctx); err != nil {
			slog.ErrorContext(ctx, "[ERROR] get lotto_history failed", slog.Any("error", err))
			return c.JSONPretty(http.StatusOK, Response{Code: 5000, Message: "failed flow 3"}, "")
		}
	case "4":
		if err := h.lottoHistoryService.TestDBUseTransaction(ctx); err != nil {
			slog.ErrorContext(ctx, "[ERROR] db use transaction failed", slog.Any("error", err))
			return c.JSONPretty(http.StatusOK, Response{Code: 5000, Message: "failed flow 4"}, "")
		}
	case "5":
		if err := h.lottoHistoryService.MigrateSpeacificLotto(ctx); err != nil {
			slog.ErrorContext(ctx, "[ERROR] migrate speacific lotto failed", slog.Any("error", err))
			return c.JSONPretty(http.StatusOK, Response{Code: 5000, Message: "failed flow 5"}, "")
		}
	case "6":
		if err := h.lottoHistoryService.UpdateRewardStatusByElasticQuery(ctx); err != nil {
			slog.ErrorContext(ctx, "[ERROR] update reward_status by elastic query failed", slog.Any("error", err))
			return c.JSONPretty(http.StatusOK, Response{Code: 5000, Message: "failed flow 6"}, "")
		}
	case "7":
		if err := h.lottoHistoryService.UpdateRewardStatus(ctx); err != nil {
			slog.ErrorContext(ctx, "[ERROR] update reward_status failed", slog.Any("error", err))
			return c.JSONPretty(http.StatusOK, Response{Code: 5000, Message: "failed flow 7"}, "")
		}
	case "8":
		if err := h.lottoHistoryService.MigrateSoldLottos(ctx); err != nil {
			slog.ErrorContext(ctx, "[ERROR] migrate sold lotto", slog.Any("error", err))
			return c.JSONPretty(http.StatusOK, Response{Code: 5000, Message: "failed flow 8"}, "")
		}
	case "9":
		if err := h.lottoHistoryService.PocBulkData(ctx); err != nil {
			slog.ErrorContext(ctx, "[ERROR] bulk update data lotto_history failed", slog.Any("error", err))
			return c.JSONPretty(http.StatusOK, Response{Code: 5000, Message: "failed flow 9"}, "")
		}
	case "10":
		if err := h.lottoHistoryService.GetMultipleLottoHistoryAndMigrateToElasticsearch(ctx); err != nil {
			slog.ErrorContext(ctx, "[ERROR] bulk update data lotto_history failed", slog.Any("error", err))
			return c.JSONPretty(http.StatusOK, Response{Code: 5000, Message: "failed flow 10"}, "")
		}
	case "11":
		if err := h.lottoHistoryService.PocElasticsearchPIT(ctx); err != nil {
			slog.ErrorContext(ctx, "[ERROR] poc pit elasticsearch lotto_history failed", slog.Any("error", err))
			return c.JSONPretty(http.StatusOK, Response{Code: 5000, Message: "failed flow 11"}, "")
		}
	default:
		slog.InfoContext(ctx, "not do anything and pass")
	}

	return c.JSONPretty(http.StatusOK, Response{Code: 200, Message: "success"}, "")
}

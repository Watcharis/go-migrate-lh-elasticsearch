package services

import (
	"context"
	"fmt"
	"log/slog"
	"watcharis/go-migrate-lotto-history-els/models"
)

func (s *lottoHistoryServices) PocBulkData(ctx context.Context) error {

	// now := time.Now()

	// index := fmt.Sprintf("%s.%s", models.INDEX_POV, now.Format(time.DateOnly))

	// // index := "lotto_history_pov.2024-02-04"

	// if err := s.lottoHistoryElasticSearch.BulkData(ctx, index); err != nil {
	// 	s.slogger.ErrorContext(ctx, "cannot bulk data to elasticsearch", slog.Any("error", err))
	// 	return err
	// }

	index := fmt.Sprintf(models.INDEX_POV, "3011-02-01")
	ids := []int{
		90124570054,
		90124570056,
	}
	if err := s.lottoHistoryElasticSearch.BulkUpdateData(ctx, index, ids, ""); err != nil {
		slog.ErrorContext(ctx, "cannot bulk update data", slog.Any("error", err))
		return err
	}

	return nil
}

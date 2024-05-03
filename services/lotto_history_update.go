package services

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
	"watcharis/go-migrate-lotto-history-els/models"
)

func (s *lottoHistoryServices) UpdateRewardStatusByElasticQuery(ctx context.Context) error {

	lottoID := 90124484109
	uuid := "4e43640c08a809f52dd95b55d4f2304d154421185268bf4b6a5dd7c381fee291"

	lotto, err := s.lottoRepository.GetLottoByID(ctx, lottoID)
	if err != nil {
		log.Printf("cannot get lotto by id [error]: %+v\n", err)
		return err
	}
	fmt.Printf("lotto_price_due : %+v\n", lotto.LottoPriceDue.Format(time.DateOnly))

	index := fmt.Sprintf(models.INDEX_POV, lotto.LottoPriceDue.Format(time.DateOnly))

	_id := int64(lottoID)
	source := `
		ctx._source.reward_status = 'test_update_reward_status_by_query'
	`

	if err := s.lottoHistoryElasticSearch.UpdateLottoHistoryRewardStatusByElasticQuery(ctx, index, uuid, _id, source); err != nil {
		log.Printf("cannot update reward_status id [error]: %+v\n", err)
		return err
	}
	return nil
}

func (s *lottoHistoryServices) UpdateRewardStatus(ctx context.Context) error {

	lottoID := 90124484109

	lotto, err := s.lottoRepository.GetLottoByID(ctx, lottoID)
	if err != nil {
		log.Printf("cannot get lotto by id [error]: %+v\n", err)
		return err
	}
	fmt.Printf("lotto_price_due : %+v\n", lotto.LottoPriceDue.Format(time.DateOnly))

	index := fmt.Sprintf(models.INDEX_POV, lotto.LottoPriceDue.Format(time.DateOnly))

	id := strconv.Itoa(lottoID)

	doc := map[string]interface{}{
		"reward_status": "test_update_reward_status",
	}

	if err := s.lottoHistoryElasticSearch.UpdateLottoHistoryRewardStatus(ctx, index, id, doc); err != nil {
		log.Printf("cannot update reward_status id [error]: %+v\n", err)
		return err
	}
	return nil
}

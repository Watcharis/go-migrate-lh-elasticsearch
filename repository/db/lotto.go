package db

import (
	"context"
	"watcharis/go-migrate-lotto-history-els/models"

	"gorm.io/gorm"
)

type lottoRepository struct {
	db *gorm.DB
}

type LottoRepository interface {
	GetLottoWithReward(ctx context.Context, uuid string) ([]models.LottoWithReward, error)
	GetLottoByID(ctx context.Context, id int) (models.Lottos, error)
	GetLottoSoldWithRewardSpeacificRoundDate(ctx context.Context, uuid string, lottoPriceDue string) ([]models.LottoWithReward, error)
}

func NewLottoRepository(db *gorm.DB) LottoRepository {
	return &lottoRepository{
		db: db,
	}
}

func (r *lottoRepository) GetLottoWithReward(ctx context.Context, uuid string) ([]models.LottoWithReward, error) {

	var lottoWithReward []models.LottoWithReward
	err := r.db.Debug().Raw(`
		select 
			lt.id, lt.lotto_number, lt.lotto_round, lt.lotto_set, lt.lotto_year, lt.lotto_item, lt.lotto_price_due, lt.lotto_price, 
			lt.scanner_uuid, lt.buyer_uuid, lt.lotto_url, lt.lotto_bc_ref, lt.lotto_bc_hash, lt.lotto_type, lt.lotto_status, 
			lt.create_at, lt.create_by, lt.update_at, lt.update_by, lt.payment_id, lt.lotto_uuid, lt.tags,
			CASE 
				WHEN rw.id IS NOT NULL THEN TRUE 
				ELSE FALSE
			END AS is_win,
			rw.reward_desc as win_desc,
			rw.reward_status
		from lottos lt
		left join rewards rw on lt.id = rw.id
		where lt.scanner_uuid = ? and lt.buyer_uuid = ? and lt.lotto_status = ?`, uuid, uuid, "transfered").
		Find(&lottoWithReward).Error

	if err != nil {
		return nil, err
	}

	return lottoWithReward, nil
}

func (r *lottoRepository) GetLottoByID(ctx context.Context, id int) (models.Lottos, error) {

	var lottos models.Lottos
	if err := r.db.Debug().Table("lottos l").Where("l.id = ?", id).First(&lottos).Error; err != nil {
		return models.Lottos{}, err
	}

	return lottos, nil
}

func (r *lottoRepository) GetLottoSoldWithRewardSpeacificRoundDate(ctx context.Context, uuid string, lottoPriceDue string) ([]models.LottoWithReward, error) {
	var lottoWithReward []models.LottoWithReward
	err := r.db.Debug().Raw(`
		select 
			lt.id, lt.lotto_number, lt.lotto_round, lt.lotto_set, lt.lotto_year, lt.lotto_item, lt.lotto_price_due, lt.lotto_price, 
			lt.scanner_uuid, lt.buyer_uuid, lt.lotto_url, lt.lotto_bc_ref, lt.lotto_bc_hash, lt.lotto_type, lt.lotto_status, 
			lt.create_at, lt.create_by, lt.update_at, lt.update_by, lt.payment_id, lt.lotto_uuid, lt.tags,
			CASE 
				WHEN rw.id IS NOT NULL THEN TRUE 
				ELSE FALSE
			END AS is_win,
			rw.reward_desc,
			rw.reward_status
		from lottos lt
		left join rewards rw on lt.id = rw.lotto_id
		WHERE lt.id in (
			select id from (
				select l.id from lottos l where l.lotto_price_due = ?
					and l.scanner_uuid = ?
					and l.lotto_status = ?
			) as temp_lotto_id
		)`, lottoPriceDue, uuid, "sold").
		Find(&lottoWithReward).Error

	if err != nil {
		return nil, err
	}

	return lottoWithReward, nil
}

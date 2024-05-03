package db

import (
	"context"
	"watcharis/go-migrate-lotto-history-els/models"

	"gorm.io/gorm"
)

type lottoHistoryRepository struct {
	db *gorm.DB
}

type LottoHistoryRepository interface {
	Begin(ctx context.Context) *gorm.DB
	Commit(ctx context.Context, tx *gorm.DB) error
	GeLottoRoundList(ctx context.Context) ([]models.LottoHistoryRound, error)
	GetLottoHistoryByLottoPriceDue(ctx context.Context, lottoPriceDue string) ([]models.LottoHistory, error)
	GetLottoHistoryTransaction(ctx context.Context, tx *gorm.DB, buyerUuid string) ([]models.LottoHistory, error)
	GetLottoHistoryWithoutTransaction(ctx context.Context, buyerUuid string) ([]models.LottoHistory, error)
	CountLottoHistoryByLottoPriceDue(ctx context.Context, lottoPriceDue string) (int, error)
	GetLottoHistoryWithLimitAndOffset(ctx context.Context, lottoPriceDue string, limit, offset int) ([]models.LottoHistory, error)
}

func NewLottoHistoryRepository(db *gorm.DB) LottoHistoryRepository {
	return &lottoHistoryRepository{
		db: db,
	}
}

func (r *lottoHistoryRepository) GeLottoRoundList(ctx context.Context) ([]models.LottoHistoryRound, error) {

	var lottoHistoryRound []models.LottoHistoryRound

	err := r.db.Debug().Table("lottos_history lh").
		Select("DISTINCT (lotto_price_due)").
		Where("lh.lotto_price_due not in ('0001-01-01 00:00:00', '0000-00-00 00:00:00')").
		Find(&lottoHistoryRound).Error
	if err != nil {
		return nil, err
	}
	return lottoHistoryRound, nil
}

func (r *lottoHistoryRepository) GetLottoHistoryByLottoPriceDue(ctx context.Context, lottoPriceDue string) ([]models.LottoHistory, error) {
	var lottoHistoryList []models.LottoHistory

	err := r.db.Debug().Table("lottos_history lh").Where("lh.lotto_price_due = ?", lottoPriceDue).Find(&lottoHistoryList).Error
	if err != nil {
		return nil, err
	}
	return lottoHistoryList, nil
}

func (r *lottoHistoryRepository) Begin(ctx context.Context) *gorm.DB {
	tx := r.db.Begin()
	return tx.WithContext(ctx)
}

func (r *lottoHistoryRepository) Commit(ctx context.Context, tx *gorm.DB) error {
	err := tx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}

func (r *lottoHistoryRepository) GetLottoHistoryTransaction(ctx context.Context, tx *gorm.DB, buyerUuid string) ([]models.LottoHistory, error) {
	var lottoHistoryList []models.LottoHistory
	err := tx.WithContext(ctx).Debug().Table("lottos_history lh").Where("lh.buyer_uuid = ?", buyerUuid).Order("lh.id desc").Limit(10).Find(&lottoHistoryList).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return lottoHistoryList, nil
}

func (r *lottoHistoryRepository) GetLottoHistoryWithoutTransaction(ctx context.Context, buyerUuid string) ([]models.LottoHistory, error) {
	var lottoHistoryList []models.LottoHistory
	err := r.db.WithContext(ctx).Debug().Table("lottos_history lh").Where("lh.buyer_uuid = ?", buyerUuid).Order("lh.id desc").Limit(10).Find(&lottoHistoryList).Error
	if err != nil {
		return nil, err
	}
	return lottoHistoryList, nil
}

func (r *lottoHistoryRepository) CountLottoHistoryByLottoPriceDue(ctx context.Context, lottoPriceDue string) (int, error) {
	var countLottoHistory int
	err := r.db.WithContext(ctx).Debug().Raw(`select count(1) from lottos_history lh where lh.lotto_price_due = ?`, lottoPriceDue).Scan(&countLottoHistory).Error
	if err != nil {
		return 0, err
	}
	return countLottoHistory, nil
}

func (r *lottoHistoryRepository) GetLottoHistoryWithLimitAndOffset(ctx context.Context, lottoPriceDue string, limit, offset int) ([]models.LottoHistory, error) {
	var lottoHistory []models.LottoHistory
	if err := r.db.WithContext(ctx).Debug().
		Table("lottos_history lh").Where("lh.lotto_price_due = ?", lottoPriceDue).Limit(limit).Offset(offset).Find(&lottoHistory).Error; err != nil {
		return nil, err
	}
	return lottoHistory, nil
}

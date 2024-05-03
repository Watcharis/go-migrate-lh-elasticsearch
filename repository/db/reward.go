package db

import (
	"context"
	"watcharis/go-migrate-lotto-history-els/models"

	"gorm.io/gorm"
)

type rewardRepository struct {
	db *gorm.DB
}

type RewardRepository interface {
	GetRewardByListLottoID(ctx context.Context, lottoID []int) ([]models.Rewards, error)
}

func NewRewardRepository(db *gorm.DB) RewardRepository {
	return &rewardRepository{
		db: db,
	}
}

func (r *rewardRepository) GetRewardByListLottoID(ctx context.Context, lottoID []int) ([]models.Rewards, error) {
	var rewardList []models.Rewards
	if err := r.db.WithContext(ctx).Table("rewards rw").
		Select(`
			rw.lotto_id,
			CASE WHEN rw.id IS NOT NULL THEN TRUE ELSE FALSE END AS is_win, 
			rw.reward_desc as win_desc, 
			rw.reward_status`).
		Where("rw.lotto_id in (?)", lottoID).Find(&rewardList).Error; err != nil {
		return nil, err
	}
	return rewardList, nil
}

func (r *rewardRepository) GetLottoIdFromRewards(ctx context.Context, buyerUuid string, limit int) ([]int, error) {

	rewardStatus := []string{"ban", "verified", "success"}

	var lottoID []int

	err := r.db.WithContext(ctx).Table("rewards rw").
		Select("rw.lotto_id").
		Where("rw.buyer_uuid = ?", buyerUuid).
		Where("(rw.claim_type = ? or rw.claim_type is null or rw.claim_type = '')", "GLO").
		Where("rw.reward_status in ?", rewardStatus).
		Limit(limit).
		Find(&lottoID).Error

	if err != nil {
		return nil, err
	}

	return lottoID, nil
}

//ปกติผมแปลง ท่านี้ time.ParseInLocation()

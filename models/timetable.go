package models

import "time"

type TimeTable struct {
	// Id                   int       `gorm:"column:id"`
	// Round                string    `gorm:"column:round"`
	RoundDate time.Time `gorm:"column:round_date"`
	// StartDate            time.Time `gorm:"column:start_date"`
	// RewardDate           time.Time `gorm:"column:reward_date"`
	// MarketOpenTime       string    `gorm:"column:market_open_time"`
	// MarketCloseTime      string    `gorm:"column:market_close_time"`
	// RewardTime           string    `gorm:"column:reward_time"`
	// ExpireNormalDate     time.Time `gorm:"column:expire_normal_date"`
	// ExpireCharityDate    time.Time `gorm:"column:expire_charity_date"`
	// IsDisable            string    `gorm:"column:is_disable"`
	// RewardStartTime      string    `gorm:"column:reward_start_time"`
	// RewardStartTimeP80   string    `gorm:"column:reward_start_time_p80"`
	// RewardTimeP80        string    `gorm:"column:reward_time_p80"`
	// CreateAt             time.Time `gorm:"column:create_at"`
	// ExpireDateChannelGLO time.Time `gorm:"column:expire_date_channel_glo"`
	// ExpireTimeChannelGLO string    `gorm:"column:expire_time_channel_glo"`
	// ExpireDateChannelKTB time.Time `gorm:"column:expire_date_channel_ktb"`
	// ExpireTimeChannelKTB string    `gorm:"column:expire_time_channel_ktb"`
}

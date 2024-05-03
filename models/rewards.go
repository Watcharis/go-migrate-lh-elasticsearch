package models

type Rewards struct {
	LottoID      int    `json:"lotto_id,omitempty" gorm:"column:lotto_id"`
	IsWin        int    `json:"is_win" gorm:"column:is_win"`
	WinDesc      string `json:"win_desc" gorm:"column:reward_desc"`
	RewardStatus string `json:"reward_status" gorm:"column:reward_status"`
}

package models

import "time"

type Lottos struct {
	ID            int       `json:"id,omitempty" gorm:"column:id"`
	LottoNumber   string    `json:"lt_number" gorm:"lotto_number"`
	LottoRound    string    `json:"lotto_round" gorm:"column:lotto_round"`
	LottoSet      string    `json:"lotto_set" gorm:"column:lotto_set"`
	LottoYear     string    `json:"lotto_year" gorm:"column:lotto_year"`
	LottoItem     string    `json:"lotto_item" gorm:"column:lotto_item"`
	LottoPriceDue time.Time `json:"lotto_price_due" gorm:"column:lotto_price_due"`
	LottoPrice    int       `json:"lotto_price" gorm:"column:lotto_price"`
	ScannerUUID   string    `json:"scanner_uuid" gorm:"column:scanner_uuid"`
	BuyerUUID     string    `json:"buyer_uuid" gorm:"buyer_uuid"`
	LottoURL      string    `json:"lotto_url" gorm:"column:lotto_url"`
	LottoBCRef    string    `json:"lotto_bc_ref" gorm:"column:lotto_bc_ref"`
	LottoBCHash   string    `json:"lotto_bc_hash" gorm:"column:lotto_bc_hash"`
	LottoType     string    `json:"lotto_type" gorm:"column:lotto_type"`
	LottoStatus   string    `json:"lotto_status" gorm:"column:lotto_status"`
	CreateAt      time.Time `json:"create_at" gorm:"column:create_at"`
	CreateBy      string    `json:"create_by" gorm:"column:create_by"`
	UpdateAt      time.Time `json:"update_at" gorm:"update_at"`
	UpdateBy      string    `json:"update_by" gorm:"update_by"`
	PaymentID     int64     `json:"payment_id" gorm:"payment_id"`
	LottoUUID     string    `json:"lotto_uuid" gorm:"column:lotto_uuid"`
	Tags          string    `json:"tags" gorm:"tags"`
}

type LottoWithReward struct {
	ID            int       `json:"id,omitempty" gorm:"column:id"`
	LottoNumber   string    `json:"lt_number" gorm:"lotto_number"`
	LottoRound    string    `json:"lotto_round" gorm:"column:lotto_round"`
	LottoSet      string    `json:"lotto_set" gorm:"column:lotto_set"`
	LottoYear     string    `json:"lotto_year" gorm:"column:lotto_year"`
	LottoItem     string    `json:"lotto_item" gorm:"column:lotto_item"`
	LottoPriceDue time.Time `json:"lotto_price_due" gorm:"column:lotto_price_due"`
	LottoPrice    int       `json:"lotto_price" gorm:"column:lotto_price"`
	ScannerUUID   string    `json:"scanner_uuid" gorm:"column:scanner_uuid"`
	BuyerUUID     string    `json:"buyer_uuid" gorm:"buyer_uuid"`
	LottoURL      string    `json:"lotto_url" gorm:"column:lotto_url"`
	LottoBCRef    string    `json:"lotto_bc_ref" gorm:"column:lotto_bc_ref"`
	LottoBCHash   string    `json:"lotto_bc_hash" gorm:"column:lotto_bc_hash"`
	LottoType     string    `json:"lotto_type" gorm:"column:lotto_type"`
	LottoStatus   string    `json:"lotto_status" gorm:"column:lotto_status"`
	CreateAt      time.Time `json:"create_at" gorm:"column:create_at"`
	CreateBy      string    `json:"create_by" gorm:"column:create_by"`
	UpdateAt      time.Time `json:"update_at" gorm:"update_at"`
	UpdateBy      string    `json:"update_by" gorm:"update_by"`
	PaymentID     int64     `json:"payment_id" gorm:"payment_id"`
	LottoUUID     string    `json:"lotto_uuid" gorm:"column:lotto_uuid"`
	Tags          string    `json:"tags" gorm:"column:tags"`
	IsWin         int       `json:"is_win" gorm:"column:is_win"`
	WinDesc       string    `json:"win_desc" gorm:"column:reward_desc"`
	RewardStatus  string    `json:"reward_status" gorm:"column:reward_status"`
}

type LottoSearch struct {
	ID []int `json:"id"`
}

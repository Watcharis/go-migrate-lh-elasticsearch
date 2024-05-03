package models

import "time"

const (
	INDEX     = "lotto_history.%s"
	INDEX_POC = "lotto_history_poc.%s"
	INDEX_POV = "lotto_history_pov.%s"

	LIMIT_LOTTO_HISTORY_QUERY_DB = 1000

	ELSATIC_SIZE_1000 = 1000
)

const (
	REDISKEY_CURRENTPERIOD = "CURRENTPERIOD"
)

const (
	DATETIME_FORMAT_PRICE_DUE_IN_DB     = "2006-01-02 15:04:05"
	DATETIME_FORMAT_PRICE_DUE_IN_REDIS  = "02012006"
	DATETIME_FORMAT_DATE_IN_RESPONSE_EN = "2 January 2006"
	DATETIME_FORMAT_DATE_IN_SLIP        = "2 January 2006 15:04:05"
	DATETIME_FORMAT_PRICE_DUE           = "2006-01-02"
	DATETIME_FORMAT_PRICE_DUE_NO_DAT    = "20060102"
)

type LottoHistory struct {
	ID            int       `gorm:"column:id" json:"id"`
	LottoNumber   string    `gorm:"column:lotto_number" json:"lotto_number"`
	LottoRound    string    `gorm:"column:lotto_round" json:"lotto_round"`
	LottoSet      string    `gorm:"column:lotto_set" json:"lotto_set"`
	LottoYear     string    `gorm:"column:lotto_year" json:"lotto_year"`
	LottoItem     string    `gorm:"column:lotto_item" json:"lotto_item"`
	LottoPriceDue time.Time `gorm:"column:lotto_price_due" json:"lotto_price_due"`
	LottoPrice    int       `gorm:"column:lotto_price" json:"lotto_price"`
	ScannerUUID   string    `gorm:"column:scanner_uuid" json:"scanner_uuid"`
	BuyerUUID     string    `gorm:"column:buyer_uuid" json:"buyer_uuid"`
	LottoURL      string    `gorm:"column:lotto_url" json:"lotto_url"`
	LottoBcRef    string    `gorm:"column:lotto_bc_ref" json:"lotto_bc_ref"`
	LottoBcHash   string    `gorm:"column:lotto_bc_hash" json:"lotto_bc_hash"`
	LottoType     string    `gorm:"column:lotto_type" json:"lotto_type"`
	LottoStatus   string    `gorm:"column:lotto_status" json:"lotto_status"`
	CreateAt      time.Time `gorm:"column:create_at" json:"create_at"`
	CreateBy      string    `gorm:"column:create_by" json:"create_by"`
	UpdateAt      time.Time `gorm:"column:update_at" json:"update_at"`
	UpdateBy      string    `gorm:"column:update_by" json:"update_by"`
	PaymentID     int       `gorm:"column:payment_id" json:"payment_id"`
	LottoUUID     string    `gorm:"column:lotto_uuid" json:"lotto_uuid"`
	Tags          string    `gorm:"column:tags" json:"tags"`
	IsWin         bool      `gorm:"column:is_win" json:"is_win"`
	WinDesc       string    `gorm:"column:win_desc" json:"win_desc"`
}

type LottoHistoryRound struct {
	LottoPriceDue time.Time `gorm:"column:lotto_price_due" json:"lotto_price_due"`
}

type LottoHistoryElasticSearchOld struct {
	ElasticId     string          `json:"-"`
	ID            int             `json:"id"`
	LottoNumber   string          `json:"lotto_number"`
	LottoRound    string          `json:"lotto_round"`
	LottoSet      string          `json:"lotto_set"`
	LottoYear     string          `json:"lotto_year"`
	LottoItem     string          `json:"lotto_item"`
	LottoPriceDue string          `json:"lotto_price_due"`
	LottoPrice    int             `json:"lotto_price"`
	ScannerUUID   string          `json:"scanner_uuid"`
	BuyerUUID     string          `json:"buyer_uuid"`
	LottoURL      string          `json:"lotto_url"`
	LottoBcRef    string          `json:"lotto_bc_ref"`
	LottoBcHash   string          `json:"lotto_bc_hash"`
	LottoType     string          `json:"lotto_type"`
	LottoStatus   string          `json:"lotto_status"`
	CreateAt      string          `json:"create_at"`
	CreateBy      string          `json:"create_by"`
	UpdateAt      string          `json:"update_at"`
	UpdateBy      string          `json:"update_by"`
	PaymentID     int             `json:"payment_id"`
	LottoUUID     string          `json:"lotto_uuid"`
	Tags          string          `json:"tags"`
	IsWin         bool            `json:"is_win"`
	WinDesc       []WinDescObject `json:"win_desc"`
}

type WinDescObject struct {
	Desc          string `json:"desc"`
	Price         string `json:"price"`
	IsFirstReward bool   `json:"is_firstreward"`
}

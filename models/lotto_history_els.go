package models

type RespElasticsearch struct {
	Took   int                      `json:"took"`
	Errors bool                     `json:"errors"`
	Items  []map[string]interface{} `json:"items"`
}

type ElasticResult struct {
	PitID    string `json:"pit_id,omitempty"`
	Took     uint64 `json:"took"`
	TimedOut bool   `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits         ElastictHits `json:"hits"`
	Aggregations struct {
		AvgRate struct {
			Val float64 `json:"value"`
		} `json:"avg_rate"`
	} `json:"aggregations"`
}

type ElastictHits struct {
	Total struct {
		Val int `json:"value"`
	} `json:"total"`
	MaxScore float32            `json:"max_score"`
	Hits     []ElastictHitsHits `json:"hits"`
}

type ElastictHitsHits struct {
	Index     string                     `json:"_index"`
	Type      string                     `json:"_type"`
	ID        string                     `json:"_id"`
	Score     float64                    `json:"_score"`
	Source    LottosHistoryElasticsearch `json:"_source"`
	Highlight map[string][]string        `json:"highlight,omitempty"`
	Sort      []int64                    `json:"sort"`
}

type ElasticResultPIT struct {
	NumFreed  int  `json:"num_freed"`
	Succeeded bool `json:"succeeded"`
}

// type LottosHistoryElasticsearchV3 struct {
// 	ID            int64      `json:"id,omitempty" gorm:"column:id"`
// 	LottoNumber   string     `json:"lotto_number" gorm:"column:lotto_number"`
// 	LottoDt1      int        `json:"lotto_dt1" gorm:"column:lotto_dt1"`
// 	LottoDt2      int        `json:"lotto_dt2" gorm:"column:lotto_dt2"`
// 	LottoDt3      int        `json:"lotto_dt3" gorm:"column:lotto_dt3"`
// 	LottoDt4      int        `json:"lotto_dt4" gorm:"column:lotto_dt4"`
// 	LottoDt5      int        `json:"lotto_dt5" gorm:"column:lotto_dt5"`
// 	LottoDt6      int        `json:"lotto_dt6" gorm:"column:lotto_dt6"`
// 	LottoRound    string     `json:"lotto_round" gorm:"column:lotto_round"`
// 	LottoSet      string     `json:"lotto_set" gorm:"column:lotto_set"`
// 	LottoYear     string     `json:"lotto_year" gorm:"column:lotto_year"`
// 	LottoItem     string     `json:"lotto_item" gorm:"column:lotto_item"`
// 	LottoPriceDue time.Time  `json:"lotto_price_due" gorm:"column:lotto_price_due"`
// 	LottoPrice    int        `json:"lotto_price" gorm:"column:lotto_price"`
// 	ScannerUUID   string     `json:"scanner_uuid" gorm:"column:scanner_uuid"`
// 	BuyerUUID     *string    `json:"buyer_uuid" gorm:"buyer_uuid"`
// 	LottoURL      string     `json:"lotto_url" gorm:"column:lotto_url"`
// 	LottoBCRef    string     `json:"lotto_bc_ref" gorm:"column:lotto_bc_ref"`
// 	LottoBCHash   string     `json:"lotto_bc_hash" gorm:"column:lotto_bc_hash"`
// 	LottoType     string     `json:"lotto_type" gorm:"column:lotto_type"`
// 	LottoStatus   string     `json:"lotto_status" gorm:"column:lotto_status"`
// 	CreateAt      *time.Time `json:"create_at" gorm:"column:create_at"`
// 	CreateBy      string     `json:"create_by" gorm:"column:create_by"`
// 	UpdateAt      *time.Time `json:"update_at" gorm:"update_at"`
// 	UpdateBy      *string    `json:"update_by" gorm:"update_by"`
// 	PaymentID     *int64     `json:"payment_id" gorm:"payment_id"`
// 	LottoUUID     *string    `json:"lotto_uuid" gorm:"column:lotto_uuid"`
// 	Tags          *string    `json:"tags" gorm:"tags"`
// 	IsWin         int        `json:"is_win" gorm:"is_win"`
// 	WinDesc       *string    `json:"win_desc" gorm:"win_desc"`
// 	RewardStatus  *string    `json:"reward_status" gorm:"reward_status"`
// }

// type LottosHistoryElasticsearch struct {
// 	ID            int64     `json:"id,omitempty" gorm:"column:id"`
// 	LottoNumber   string    `json:"lotto_number" gorm:"column:lotto_number"`
// 	LottoDt1      int       `json:"lotto_dt1" gorm:"column:lotto_dt1"`
// 	LottoDt2      int       `json:"lotto_dt2" gorm:"column:lotto_dt2"`
// 	LottoDt3      int       `json:"lotto_dt3" gorm:"column:lotto_dt3"`
// 	LottoDt4      int       `json:"lotto_dt4" gorm:"column:lotto_dt4"`
// 	LottoDt5      int       `json:"lotto_dt5" gorm:"column:lotto_dt5"`
// 	LottoDt6      int       `json:"lotto_dt6" gorm:"column:lotto_dt6"`
// 	LottoRound    string    `json:"lotto_round" gorm:"column:lotto_round"`
// 	LottoSet      string    `json:"lotto_set" gorm:"column:lotto_set"`
// 	LottoYear     string    `json:"lotto_year" gorm:"column:lotto_year"`
// 	LottoItem     string    `json:"lotto_item" gorm:"column:lotto_item"`
// 	LottoPriceDue time.Time `json:"lotto_price_due" gorm:"column:lotto_price_due"`
// 	LottoPrice    int       `json:"lotto_price" gorm:"column:lotto_price"`
// 	ScannerUUID   string    `json:"scanner_uuid" gorm:"column:scanner_uuid"`
// 	BuyerUUID     string    `json:"buyer_uuid" gorm:"buyer_uuid"`
// 	LottoURL      string    `json:"lotto_url" gorm:"column:lotto_url"`
// 	LottoBCRef    string    `json:"lotto_bc_ref" gorm:"column:lotto_bc_ref"`
// 	LottoBCHash   string    `json:"lotto_bc_hash" gorm:"column:lotto_bc_hash"`
// 	LottoType     string    `json:"lotto_type" gorm:"column:lotto_type"`
// 	LottoStatus   string    `json:"lotto_status" gorm:"column:lotto_status"`
// 	CreateAt      time.Time `json:"create_at" gorm:"column:create_at"`
// 	CreateBy      string    `json:"create_by" gorm:"column:create_by"`
// 	UpdateAt      time.Time `json:"update_at" gorm:"update_at"`
// 	UpdateBy      string    `json:"update_by" gorm:"update_by"`
// 	PaymentID     int64     `json:"payment_id" gorm:"payment_id"`
// 	LottoUUID     string    `json:"lotto_uuid" gorm:"column:lotto_uuid"`
// 	Tags          string    `json:"tags" gorm:"tags"`
// 	IsWin         int       `json:"is_win" gorm:"is_win"`
// 	WinDesc       string    `json:"win_desc" gorm:"win_desc"`
// 	RewardStatus  string    `json:"reward_status" gorm:"reward_status"`
// }

type LottosHistoryElasticsearch struct {
	ID               int64     `json:"id,omitempty"`
	LottoNumber      string    `json:"lotto_number"`
	LottoDt1         int       `json:"lotto_dt1"`
	LottoDt2         int       `json:"lotto_dt2"`
	LottoDt3         int       `json:"lotto_dt3"`
	LottoDt4         int       `json:"lotto_dt4"`
	LottoDt5         int       `json:"lotto_dt5"`
	LottoDt6         int       `json:"lotto_dt6"`
	LottoRound       string    `json:"lotto_round"`
	LottoSet         string    `json:"lotto_set"`
	LottoYear        string    `json:"lotto_year"`
	LottoItem        string    `json:"lotto_item"`
	LottoPriceDue    string    `json:"lotto_price_due"`
	LottoPrice       int       `json:"lotto_price"`
	ScannerUUID      string    `json:"scanner_uuid"`
	BuyerUUID        string    `json:"buyer_uuid"`
	LottoURL         string    `json:"lotto_url"`
	LottoBCRef       string    `json:"lotto_bc_ref"`
	LottoBCHash      string    `json:"lotto_bc_hash"`
	LottoType        string    `json:"lotto_type"`
	LottoStatus      string    `json:"lotto_status"`
	PurchaseDateTime string    `json:"purchase_datetime"`
	PurchaseDate     string    `json:"purchase_date"`
	CreateAt         string    `json:"create_at"`
	CreateBy         string    `json:"create_by"`
	UpdateAt         string    `json:"update_at"`
	UpdateBy         string    `json:"update_by"`
	PaymentID        int64     `json:"payment_id"`
	LottoUUID        string    `json:"lotto_uuid"`
	Tags             []int     `json:"tags"`
	IsWin            int       `json:"is_win"`
	WinDesc          []WinDesc `json:"win_desc"`
	RewardStatus     string    `json:"reward_status"`
}

type WinDesc struct {
	Description   string `json:"description"`
	Reward        string `json:"reward"`
	IsFirstReward bool   `json:"is_firstreward"`
}

package services

func validateLottoPriceDue(lottoRound []string, lottoPriceDue string) bool {
	for _, lr := range lottoRound {
		if lr == lottoPriceDue {

			return true
		}
	}
	return false
}

func IsEmptyString(v string) bool {
	return v == ""
}

func LottoHistoryElasticSearchMapping() string {
	return `{
		"mappings": {
			"dynamic": false,
			"properties": {
				"id": {
					"type": "long"
				},
				"lotto_price": {
					"type": "long"
				},
				"lotto_uuid": {
					"type": "keyword"
				},
				"scanner_uuid": {
					"type": "keyword"
				},
				"buyer_uuid": {
					"type": "keyword"
				},
				"lotto_url": {
					"type": "keyword"
				},
				"lotto_price_due": {
					"type": "date",
					"format": "yyyy-MM-dd"
				},
				"lotto_number": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword"
						}
					}
				},
				"lotto_dt1": {
					"type": "long"
				},
				"lotto_dt2": {
					"type": "long"
				},
				"lotto_dt3": {
					"type": "long"
				},
				"lotto_dt4": {
					"type": "long"
				},
				"lotto_dt5": {
					"type": "long"
				},
				"lotto_dt6": {
					"type": "long"
				},
				"lotto_round": {
					"type": "keyword"
				},
				"lotto_set": {
					"type": "keyword"
				},
				"lotto_year": {
					"type": "keyword"
				},
				"lotto_item": {
					"type": "keyword"
				},
				"lotto_bc_ref": {
					"type": "text"
				},
				"lotto_bc_hash": {
					"type": "keyword"
				},
				"lotto_type": {
					"type": "keyword"
				},
				"lotto_status": {
					"type": "keyword"
				},
				"payment_id": {
					"type": "long"
				},
				"tags": {
					"type": "long"
				},
				"is_win": {
					"type": "integer"
				},
				"win_desc": {
					"type": "nested",
					"properties": {
						"description": {
							"type": "text"
						},
						"reward": {
							"type": "scaled_float",
							"scaling_factor": 100
						},
						"is_firstreward": {
							"type": "boolean"
						}
					}
				},
				"reward_status": {
					"type": "keyword"
				},
				"purchase_date": {
					"type": "date",
					"format": "yyyy-MM-dd"
				},
				"purchase_datetime": {
					"type": "date",
					"format": "yyyy-MM-dd HH:mm:ss"
				},
				"create_by": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword"
						}
					}
				},
				"update_by": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword"
						}
					}
				},
				"create_at": {
					"type": "date",
					"format": "yyyy-MM-dd HH:mm:ss"
				},
				"update_at": {
					"type": "date",
					"format": "yyyy-MM-dd HH:mm:ss"
				}
			}
		}
	}`
}

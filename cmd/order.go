package wotwhb

import (
	"time"
)

type HbOrder struct {
	AmountSpent  float64     `json:"amount_spent"`
	Claimed      bool        `json:"claimed"`
	Created      time.Time   `json:"created"`
	Currency     string      `json:"currency"`
	Gamekey      string      `json:"gamekey"`
	IsGiftee     bool        `json:"is_giftee"`
	MissedCredit interface{} `json:"missed_credit"`
	PathIds      interface{} `json:"path_ids"`
	Product      interface{} `json:"product"`
	Subproducts  interface{} `json:"subproducts"`
	Total        float64     `json:"total"`
	Uid          string      `json:"uid"`
}

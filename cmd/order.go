package wotwhb

import (
	"time"
)

type HbPathIds []string

type HbProduct struct {
	Category           string      `json:"category"`
	EmptyTpkds         interface{} `json:"empty_tpkds"`
	HumanName          string      `json:"human_name"`
	MachineName        string      `json:"machine_name"`
	PartialGiftEnabled bool        `json:"partial_gift_enabled"`
	PostPurchaseText   string      `json:"post_purchase_text"`
}

type HbPayee struct {
	HumanName   string `json:"human_name"`
	MachineName string `json:"machine_name"`
}

type HbSubProduct struct {
	CustomDownloadPageBox string      `json:"custom_download_page_box_html"`
	Downloads             interface{} `json:"downloads"`
	HumanName             string      `json:"human_name"`
	Icon                  string      `json:"icon"`
	LibraryFamilyName     string      `json:"library_family_name"`
	MachineName           string      `json:"machine_name"`
	Payee                 HbPayee     `json:"payee"`
	Url                   string      `json:"url"`
}

type HbOrder struct {
	AmountSpent  float64        `json:"amount_spent"`
	Claimed      bool           `json:"claimed"`
	Created      time.Time      `json:"created"`
	Currency     string         `json:"currency"`
	Gamekey      string         `json:"gamekey"`
	IsGiftee     bool           `json:"is_giftee"`
	MissedCredit interface{}    `json:"missed_credit"`
	PathIds      HbPathIds      `json:"path_ids"`
	Product      HbProduct      `json:"product"`
	SubProducts  []HbSubProduct `json:"subproducts"`
	Total        float64        `json:"total"`
	Uid          string         `json:"uid"`
}

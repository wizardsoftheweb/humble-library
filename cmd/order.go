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

type HbSubProductDownload struct {
	AndroidAppOnly        bool        `json:"android_app_only"`
	DownloadIdentifer     string      `json:"download_identifier"`
	DownloadStruct        interface{} `json:"download_struct"`
	DownloadVersionNumber int         `json:"download_version_number"`
	MachineName           string      `json:"machine_name"`
	OptionsDict           interface{} `json:"options_dict"`
	Platform              string      `json:"platform"`
}

type HbSubProduct struct {
	CustomDownloadPageBox string                 `json:"custom_download_page_box_html"`
	Downloads             []HbSubProductDownload `json:"downloads"`
	HumanName             string                 `json:"human_name"`
	Icon                  string                 `json:"icon"`
	LibraryFamilyName     string                 `json:"library_family_name"`
	MachineName           string                 `json:"machine_name"`
	Payee                 HbPayee                `json:"payee"`
	Url                   string                 `json:"url"`
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

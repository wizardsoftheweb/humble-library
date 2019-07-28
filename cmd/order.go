package wotwhb

import (
	"time"
)

// jq '[.[] | select(0 < (.subproducts | length)).subproducts[] | select(0 < (.downloads | length)).downloads[] | select(0 < (.download_struct | length)).download_struct[] | select(.asm_config != null).asm_config | keys] | flatten | unique' all-orders.json
type HbSubProductDownloadDownloadStructAsmConfig struct {
	CloudMountPoint string `json:"cloudMountPoint"`
	DisplayItem     string `json:"display_item"`
	WarnCrash       bool   `json:"warnCrash"`
}

// q '[.[] | select(0 < (.subproducts | length)).subproducts[] | select(0 < (.downloads | length)).downloads[] | select(0 < (.download_struct | length)).download_struct[] | select(.asm_manifest != null).asm_manifest | keys] | flatten | unique' all-orders.json
type HbSubProductDownloadDownloadStructAsmManifest interface{}

// jq '[.[] | select(0 < (.subproducts | length)).subproducts[] | select(0 < (.downloads | length)).downloads[] | select(0 < (.download_struct | length)).download_struct[] | select(.url != null).url | keys] | flatten | unique' all-orders.json
type HbSubProductDownloadDownloadStructUrl struct {
	BitTorrent string `json:"bittorrent"`
	Web        string `json:"web"`
}

// jq '[.[].subproducts[].downloads[].download_struct[] | keys] | flatten | unique' all-orders.json
type HbSubProductDownloadDownloadStruct struct {
	Arch             string                                        `json:"arch"`
	AsmConfig        HbSubProductDownloadDownloadStructAsmConfig   `json:"asm_config"`
	AsmManifest      HbSubProductDownloadDownloadStructAsmManifest `json:"asm_manifest"`
	ExternalLink     string                                        `json:"external_link"`
	FileSize         int64                                         `json:"file_size"`
	ForceDownload    bool                                          `json:"force_download"`
	HdStreamUrl      string                                        `json:"hd_stream_url"`
	HumanSize        string                                        `json:"human_size"`
	KindleFriendly   bool                                          `json:"kindle_friendly"`
	Md5              string                                        `json:"md5"`
	Name             string                                        `json:"name"`
	SdStreamUrl      string                                        `json:"sd_stream_url"`
	Sha1             string                                        `json:"sha1"`
	Small            int                                           `json:"small"`
	Timestamp        int64                                         `json:"timestamp"`
	Timetstamp       int64                                         `json:"timetstamp"`
	UploadedAt       time.Time                                     `json:"uploaded_at"`
	Url              HbSubProductDownloadDownloadStructUrl         `json:"url"`
	UsesKindleSender bool                                          `json:"uses_kindle_sender"`
}

// jq '[.[] | select(0 < (.subproducts | length)).subproducts[] | select(0 < (.downloads | length)).downloads[] | select(.options_dict != null).options_dict | keys] | flatten | unique' all-orders.json
type HbSubProductDownloadOptionsDict struct {
	Is64BitToggle int `json:"is64bittoggle"`
}

// jq '[.[].subproducts[].downloads[] | keys] | flatten | unique' all-orders.json
type HbSubProductDownload struct {
	AndroidAppOnly        bool                                 `json:"android_app_only"`
	DownloadIdentifer     string                               `json:"download_identifier"`
	DownloadStruct        []HbSubProductDownloadDownloadStruct `json:"download_struct"`
	DownloadVersionNumber int                                  `json:"download_version_number"`
	MachineName           string                               `json:"machine_name"`
	OptionsDict           HbSubProductDownloadOptionsDict      `json:"options_dict"`
	Platform              string                               `json:"platform"`
}

// jq '[.[].subproducts[].payee | keys] | flatten | unique' all-orders.json
type HbPayee struct {
	HumanName   string `json:"human_name"`
	MachineName string `json:"machine_name"`
}

type HbPathIds []string

// jq '[.[].product | keys] | flatten | unique' all-orders.json
type HbProduct struct {
	Category string `json:"category"`
	// Haven't found a pattern to this yet
	EmptyTpkds          interface{} `json:"empty_tpkds"`
	HumanName           string      `json:"human_name"`
	MachineName         string      `json:"machine_name"`
	PartialGiftEnabled  bool        `json:"partial_gift_enabled"`
	PostPurchaseText    string      `json:"post_purchase_text"`
	SubscriptionCredits int         `json:"subscription_credits"`
}

// jq '[.[].subproducts[] | keys] | flatten | unique' all-orders.json
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

// jq '[.[] | keys] | flatten | unique' all-orders.json
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

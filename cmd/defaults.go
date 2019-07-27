package wotwhb

import (
	"path/filepath"
)

const (
	cookieFileBasename       = "cookies.json"
	orderKeyListFileBasename = "order-keys.json"
	allOrdersFileBasename    = "all-orders.json"
)

var (
	configDirectory   = filepath.Clean(filepath.Join("~", ".config", "wotw", "humblebundle"))
	downloadDirectory = filepath.Clean(filepath.Join("~", "Downloads"))
)

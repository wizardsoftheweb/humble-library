package wotwhb

import (
	"os"
	"path/filepath"
)

const (
	cookieFileBasename       = "cookies.json"
	orderKeyListFileBasename = "order-keys.json"
	allOrdersFileBasename    = "all-orders.json"
)

var (
	homeDirectory     = os.Getenv("HOME")
	configDirectory   = filepath.Clean(filepath.Join(homeDirectory, ".config", "wotw", "humblebundle"))
	downloadDirectory = filepath.Clean(filepath.Join(homeDirectory, "Downloads"))
)

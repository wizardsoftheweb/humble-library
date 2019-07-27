package wotwhb

import (
	"os"
)

func ensureDirectoryExists(directory string) {
	_ = os.MkdirAll(directory, os.ModePerm)
}

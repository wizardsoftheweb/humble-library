package wotwhb

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func ensureDirectoryExists(directory string) {
	_ = os.MkdirAll(directory, os.ModePerm)
}

func writeJsonToFile(contents interface{}, fileName string) {
	fileContents, err := json.Marshal(contents)
	fatalCheck(err)
	err = ioutil.WriteFile(
		fileName,
		fileContents,
		0644,
	)
	fatalCheck(err)
}

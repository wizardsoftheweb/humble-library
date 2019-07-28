package wotwhb

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"

	"github.com/spyzhov/ajson"
)

func queryOrderList(printer CanPrint) []byte {
	client, jar := buildSession()
	return getResource(printer, client, jar, keysResource, url.Values{}, nil)
}

func parseRawOrderList(rawResponse []byte) []string {
	nodes, err := ajson.JSONPath(rawResponse, "$..gamekey")
	fatalCheck(err)
	keys := make([]string, len(nodes))
	for index, node := range nodes {
		keys[index] = node.MustString()
	}
	return keys
}

func updateOrderList(printer CanPrint) []string {
	rawResponse := queryOrderList(printer)
	keys := parseRawOrderList(rawResponse)
	writeJsonToFile(keys, filepath.Join(ConfigDirectoryFlagValue, orderKeyListFileBasename))
	return keys
}

func loadSavedOrderList() []string {
	contents, err := ioutil.ReadFile(filepath.Join(ConfigDirectoryFlagValue, orderKeyListFileBasename))
	fatalCheck(err)
	root, err := ajson.Unmarshal(contents)
	fatalCheck(err)
	fmt.Println(root)
	return []string{"cool"}
}

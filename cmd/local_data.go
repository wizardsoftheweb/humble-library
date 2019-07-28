package wotwhb

import (
	"net/url"
	"path/filepath"

	"github.com/spyzhov/ajson"
)

func queryOrderList(printer CanPrint) []byte {
	client, jar := buildSession()
	return getResource(printer, client, jar, keysResource, url.Values{}, nil)
}

func parseOrderList(rawResponse []byte) []string {
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
	keys := parseOrderList(rawResponse)
	writeJsonToFile(keys, filepath.Join(ConfigDirectoryFlagValue, orderKeyListFileBasename))
	return keys
}

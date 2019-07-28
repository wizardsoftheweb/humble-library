package wotwhb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"

	cookiejar "github.com/juju/persistent-cookiejar"
	"github.com/spyzhov/ajson"
)

func queryKeyList(printer CanPrint) []byte {
	client, jar := buildSession()
	return getResource(printer, client, jar, keysResource, url.Values{}, nil)
}

func parseRawKeyList(rawResponse []byte) []string {
	nodes, err := ajson.JSONPath(rawResponse, "$..gamekey")
	fatalCheck(err)
	keys := make([]string, len(nodes))
	for index, node := range nodes {
		keys[index] = node.MustString()
	}
	return keys
}

func updateKeyList(printer CanPrint) []string {
	rawResponse := queryKeyList(printer)
	keys := parseRawKeyList(rawResponse)
	writeJsonToFile(keys, filepath.Join(ConfigDirectoryFlagValue, orderKeyListFileBasename))
	return keys
}

func loadSavedKeyList() []string {
	contents, err := ioutil.ReadFile(filepath.Join(ConfigDirectoryFlagValue, orderKeyListFileBasename))
	fatalCheck(err)
	var keys []string
	err = json.Unmarshal(contents, &keys)
	fatalCheck(err)
	return keys
}

func queryIndividualOrder(printer CanPrint, client HttpClient, jar *cookiejar.Jar, key string) []byte {
	return getResource(printer, client, jar, orderResource+key, url.Values{}, nil)
}

func queryAllOrders(printer CanPrint, keys []string) []map[string]interface{} {
	client, jar := buildSession()
	size := len(keys)
	allOrders := make([]map[string]interface{}, size)
	for index, key := range keys {
		var parsedOrder map[string]interface{}
		switch {
		case 0 == index%50:
			Logger.Info(fmt.Sprintf("Queried %d out of %d", index, size))
			break
		case 0 == index%10:
			Logger.Trace(fmt.Sprintf("Queried %d out of %d", index, size))
			break
		case 0 == index%5:
			Logger.Debug(fmt.Sprintf("Queried %d out of %d", index, size))
			break
		}
		rawOrder := queryIndividualOrder(printer, client, jar, key)
		err := json.Unmarshal(rawOrder, &parsedOrder)
		fatalCheck(err)
		allOrders[index] = parsedOrder
	}
	return allOrders
}

func updateAllOrders(printer CanPrint, keys []string) {
	rawOrders := queryAllOrders(printer, keys)
	writeJsonToFile(rawOrders, filepath.Join(ConfigDirectoryFlagValue, allOrdersFileBasename))
}

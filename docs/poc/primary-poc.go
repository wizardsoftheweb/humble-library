package poc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"

	// "net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	cookiejar "github.com/juju/persistent-cookiejar"
	"github.com/sirupsen/logrus"
	"github.com/spyzhov/ajson"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

const (
	baseUrl           = "https://hr-humblebundle.appspot.com"
	loginResource     = "/processlogin"
	keysResource      = "/api/v1/user/order"
	orderResource     = "/api/v1/order/"
	baseDomain        = "hr-humblebundle.appspot.com"
	cookieFile        = "cookies.json"
	orderKeyListFile  = "order-keys.json"
	allOrdersFile     = "all-orders.json"
	downloadDirectory = "downloads"
)

func BootstrapLogger(verbosityLevel int) {
	formatter := &prefixed.TextFormatter{
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
		QuoteEmptyFields: true,
		ForceFormatting:  true,
		ForceColors:      true,
	}
	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle:     "cyan",
		TimestampStyle:  "black+h",
		DebugLevelStyle: "blue+h:",
		InfoLevelStyle:  "green+h:",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "red+b",
		PanicLevelStyle: "red+B",
	})
	logrus.SetFormatter(formatter)
	switch {
	case -2 >= verbosityLevel:
		logrus.SetLevel(logrus.PanicLevel)
		break
	case -1 == verbosityLevel:
		logrus.SetLevel(logrus.FatalLevel)
		break
	case 0 == verbosityLevel:
		logrus.SetLevel(logrus.ErrorLevel)
		break
	case 1 == verbosityLevel:
		logrus.SetLevel(logrus.WarnLevel)
		break
	case 2 == verbosityLevel:
		logrus.SetLevel(logrus.InfoLevel)
		break
	case 3 == verbosityLevel:
		logrus.SetLevel(logrus.TraceLevel)
		break
	default:
		logrus.SetLevel(logrus.DebugLevel)
		break
	}

}

type GuardNeededResponse struct {
	Required bool     `json:"humble_guard_required"`
	Code     []string `json:"skip_code"`
}

type SuccessfulLoginResponse struct {
	UserTerms interface{} `json:"user_terms_opt_in_data"`
	Goto      string      `json:"goto"`
}

type OrderObject struct {
	Key string `json:"gamekey"`
}

type OrderResponse []OrderObject

func getInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", prompt)
	input, err := reader.ReadString('\n')
	if nil != err {
		logrus.Fatal(err)
	}
	return strings.TrimSpace(input)
}

func parseResponseBody(body io.ReadCloser) string {
	contents, err := ioutil.ReadAll(body)
	if nil != err {
		logrus.Fatal(err)
	}
	return string(contents)
}

func sanitizeCookieString(cookieString string) string {
	cleanedValue := strings.ReplaceAll(cookieString, `\075`, "=")
	cleanedValue = strings.ReplaceAll(cleanedValue, `"`, "")
	return cleanedValue
}

func parseRawCookie(rawCookie string) *http.Cookie {
	explodedCookie := strings.Split(rawCookie, ";")
	splitValue := strings.Split(explodedCookie[0], "=")
	newCookie := &http.Cookie{
		Name:   splitValue[0],
		Value:  sanitizeCookieString(splitValue[1]),
		Raw:    sanitizeCookieString(rawCookie),
		Domain: baseDomain,
	}
	for _, value := range explodedCookie[1:] {
		splitValue := strings.Split(strings.TrimSpace(value), "=")
		switch strings.ToLower(splitValue[0]) {
		case "max-age":
			newCookie.MaxAge, _ = strconv.Atoi(splitValue[1])
		case "path":
			newCookie.Path = splitValue[1]
		case "secure":
			newCookie.Secure = true
		case "httponly":
			newCookie.HttpOnly = true
		case "domain":
			newCookie.Domain = splitValue[1]
		}
	}
	return newCookie
}

func updateCookies(jar *cookiejar.Jar, response *http.Response) {
	rawCsrfCookies := response.Header["Set-Cookie"]
	parsedCsrfCookies := make([]*http.Cookie, len(rawCsrfCookies))
	for index, rawCookie := range rawCsrfCookies {
		parsedCsrfCookies[index] = parseRawCookie(rawCookie)
	}
	jar.SetCookies(
		response.Request.URL,
		append(
			jar.Cookies(response.Request.URL),
			parsedCsrfCookies...,
		),
	)
	err := jar.Save()
	if nil != err {
		logrus.Fatal(err)
	}
}

func makeCaptchaRequest(client *http.Client, jar *cookiejar.Jar, csrfCookie *http.Cookie) string {
	data := url.Values{}
	data.Set("ajax", "true")
	data.Set("username", os.Getenv("HB_USERNAME"))
	data.Set("password", os.Getenv("HB_PASSWORD"))
	data.Set("recaptcha_challenge_field", "")
	data.Set("recaptcha_response_field", getInput("Enter Recaptcha"))
	return getResource(client, jar, loginResource, data, csrfCookie)
}

func discoverSkipCode(guardResponse string) string {
	var parsedGuardResponse GuardNeededResponse
	err := json.Unmarshal([]byte(guardResponse), &parsedGuardResponse)
	if nil != err {
		logrus.Fatal(err)
	}
	return parsedGuardResponse.Code[0]
}

func makeGuardRequest(client *http.Client, jar *cookiejar.Jar, csrfCookie *http.Cookie, skipCode string) string {
	data := url.Values{}
	data.Set("ajax", "true")
	data.Set("username", os.Getenv("HB_USERNAME"))
	data.Set("password", os.Getenv("HB_PASSWORD"))
	data.Set("captcha-skip-code", skipCode)
	data.Set("guard", getInput("Enter guard code"))
	return getResource(client, jar, loginResource, data, csrfCookie)
}

func wasLoginSuccessful(responseBody string) bool {
	var parsedResponse SuccessfulLoginResponse
	err := json.Unmarshal([]byte(responseBody), &parsedResponse)
	if nil != err {
		return false
	}
	return "/home" == parsedResponse.Goto
}

func authenticate(client *http.Client, jar *cookiejar.Jar, csrfCookie *http.Cookie) {
	initialResponse := makeCaptchaRequest(client, jar, csrfCookie)
	skipCode := discoverSkipCode(initialResponse)
	finalResponse := makeGuardRequest(client, jar, csrfCookie, skipCode)
	if !wasLoginSuccessful(finalResponse) {
		logrus.Fatal("unable to log in")
	}
}

func getResource(client *http.Client, jar *cookiejar.Jar, resource string, data url.Values, csrfCookie *http.Cookie) string {
	var request *http.Request
	var err error
	if 0 == len(data) {
		request, err = http.NewRequest("GET", baseUrl+resource, nil)
	} else {
		request, err = http.NewRequest("POST", baseUrl+resource, strings.NewReader(data.Encode()))
	}
	if nil != err {
		logrus.Fatal(err)
	}
	if 0 < len(data) {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if nil != csrfCookie {
		request.Header.Set("csrf-prevention-token", csrfCookie.Value)
	}
	request.Header.Set("X-Requested-By", "hb_android_app")
	response, err := client.Do(request)
	if nil != err {
		logrus.Fatal(err)
	}
	logrus.Debug(
		fmt.Sprintf(
			"%s->%s status: %s\n",
			request.URL.Path,
			response.Request.URL.Path,
			response.Status,
		),
	)
	body := parseResponseBody(response.Body)
	err = response.Body.Close()
	if nil != err {
		logrus.Fatal(err)
	}
	if request.URL.Path == response.Request.URL.Path {
		updateCookies(jar, response)
	} else {
		for _, cookie := range response.Cookies() {
			if "csrf_cookie" == cookie.Name {
				jar.SetCookies(
					response.Request.URL,
					append(
						jar.Cookies(response.Request.URL),
						cookie,
					),
				)
				authenticate(client, jar, cookie)
				return getResource(client, jar, resource, data, csrfCookie)
			}
		}
		logrus.Fatal("No CSRF Token")
	}
	return body
}

func parseOrderKeyList(responseBody string) []string {
	var orderResponse OrderResponse
	err := json.Unmarshal([]byte(responseBody), &orderResponse)
	if nil != err {
		logrus.Fatal(err)
	}
	orderKeys := make([]string, len(orderResponse))
	for index, value := range orderResponse {
		orderKeys[index] = value.Key
	}
	return orderKeys
}

func writeJsonToFile(contents interface{}, fileName string) {
	fileContents, err := json.Marshal(contents)
	if nil != err {
		logrus.Fatal(err)
	}
	err = ioutil.WriteFile(
		fileName,
		fileContents,
		0644,
	)
	if nil != err {
		logrus.Fatal(err)
	}
}

func updateOrderKeyList(client *http.Client, jar *cookiejar.Jar) []string {
	rawResponse := getResource(client, jar, keysResource, url.Values{}, nil)
	orderKeys := parseOrderKeyList(rawResponse)
	writeJsonToFile(orderKeys, orderKeyListFile)
	return orderKeys
}

func updateAllOrders(client *http.Client, jar *cookiejar.Jar, orderKeys []string) []map[string]interface{} {
	allOrders := make([]map[string]interface{}, len(orderKeys))
	for index, key := range orderKeys {
		if 0 == index%10 {
			logrus.Info(fmt.Sprintf("Parsed %d of %d orders", index, len(orderKeys)))
		}
		allOrders[index] = updateIndividualOrder(client, jar, key)
	}
	writeJsonToFile(allOrders, allOrdersFile)
	return allOrders
}

func updateIndividualOrder(client *http.Client, jar *cookiejar.Jar, key string) map[string]interface{} {
	rawResponse := getResource(client, jar, orderResource+key, url.Values{}, nil)
	var result map[string]interface{}
	err := json.Unmarshal([]byte(rawResponse), &result)
	if nil != err {
		logrus.Fatal(err)
	}
	return result
}

func loadAllOrdersAsNodes() *ajson.Node {
	contents, err := ioutil.ReadFile(allOrdersFile)
	if nil != err {
		logrus.Fatal(err)
	}
	// return contents
	root, err := ajson.Unmarshal(contents)
	if nil != err {
		logrus.Fatal(err)
	}
	return root
}

func loadAllOrdersAsBytes() []byte {
	contents, err := ioutil.ReadFile(allOrdersFile)
	if nil != err {
		logrus.Fatal(err)
	}
	return contents
}

type FileTypeReturn struct {
	HumanName   string
	MachineName string
	FileUrl     string
	FileSize    int64
}

func NewFileTypeReturn(humanName, machinName, fileUrl string, fileSize int64) *FileTypeReturn {
	return &FileTypeReturn{
		HumanName:   humanName,
		MachineName: machinName,
		FileUrl:     fileUrl,
		FileSize:    fileSize,
	}
}

func humanReadableFileSize(size int64) string {
	fmt.Println(size)
	const unit = 1000
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	dividend, index := int64(unit), 0
	for factor := size / unit; factor >= unit; factor /= unit {
		dividend *= unit
		index++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(dividend), "kMGTPE"[index])
}

func searchForFileType(fileType string) []*FileTypeReturn {
	allOrders := loadAllOrdersAsNodes()
	root := allOrders.MustArray()
	totalSize := int64(0)
	var results []*FileTypeReturn
	for _, orders := range root {
		if orders.IsObject() && orders.HasKey("subproducts") && 0 < len(orders.MustKey("subproducts").MustArray()) {
			subproducts := orders.MustKey("subproducts").MustArray()
			for _, subproduct := range subproducts {
				if subproduct.IsObject() && subproduct.HasKey("human_name") && subproduct.HasKey("machine_name") && subproduct.HasKey("downloads") && 0 < len(subproduct.MustKey("downloads").MustArray()) {
					downloads := subproduct.MustKey("downloads").MustArray()
					for _, download := range downloads {
						if download.IsObject() && download.HasKey("download_struct") && 0 < len(download.MustKey("download_struct").MustArray()) {
							downloadStructs := download.MustKey("download_struct").MustArray()
							for _, downloadStruct := range downloadStructs {
								if downloadStruct.IsObject() && downloadStruct.HasKey("name") && fileType == downloadStruct.MustKey("name").MustString() && downloadStruct.HasKey("file_size") && downloadStruct.MustKey("file_size").IsNumeric() {
									if downloadStruct.HasKey("url") && downloadStruct.MustKey("url").IsObject() && downloadStruct.MustKey("url").HasKey("web") {
										humanName := subproduct.MustKey("human_name").MustString()
										machineName := subproduct.MustKey("machine_name").MustString()
										fileSize := int64(downloadStruct.MustKey("file_size").MustNumeric())
										results = append(
											results,
											NewFileTypeReturn(
												humanName,
												fmt.Sprintf("%s.%s", machineName, strings.ToLower(fileType)),
												strings.ReplaceAll(downloadStruct.MustKey("url").MustKey("web").MustString(), `\u0026`, "&"),
												fileSize,
											),
										)
										totalSize += fileSize
									}
								}
							}
						}
					}
				}
			}
		}
	}
	logrus.Info(fmt.Sprintf("Total file size is %s", humanReadableFileSize(totalSize)))
	return results
}

func uniqueStringSlice(stringSlice []string) []string {
	existingStrings := make(map[string]bool)
	uniqueStrings := []string{}
	for _, element := range stringSlice {
		if _, ok := existingStrings[element]; !ok {
			existingStrings[element] = true
			uniqueStrings = append(uniqueStrings, element)
		}
	}
	return uniqueStrings
}

func discoverFileTypes() {
	allOrders := loadAllOrdersAsBytes()
	nodes, err := ajson.JSONPath(allOrders, "$..subproducts..downloads..download_struct..name")
	if nil != err {
		logrus.Fatal(err)
	}
	existingStrings := make(map[string]bool)
	var uniqueStrings []string
	for _, node := range nodes {
		element := node.MustString()
		if _, ok := existingStrings[element]; !ok {
			existingStrings[element] = true
			uniqueStrings = append(uniqueStrings, element)
		}
	}
	sort.Strings(uniqueStrings)
	fmt.Println(uniqueStrings)
}

type DownloadProgress struct {
	Downloaded int64
	Total      int64
}

func (d *DownloadProgress) Write(contents []byte) (int, error) {
	chunkSize := len(contents)
	d.Downloaded += int64(chunkSize)
	d.Update()
	return chunkSize, nil
}

func (d *DownloadProgress) Update() {
	fmt.Printf("\r%s", strings.Repeat(" ", 40))
	doneRatio := float64(d.Downloaded) / float64(d.Total)
	fmt.Printf(
		"\r%3.f%% [%s%s]",
		doneRatio*100,
		strings.Repeat(`█`, int(math.Floor(doneRatio*20))),
		strings.Repeat(` `, int(20-math.Floor(doneRatio*20))),
	)
}

func downloadFile(fileInfo *FileTypeReturn) {
	logrus.Info(fmt.Sprintf("Downloading %s", fileInfo.HumanName))
	fileName := filepath.Join(downloadDirectory, fileInfo.MachineName)
	output, err := os.Create(fileName + ".download")
	if nil != err {
		logrus.Fatal(err)
	}
	defer (func() { _ = output.Close() })()
	response, err := http.Get(fileInfo.FileUrl)
	if nil != err {
		logrus.Fatal(err)
	}
	defer (func() { _ = response.Body.Close() })()
	progress := &DownloadProgress{
		Total: fileInfo.FileSize,
	}
	_, err = io.Copy(output, io.TeeReader(response.Body, progress))
	if nil != err {
		logrus.Fatal(err)
	}
	fmt.Print("\n")
	err = os.Rename(fileName+".download", fileName)
	if nil != err {
		logrus.Fatal(err)
	}
}

func ensureDownloadDirectoryExists() {
	_ = os.MkdirAll(downloadDirectory, os.ModePerm)
}

// func main() {
// 	BootstrapLogger(4)
// 	ensureDownloadDirectoryExists()
// 	_ = godotenv.Load()
// 	discoverFileTypes()
// 	// results := searchForFileType("EPUB")
// 	// for _, value := range results[0:1] {
// 	// 	downloadFile(value)
// 	// }
// 	// options := cookiejar.Options{
// 	// 	PublicSuffixList: publicsuffix.List,
// 	// 	Filename:         "cookies.json",
// 	// }
// 	// jar, err := cookiejar.New(&options)
// 	// if nil != err {
// 	// 	logrus.Fatal(err)
// 	// }
// 	// fmt.Println(jar)
// 	// client := &http.Client{
// 	// 	Jar: jar,
// 	// }
// 	// orderKeys := updateOrderKeyList(client, jar)
// 	// updateAllOrders(client, jar, orderKeys)
// 	// updateIndividualOrder(client, jar, orderKeys[0])
// }

package wotwhb

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	cookiejar "github.com/juju/persistent-cookiejar"
	"github.com/spyzhov/ajson"
)

type HttpClient interface {
	Do(request *http.Request) (*http.Response, error)
}

type ResponseGuardNeeded struct {
	Required bool     `json:"humble_guard_required"`
	Code     []string `json:"skip_code"`
}

type ResponseLoginSuccessful struct {
	UserTerms interface{} `json:"user_terms_opt_in_data"`
	Goto      string      `json:"goto"`
}

func createNewRequest(resource string, data url.Values, csrfCookie *http.Cookie) *http.Request {
	var request *http.Request
	var err error
	if 0 == len(data) {
		request, err = http.NewRequest("GET", baseUrl+resource, nil)
	} else {
		request, err = http.NewRequest("POST", baseUrl+resource, strings.NewReader(data.Encode()))
	}
	fatalCheck(err)
	if 0 < len(data) {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	request.Header.Set("X-Requested-By", "hb_android_app")
	if nil != csrfCookie {
		request.Header.Set("csrf-prevention-token", csrfCookie.Value)
	}
	return request
}

func parseResponseBody(body io.Reader) []byte {
	contents, err := ioutil.ReadAll(body)
	fatalCheck(err)
	return contents
}

func executeRequest(client HttpClient, request *http.Request) (*http.Response, []byte) {
	response, err := client.Do(request)
	fatalCheck(err)
	Logger.Debug(
		fmt.Sprintf(
			"%s->%s status: %s\n",
			request.URL.Path,
			response.Request.URL.Path,
			response.Status,
		),
	)
	body := parseResponseBody(response.Body)
	err = response.Body.Close()
	fatalCheck(err)
	return response, body
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
	fatalCheck(err)
}

func loginWithRecaptcha(printer CanPrint, client HttpClient, jar *cookiejar.Jar, csrfCookie *http.Cookie) (string, string, []byte) {
	username := getInput(printer, usernameLongPrompt, usernameShortPrompt)
	password := getInput(printer, passwordLongPrompt, passwordShortPrompt)
	data := url.Values{}
	data.Set("ajax", "true")
	data.Set("username", username)
	data.Set("password", password)
	data.Set("recaptcha_challenge_field", "")
	data.Set("recaptcha_response_field", getInput(printer, recaptchaLongPrompt, recaptchaShortPrompt))
	return username, password, getResource(printer, client, jar, loginResource, data, csrfCookie)
}

func loginWithGuard(printer CanPrint, client HttpClient, jar *cookiejar.Jar, csrfCookie *http.Cookie, username, password, skipCode string) []byte {
	data := url.Values{}
	data.Set("ajax", "true")
	data.Set("username", username)
	data.Set("password", password)
	data.Set("captcha-skip-code", skipCode)
	data.Set("guard", getInput(printer, guardLongPrompt, guardShortPrompt))
	return getResource(printer, client, jar, loginResource, data, csrfCookie)
}

func discoverCsrfCookie(response *http.Response, jar *cookiejar.Jar) *http.Cookie {
	for _, cookie := range response.Cookies() {
		if "csrf_cookie" == cookie.Name {
			jar.SetCookies(
				response.Request.URL,
				append(
					jar.Cookies(response.Request.URL),
					cookie,
				),
			)
			return cookie
		}
	}
	fatalCheck(errors.New("unable to find CSRF Token"))
	return nil
}

func authenticate(printer CanPrint, client HttpClient, jar *cookiejar.Jar, csrfCookie *http.Cookie) {
	username, password, guardResponse := loginWithRecaptcha(printer, client, jar, csrfCookie)
	guardRoot, err := ajson.Unmarshal(guardResponse)
	fatalCheck(err)
	skipCode := guardRoot.MustKey("skip_code").MustArray()[0].MustString()
	finalResponse := loginWithGuard(printer, client, jar, csrfCookie, username, password, skipCode)
	finalRoot, err := ajson.Unmarshal(finalResponse)
	fatalCheck(err)
	finalRoot.MustKey("goto").MustString()
}

func getResource(printer CanPrint, client HttpClient, jar *cookiejar.Jar, resource string, data url.Values, csrfCookie *http.Cookie) []byte {
	request := createNewRequest(resource, data, csrfCookie)
	response, body := executeRequest(client, request)
	if request.URL.Path == response.Request.URL.Path {
		updateCookies(jar, response)
	} else {
		csrfCookie = discoverCsrfCookie(response, jar)
		authenticate(printer, client, jar, csrfCookie)
		return getResource(printer, client, jar, resource, data, csrfCookie)
	}
	return body
}

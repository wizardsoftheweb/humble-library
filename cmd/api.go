package wotwhb

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	cookiejar "github.com/juju/persistent-cookiejar"
	"github.com/sirupsen/logrus"
)

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

func parseResponseBody(body io.ReadCloser) []byte {
	contents, err := ioutil.ReadAll(body)
	fatalCheck(err)
	return contents
}

func executeRequest(client *http.Client, request *http.Request) (*http.Response, []byte) {
	response, err := client.Do(request)
	fatalCheck(err)
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
	fatalCheck(err)
	return response, body
}

func updateCookies(jar *cookiejar.Jar, response *http.Response) {}

func loginWithRecaptcha(client *http.Client, jar *cookiejar.Jar, csrfCookie *http.Cookie) string {
	return ""
}

func loginWithGuard(client *http.Client, jar *cookiejar.Jar, csrfCookie *http.Cookie, skipCode string) string {
	return ""
}

func authenticate(client *http.Client, jar *cookiejar.Jar, csrfCookie *http.Cookie) {}

func getResource(client *http.Client, jar *cookiejar.Jar, resource string, data url.Values, csrfCookie *http.Cookie) []byte {
	request := createNewRequest(resource, data, csrfCookie)
	response, body := executeRequest(client, request)
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
		Logger.Fatal("No CSRF Token")
	}
	return body
}

package wotwhb

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	cookiejar "github.com/juju/persistent-cookiejar"
	. "gopkg.in/check.v1"
)

type ApiSuite struct {
	BaseSuite
	HttpClient HttpClient
	BadClient  HttpClient
	resource   string
	jar        *cookiejar.Jar
	postData   url.Values
	csrfCookie *http.Cookie
}

const testBody = `{"test":"body"}`
const rawCookie = "test=value; Max-Age=1564282717; Path=/; Secure; HttpOnly; Domain=qqq"
const recaptchaInputTestResults = `test1
test2
test3
test3
`

func (s *ApiSuite) Printf(format string, args ...interface{}) {}

type GoodHttpClientMock struct{}

func (h GoodHttpClientMock) Do(request *http.Request) (*http.Response, error) {
	response := &http.Response{
		Request:    request,
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader(testBody)),
	}
	return response, nil
}

type BadHttpClientMock struct{}

func (h BadHttpClientMock) Do(request *http.Request) (*http.Response, error) {
	response := &http.Response{
		Request: &http.Request{
			URL: &url.URL{
				Path: "wrong",
			},
		},
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader(testBody)),
		Header:     http.Header{},
	}
	return response, nil
}

var _ = Suite(&ApiSuite{})

func (s *ApiSuite) SetUpTest(c *C) {
	s.resource = loginResource
	s.postData = url.Values{}
	s.postData.Set("test", "value")
	s.csrfCookie = &http.Cookie{
		Name:  "csrf_cookie",
		Value: "qqq",
	}
	s.jar, _ = cookiejar.New(&cookiejar.Options{})
	s.HttpClient = GoodHttpClientMock{}
	s.BadClient = BadHttpClientMock{}
	inputReader = (func() io.Reader { return strings.NewReader(recaptchaInputTestResults) })()
}

func (s *ApiSuite) TestCreateNewGetRequest(c *C) {
	request := createNewRequest(s.resource, url.Values{}, nil)
	c.Assert(request.Method, Equals, "GET")
}

func (s *ApiSuite) TestCreateNewPostRequest(c *C) {
	request := createNewRequest(s.resource, s.postData, nil)
	c.Assert(request.Method, Equals, "POST")
	request = createNewRequest(s.resource, s.postData, s.csrfCookie)
	value := request.Header.Get("csrf-prevention-token")
	c.Assert(value, Equals, s.csrfCookie.Value)
}

func (s *ApiSuite) TestParseBody(c *C) {
	results := parseResponseBody(strings.NewReader(testBody))
	c.Assert(results, DeepEquals, []byte(testBody))
}

func (s *ApiSuite) TestExecuteResponse(c *C) {
	request := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Path: s.resource,
		},
	}
	response, body := executeRequest(s.HttpClient, request)
	c.Assert(body, DeepEquals, []byte(testBody))
	c.Assert(response.Request.URL.Path, Equals, request.URL.Path)
}

func (s *ApiSuite) TestDiscoverCsrfCookie(c *C) {
	response := &http.Response{
		Request: &http.Request{
			URL: &url.URL{
				Path: s.resource,
			},
		},
		Header: http.Header{},
	}
	c.Assert(func() {
		discoverCsrfCookie(response, s.jar)
	},
		PanicMatches,
		".*CSRF.*",
	)
	response.Header.Set("Set-Cookie", s.csrfCookie.String())
	result := discoverCsrfCookie(response, s.jar)
	c.Assert(result.Value, Equals, s.csrfCookie.Value)
}

func (s *ApiSuite) TestLoginWithRecaptcha(c *C) {
	c.Assert(
		func() {
			loginWithRecaptcha(s, s.HttpClient, s.jar, s.csrfCookie)
		},
		PanicMatches,
		".*EOF.*",
	)
}

func (s *ApiSuite) TestLoginWithGuard(c *C) {
	result := loginWithGuard(s, s.HttpClient, s.jar, s.csrfCookie, "qqq", "qqq", "qqq")
	c.Assert(result, DeepEquals, []byte(testBody))
}

func (s *ApiSuite) TestAuthenticate(c *C) {
	c.Assert(
		func() {
			authenticate(s, s.HttpClient, s.jar, s.csrfCookie)
		},
		PanicMatches,
		".*EOF.*",
	)
}

func (s *ApiSuite) TestGetResource(c *C) {
	result := getResource(s, s.HttpClient, s.jar, s.resource, s.postData, s.csrfCookie)
	c.Assert(result, NotNil)
	c.Assert(func() {
		getResource(s, s.BadClient, s.jar, s.resource, s.postData, s.csrfCookie)
	},
		PanicMatches,
		".*CSRF.*",
	)
}

func (s *ApiSuite) TestSanitizeCookie(c *C) {
	c.Assert(
		sanitizeCookieString(`"qqq\075"`),
		Equals,
		"qqq=",
	)
}

func (s *ApiSuite) TestParseRawCookie(c *C) {
	parsedCookie := parseRawCookie(rawCookie)
	c.Assert(parsedCookie.Name, Equals, "test")
	c.Assert(parsedCookie.Value, Equals, "value")
	c.Assert(parsedCookie.MaxAge, Equals, 1564282717)
	c.Assert(parsedCookie.Path, Equals, "/")
	c.Assert(parsedCookie.Secure, Equals, true)
	c.Assert(parsedCookie.HttpOnly, Equals, true)
	c.Assert(parsedCookie.Domain, Equals, "qqq")
}

func (s *ApiSuite) TestUpdateCookies(c *C) {
	response := &http.Response{
		Request: &http.Request{
			URL: &url.URL{
				Scheme: "https",
				Host:   baseDomain,
				Path:   s.resource,
			},
		},
		Header: http.Header{},
	}
	response.Header.Set("Set-Cookie", s.csrfCookie.String())
	updateCookies(s.jar, response)
	c.Assert(len(s.jar.AllCookies()), Equals, 1)
}

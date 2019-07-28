package wotwhb

import (
	"net/http"
	"net/url"
	"strings"

	. "gopkg.in/check.v1"
)

type ApiSuite struct {
	BaseSuite
	resource   string
	postData   url.Values
	csrfCookie *http.Cookie
}

const testBody = `{"test":"body"}`

var _ = Suite(&ApiSuite{})

func (s *ApiSuite) SetUpTest(c *C) {
	s.resource = loginResource
	s.postData = url.Values{}
	s.postData.Set("test", "value")
	s.csrfCookie = &http.Cookie{
		Name:  "csrf_cookie",
		Value: "qqq",
	}
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

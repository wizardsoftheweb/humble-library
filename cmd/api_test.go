package wotwhb

import (
	"fmt"
	"net/http"
	"net/url"

	. "gopkg.in/check.v1"
)

type ApiSuite struct {
	BaseSuite
	resource   string
	postData   url.Values
	csrfCookie *http.Cookie
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
}

func (s *ApiSuite) TestCreateNewGetRequest(c *C) {
	request := createNewRequest(s.resource, url.Values{}, nil)
	c.Assert(request.Method, Equals, "GET")
}

func (s *ApiSuite) TestCreateNewPostRequest(c *C) {
	request := createNewRequest(s.resource, s.postData, nil)
	c.Assert(request.Method, Equals, "POST")
	request = createNewRequest(s.resource, s.postData, s.csrfCookie)
	fmt.Println(request.Cookies())
	value := request.Header.Get("csrf-prevention-token")
	c.Assert(value, Equals, s.csrfCookie.Value)
}

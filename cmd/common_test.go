package wotwhb

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type BaseSuite struct {
	WorkingDir string
	HttpClient HttpClient
}

var _ = Suite(&BaseSuite{})

func (s *BaseSuite) SetUpSuite(c *C) {
	s.WorkingDir = c.MkDir()
	ConfigDirectoryFlagValue = filepath.Join(s.WorkingDir, "config")
	DownloadDirectoryFlagValue = filepath.Join(s.WorkingDir, "downloads")
	fatalHandler = func(args ...interface{}) { panic(args[0]) }
	s.HttpClient = HttpClientMock{}
}

type HttpClientMock struct{}

func (h HttpClientMock) Do(request *http.Request) (*http.Response, error) {
	response := &http.Response{
		Request:    request,
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader(testBody)),
	}
	return response, nil
}

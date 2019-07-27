package wotwhb

import (
	"net/http"
	"path/filepath"

	cookiejar "github.com/juju/persistent-cookiejar"
	"golang.org/x/net/publicsuffix"
)

func BootstrapConfig(configPath, downloadPath string) {
	ensureDirectoryExists(configPath)
	ensureDirectoryExists(downloadPath)
}

func buildSession() (*http.Client, *cookiejar.Jar) {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
		Filename:         filepath.Join(ConfigDirectoryFlagValue, cookieFileBasename),
	}
	jar, err := cookiejar.New(&options)
	if nil != err {
		Logger.Fatal(err)
	}
	Logger.Trace(jar.AllCookies())
	client := &http.Client{
		Jar: jar,
	}
	return client, jar
}

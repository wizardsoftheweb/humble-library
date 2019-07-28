package wotwhb

import (
	"path/filepath"
)

const (
	cookieFileBasename       = "cookies.json"
	orderKeyListFileBasename = "order-keys.json"
	allOrdersFileBasename    = "all-orders.json"
)

var (
	configDirectory   = filepath.Clean(filepath.Join("~", ".config", "wotw", "humblebundle"))
	downloadDirectory = filepath.Clean(filepath.Join("~", "Downloads"))
)

const (
	usernameLongPrompt   = "Please enter your Humble Username/Email."
	usernameShortPrompt  = "Username"
	passwordLongPrompt   = "Please enter your Humble Password."
	passwordShortPrompt  = "Password"
	recaptchaLongPrompt  = "Please visit https://www.schiff.io/projects/humble-bundle-api#getting-a-captcha to learn how to get a Recaptcha code. Once you've got one, enter it here."
	recaptchaShortPrompt = "Recaptcha Code"
	guardLongPrompt      = "Please check your email for a Humble Bundle account protection code."
	guardShortPrompt     = "Code"
)

package wotwhb

import (
	"os"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	formatter = &prefixed.TextFormatter{
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
		QuoteEmptyFields: true,
		ForceFormatting:  true,
		ForceColors:      true,
	}
	formatterColorScheme = &prefixed.ColorScheme{
		PrefixStyle:     "cyan",
		TimestampStyle:  "black+h",
		DebugLevelStyle: "blue+h:",
		InfoLevelStyle:  "green+h:",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "red+b",
		PanicLevelStyle: "red+B",
	}
	Logger = &logrus.Logger{
		Out:       os.Stderr,
		Formatter: formatter,
	}
)

func setLoggerLevel(verbosityLevel int) {
	switch {
	case -2 >= verbosityLevel:
		Logger.SetLevel(logrus.PanicLevel)
		break
	case -1 == verbosityLevel:
		Logger.SetLevel(logrus.FatalLevel)
		break
	case 0 == verbosityLevel:
		Logger.SetLevel(logrus.ErrorLevel)
		break
	case 1 == verbosityLevel:
		Logger.SetLevel(logrus.WarnLevel)
		break
	case 2 == verbosityLevel:
		Logger.SetLevel(logrus.InfoLevel)
		break
	case 3 == verbosityLevel:
		Logger.SetLevel(logrus.TraceLevel)
		break
	default:
		Logger.SetLevel(logrus.DebugLevel)
		break
	}
}

func BootstrapLogger(verbosityLevel int) {
	formatter.SetColorScheme(formatterColorScheme)
	setLoggerLevel(verbosityLevel)
}

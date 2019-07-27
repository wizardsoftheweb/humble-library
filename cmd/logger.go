package wotwhb

import (
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
)

func setLoggerLevel(verbosityLevel int) {
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

func BootstrapLogger(verbosityLevel int) {
	formatter.SetColorScheme(formatterColorScheme)
	logrus.SetFormatter(formatter)
	setLoggerLevel(verbosityLevel)
}

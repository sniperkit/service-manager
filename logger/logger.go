package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
)

const RootPackageName = "root"

var (
	mux       sync.RWMutex
	formatter logrus.Formatter
	loggers   = make(map[string]*logrus.Logger)

	defaultLogLevel = logrus.InfoLevel
)

// Configuration is the logging configuration.
// Example yml configuration is the following
//  log:
//    format: json
//    level:
//      root: error
//      github.com/Peripli/service-manager/server: debug
type Configuration struct {
	Format    string            `mapstructure:"format"`
	LogLevels map[string]string `mapstructure:"level"`
}

func Get(packageName string) *logrus.Logger {
	packageLowercase := strings.ToLower(packageName)
	mux.Lock()
	defer mux.Unlock()

	logger, exists := loggers[packageLowercase]
	if exists {
		return logger
	}
	logger = newLogger(defaultLogLevel)
	loggers[packageLowercase] = logger
	return logger
}

func Init(configuration Configuration) {
	formatter = formatterForName(configuration.Format)
	loggers = initializeLoggers(configuration.LogLevels)
	defaultLogLevel = parseLogLevel(configuration.LogLevels[RootPackageName])
}

func initializeLoggers(packageLogLevels map[string]string) map[string]*logrus.Logger {
	result := make(map[string]*logrus.Logger)
	for packageName, logLevel := range packageLogLevels {
		packageLowercase := strings.ToLower(packageName)
		levelForPackage := parseLogLevel(logLevel)
		if configuredLogger, dup := result[packageLowercase]; dup {
			if configuredLogger.Level != levelForPackage {
				panic(fmt.Sprintf("Ambiguous log level configuration for package %s ", packageName))
			}
		} else {
			result[packageLowercase] = newLogger(levelForPackage)
		}
	}
	return result
}

func formatterForName(formatName string) logrus.Formatter {
	switch formatName {
	case "json":
		return &logrus.JSONFormatter{}
	case "text":
		fallthrough
	default:
		return &logrus.TextFormatter{}
	}
}

func parseLogLevel(logLevel string) logrus.Level {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		fmt.Printf("Cannot parse log level %s. Falling back to %s...", logLevel, defaultLogLevel)
		level = defaultLogLevel
	}
	return level
}

func newLogger(level logrus.Level) *logrus.Logger {
	return &logrus.Logger{
		Out:       os.Stderr,
		Level:     level,
		Hooks:     make(logrus.LevelHooks),
		Formatter: formatter,
	}
}

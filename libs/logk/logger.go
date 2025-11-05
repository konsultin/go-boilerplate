package logk

import (
	"fmt"
	stdLog "log"
	"os"
	"sync"

	"github.com/Konsultin/project-goes-here/libs/logk/level"
	logkOption "github.com/Konsultin/project-goes-here/libs/logk/option"
)

type Logger interface {
	// Fatal must write an error, message that explaining the error and where it's occurred in FATAL level.
	Fatal(msg string, options ...logkOption.SetterFunc)

	// Fatalf must write a formatted message and where it's occurred in FATAL level.
	Fatalf(format string, args ...interface{})

	// Error must write an error, message that explaining the error and where it's occurred in ERROR level.
	Error(msg string, options ...logkOption.SetterFunc)

	// Errorf must write a formatted message and where it's occurred in ERROR level.
	Errorf(format string, args ...interface{})

	// Warn must write a message in WARN level.
	Warn(msg string, options ...logkOption.SetterFunc)

	// Warnf must write a formatted message in WARN level.
	Warnf(format string, args ...interface{})

	// Info must write a message in INFO level.
	Info(msg string, options ...logkOption.SetterFunc)

	// Infof must write a formatted message in INFO level.
	Infof(format string, args ...interface{})

	// Debug must write a message in DEBUG level.
	Debug(msg string, options ...logkOption.SetterFunc)

	// Debugf must write a formatted message in DEBUG level.
	Debugf(format string, args ...interface{})

	// Trace must write a message in TRACE level.
	Trace(msg string, options ...logkOption.SetterFunc)

	// Tracef must write a formatted message in TRACE level.
	Tracef(format string, args ...interface{})

	// NewChild must create a child logger and inherit level, writer and other flags
	// only option such as namespace could be overridden
	NewChild(args ...logkOption.SetterFunc) Logger
}

var log Logger
var logMutex sync.RWMutex

// Get retrieve logger instance and will fallback to StdLogger if no logger registered
func Get() Logger {
	// If log is nil, initiate standard logger
	if log == nil {
		// Get logger from env
		logLevelStr, _ := os.LookupEnv(EnvLogLevel)
		logLevel := level.Parse(logLevelStr)

		// Get logger prefix
		namespace, _ := os.LookupEnv(EnvLogNamespace)

		// Init standard logger
		p := NewStdLogPrinter(os.Stdout, stdLog.LstdFlags)
		l := NewStdLogger(p, logkOption.Level(logLevel), logkOption.WithNamespace(namespace))

		// Register logger
		Register(l)
		log.Trace("No logger found. StdLogger initiated")
	}
	return log
}

func NewChild(args ...logkOption.SetterFunc) Logger {
	// Get parent logger
	logger := Get()
	return logger.NewChild(args...)
}

// Register a logger implementation instance
func Register(l Logger) {
	// If logger is nil, return error
	if l == nil {
		panic(fmt.Errorf("%s: logger to be registered is nil", pkgName))
	}

	// Set logger
	logMutex.Lock()
	defer logMutex.Unlock()
	log = l
}

// Clear logger implementation instance
func Clear() {
	// Set logger
	logMutex.Lock()
	defer logMutex.Unlock()
	log = nil
}

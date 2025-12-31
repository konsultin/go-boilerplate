package logk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	stdLog "log"
	"os"

	logkContext "github.com/konsultin/project-goes-here/libs/logk/context"
	"github.com/konsultin/project-goes-here/libs/logk/level"
	logkOption "github.com/konsultin/project-goes-here/libs/logk/option"
)

var stdLevelPrefix = map[level.LogLevel]string{
	level.Fatal: "[FATAL] > ",
	level.Error: "[ERROR] > ",
	level.Warn:  "[WARN]  > ",
	level.Info:  "[INFO]  > ",
	level.Debug: "[DEBUG] > ",
	level.Trace: "[TRACE] > ",
}

type StdLogger struct {
	level     level.LogLevel
	printer   Printer
	namespace string
	ctx       context.Context
}

func (l *StdLogger) Fatal(msg string, args ...logkOption.SetterFunc) {
	l.print(level.Fatal, msg, logkOption.Evaluate(args))
}

func (l *StdLogger) Fatalf(format string, args ...interface{}) {
	l.print(level.Fatal, format, logkOption.NewFormatOptions(args...))
}

func (l *StdLogger) Error(msg string, args ...logkOption.SetterFunc) {
	l.print(level.Error, msg, logkOption.Evaluate(args))
}

func (l *StdLogger) Errorf(format string, args ...interface{}) {
	l.print(level.Error, format, logkOption.NewFormatOptions(args...))
}

func (l *StdLogger) Warn(msg string, args ...logkOption.SetterFunc) {
	l.print(level.Warn, msg, logkOption.Evaluate(args))
}

func (l *StdLogger) Warnf(format string, args ...interface{}) {
	l.print(level.Warn, format, logkOption.NewFormatOptions(args...))
}

func (l *StdLogger) Info(msg string, args ...logkOption.SetterFunc) {
	l.print(level.Info, msg, logkOption.Evaluate(args))
}

func (l *StdLogger) Infof(format string, args ...interface{}) {
	l.print(level.Info, format, logkOption.NewFormatOptions(args...))
}

func (l *StdLogger) Debug(msg string, args ...logkOption.SetterFunc) {
	l.print(level.Debug, msg, logkOption.Evaluate(args))
}

func (l *StdLogger) Debugf(format string, args ...interface{}) {
	l.print(level.Debug, format, logkOption.NewFormatOptions(args...))
}

func (l *StdLogger) Trace(msg string, args ...logkOption.SetterFunc) {
	l.print(level.Trace, msg, logkOption.Evaluate(args))
}

func (l *StdLogger) Tracef(format string, args ...interface{}) {
	l.print(level.Trace, format, logkOption.NewFormatOptions(args...))
}

func (l *StdLogger) NewChild(args ...logkOption.SetterFunc) Logger {
	options := logkOption.Evaluate(args)

	// Override namespace if option is set
	namespace, _ := logkOption.GetString(options, logkOption.NamespaceKey)

	// If not set and parent has namespace, then use parent namespace
	if namespace == "" && l.namespace != "" {
		args = append(args, logkOption.WithNamespace(l.namespace))
	}

	// Override level arguments
	args = append(args, logkOption.Level(l.level))

	// Initiate new logger
	cl := NewStdLogger(l.printer, args...)

	// Set context if available
	if ctx := options.Context; ctx != nil {
		cl.ctx = ctx
	}

	return cl
}

func (l *StdLogger) print(outLevel level.LogLevel, msg string, options *logkOption.Options) {
	// if output level is greater than log level, don't print
	if outLevel > l.level {
		return
	}

	// Inject context if not set
	if l.ctx != nil && options.Context == nil {
		options.Context = l.ctx
	}

	l.printer.Print(l.namespace, outLevel, msg, options)
}

func NewStdLogger(printer Printer, args ...logkOption.SetterFunc) *StdLogger {
	// Init standard logger instance
	l := StdLogger{}

	// Evaluate options
	o := logkOption.Evaluate(args)

	// Set level
	l.level = o.Level

	// Get namespace
	if namespace, _ := logkOption.GetString(o, logkOption.NamespaceKey); namespace != "" {
		l.namespace = namespace
	}

	// Get context
	if ctx := o.Context; ctx != nil {
		l.ctx = ctx
	}

	// Init printer if nil
	if printer == nil {
		l.printer = NewStdLogPrinter(os.Stdout, stdLog.LstdFlags)
	} else {
		l.printer = printer
	}

	return &l
}

func NewStdLogPrinter(out io.Writer, flag int) *stdLogPrinter {
	// If writer is nil, set default writer to Stdout
	if out == nil {
		out = os.Stdout
	}

	// Init log.Logger
	writer := stdLog.New(out, "", flag)

	return &stdLogPrinter{writer: writer}
}

type stdLogPrinter struct {
	writer *stdLog.Logger
}

func (s *stdLogPrinter) Print(namespace string, lv level.LogLevel, msg string, options *logkOption.Options) {
	writer := s.writer

	// Generate prefix
	prefix := stdLevelPrefix[lv]

	// Append namespace
	if namespace != "" {
		prefix = fmt.Sprintf("%s(%s) ", prefix, namespace)
	}

	// If options is existed
	// If formatted arguments is available, then print as formatted
	fmtArgs := options.FmtArgs
	if len(fmtArgs) > 0 {
		writer.Printf(prefix+msg+"\n", fmtArgs...)
	} else {
		writer.Printf("%s%s\n", prefix, msg)
	}

	// Get request id
	if reqId := logkContext.GetRequestId(options.Context); reqId != "" {
		writer.Printf("  > Request ID: %s\n", reqId)
	}

	// If error exists, then print error
	logErr := logkOption.GetError(options, logkOption.ErrorKey)
	if logErr != nil && lv <= level.Error {
		writer.Printf("  > Error: %s\n", logErr)
	}

	meta := options.Metadata
	if meta != nil && len(meta) > 0 {
		// Serialize to json
		metadata, err := json.Marshal(meta)
		// If not error, then print
		if err == nil {
			writer.Printf("  > Metadata: %s\n", metadata)
		}
	}
}

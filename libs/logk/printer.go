package logk

import (
	"github.com/konsultin/project-goes-here/libs/logk/level"
	logkOption "github.com/konsultin/project-goes-here/libs/logk/option"
)

type Printer interface {
	Print(namespace string, outLevel level.LogLevel, msg string, options *logkOption.Options)
}

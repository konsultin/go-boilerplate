package logk

import (
	"github.com/Konsultin/project-goes-here/libs/logk/level"
	logkOption "github.com/Konsultin/project-goes-here/libs/logk/option"
)

type Printer interface {
	Print(namespace string, outLevel level.LogLevel, msg string, options *logkOption.Options)
}

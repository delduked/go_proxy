package logger

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

var Log = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller:    true,
	ReportTimestamp: true,
	TimeFormat:      time.Kitchen,
	Prefix:          "Proxy Server",
})

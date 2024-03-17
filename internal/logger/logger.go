package logger

import (
	"os"
	"time"

	a "github.com/charmbracelet/log"
)

var Log = a.NewWithOptions(os.Stderr, a.Options{
	ReportCaller:    true,
	ReportTimestamp: true,
	TimeFormat:      time.Kitchen,
})

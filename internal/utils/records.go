package utils

import (
	"strings"
	"errors"
	"prx/internal/logger"
)

var RedirectRecords = make(map[string]string)

type RedirectFlag []string

func (f *RedirectFlag) String() string {
	return strings.Join(*f, ", ")
}

func (f *RedirectFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func ParseRedirectRecords(records RedirectFlag) error {
	for _, record := range records {
		parts := strings.SplitN(record, "=", 2)
		if len(parts) != 2 {
			logger.Log.Info("Skipping invalid record format:", "record", record)
			continue
		}
		RedirectRecords[parts[0]] = parts[1]
	}
	if len(RedirectRecords) == 0 {
		return errors.New("no valid redirect records provided")
	}
	return nil
}

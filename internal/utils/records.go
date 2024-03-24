package utils

import (
  "strings"
)

var RedirectRecords = make(map[string]string)

type RecordsFlag []string

func (i *RecordsFlag) String() string {
  return strings.Join(*i, ", ")
}

func (i *RecordsFlag) Set(value string) error {
  *i = append(*i, value)
  return nil
}

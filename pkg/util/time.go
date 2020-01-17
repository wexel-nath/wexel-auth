package util

import (
	"regexp"
	"time"
)

const (
	format = "Monday 2 Jan 06 3:04pm MST"

	whitespaceRegexString = `\s+`
)

var (
	whitespaceRegex = regexp.MustCompile(whitespaceRegexString)
)

func FormatTime(t time.Time) string {
	return t.In(time.Local).Format(format)
}

func StripWhitespace(s string) string {
	return whitespaceRegex.ReplaceAllString(s, "")
}

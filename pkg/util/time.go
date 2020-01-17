package util

import (
	"time"
)

const format = "Monday 2 Jan 06 3:04pm MST"

func FormatTime(t time.Time) string {
	return t.In(time.Local).Format(format)
}

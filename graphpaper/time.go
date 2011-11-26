package graphpaper

import (
	"time"
)

func FormatTimeLocal(nanoSeconds int64) string {
	seconds := nanoSeconds / 1000000000
	return time.SecondsToLocalTime(seconds).Format("2006-01-02 15:04:05Z07:00")
}

func FormatTime(nanoSeconds int64, utc bool, format string) string {
	seconds := nanoSeconds / 1000000000
	if utc {
		return time.SecondsToUTC(seconds).Format(format)
	} else {
		return time.SecondsToLocalTime(seconds).Format(format)
	}
	panic("unreachable code")
}

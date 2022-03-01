package utime

import "time"

const (
	FormatTimeSecond = "2006-01-02 15:04:05"
	FormatTimeDay    = "2006-01-02"
)

func StrSecond() string {
	return time.Now().Format(FormatTimeSecond)
}

func StrDay() string {
	return time.Now().Format(FormatTimeDay)
}

func StrDayByTime(t time.Time) string {
	return t.Format(FormatTimeDay)
}

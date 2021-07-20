package utils

import (
	"fmt"
	"time"
)

func ConvertIntervalToMins(interval string) (mins int, err error) {
	switch interval {
	case "15m":
		mins = 15
	case "30m":
		mins = 30
	case "1h":
		mins = 60
	case "2h":
		mins = 120
	case "4h":
		mins = 240
	case "1d":
		mins = 1440
	default:
		err = fmt.Errorf("interval '%s' not supported", interval)
	}
	return
}

// Get time block based on the length and time point
// For example, if time point is 10:00 and lenght is 4h, the time block would be between 08:00 and 12:00
func GetTimeBlockByLength(t time.Time, lengthMins int) (tStart time.Time, tEnd time.Time, err error) {
	// Get the minutes of time
	mins := t.Hour()*60 + t.Minute()
	timeBlocks := mins / lengthMins
	tStart = time.Date(t.Year(), t.Month(), t.Day(), 0, timeBlocks*lengthMins, 0, 0, time.UTC)
	tEnd = time.Date(t.Year(), t.Month(), t.Day(), 0, (timeBlocks+1)*lengthMins, 0, 0, time.UTC)
	return
}

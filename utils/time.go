package utils

import "fmt"

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

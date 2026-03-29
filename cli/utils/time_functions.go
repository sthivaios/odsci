package utils

import "fmt"

func TimeString(seconds int64) string {
	if (seconds > 60) {
		if (seconds % 60 == 0) {
			return fmt.Sprintf("%dm", seconds/60)
		} else {
			return fmt.Sprintf("%dm %ds", seconds/60, seconds % 60)
		}
	} else {
		return fmt.Sprintf("%ds", seconds)
	}
}

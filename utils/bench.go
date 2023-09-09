package utils

import (
	"fmt"
	"time"
)

func Trace(s string) (string, time.Time) {
	return s, time.Now()
}

func Un(s string, startTime time.Time) {
	endTime := time.Now()
	fmt.Println("[benchmark] #"+s, ": ", endTime.Sub(startTime))
}

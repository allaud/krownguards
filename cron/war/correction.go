package war

import (
	"math"
	"time"
	"ws/game"
)

func delayCorrection(startTime time.Time) float64 {
	endTime := time.Now()
	computationTimeMs := math.Ceil(float64(endTime.Sub(startTime).Nanoseconds() / 1000000))
	idealDelayMs := 1000 * game.WAR_STEP_DELAY
	delayMs := math.Max(10, idealDelayMs-computationTimeMs)
	return delayMs
}

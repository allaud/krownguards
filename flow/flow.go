package flow

import (
	//"fmt"
	"time"
	"ws/datatypes"
	"ws/models"
)

var Messagebox chan datatypes.Signal

func Run(state *models.State) {
	for {
		signal := <-Messagebox

		if signal.Duration == 0 {
			signal.Fn()
			continue
		}

		time.AfterFunc(signal.Duration, func() {
			if state.Pause == 0 || signal.Type == "pause" {
				signal.Duration = 0
			}

			Messagebox <- signal
		})
	}
}

func init() {
	Messagebox = make(chan datatypes.Signal, 100000)
}

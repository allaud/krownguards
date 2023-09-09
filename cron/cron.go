package cron

import (
	"fmt"
	"strconv"
	"time"
	"ws/datatypes"
	"ws/flow"
	"ws/game"
	"ws/models"
	"ws/utils"
)

func Run(state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	if game.COUNTDOWN > 0 {
		fmt.Println(game.COUNTDOWN)
		utils.SendNotification(state, channels, strconv.Itoa(game.COUNTDOWN), 0)
		utils.SendAll(channels, &datatypes.MessageData{
			"type":  "state",
			"state": state,
		}, 0)
		game.COUNTDOWN -= 1
		flow.Messagebox <- datatypes.Signal{
			Type:     "cron",
			Duration: time.Duration(game.TIMER_DELAY) * time.Second,
			Fn:       func() { Run(state, links, channels) },
		}
	} else {
		state.Phase = game.Building

		flow.Messagebox <- datatypes.Signal{
			Type:     "cron",
			Duration: time.Duration(game.TIMER_DELAY) * time.Second,
			Fn:       func() { updateTimer(state, links, channels) },
		}
	}
}

func ConnectDelay(state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	if state.WaitFor == 0 {
		return
	}
	if game.CONNECT_DELAY < 0 {
		utils.EndGame(state)
	}
	fmt.Println("Waiting for players:", game.CONNECT_DELAY)
	game.CONNECT_DELAY -= 1
	flow.Messagebox <- datatypes.Signal{
		Type:     "cron",
		Duration: time.Duration(game.TIMER_DELAY) * time.Second,
		Fn:       func() { ConnectDelay(state, links, channels) },
	}
}

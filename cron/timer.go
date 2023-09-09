package cron

import (
	"fmt"
	"strconv"
	"time"
	"ws/bot"
	"ws/cron/war"
	"ws/datatypes"
	"ws/flow"
	"ws/game"
	"ws/game/units"
	"ws/models"
	"ws/utils"
)

func updateTimer(state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	state.Servertime = state.Servertime + 1

	// regenerate king's hp
	for _, king := range state.Kings {
		king.Hp = king.Hp + king.HpReg
		if king.Hp > king.MaxHp {
			king.Hp = king.MaxHp
		}
	}

	// add stone to players
	if state.Servertime%game.STONE_INCOME_STEP == 0 {
		for _, player := range state.Players {
			player.Stone += player.StoneGrade.StoneInc
		}
	}

	for _, player := range state.Players {
		// cooldown summons
		for _, summon := range player.SummonUnits {
			if summon.Current >= summon.Cap {
				continue
			}
			summon.Timer = summon.Timer + 1
			if summon.Timer >= summon.Cooldown {
				summon.Current = summon.Current + 1
				if summon.Current > summon.Cap {
					summon.Current = summon.Cap
				}
				summon.Timer = 0
			}
		}
		// progress upgrades
		if player.FoodGrade.InProgress {
			player.FoodGrade.Timer = player.FoodGrade.Timer + 1
			if player.FoodGrade.Timer >= player.FoodGrade.Cooldown {
				player.FoodGrade.InProgress = false
				player.FoodGrade.Timer = 0
				player.Food[1] = player.Food[1] + game.FOOD_GRADE_STEP
				/*utils.SendAll(channels, &datatypes.MessageData{
					"type":  "upgrade_farm!",
					"state": state,
				}, 0)*/
			}
		}
		if player.StoneGrade.InProgress {
			player.StoneGrade.Timer = player.StoneGrade.Timer + 1
			if player.StoneGrade.Timer >= player.StoneGrade.Cooldown {
				player.StoneGrade.InProgress = false
				player.StoneGrade.Timer = 0
				player.StoneGrade.StoneInc = game.STONE_GRADE[player.StoneGrade.CurrGrade][2]
				player.StoneGrade.CurrGrade += 1
				player.StoneGrade.Price = [2]int{game.STONE_GRADE[player.StoneGrade.CurrGrade][0], game.STONE_GRADE[player.StoneGrade.CurrGrade][1]}
				/*utils.SendAll(channels, &datatypes.MessageData{
					"type":  "upgrade_stone!",
					"state": state,
				}, 0)*/
			}
		}
	}

	// change timer
	if state.Phase == game.Building {
		if state.WaveCountdown == game.SWITCH_WAR_DELAY {
			wave := units.Waves[state.Wave]
			before := strconv.Itoa(game.SWITCH_WAR_DELAY)
			arenaNext := false
			if len(state.ArenaWaves) > 0 {
				if state.ArenaWaves[0] == state.Wave {
					arenaNext = true
				}
			}
			if arenaNext {
				utils.SendNotification(state, channels, "Arena battle in "+before+" seconds", 0)
			} else {
				utils.SendNotification(state, channels, "Wave #"+strconv.Itoa(state.Wave)+
					" ("+wave.Unit+" x "+strconv.Itoa(wave.Count)+") in "+before+" seconds", 0)
			}
		}
		state.WaveCountdown = state.WaveCountdown - 1
		utils.SendAll(channels, &datatypes.MessageData{
			"type":  "state",
			"state": state,
		}, 0)
	} else {
		state.WaveCountdown += 1
	}

	fmt.Println("time>", state.Servertime)

	bot.BotsStep(state, links, channels)

	if state.WaveCountdown <= 0 && state.Phase == game.Building {
		flow.Messagebox <- datatypes.Signal{
			Type: "cron",
			Fn:   func() { war.SwitchWar(state, links, channels) },
		}
	}

	flow.Messagebox <- datatypes.Signal{
		Type:     "cron",
		Duration: 1 * time.Second,
		Fn:       func() { updateTimer(state, links, channels) },
	}
}

func PauseTimer(state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	if state.Pause == 0 {
		return
	}
	player, ok := state.Players[state.Pause]
	if ok {
		player.PauseCap -= 1
		fmt.Println(player.Name, "pause cap left", player.PauseCap)
	}
	utils.SendAll(channels, &datatypes.MessageData{
		"type":  "state",
		"state": state,
	}, 0)
	flow.Messagebox <- datatypes.Signal{
		Type:     "pause",
		Duration: 1 * time.Second,
		Fn:       func() { PauseTimer(state, links, channels) },
	}
}

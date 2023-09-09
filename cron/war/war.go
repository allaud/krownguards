package war

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
	"ws/abilities"
	"ws/datatypes"
	"ws/flow"
	"ws/game"
	"ws/game/units"
	"ws/models"
	"ws/utils"
)

func SwitchWar(state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	defer utils.Un(utils.Trace("switch_war"))
	abilities.ResetProjId()
	ClearDelayedLeavers(state)

	for ind, number := range state.ArenaWaves {
		if state.Wave == number {
			state.ArenaWaves = append(state.ArenaWaves[:ind], state.ArenaWaves[ind+1:]...)
			SwitchArena(state, links, channels)
			return
		}
	}

	fmt.Println("switch war!!!")
	state.Phase = game.War
	originalUnits := map[int][]models.Unit{}

	for slot, player := range state.Players {
		for _, unit := range state.PlayerUnits[slot] {
			// saving unit positions for build phase
			SaveOriginalUnit(originalUnits, unit)
		}

		for index := 0; index < units.Waves[state.Wave].Count; index++ {
			waveUnit, exists := units.UNITS_BY_NAME[units.Waves[state.Wave].Unit]
			if !exists {
				return
			}
			waveUnit.Id = utils.UniqId()
			waveUnit.Slot = slot
			waveUnit.Waypoint = 1
			waveUnit.Action = game.Idle
			utils.RandomSpawn(&waveUnit, state.WaveUnits[slot], game.SpawnSlots)
			waveUnit.Abilities = utils.CopyAbilities(waveUnit.Abilities)
			waveUnit.Affected = make(map[string]*models.Effect)
			state.WaveUnits[slot] = append(state.WaveUnits[slot], &waveUnit)
		}

		enemySlots := []int{}
		for eSlot, _ := range state.Players {
			if utils.Side(slot) != utils.Side(eSlot) {
				enemySlots = append(enemySlots, eSlot)
			}
		}
		if len(enemySlots) == 0 {
			if utils.Side(slot) == game.West {
				enemySlots = append(enemySlots, 5)
			} else {
				enemySlots = append(enemySlots, 1)
			}
		}
		//fmt.Println("Player", player.Name, "have", player.IncomeUnits, "at start of wave")
		for _, incUnit := range player.IncomeUnits {
			incUnit.Slot = enemySlots[rand.Intn(len(enemySlots))]
			utils.RandomSpawn(incUnit, state.WaveUnits[incUnit.Slot], game.SpawnSlots)
			state.WaveUnits[incUnit.Slot] = append(state.WaveUnits[incUnit.Slot], incUnit)
		}
		//clearing player's income box
		player.IncomeUnits = []*models.Unit{}
	}
	//Count number of enemies each side and units leaked by each player
	RefreshWarStatistics(state)
	stateTemp := utils.Deepcopy(state).(*models.State)

	utils.SendAll(channels, &datatypes.MessageData{
		"type":  "war!",
		"state": state,
	}, 0)

	flow.Messagebox <- datatypes.Signal{
		Type:     "cron",
		Duration: time.Duration(1000*game.WAR_STEP_DELAY) * time.Millisecond,
		Fn:       func() { WarStep(stateTemp, state, links, &originalUnits, channels) },
	}
}

func WarStep(stateTemp, state *models.State, links *models.Links, originalUnits *map[int][]models.Unit, channels datatypes.ChannelsMap) {
	//defer utils.Un(utils.Trace("war_step"))
	startTime := time.Now()
	// check if game over
	for _, king := range state.Kings {
		if king.Action == game.Dead {
			state.Phase = game.GameOver

			utils.SendAll(channels, &datatypes.MessageData{
				"type":  "end!",
				"state": state,
			}, 0)

			utils.EndGame(state)
			return
		}
	}
	// check if wave in progress
	if !IsWarFinished(state) {
		CombatStep(stateTemp, state, links, true)
		CombatStep(stateTemp, state, links, false)
		abilities.Step(state, links)
		//Count number of enemies each side and units leaked by each player
		RefreshWarStatistics(state)
		stateTemp = utils.Deepcopy(state).(*models.State)
		utils.SendAll(channels, &datatypes.MessageData{
			"type":  "war!",
			"state": state,
		}, 0)

		flow.Messagebox <- datatypes.Signal{
			Type:     "cron",
			Duration: time.Duration(delayCorrection(startTime)) * time.Millisecond,
			Fn:       func() { WarStep(stateTemp, state, links, originalUnits, channels) },
		}
		return
	}
	// wave is over
	fmt.Println("war is over")
	state.Phase = game.Building
	state.WaveCountdown = game.SWITCH_WAR_DELAY
	state.Wave += 1
	// upgrade available summons
	if _, ok := units.Summons[state.Wave]; ok {
		for _, player := range state.Players {
			for k, v := range units.Summons[state.Wave] {
				tempV := *v
				player.SummonUnits[k] = &tempV
			}
		}
	}

	// reset player's units
	for pSlot, _ := range state.PlayerUnits {
		state.PlayerUnits[pSlot] = []*models.Unit{}
	}
	// restore units to build state
	for pSlot, pUnits := range *originalUnits {
		for _, unit := range pUnits {
			unitTemp := unit
			state.PlayerUnits[pSlot] = append(state.PlayerUnits[pSlot], &unitTemp)
		}
	}
	// get gold from all sources
	for pSlot, _ := range state.Players {
		state.Players[pSlot].Gold += state.Players[pSlot].Income
		state.Players[pSlot].Gold += units.Waves[state.Wave-1].Bounty
		state.Players[pSlot].Gold += (state.Kings[0].Bounty + state.Kings[1].Bounty) / len(state.Players)
	}
	// reset kings
	for index, king := range state.Kings {
		king.Action = game.Idle
		king.Bounty = 0
		king.Coords = game.KingCoords[index]
		king.Dir = game.KingDir[index]
	}
	// reset wave
	for wSlot, _ := range state.WaveUnits {
		state.WaveUnits[wSlot] = []*models.Unit{}
	}
	// reset Projectiles
	state.Projectiles = []*models.Projectile{}
	links.ProjLinks = []*models.ProjLink{}

	utils.SendAll(channels, &datatypes.MessageData{
		"type":  "war!",
		"state": state,
	}, 0)

	//flow.Messagebox <- datatypes.Signal{
	//	Type:     "cron",
	//	Duration: time.Duration(game.SWITCH_WAR_DELAY) * time.Second,
	//	Fn:       func() { SwitchWar(state, links, channels) },
	//}
}

func SwitchArena(state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	fmt.Println("switch arena!!!")
	state.Phase = game.Arena
	originalUnits := map[int][]models.Unit{}

	// count teams
	westSide := []int{}
	eastSide := []int{}
	for slot, _ := range state.Players {
		if utils.Side(slot) == game.West {
			westSide = append(westSide, slot)
		} else {
			eastSide = append(eastSide, slot)
		}
	}
	sort.Ints(westSide)
	sort.Ints(eastSide)

	for slot, _ := range state.Players {
		// get team size and player position in team
		plIndex := 0
		plSide := []int{}
		if utils.Side(slot) == game.West {
			plSide = westSide
		} else {
			plSide = eastSide
		}
		for sIndex, sSlot := range plSide {
			if slot == sSlot {
				plIndex = sIndex + 1
			}
		}
		max := len(plSide)
		// get player slots
		spawnSlotsValue := *game.SpawnSlots
		spawnSlot := spawnSlotsValue[game.ArenaCoef+slot/game.SideCoef]
		buildSlotsValue := *game.BuildSlots
		buildSlot := buildSlotsValue[slot]
		spawnCenter := 0.5 * (spawnSlot.YRange[1] + spawnSlot.YRange[0])
		halfSize := 0.5 * (buildSlot.XRange[1] - buildSlot.XRange[0])
		slotCenter := spawnCenter - halfSize*float64(2*plIndex-max-1)
		personalSlot := datatypes.Size{
			XRange: spawnSlot.XRange,
			YRange: [2]float64{
				slotCenter - halfSize,
				slotCenter + halfSize,
			},
		}

		for _, unit := range state.PlayerUnits[slot] {
			// saving unit positions for build phase
			SaveOriginalUnit(originalUnits, unit)
			unit.Slot = game.ArenaCoef + slot/game.SideCoef
			// new samurai spawn
			dX := unit.Coords[0] - buildSlot.XRange[0]
			dY := unit.Coords[1] - buildSlot.YRange[0]
			switch slot {
			case 1, 2, 7, 8:
				unit.Coords[0] = personalSlot.XRange[0] + dY
				unit.Coords[1] = personalSlot.YRange[1] - dX
			case 3, 4, 5, 6:
				unit.Coords[0] = personalSlot.XRange[1] - dY
				unit.Coords[1] = personalSlot.YRange[0] + dX
			}
			//unit.Coords[0] = personalSlot.XRange[0] + unit.Coords[0] - buildSlot.XRange[0]
			//unit.Coords[1] = personalSlot.YRange[0] + unit.Coords[1] - buildSlot.YRange[0]
			//state.PlayerUnits[slot] => GetUnitsBySlots(state, SideSlots(slot), true)
			//utils.RandomSpawn(unit, GetUnitsBySlots(state, SideSlots(slot), true), game.SpawnSlots)
		}
	}
	// spawn summons after all units spawned
	for slot, player := range state.Players {
		for _, incUnit := range player.IncomeUnits {
			incUnit.Slot = game.ArenaCoef + slot/game.SideCoef
			//state.PlayerUnits[slot] => GetUnitsBySlots(state, SideSlots(slot), true)
			utils.RandomSpawn(incUnit, GetUnitsBySlots(state, SideSlots(slot), true), game.SpawnSlots)
			state.PlayerUnits[slot] = append(state.PlayerUnits[slot], incUnit)
		}
		// clearing player's income box
		player.IncomeUnits = []*models.Unit{}
	}
	// Add bonus turtles REPLACE CUSTOM THINGS ("Warturtle", 5)
	if (len(westSide) != 0) && (len(eastSide) != 0) && (len(westSide) != len(eastSide)) {
		limit := 0
		lesserSide := []int{}

		if len(westSide) > len(eastSide) {
			limit = (len(westSide) - len(eastSide)) * (state.Wave / 5)
			lesserSide = eastSide
		} else {
			limit = (len(eastSide) - len(westSide)) * (state.Wave / 5)
			lesserSide = westSide
		}

		index := 0
		for index < limit {
			arenaUnit, exists := units.UNITS_BY_NAME["Warturtle"]
			if !exists {
				return
			}
			arenaUnit.Id = utils.UniqId()
			arenaSlot := lesserSide[rand.Intn(len(lesserSide))]
			arenaUnit.Slot = game.ArenaCoef + arenaSlot/game.SideCoef
			arenaUnit.Waypoint = 1
			arenaUnit.Action = game.Idle
			//state.PlayerUnits[slot] => GetUnitsBySlots(state, SideSlots(slot), true)
			utils.RandomSpawn(&arenaUnit, GetUnitsBySlots(state, SideSlots(arenaSlot), true), game.SpawnSlots)
			arenaUnit.Affected = make(map[string]*models.Effect)
			state.PlayerUnits[arenaSlot] = append(state.PlayerUnits[arenaSlot], &arenaUnit)
			index = index + 1
		}
	}

	RefreshArenaStatistics(state)
	stateTemp := utils.Deepcopy(state).(*models.State)

	utils.SendAll(channels, &datatypes.MessageData{
		"type":  "arena!",
		"state": state,
	}, 0)

	flow.Messagebox <- datatypes.Signal{
		Type:     "cron",
		Duration: time.Duration(1000*game.WAR_STEP_DELAY) * time.Millisecond,
		Fn:       func() { ArenaStep(stateTemp, state, links, &originalUnits, channels) },
	}
}

func ArenaStep(stateTemp, state *models.State, links *models.Links, originalUnits *map[int][]models.Unit, channels datatypes.ChannelsMap) {
	startTime := time.Now()
	westDefeated := IsSideDefeated(state, game.West)
	eastDefeated := IsSideDefeated(state, game.East)
	// check if arena in progress
	if !westDefeated && !eastDefeated {
		if state.WaveCountdown > 2 {
			CombatStep(stateTemp, state, links, true)
			CombatStep(stateTemp, state, links, false)
			abilities.Step(state, links)
			RefreshArenaStatistics(state)
		}
		stateTemp = utils.Deepcopy(state).(*models.State)
		utils.SendAll(channels, &datatypes.MessageData{
			"type":  "arena!",
			"state": state,
		}, 0)

		flow.Messagebox <- datatypes.Signal{
			Type:     "cron",
			Duration: time.Duration(delayCorrection(startTime)) * time.Millisecond,
			Fn:       func() { ArenaStep(stateTemp, state, links, originalUnits, channels) },
		}
		return
	}
	// arena is over
	fmt.Println("arena is over")
	state.Phase = game.Building
	state.WaveCountdown = game.SWITCH_WAR_DELAY
	state.WestAlive = 0
	state.EastAlive = 0
	winnerSide := "noone"
	if westDefeated && eastDefeated {
		utils.SendNotification(state, channels, "It's a draw!", 0)
	} else if westDefeated {
		winnerSide = game.East
		utils.SendNotification(state, channels, "Congratulations to the East!", 0)
	} else if eastDefeated {
		winnerSide = game.West
		utils.SendNotification(state, channels, "Congratulations to the West!", 0)
	}
	// reset player's units
	for pSlot, _ := range state.PlayerUnits {
		state.PlayerUnits[pSlot] = []*models.Unit{}
	}
	// restore units to build state
	for pSlot, pUnits := range *originalUnits {
		for _, unit := range pUnits {
			unitTemp := unit
			state.PlayerUnits[pSlot] = append(state.PlayerUnits[pSlot], &unitTemp)
		}
	}
	// get gold from all sources
	for pSlot, _ := range state.Players {
		if utils.Side(pSlot) == winnerSide {
			state.Players[pSlot].Gold += 10 * (state.Wave - 1)
		}
	}

	// reset Projectiles
	state.Projectiles = []*models.Projectile{}
	links.ProjLinks = []*models.ProjLink{}

	utils.SendAll(channels, &datatypes.MessageData{
		"type":  "arena!",
		"state": state,
	}, 0)

	//flow.Messagebox <- datatypes.Signal{
	//	Type:     "cron",
	//	Duration: time.Duration(game.SWITCH_WAR_DELAY) * time.Second,
	//	Fn:       func() { SwitchWar(state, links, channels) },
	//}
}

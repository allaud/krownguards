package utils

import (
	"github.com/fatih/structs"
	"ws/datatypes"
	"ws/game"
	"ws/game/units"
	"ws/models"
)

var FULL_TRANSFORMATIONS = []func(data *datatypes.MessageData){
	deepcopy,
	skipWaitFor,
	stripServiceFields,
	stripProjectiles,
	gen_diff,
}

var SIMPLE_TRANSFORMATIONS = []func(data *datatypes.MessageData){
	deepcopy,
	skipWaitFor,
}

func preprocess(data *datatypes.MessageData, transforms []func(data *datatypes.MessageData)) *datatypes.MessageData {
	for _, transform := range transforms {
		transform(data)
	}
	return data
}

func deepcopy(data *datatypes.MessageData) {
	state := (*data)["state"].(*models.State)
	(*data)["state"] = Deepcopy(state).(*models.State)
}

func skipWaitFor(data *datatypes.MessageData) {
	state := (*data)["state"].(*models.State)
	state.WaitFor = 0
}

func stripServiceFields(data *datatypes.MessageData) {
	state := (*data)["state"].(*models.State)
	for _, player := range state.Players {
		player.Score = float64(int(player.Score))
	}
	for _, pUnits := range state.PlayerUnits {
		for _, unit := range pUnits {
			//unit.Projectile = ""
			//unit.ProjSpeed = 0
			//unit.ProjTraj = ""
			unit.AspdTimer = 0
			unit.Size = 0
			unit.TargetCoords = [2]float64{0, 0}
			unit.TargetSign = 0
			//unit.Slot = 0
			unit.Waypoint = 0
			//unit.BuildWave = 0
			unit.Killer = 0
			unit.Atk[0] = float64(int(unit.Atk[0]))
			unit.Atk[1] = float64(int(unit.Atk[1]))
			unit.Hp = float64(int(unit.Hp))
			unit.MaxHp = float64(int(unit.MaxHp))
			unit.Mp = float64(int(unit.Mp))
			unit.MaxMp = float64(int(unit.MaxMp))
			unit.Abilities = CopyAbilities(units.UNITS_BY_NAME[unit.Name].Abilities)
		}
	}
	for _, wUnits := range state.WaveUnits {
		for _, unit := range wUnits {
			//unit.Projectile = ""
			//unit.ProjSpeed = 0
			//unit.ProjTraj = ""
			unit.AspdTimer = 0
			unit.Size = 0
			unit.TargetCoords = [2]float64{0, 0}
			unit.TargetSign = 0
			//unit.Slot = 0
			unit.Waypoint = 0
			//unit.BuildWave = 0
			unit.Killer = 0
			unit.Atk[0] = float64(int(unit.Atk[0]))
			unit.Atk[1] = float64(int(unit.Atk[1]))
			unit.Hp = float64(int(unit.Hp))
			unit.MaxHp = float64(int(unit.MaxHp))
			unit.Mp = float64(int(unit.Mp))
			unit.MaxMp = float64(int(unit.MaxMp))
			unit.Abilities = CopyAbilities(units.UNITS_BY_NAME[unit.Name].Abilities)
		}
	}
	for _, unit := range state.Kings {
		//unit.Projectile = ""
		//unit.ProjSpeed = 0
		//unit.ProjTraj = ""
		unit.AspdTimer = 0
		unit.TargetCoords = [2]float64{0, 0}
		unit.TargetSign = 0
		unit.Atk[0] = float64(int(unit.Atk[0]))
		unit.Atk[1] = float64(int(unit.Atk[1]))
		unit.Hp = float64(int(unit.Hp))
		unit.MaxHp = float64(int(unit.MaxHp))
		unit.Mp = float64(int(unit.Mp))
		unit.MaxMp = float64(int(unit.MaxMp))
		unit.Abilities = CopyAbilities(units.King.Abilities)
	}
}

func stripProjectiles(data *datatypes.MessageData) {
	state := (*data)["state"].(*models.State)
	index := len(state.Projectiles) - 1
	for index >= 0 {
		if state.Projectiles[index].Name == game.Invisible || state.Projectiles[index].Name == game.Melee {
			state.Projectiles = append(state.Projectiles[:index], state.Projectiles[index+1:]...)
		}
		index = index - 1
	}
}

var gen_diff = func() func(data_p *datatypes.MessageData) {
	//var StartState = map[string]interface{}{}
	var PrevState *models.State

	return func(data_p *datatypes.MessageData) {
		//defer Un(Trace("get_diff"))
		data := *data_p
		StateMap := structs.Map(data["state"])

		var StateDiff interface{}
		if PrevState != nil {
			StateDiff = data["state"].(*models.State).Diff(PrevState)
		} else {
			StateDiff = StateMap
		}

		//StateDiff := diff(StartState, StateMap)
		//fmt.Println(data["state"].(*models.State).Diff(data["state"].(*models.State)))
		//fmt.Println("DIFF> ", StateDiff)
		//StartState = Deepcopy(StateMap).(map[string]interface{})
		PrevState = Deepcopy(data["state"]).(*models.State)
		data["state"] = StateDiff
	}
}()

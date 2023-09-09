package war

import (
	//"fmt"
	"ws/abilities"
	"ws/game"
	"ws/game/units"
	"ws/models"
	"ws/utils"
	"ws/utils/status"
)

func ClearDelayedLeavers(state *models.State) {
	if len(state.Players) <= len(state.Usernames) {
		return
	}
	for slot, player := range state.Players {
		if utils.ContainsString(state.Usernames, player.Name) {
			continue
		}
		if utils.ContainsString(game.BotNames, player.Name) {
			continue
		}
		delete(state.Players, slot)
		delete(state.PlayerUnits, slot)
	}
}

func SaveOriginalUnit(originalUnits map[int][]models.Unit, unit *models.Unit) {
	utils.DampPrice(unit)
	tempUnitValue := *unit
	tempUnitValue.Abilities = utils.CopyAbilities(tempUnitValue.Abilities)
	tempUnitValue.Affected = make(map[string]*models.Effect)
	originalUnits[unit.Player] = append(originalUnits[unit.Player], tempUnitValue)
}

func RefreshWarStatistics(state *models.State) {
	state.WestAlive = 0
	state.EastAlive = 0
	for slot, wave := range state.WaveUnits {
		count := 0
		for _, unit := range wave {
			if !status.IsMarkedDead(unit) {
				count += 1
			}
		}
		if utils.Side(slot) == game.West {
			state.WestAlive += count
		} else {
			state.EastAlive += count
		}
	}
}

func RefreshArenaStatistics(state *models.State) {
	state.WestAlive = 0
	state.EastAlive = 0
	for slot, units := range state.PlayerUnits {
		count := 0
		for _, unit := range units {
			if !status.IsMarkedDead(unit) {
				count += 1
			}
		}
		if utils.Side(slot) == game.West {
			state.EastAlive += count
		} else {
			state.WestAlive += count
		}
	}
}

func GetScore(state *models.State, unit *models.Unit, player int, gain bool) {
	if _, ok := state.Players[player]; !ok {
		return
	}
	if status.IsSummon(unit) {
		return
	}
	plValue := float64(state.Players[player].Value)
	recValue := units.Waves[state.Wave].RecValue
	score := game.UnitScore
	// check if boss wave
	if state.Wave%10 == 0 {
		score = 10 * score
	}
	switch gain {
	// getting score for killing wave unit
	case true:
		if recValue != 0 && plValue < recValue {
			score = score * (2*recValue - plValue) / recValue
		}
		state.Players[player].Score += score
	// losing score for leaking wave unit
	case false:
		if recValue != 0 && plValue > recValue {
			score = score * plValue / recValue
		}
		state.Players[player].Score -= score
		state.Players[player].Leaked += 1
	}
}

func DampBounty(state *models.State, unit *models.Unit) {
	if _, ok := state.Players[unit.Slot]; !ok {
		return
	}
	if status.IsSummon(unit) {
		return
	}
	dampCoef := 2*float64(state.Players[unit.Slot].Value)/units.Waves[state.Wave].RecValue - 1
	switch {
	case dampCoef < 0:
		dampCoef = 0
	case dampCoef > 1:
		dampCoef = 1
	}
	unit.Bounty = int(float64(unit.Bounty) * dampCoef)
}

func TriggerDeath(stateTemp, state *models.State, unitTemp, unit *models.Unit, options *models.AbOptions) {
	if status.IsMarkedDead(unitTemp) {
		return
	}
	unit.Action = game.Dead
	unit.Affected = map[string]*models.Effect{}
	// apply on_death skills
	options.Unit = unit
	options.State = state
	options.StateCurrent = stateTemp
	abilities.AddByType(unitTemp, abilities.ABILITIES, options, game.OnDeath)
	RefreshAttributes(state, unit)
	// if unit's killed by player
	if unitTemp.Killer >= 1 && unitTemp.Killer <= 8 {
		if _, ok := state.Players[unitTemp.Killer]; !ok {
			return
		}
		//fmt.Println("Player", unit.Killer, "got", unit.Bounty)
		state.Players[unitTemp.Killer].Gold += unitTemp.Bounty
		// Getting score for killing unit
		GetScore(state, unitTemp, unitTemp.Killer, true)
		return
	}
	// if unit's killed by King
	if unitTemp.Killer/game.KingCoef == 1 {
		state.Kings[unitTemp.Killer%game.KingCoef].Bounty += unitTemp.Bounty
		//fmt.Println("King", unit.Killer%game.KingCoef, "got", unit.Bounty)
	}
}

func Attack(state *models.State, links *models.Links, unitTemp, target, unit *models.Unit, options *models.AbOptions) {
	unit.Action = game.Wait
	unit.AspdTimer += 1
	unit.Dir = target.Coords
	if unit.AspdTimer == HalfAspd(unitTemp) {
		unit.Action = game.Attack
	}
	if unit.AspdTimer < unitTemp.Aspd {
		return
	}
	// begin attack swing
	options.Unit = unit
	options.Target = target
	options.Atk[0] = unitTemp.Atk[0]
	options.Atk[1] = unitTemp.Atk[1]
	// check is unit have replacer ability instead of attack
	if _, ok := unitTemp.Abilities[game.AttackReplacer]; ok {
		abilities.AddByType(unitTemp, abilities.ABILITIES, options, game.AttackReplacer)
		SetDelay(unit, 0)
		return
	}
	// create attack projectile for ranged unit
	if unitTemp.Projectile != game.Melee {
		abilities.Create(state, links, unit, target, options, game.Attack)
		SetDelay(unit, 0)
		return
	}
	// set attack hit for melee unit
	abilities.AttackHit(unit, target, options)
	SetDelay(unit, 0)
}

func RegenSources(unit *models.Unit) {
	if unit.HpReg != 0 && !status.IsKing(unit) {
		unit.Hp = unit.Hp + unit.HpReg
	}
	if unit.Hp > unit.MaxHp {
		unit.Hp = unit.MaxHp
	}
	if unit.MaxMp == 0 || unit.MpReg == 0 {
		return
	}
	unit.Mp = unit.Mp + unit.MpReg
	if unit.Mp > unit.MaxMp {
		unit.Mp = unit.MaxMp
	}
}

func RefreshAttributes(state *models.State, unit *models.Unit) {
	unitProto := models.Unit{}
	if status.IsKing(unit) {
		unitProto = units.King
		atk := state.KingAttrs[unit.Slot-game.KingCoef].AtkGrade
		unitProto.Atk[0] += float64(atk.CurrGrade) * atk.GradeStep
		unitProto.Atk[1] += float64(atk.CurrGrade) * atk.GradeStep
		maxHp := state.KingAttrs[unit.Slot-game.KingCoef].MaxHpGrade
		unitProto.MaxHp += float64(maxHp.CurrGrade) * maxHp.GradeStep
		unitProto.Def += float64(maxHp.CurrGrade / 4)
		hpReg := state.KingAttrs[unit.Slot-game.KingCoef].HpRegGrade
		unitProto.HpReg += float64(hpReg.CurrGrade) * hpReg.GradeStep
	} else {
		unitProto = units.UNITS_BY_NAME[unit.Name]
	}
	unit.Atk = unitProto.Atk
	unit.Aspd = unitProto.Aspd
	unit.Def = unitProto.Def
	unit.Movespeed = unitProto.Movespeed
	unit.MaxHp = unitProto.MaxHp
	unit.MpReg = unitProto.MpReg
	unit.Abilities = utils.CopyAbilities(unitProto.Abilities)
}

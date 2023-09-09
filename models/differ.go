package models

func (current Upgrade) Diff(old Upgrade) map[string]interface{} {
	result := map[string]interface{}{}
	
	if old.GradeStep != current.GradeStep {
		result["GradeStep"] = current.GradeStep
	}

	if old.MaxGrade != current.MaxGrade {
		result["MaxGrade"] = current.MaxGrade
	}

	if old.CurrGrade != current.CurrGrade {
		result["CurrGrade"] = current.CurrGrade
	}

	if old.Price != current.Price {
		result["Price"] = current.Price
	}

	return result
}

func (current *KingUpgrades) Diff(old *KingUpgrades) map[string]interface{} {
	result := map[string]interface{}{}
	
	diffHpRegGrade := current.HpRegGrade.Diff(old.HpRegGrade)
	if len(diffHpRegGrade) > 0 {
		result["HpRegGrade"] = diffHpRegGrade
	}

	diffAtkGrade := current.AtkGrade.Diff(old.AtkGrade)
	if len(diffAtkGrade) > 0 {
		result["AtkGrade"] = diffAtkGrade
	}

	diffMaxHpGrade := current.MaxHpGrade.Diff(old.MaxHpGrade)
	if len(diffMaxHpGrade) > 0 {
		result["MaxHpGrade"] = diffMaxHpGrade
	}

	return result
}

func (current *State) Diff(old *State) interface{} {
	result := map[string]interface{}{}
	delta := 0

	if old.Id != current.Id {
		result["Id"] = current.Id
	}

	if old.Pause != current.Pause {
		result["Pause"] = current.Pause
	}

	if old.Servertime != current.Servertime {
		result["Servertime"] = current.Servertime
	}

	if old.WaveCountdown != current.WaveCountdown {
		result["WaveCountdown"] = current.WaveCountdown
	}

	if old.WaitFor != current.WaitFor {
		result["WaitFor"] = current.WaitFor
	}

	if old.Phase != current.Phase {
		result["Phase"] = current.Phase
	}

	if old.Wave != current.Wave {
		result["Wave"] = current.Wave
	}

	resultArenaWaves := make([]interface{}, len(current.ArenaWaves))
	delta = 0

	if len(current.ArenaWaves) == 0 || len(old.ArenaWaves) == 0 {
		if len(current.ArenaWaves) == 0 && len(old.ArenaWaves) != 0 {
			result["ArenaWaves"] = []interface{}{}
		} else if len(current.ArenaWaves) != 0 && len(old.ArenaWaves) == 0 {
			result["ArenaWaves"] = current.ArenaWaves
		}
	} else {
		for index, _ := range current.ArenaWaves {
			if index >= len(old.ArenaWaves) {
				resultArenaWaves[index] = current.ArenaWaves[index]
				delta = delta + 1
				continue
			}
			if old.ArenaWaves[index] != current.ArenaWaves[index] {
				resultArenaWaves[index] = current.ArenaWaves[index]
				delta = delta + 1
			} else {
				resultArenaWaves[index] = "__"
			}
		}
		if delta > 0 || len(current.ArenaWaves) != len(old.ArenaWaves) {
			result["ArenaWaves"] = resultArenaWaves
		}
	}

	if old.WestAlive != current.WestAlive {
		result["WestAlive"] = current.WestAlive
	}

	if old.EastAlive != current.EastAlive {
		result["EastAlive"] = current.EastAlive
	}

	if old.Port != current.Port {
		result["Port"] = current.Port
	}

	if old.LobbyAddr != current.LobbyAddr {
		result["LobbyAddr"] = current.LobbyAddr
	}

	resultUsernames := make([]interface{}, len(current.Usernames))
	delta = 0

	if len(current.Usernames) == 0 || len(old.Usernames) == 0 {
		if len(current.Usernames) == 0 && len(old.Usernames) != 0 {
			result["Usernames"] = []interface{}{}
		} else if len(current.Usernames) != 0 && len(old.Usernames) == 0 {
			result["Usernames"] = current.Usernames
		}
	} else {
		for index, _ := range current.Usernames {
			if index >= len(old.Usernames) {
				resultUsernames[index] = current.Usernames[index]
				delta = delta + 1
				continue
			}
			if old.Usernames[index] != current.Usernames[index] {
				resultUsernames[index] = current.Usernames[index]
				delta = delta + 1
			} else {
				resultUsernames[index] = "__"
			}
		}
		if delta > 0 || len(current.Usernames) != len(old.Usernames) {
			result["Usernames"] = resultUsernames
		}
	}

	resultPlayers := map[int]interface{}{}
	for key, val := range current.Players {
		if oldPlayers, ok := old.Players[key]; ok {
			diff := val.Diff(oldPlayers)
			if len(diff) > 0 {
				resultPlayers[key] = diff
			}
		} else {
			resultPlayers[key] = val
		}
	}
	for key, _ := range old.Players {
		if _, ok := current.Players[key]; !ok {
			resultPlayers[key] = "_del_"
		}
	}
	if len(resultPlayers) > 0 {
		result["Players"] = resultPlayers
	}

	resultPlayerUnits := map[int]interface{}{}
	for key, val := range current.PlayerUnits {
		if oldPlayerUnits, ok := old.PlayerUnits[key]; ok {
			diff := UnitArrayDiff(val, oldPlayerUnits)
			if _, ok := diff.(bool); !ok {
				resultPlayerUnits[key] = diff
			}
		} else {
			resultPlayerUnits[key] = val
		}
	}
	for key, _ := range old.PlayerUnits {
		if _, ok := current.PlayerUnits[key]; !ok {
			resultPlayerUnits[key] = "_del_"
		}
	}
	if len(resultPlayerUnits) > 0 {
		result["PlayerUnits"] = resultPlayerUnits
	}

	resultWaveUnits := map[int]interface{}{}
	for key, val := range current.WaveUnits {
		if oldWaveUnits, ok := old.WaveUnits[key]; ok {
			diff := UnitArrayDiff(val, oldWaveUnits)
			if _, ok := diff.(bool); !ok {
				resultWaveUnits[key] = diff
			}
		} else {
			resultWaveUnits[key] = val
		}
	}
	for key, _ := range old.WaveUnits {
		if _, ok := current.WaveUnits[key]; !ok {
			resultWaveUnits[key] = "_del_"
		}
	}
	if len(resultWaveUnits) > 0 {
		result["WaveUnits"] = resultWaveUnits
	}

	resultKings := make([]interface{}, len(current.Kings))
	delta = 0

	if len(current.Kings) == 0 || len(old.Kings) == 0 {
		if len(current.Kings) == 0 && len(old.Kings) != 0 {
			result["Kings"] = []interface{}{}
		} else if len(current.Kings) != 0 && len(old.Kings) == 0 {
			result["Kings"] = current.Kings
		}
	} else {
		for index, _ := range current.Kings {
			if index >= len(old.Kings) {
				resultKings[index] = current.Kings[index]
				delta = delta + 1
				continue
			}
			diff := current.Kings[index].Diff(old.Kings[index])
			if len(diff) > 0 {
				resultKings[index] = diff
				delta = delta + 1
			} else {
				resultKings[index] = "__"
			}
		}
		if delta > 0 || len(current.Kings) != len(old.Kings) {
			result["Kings"] = resultKings
		}
	}

	resultKingAttrs := make([]interface{}, len(current.KingAttrs))
	delta = 0

	if len(current.KingAttrs) == 0 || len(old.KingAttrs) == 0 {
		if len(current.KingAttrs) == 0 && len(old.KingAttrs) != 0 {
			result["KingAttrs"] = []interface{}{}
		} else if len(current.KingAttrs) != 0 && len(old.KingAttrs) == 0 {
			result["KingAttrs"] = current.KingAttrs
		}
	} else {
		for index, _ := range current.KingAttrs {
			if index >= len(old.KingAttrs) {
				resultKingAttrs[index] = current.KingAttrs[index]
				delta = delta + 1
				continue
			}
			diff := current.KingAttrs[index].Diff(old.KingAttrs[index])
			if len(diff) > 0 {
				resultKingAttrs[index] = diff
				delta = delta + 1
			} else {
				resultKingAttrs[index] = "__"
			}
		}
		if delta > 0 || len(current.KingAttrs) != len(old.KingAttrs) {
			result["KingAttrs"] = resultKingAttrs
		}
	}

	resultProjectiles := make([]interface{}, len(current.Projectiles))
	delta = 0

	if len(current.Projectiles) == 0 || len(old.Projectiles) == 0 {
		if len(current.Projectiles) == 0 && len(old.Projectiles) != 0 {
			result["Projectiles"] = []interface{}{}
		} else if len(current.Projectiles) != 0 && len(old.Projectiles) == 0 {
			result["Projectiles"] = current.Projectiles
		}
	} else {
		for index, _ := range current.Projectiles {
			if index >= len(old.Projectiles) {
				resultProjectiles[index] = current.Projectiles[index]
				delta = delta + 1
				continue
			}
			diff := current.Projectiles[index].Diff(old.Projectiles[index])
			if len(diff) > 0 {
				resultProjectiles[index] = diff
				delta = delta + 1
			} else {
				resultProjectiles[index] = "__"
			}
		}
		if delta > 0 || len(current.Projectiles) != len(old.Projectiles) {
			result["Projectiles"] = resultProjectiles
		}
	}

	resultNotifications := make([]interface{}, len(current.Notifications))
	delta = 0

	if len(current.Notifications) == 0 || len(old.Notifications) == 0 {
		if len(current.Notifications) == 0 && len(old.Notifications) != 0 {
			result["Notifications"] = []interface{}{}
		} else if len(current.Notifications) != 0 && len(old.Notifications) == 0 {
			result["Notifications"] = current.Notifications
		}
	} else {
		for index, _ := range current.Notifications {
			if index >= len(old.Notifications) {
				resultNotifications[index] = current.Notifications[index]
				delta = delta + 1
				continue
			}
			if old.Notifications[index] != current.Notifications[index] {
				resultNotifications[index] = current.Notifications[index]
				delta = delta + 1
			} else {
				resultNotifications[index] = "__"
			}
		}
		if delta > 0 || len(current.Notifications) != len(old.Notifications) {
			result["Notifications"] = resultNotifications
		}
	}

	if len(result) == 0 {
		return false
	}
	return result
}

func (current StoneUpgrade) Diff(old StoneUpgrade) map[string]interface{} {
	result := map[string]interface{}{}
	delta := 0

	if old.StoneInc != current.StoneInc {
		result["StoneInc"] = current.StoneInc
	}

	if old.MaxGrade != current.MaxGrade {
		result["MaxGrade"] = current.MaxGrade
	}

	if old.CurrGrade != current.CurrGrade {
		result["CurrGrade"] = current.CurrGrade
	}

	if old.InProgress != current.InProgress {
		result["InProgress"] = current.InProgress
	}

	if old.Timer != current.Timer {
		result["Timer"] = current.Timer
	}

	if old.Cooldown != current.Cooldown {
		result["Cooldown"] = current.Cooldown
	}

	resultPrice := make([]interface{}, len(current.Price))
	delta = 0

	if len(current.Price) == 0 || len(old.Price) == 0 {
		if len(current.Price) == 0 && len(old.Price) != 0 {
			result["Price"] = []interface{}{}
		} else if len(current.Price) != 0 && len(old.Price) == 0 {
			result["Price"] = current.Price
		}
	} else {
		for index, _ := range current.Price {
			if index >= len(old.Price) {
				resultPrice[index] = current.Price[index]
				delta = delta + 1
				continue
			}
			if old.Price[index] != current.Price[index] {
				resultPrice[index] = current.Price[index]
				delta = delta + 1
			} else {
				resultPrice[index] = "__"
			}
		}
		if delta > 0 || len(current.Price) != len(old.Price) {
			result["Price"] = resultPrice
		}
	}

	return result
}

func (current FoodUpgrade) Diff(old FoodUpgrade) map[string]interface{} {
	result := map[string]interface{}{}
	
	if old.InProgress != current.InProgress {
		result["InProgress"] = current.InProgress
	}

	if old.Timer != current.Timer {
		result["Timer"] = current.Timer
	}

	if old.Cooldown != current.Cooldown {
		result["Cooldown"] = current.Cooldown
	}

	return result
}

func (current *Player) Diff(old *Player) map[string]interface{} {
	result := map[string]interface{}{}
	delta := 0

	if old.Name != current.Name {
		result["Name"] = current.Name
	}

	if old.Gold != current.Gold {
		result["Gold"] = current.Gold
	}

	if old.Stone != current.Stone {
		result["Stone"] = current.Stone
	}

	resultFood := make([]interface{}, len(current.Food))
	delta = 0

	if len(current.Food) == 0 || len(old.Food) == 0 {
		if len(current.Food) == 0 && len(old.Food) != 0 {
			result["Food"] = []interface{}{}
		} else if len(current.Food) != 0 && len(old.Food) == 0 {
			result["Food"] = current.Food
		}
	} else {
		for index, _ := range current.Food {
			if index >= len(old.Food) {
				resultFood[index] = current.Food[index]
				delta = delta + 1
				continue
			}
			if old.Food[index] != current.Food[index] {
				resultFood[index] = current.Food[index]
				delta = delta + 1
			} else {
				resultFood[index] = "__"
			}
		}
		if delta > 0 || len(current.Food) != len(old.Food) {
			result["Food"] = resultFood
		}
	}

	if old.Income != current.Income {
		result["Income"] = current.Income
	}

	diffStoneGrade := current.StoneGrade.Diff(old.StoneGrade)
	if len(diffStoneGrade) > 0 {
		result["StoneGrade"] = diffStoneGrade
	}

	diffFoodGrade := current.FoodGrade.Diff(old.FoodGrade)
	if len(diffFoodGrade) > 0 {
		result["FoodGrade"] = diffFoodGrade
	}

	if old.Guild != current.Guild {
		result["Guild"] = current.Guild
	}

	resultGuildsList := make([]interface{}, len(current.GuildsList))
	delta = 0

	if len(current.GuildsList) == 0 || len(old.GuildsList) == 0 {
		if len(current.GuildsList) == 0 && len(old.GuildsList) != 0 {
			result["GuildsList"] = []interface{}{}
		} else if len(current.GuildsList) != 0 && len(old.GuildsList) == 0 {
			result["GuildsList"] = current.GuildsList
		}
	} else {
		for index, _ := range current.GuildsList {
			if index >= len(old.GuildsList) {
				resultGuildsList[index] = current.GuildsList[index]
				delta = delta + 1
				continue
			}
			if old.GuildsList[index] != current.GuildsList[index] {
				resultGuildsList[index] = current.GuildsList[index]
				delta = delta + 1
			} else {
				resultGuildsList[index] = "__"
			}
		}
		if delta > 0 || len(current.GuildsList) != len(old.GuildsList) {
			result["GuildsList"] = resultGuildsList
		}
	}

	if old.Resets != current.Resets {
		result["Resets"] = current.Resets
	}

	resultResetPrice := make([]interface{}, len(current.ResetPrice))
	delta = 0

	if len(current.ResetPrice) == 0 || len(old.ResetPrice) == 0 {
		if len(current.ResetPrice) == 0 && len(old.ResetPrice) != 0 {
			result["ResetPrice"] = []interface{}{}
		} else if len(current.ResetPrice) != 0 && len(old.ResetPrice) == 0 {
			result["ResetPrice"] = current.ResetPrice
		}
	} else {
		for index, _ := range current.ResetPrice {
			if index >= len(old.ResetPrice) {
				resultResetPrice[index] = current.ResetPrice[index]
				delta = delta + 1
				continue
			}
			if old.ResetPrice[index] != current.ResetPrice[index] {
				resultResetPrice[index] = current.ResetPrice[index]
				delta = delta + 1
			} else {
				resultResetPrice[index] = "__"
			}
		}
		if delta > 0 || len(current.ResetPrice) != len(old.ResetPrice) {
			result["ResetPrice"] = resultResetPrice
		}
	}

	resultAvailableUnits := make([]interface{}, len(current.AvailableUnits))
	delta = 0

	if len(current.AvailableUnits) == 0 || len(old.AvailableUnits) == 0 {
		if len(current.AvailableUnits) == 0 && len(old.AvailableUnits) != 0 {
			result["AvailableUnits"] = []interface{}{}
		} else if len(current.AvailableUnits) != 0 && len(old.AvailableUnits) == 0 {
			result["AvailableUnits"] = current.AvailableUnits
		}
	} else {
		for index, _ := range current.AvailableUnits {
			if index >= len(old.AvailableUnits) {
				resultAvailableUnits[index] = current.AvailableUnits[index]
				delta = delta + 1
				continue
			}
			if old.AvailableUnits[index] != current.AvailableUnits[index] {
				resultAvailableUnits[index] = current.AvailableUnits[index]
				delta = delta + 1
			} else {
				resultAvailableUnits[index] = "__"
			}
		}
		if delta > 0 || len(current.AvailableUnits) != len(old.AvailableUnits) {
			result["AvailableUnits"] = resultAvailableUnits
		}
	}

	resultSummonUnits := map[string]interface{}{}
	for key, val := range current.SummonUnits {
		if oldSummonUnits, ok := old.SummonUnits[key]; ok {
			diff := val.Diff(oldSummonUnits)
			if len(diff) > 0 {
				resultSummonUnits[key] = diff
			}
		} else {
			resultSummonUnits[key] = val
		}
	}
	for key, _ := range old.SummonUnits {
		if _, ok := current.SummonUnits[key]; !ok {
			resultSummonUnits[key] = "_del_"
		}
	}
	if len(resultSummonUnits) > 0 {
		result["SummonUnits"] = resultSummonUnits
	}

	resultIncomeUnits := make([]interface{}, len(current.IncomeUnits))
	delta = 0

	if len(current.IncomeUnits) == 0 || len(old.IncomeUnits) == 0 {
		if len(current.IncomeUnits) == 0 && len(old.IncomeUnits) != 0 {
			result["IncomeUnits"] = []interface{}{}
		} else if len(current.IncomeUnits) != 0 && len(old.IncomeUnits) == 0 {
			result["IncomeUnits"] = current.IncomeUnits
		}
	} else {
		for index, _ := range current.IncomeUnits {
			if index >= len(old.IncomeUnits) {
				resultIncomeUnits[index] = current.IncomeUnits[index]
				delta = delta + 1
				continue
			}
			diff := current.IncomeUnits[index].Diff(old.IncomeUnits[index])
			if len(diff) > 0 {
				resultIncomeUnits[index] = diff
				delta = delta + 1
			} else {
				resultIncomeUnits[index] = "__"
			}
		}
		if delta > 0 || len(current.IncomeUnits) != len(old.IncomeUnits) {
			result["IncomeUnits"] = resultIncomeUnits
		}
	}

	if old.Value != current.Value {
		result["Value"] = current.Value
	}

	if old.Leaked != current.Leaked {
		result["Leaked"] = current.Leaked
	}

	if old.Score != current.Score {
		result["Score"] = current.Score
	}

	if old.PauseCap != current.PauseCap {
		result["PauseCap"] = current.PauseCap
	}

	return result
}

func (current *Projectile) Diff(old *Projectile) map[string]interface{} {
	result := map[string]interface{}{}
	delta := 0

	if old.Id != current.Id {
		result["Id"] = current.Id
	}

	if old.Name != current.Name {
		result["Name"] = current.Name
	}

	if old.Movespeed != current.Movespeed {
		result["Movespeed"] = current.Movespeed
	}

	if old.Trajectory != current.Trajectory {
		result["Trajectory"] = current.Trajectory
	}

	resultCoords := make([]interface{}, len(current.Coords))
	delta = 0

	if len(current.Coords) == 0 || len(old.Coords) == 0 {
		if len(current.Coords) == 0 && len(old.Coords) != 0 {
			result["Coords"] = []interface{}{}
		} else if len(current.Coords) != 0 && len(old.Coords) == 0 {
			result["Coords"] = current.Coords
		}
	} else {
		for index, _ := range current.Coords {
			if index >= len(old.Coords) {
				resultCoords[index] = current.Coords[index]
				delta = delta + 1
				continue
			}
			if old.Coords[index] != current.Coords[index] {
				resultCoords[index] = current.Coords[index]
				delta = delta + 1
			} else {
				resultCoords[index] = "__"
			}
		}
		if delta > 0 || len(current.Coords) != len(old.Coords) {
			result["Coords"] = resultCoords
		}
	}

	resultDir := make([]interface{}, len(current.Dir))
	delta = 0

	if len(current.Dir) == 0 || len(old.Dir) == 0 {
		if len(current.Dir) == 0 && len(old.Dir) != 0 {
			result["Dir"] = []interface{}{}
		} else if len(current.Dir) != 0 && len(old.Dir) == 0 {
			result["Dir"] = current.Dir
		}
	} else {
		for index, _ := range current.Dir {
			if index >= len(old.Dir) {
				resultDir[index] = current.Dir[index]
				delta = delta + 1
				continue
			}
			if old.Dir[index] != current.Dir[index] {
				resultDir[index] = current.Dir[index]
				delta = delta + 1
			} else {
				resultDir[index] = "__"
			}
		}
		if delta > 0 || len(current.Dir) != len(old.Dir) {
			result["Dir"] = resultDir
		}
	}

	if old.Ability != current.Ability {
		result["Ability"] = current.Ability
	}

	return result
}

func (current *Effect) Diff(old *Effect) map[string]interface{} {
	result := map[string]interface{}{}
	
	if old.Owner != current.Owner {
		result["Owner"] = current.Owner
	}

	if old.Duration != current.Duration {
		result["Duration"] = current.Duration
	}

	return result
}

func (current *Unit) Diff(old *Unit) map[string]interface{} {
	result := map[string]interface{}{}
	delta := 0

	if old.Id != current.Id {
		result["Id"] = current.Id
	}

	if old.Name != current.Name {
		result["Name"] = current.Name
	}

	if old.Projectile != current.Projectile {
		result["Projectile"] = current.Projectile
	}

	if old.ProjSpeed != current.ProjSpeed {
		result["ProjSpeed"] = current.ProjSpeed
	}

	if old.ProjTraj != current.ProjTraj {
		result["ProjTraj"] = current.ProjTraj
	}

	if old.Player != current.Player {
		result["Player"] = current.Player
	}

	resultAtk := make([]interface{}, len(current.Atk))
	delta = 0

	if len(current.Atk) == 0 || len(old.Atk) == 0 {
		if len(current.Atk) == 0 && len(old.Atk) != 0 {
			result["Atk"] = []interface{}{}
		} else if len(current.Atk) != 0 && len(old.Atk) == 0 {
			result["Atk"] = current.Atk
		}
	} else {
		for index, _ := range current.Atk {
			if index >= len(old.Atk) {
				resultAtk[index] = current.Atk[index]
				delta = delta + 1
				continue
			}
			if old.Atk[index] != current.Atk[index] {
				resultAtk[index] = current.Atk[index]
				delta = delta + 1
			} else {
				resultAtk[index] = "__"
			}
		}
		if delta > 0 || len(current.Atk) != len(old.Atk) {
			result["Atk"] = resultAtk
		}
	}

	if old.AtkRange != current.AtkRange {
		result["AtkRange"] = current.AtkRange
	}

	if old.Aspd != current.Aspd {
		result["Aspd"] = current.Aspd
	}

	if old.AspdTimer != current.AspdTimer {
		result["AspdTimer"] = current.AspdTimer
	}

	if old.AtkType != current.AtkType {
		result["AtkType"] = current.AtkType
	}

	if old.Def != current.Def {
		result["Def"] = current.Def
	}

	if old.DefType != current.DefType {
		result["DefType"] = current.DefType
	}

	if old.Movespeed != current.Movespeed {
		result["Movespeed"] = current.Movespeed
	}

	if old.MaxHp != current.MaxHp {
		result["MaxHp"] = current.MaxHp
	}

	if old.Hp != current.Hp {
		result["Hp"] = current.Hp
	}

	if old.HpReg != current.HpReg {
		result["HpReg"] = current.HpReg
	}

	if old.MaxMp != current.MaxMp {
		result["MaxMp"] = current.MaxMp
	}

	if old.Mp != current.Mp {
		result["Mp"] = current.Mp
	}

	if old.MpReg != current.MpReg {
		result["MpReg"] = current.MpReg
	}

	if old.Price != current.Price {
		result["Price"] = current.Price
	}

	if old.Food != current.Food {
		result["Food"] = current.Food
	}

	if old.Bounty != current.Bounty {
		result["Bounty"] = current.Bounty
	}

	if old.Size != current.Size {
		result["Size"] = current.Size
	}

	if old.Guild != current.Guild {
		result["Guild"] = current.Guild
	}

	if old.Tier != current.Tier {
		result["Tier"] = current.Tier
	}

	resultCoords := make([]interface{}, len(current.Coords))
	delta = 0

	if len(current.Coords) == 0 || len(old.Coords) == 0 {
		if len(current.Coords) == 0 && len(old.Coords) != 0 {
			result["Coords"] = []interface{}{}
		} else if len(current.Coords) != 0 && len(old.Coords) == 0 {
			result["Coords"] = current.Coords
		}
	} else {
		for index, _ := range current.Coords {
			if index >= len(old.Coords) {
				resultCoords[index] = current.Coords[index]
				delta = delta + 1
				continue
			}
			if old.Coords[index] != current.Coords[index] {
				resultCoords[index] = current.Coords[index]
				delta = delta + 1
			} else {
				resultCoords[index] = "__"
			}
		}
		if delta > 0 || len(current.Coords) != len(old.Coords) {
			result["Coords"] = resultCoords
		}
	}

	resultDir := make([]interface{}, len(current.Dir))
	delta = 0

	if len(current.Dir) == 0 || len(old.Dir) == 0 {
		if len(current.Dir) == 0 && len(old.Dir) != 0 {
			result["Dir"] = []interface{}{}
		} else if len(current.Dir) != 0 && len(old.Dir) == 0 {
			result["Dir"] = current.Dir
		}
	} else {
		for index, _ := range current.Dir {
			if index >= len(old.Dir) {
				resultDir[index] = current.Dir[index]
				delta = delta + 1
				continue
			}
			if old.Dir[index] != current.Dir[index] {
				resultDir[index] = current.Dir[index]
				delta = delta + 1
			} else {
				resultDir[index] = "__"
			}
		}
		if delta > 0 || len(current.Dir) != len(old.Dir) {
			result["Dir"] = resultDir
		}
	}

	resultTargetCoords := make([]interface{}, len(current.TargetCoords))
	delta = 0

	if len(current.TargetCoords) == 0 || len(old.TargetCoords) == 0 {
		if len(current.TargetCoords) == 0 && len(old.TargetCoords) != 0 {
			result["TargetCoords"] = []interface{}{}
		} else if len(current.TargetCoords) != 0 && len(old.TargetCoords) == 0 {
			result["TargetCoords"] = current.TargetCoords
		}
	} else {
		for index, _ := range current.TargetCoords {
			if index >= len(old.TargetCoords) {
				resultTargetCoords[index] = current.TargetCoords[index]
				delta = delta + 1
				continue
			}
			if old.TargetCoords[index] != current.TargetCoords[index] {
				resultTargetCoords[index] = current.TargetCoords[index]
				delta = delta + 1
			} else {
				resultTargetCoords[index] = "__"
			}
		}
		if delta > 0 || len(current.TargetCoords) != len(old.TargetCoords) {
			result["TargetCoords"] = resultTargetCoords
		}
	}

	if old.TargetSign != current.TargetSign {
		result["TargetSign"] = current.TargetSign
	}

	if old.Slot != current.Slot {
		result["Slot"] = current.Slot
	}

	if old.Waypoint != current.Waypoint {
		result["Waypoint"] = current.Waypoint
	}

	if old.SellPrice != current.SellPrice {
		result["SellPrice"] = current.SellPrice
	}

	if old.Action != current.Action {
		result["Action"] = current.Action
	}

	if old.Killer != current.Killer {
		result["Killer"] = current.Killer
	}

	resultUpgrades := make([]interface{}, len(current.Upgrades))
	delta = 0

	if len(current.Upgrades) == 0 || len(old.Upgrades) == 0 {
		if len(current.Upgrades) == 0 && len(old.Upgrades) != 0 {
			result["Upgrades"] = []interface{}{}
		} else if len(current.Upgrades) != 0 && len(old.Upgrades) == 0 {
			result["Upgrades"] = current.Upgrades
		}
	} else {
		for index, _ := range current.Upgrades {
			if index >= len(old.Upgrades) {
				resultUpgrades[index] = current.Upgrades[index]
				delta = delta + 1
				continue
			}
			if old.Upgrades[index] != current.Upgrades[index] {
				resultUpgrades[index] = current.Upgrades[index]
				delta = delta + 1
			} else {
				resultUpgrades[index] = "__"
			}
		}
		if delta > 0 || len(current.Upgrades) != len(old.Upgrades) {
			result["Upgrades"] = resultUpgrades
		}
	}

	resultAbilities := map[string]interface{}{}
	for key, val := range current.Abilities {
		if oldAbilities, ok := old.Abilities[key]; ok {
			diff := stringArrayDiff(val, oldAbilities)
			if _, ok := diff.(bool); !ok {
				resultAbilities[key] = diff
			}
		} else {
			resultAbilities[key] = val
		}
	}
	for key, _ := range old.Abilities {
		if _, ok := current.Abilities[key]; !ok {
			resultAbilities[key] = "_del_"
		}
	}
	if len(resultAbilities) > 0 {
		result["Abilities"] = resultAbilities
	}

	resultAffected := map[string]interface{}{}
	for key, val := range current.Affected {
		if oldAffected, ok := old.Affected[key]; ok {
			diff := val.Diff(oldAffected)
			if len(diff) > 0 {
				resultAffected[key] = diff
			}
		} else {
			resultAffected[key] = val
		}
	}
	for key, _ := range old.Affected {
		if _, ok := current.Affected[key]; !ok {
			resultAffected[key] = "_del_"
		}
	}
	if len(resultAffected) > 0 {
		result["Affected"] = resultAffected
	}

	return result
}

func (current *Summon) Diff(old *Summon) map[string]interface{} {
	result := map[string]interface{}{}
	
	if old.Current != current.Current {
		result["Current"] = current.Current
	}

	if old.Cap != current.Cap {
		result["Cap"] = current.Cap
	}

	if old.Timer != current.Timer {
		result["Timer"] = current.Timer
	}

	if old.Cooldown != current.Cooldown {
		result["Cooldown"] = current.Cooldown
	}

	return result
}

func stringArrayDiff(current, old []string) interface{} {
	result := make([]interface{}, len(current))
	delta := 0

	if len(current) == 0 || len(old) == 0 {
		if len(current) == 0 && len(old) != 0 {
			return []interface{}{}
		} else if len(current) != 0 && len(old) == 0 {
			return current
		} else {
			return false
		}
	} else {
		for index, _ := range current {
			if index >= len(old) {
				result[index] = current[index]
				delta = delta + 1
				continue
			}
			if old[index] != current[index] {
				result[index] = current[index]
				delta = delta + 1
			} else {
				result[index] = "__"
			}
		}
		if delta == 0 && len(current) == len(old) {
			return false
		}
	}
	return result
}

func UnitArrayDiff(current, old []*Unit) interface{} {
	result := make([]interface{}, len(current))
	delta := 0

	if len(current) == 0 || len(old) == 0 {
		if len(current) == 0 && len(old) != 0 {
			return []interface{}{}
		} else if len(current) != 0 && len(old) == 0 {
			return current
		} else {
			return false
		}
	} else {
		for index, _ := range current {
			if index >= len(old) {
				result[index] = current[index]
				delta = delta + 1
				continue
			}
			diff := current[index].Diff(old[index])
			if len(diff) > 0 {
				result[index] = diff
				delta = delta + 1
			} else {
				result[index] = "__"
			}
		}
		if delta == 0 && len(current) == len(old) {
			return false
		}
	}
	return result
}

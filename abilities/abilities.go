package abilities

import (
	"fmt"
	"math"
	"ws/game"
	"ws/game/units"
	"ws/models"
	"ws/utils"
	"ws/utils/status"
)

// ADDERS

var AddModifierKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"duration": 0.0,
		"rad":      0.0,
	}
	kwargs = override(kw, kwargs)
	rad := kwargs["rad"].(float64)
	name := kwargs["name"].(string)
	duration := kwargs["duration"].(int)

	return func(options *models.AbOptions) {
		unit := options.Unit
		target := options.Target

		if status.IsMarkedDead(target) {
			return
		}

		if (utils.XYRange(unit.Coords, target.Coords) - target.Size) > rad {
			return
		}

		target.Affected[name] = &models.Effect{
			Owner:    unit.Player,
			Duration: duration,
		}
		//fmt.Println(name, "add", unit.Name, unit.Id, "=>", target.Name, target.Id, "for time: ", duration)
	}
}

var AddOnHitEffectKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"duration": 0.0,
	}
	kwargs = override(kw, kwargs)
	name := kwargs["name"].(string)
	duration := kwargs["duration"].(int)

	return func(options *models.AbOptions) {
		unit := options.Unit
		target := options.Target

		target.Affected[name] = &models.Effect{
			Owner:    unit.Player,
			Duration: duration,
		}
		//fmt.Println(name, "add", unit.Name, unit.Id, "=>", target.Name, target.Id, "for time: ", duration)
	}
}

var AddOnHitStepEffectKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"duration": 0.0,
		"step":     0.0,
	}
	kwargs = override(kw, kwargs)
	name := kwargs["name"].(string)
	duration := kwargs["duration"].(int)
	step := kwargs["step"].(int)

	return func(options *models.AbOptions) {
		unit := options.Unit
		target := options.Target

		if _, ok := target.Affected[name]; ok {
			target.Affected[name].Owner = unit.Player
			target.Affected[name].Duration = duration + target.Affected[name].Duration%step
			return
		}

		target.Affected[name] = &models.Effect{
			Owner:    unit.Player,
			Duration: duration,
		}
	}
}

var AddStepModifierKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"duration": 0.0,
		"rad":      0.0,
	}
	kwargs = override(kw, kwargs)
	rad := kwargs["rad"].(float64)
	name := kwargs["name"].(string)
	duration := kwargs["duration"].(int)

	return func(options *models.AbOptions) {
		unit := options.Unit
		target := options.Target

		if status.IsMarkedDead(target) {
			return
		}

		if (utils.XYRange(unit.Coords, target.Coords) - target.Size) > rad {
			return
		}

		if _, ok := target.Affected[name]; !ok {
			target.Affected[name] = &models.Effect{
				Owner:    unit.Player,
				Duration: duration,
			}
		}
		//fmt.Println(name, "add", unit.Name, "::", unit.Id, "=>", target.Name, "::", target.Id, "for time: ", duration)
	}
}

var AddRangedModifierKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"duration": 0.0,
		"rad":      0.0,
	}
	kwargs = override(kw, kwargs)
	rad := kwargs["rad"].(float64)
	name := kwargs["name"].(string)
	duration := kwargs["duration"].(int)

	return func(options *models.AbOptions) {
		unit := options.Unit
		target := options.Target

		if status.IsMarkedDead(target) {
			return
		}

		if !status.IsRanged(target) {
			return
		}

		if (utils.XYRange(unit.Coords, target.Coords) - target.Size) > rad {
			return
		}

		target.Affected[name] = &models.Effect{
			Owner:    unit.Player,
			Duration: duration,
		}
		//fmt.Println(name, "add", unit.Name, unit.Id, "=>", target.Name, target.Id, "for time: ", duration)
	}
}

var AddMeleeModifierKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"duration": 0.0,
		"rad":      0.0,
	}
	kwargs = override(kw, kwargs)
	rad := kwargs["rad"].(float64)
	name := kwargs["name"].(string)
	duration := kwargs["duration"].(int)

	return func(options *models.AbOptions) {
		unit := options.Unit
		target := options.Target

		if status.IsMarkedDead(target) {
			return
		}

		if status.IsRanged(target) {
			return
		}

		if (utils.XYRange(unit.Coords, target.Coords) - target.Size) > rad {
			return
		}

		target.Affected[name] = &models.Effect{
			Owner:    unit.Player,
			Duration: duration,
		}
		//fmt.Println(name, "add", unit.Name, unit.Id, "=>", target.Name, target.Id, "for time: ", duration)
	}
}

var AddProcStatusKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"chance":   0.2,
		"duration": 1.0,
		"status":   "stun",
		"dmg":      20.0,
	}
	kwargs = override(kw, kwargs)
	chance := kwargs["chance"].(float64)
	statusName := kwargs["status"].(string)
	duration := kwargs["duration"].(int)
	dmg := kwargs["dmg"].(float64)

	return func(options *models.AbOptions) {
		unit := options.Unit
		target := options.Target

		if utils.Uniform(0, 1) > chance {
			return
		}
		target.Hp -= dmg
		status.SetKiller(unit.Player, target)
		RefreshStatus(target, unit, statusName, duration)
	}
}

var AddDivineHymnKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"castTime": 0.1,
		"rad":      3.0,
		"cost":     15.0,
		"duration": 5.0,
	}
	kwargs = override(kw, kwargs)
	name := kwargs["name"].(string)
	rad := kwargs["rad"].(float64)
	cost := kwargs["cost"].(float64)
	duration := kwargs["duration"].(int)
	castTime := kwargs["castTime"].(float64)
	return func(options *models.AbOptions) {
		unit := options.Unit
		unitCurrent := options.UnitCurrent
		allies := options.Allies

		if unitCurrent.Mp < cost || status.IsCasting(unit) {
			return
		}

		unit.Mp = unit.Mp - cost
		unit.Action = game.Cast
		unit.AspdTimer = castTime

		chance := utils.Uniform(0, 1)
		switch {
		case chance < 0.33:
			name = "Hymn of blaze"
		case chance >= 0.33 && chance < 0.66:
			name = "Hymn of pleading"
		case chance >= 0.66 && chance < 0.99:
			name = "Hymn of magnificence"
		default:
			name = "Hymn of hymns"
		}
		for _, ally := range allies {
			if status.IsMarkedDead(ally) {
				continue
			}
			if (utils.XYRange(unit.Coords, ally.Coords) - ally.Size) > rad {
				continue
			}

			if name != "Hymn of hymns" {
				ally.Affected[name] = &models.Effect{
					Owner:    unit.Player,
					Duration: duration,
				}
				continue
			}
			ally.Affected["Hymn of blaze"] = &models.Effect{
				Owner:    unit.Player,
				Duration: duration,
			}
			ally.Affected["Hymn of pleading"] = &models.Effect{
				Owner:    unit.Player,
				Duration: duration,
			}
			ally.Affected["Hymn of magnificence"] = &models.Effect{
				Owner:    unit.Player,
				Duration: duration,
			}
			fmt.Println(name, ally.Name, ally.Id)
		}
	}
}

var MultiAttackKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"limit": 1,
	}
	kwargs = override(kw, kwargs)
	limit := kwargs["limit"].(int)

	return func(options *models.AbOptions) {
		unit := options.Unit
		enemies := options.Enemies
		state := options.State
		links := options.Links
		hits := 0

		for _, enemy := range enemies {
			if hits >= limit {
				return
			}
			if status.IsMarkedDead(enemy) {
				continue
			}
			if utils.XYRange(unit.Coords, enemy.Coords)-unit.Size-enemy.Size > unit.AtkRange*1.2 {
				continue
			}
			/*if !status.CanAttack(unit, enemy) {
				continue
			}*/
			hits += 1
			optionsTemp := *options
			optionsTemp.Target = enemy
			Create(state, links, unit, enemy, &optionsTemp, game.Attack)
		}
	}
}

var BouncingAttackKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"limit": 1,
		"coef":  1.0,
	}
	kwargs = override(kw, kwargs)
	limit := kwargs["limit"].(int)
	coef := kwargs["coef"].(float64)

	return func(options *models.AbOptions) {
		unit := options.Unit
		target := options.Target
		enemies := options.Enemies
		state := options.State
		links := options.Links
		//AttackHit(unit, target, options)
		options.Limit -= 1
		if options.Limit == -1 {
			options.Limit = limit
		}
		if options.Limit <= 0 {
			return
		}
		options.Atk[0] = coef * options.Atk[0]
		options.Atk[1] = coef * options.Atk[1]
		atkRange := math.Pow(coef, float64(limit-options.Limit+1)) * unit.AtkRange
		//fmt.Println("dmg", options.Atk[0], options.Atk[1], "limit", options.Limit)
		for _, enemy := range enemies {
			if target.Id == enemy.Id || status.IsMarkedDead(enemy) {
				continue
			}
			if (utils.XYRange(target.Coords, enemy.Coords) - enemy.Size) > atkRange {
				continue
			}
			options.Target = enemy
			options.ProjCoords = target.Coords
			Create(state, links, unit, enemy, options, game.Attack)
			return
		}
	}
}

var SplinterAttackKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"rad":   0.0,
		"coef":  0.0,
		"limit": 0,
	}
	kwargs = override(kw, kwargs)
	rad := kwargs["rad"].(float64)
	coef := kwargs["coef"].(float64)
	limit := kwargs["limit"].(int)

	return func(options *models.AbOptions) {
		unit := options.Unit
		target := options.Target
		enemies := options.Enemies
		atk := utils.Uniform(options.Atk[0], options.Atk[1])
		hits := 0
		for _, enemy := range enemies {
			if hits >= limit {
				return
			}
			if target.Id == enemy.Id || status.IsMarkedDead(enemy) {
				continue
			}
			if (utils.XYRange(target.Coords, enemy.Coords) - enemy.Size) > rad {
				continue
			}
			hits += 1
			enemy.Hp = enemy.Hp - atk*coef*status.DefCoef(unit, enemy)
			//fmt.Println(enemy.Name, "lost", atk*coef, "hp")
			status.SetKiller(unit.Player, enemy)
		}
	}
}

var CastOnAttackKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"name":       "Spell blast",
		"cost":       10.0,
		"projectile": "spell_blast",
	}
	kwargs = override(kw, kwargs)
	name := kwargs["name"].(string)
	cost := kwargs["cost"].(float64)
	projectile := kwargs["projectile"].(string)
	return func(options *models.AbOptions) {
		unit := options.Unit
		unitCurrent := options.UnitCurrent
		target := options.Target
		state := options.State
		links := options.Links

		Create(state, links, unit, target, options, game.Attack)
		if unitCurrent.Mp < cost {
			return
		}
		unit.Mp -= cost
		optionsTemp := *options
		optionsTemp.Projectile = projectile
		optionsTemp.ProjSpeed = unit.ProjSpeed
		optionsTemp.ProjTraj = unit.ProjTraj
		Create(state, links, unit, target, &optionsTemp, name)
	}
}

var SpreadMissilesKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"name":       "Solar spike",
		"rad":        2.75,
		"limit":      2,
		"projectile": "burning_ember",
		"projspeed":  3.0,
		"projtraj":   "",
	}
	kwargs = override(kw, kwargs)
	name := kwargs["name"].(string)
	cost := kwargs["cost"].(float64)
	rad := kwargs["rad"].(float64)
	limit := kwargs["limit"].(int)
	projectile := kwargs["projectile"].(string)
	projSpeed := kwargs["projspeed"].(float64)
	projTraj := kwargs["projtraj"].(string)

	return func(options *models.AbOptions) {
		unit := options.Unit
		unitCurrent := options.UnitCurrent
		enemies := options.Enemies
		state := options.State
		links := options.Links

		if unitCurrent.Mp < cost {
			return
		}
		unit.Mp -= cost
		hits := 0
		for _, enemy := range enemies {
			if hits >= limit {
				return
			}
			if status.IsMarkedDead(enemy) {
				continue
			}
			if (utils.XYRange(unit.Coords, enemy.Coords) - unit.Size - enemy.Size) > rad {
				continue
			}
			hits += 1
			optionsTemp := *options
			optionsTemp.Target = enemy
			optionsTemp.Projectile = projectile
			optionsTemp.ProjSpeed = projSpeed
			optionsTemp.ProjTraj = projTraj
			Create(state, links, unit, enemy, &optionsTemp, name)
		}
	}
}

var ApplyDamageKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"dmg": 50.0,
	}
	kwargs = override(kw, kwargs)
	dmg := kwargs["dmg"].(float64)

	return func(options *models.AbOptions) {
		unit := options.Unit
		target := options.Target
		target.Hp -= dmg
		status.SetKiller(unit.Player, target)
	}
}

var AddSolarSpikeKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"name":       "Solar spike",
		"rad":        2.75,
		"limit":      2,
		"projectile": "solar_spike",
		"projspeed":  3.0,
		"projtraj":   "pohui",
	}
	kwargs = override(kw, kwargs)
	name := kwargs["name"].(string)
	rad := kwargs["rad"].(float64)
	limit := kwargs["limit"].(int)
	projectile := kwargs["projectile"].(string)
	projSpeed := kwargs["projspeed"].(float64)
	projTraj := kwargs["projtraj"].(string)
	return func(options *models.AbOptions) {
		unit := options.Unit
		target := options.Target
		enemies := options.Enemies
		state := options.State
		links := options.Links
		hits := 0

		for _, enemy := range enemies {
			if hits >= limit {
				return
			}
			if target.Id == enemy.Id || status.IsMarkedDead(enemy) {
				continue
			}
			if (utils.XYRange(target.Coords, enemy.Coords) - enemy.Size) > rad {
				continue
			}
			hits += 1

			optionsTemp := *options
			optionsTemp.Target = enemy
			optionsTemp.Projectile = projectile
			optionsTemp.ProjSpeed = projSpeed
			optionsTemp.ProjTraj = projTraj
			optionsTemp.ProjCoords = target.Coords
			Create(state, links, unit, enemy, &optionsTemp, name)
		}
	}
}

var AddRocketBarrageKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"name":       "Rocket barrage",
		"abRange":    2.75,
		"cost":       10.0,
		"projectile": "pocket rocket",
		"projspeed":  3.0,
		"projtraj":   "pohui",
		"castTime":   0.3,
	}
	kwargs = override(kw, kwargs)
	name := kwargs["name"].(string)
	abRange := kwargs["abRange"].(float64)
	cost := kwargs["cost"].(float64)
	projectile := kwargs["projectile"].(string)
	projSpeed := kwargs["projspeed"].(float64)
	projTraj := kwargs["projtraj"].(string)
	castTime := kwargs["castTime"].(float64)
	return func(options *models.AbOptions) {
		unit := options.Unit
		unitCurrent := options.UnitCurrent
		target := options.Target
		state := options.State
		links := options.Links

		if unitCurrent.Mp < cost || status.IsCasting(unit) {
			return
		}

		if utils.XYRange(unitCurrent.Coords, target.Coords) > abRange {
			return
		}

		unit.Mp = unit.Mp - cost
		unit.Action = game.Cast
		unit.AspdTimer = castTime

		options.Projectile = projectile
		options.ProjSpeed = projSpeed
		options.ProjTraj = projTraj
		Create(state, links, unit, target, options, name)
	}
}

var ApplyRocketBarrageKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"name":     "Rocket barrage",
		"rad":      0.85,
		"dmg":      10.0,
		"duration": 1.0,
	}
	kwargs = override(kw, kwargs)
	rad := kwargs["rad"].(float64)
	dmg := kwargs["dmg"].(float64)
	duration := kwargs["duration"].(int)

	return func(options *models.AbOptions) {
		unit := options.Unit
		target := options.Target
		enemies := options.Enemies

		for _, enemy := range enemies {
			if status.IsMarkedDead(enemy) {
				continue
			}
			if (utils.XYRange(target.Coords, enemy.Coords) - enemy.Size) > rad {
				continue
			}
			enemy.Hp = enemy.Hp - dmg
			status.SetKiller(unit.Player, enemy)
			RefreshStatus(enemy, unit, "stun", duration)
		}
	}
}

var ReflectMeleeKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"coef": 0.0,
	}
	kwargs = override(kw, kwargs)
	coef := kwargs["coef"].(float64)

	return func(options *models.AbOptions) {
		unit := options.Unit
		if status.IsRanged(unit) {
			return
		}
		target := options.Target

		unit.Hp = unit.Hp - options.Dmg*coef
		status.SetKiller(target.Player, unit)
	}
}

var AtkTypeModifierKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"coef":    0.0,
		"atkType": "",
	}
	kwargs = override(kw, kwargs)
	coef := kwargs["coef"].(float64)
	atkType := kwargs["atkType"].(string)

	return func(options *models.AbOptions) {
		unit := options.Unit

		if unit.AtkType != atkType {
			return
		}
		options.Dmg = options.Dmg * (1 + coef)
	}
}

var ProcDmgModifierKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"chance": 0.0,
		"value":  0.0,
	}
	kwargs = override(kw, kwargs)
	chance := kwargs["chance"].(float64)
	value := kwargs["value"].(float64)

	return func(options *models.AbOptions) {
		if utils.Uniform(0, 1) > chance {
			return
		}
		options.Dmg = options.Dmg + value
		if options.Dmg < 0 {
			options.Dmg = 0
		}
	}
}

var ManaShieldKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"name":     "Illusive barrier",
		"cost":     1.0,
		"modifier": 0.0,
		"duration": 1.0,
	}
	kwargs = override(kw, kwargs)
	name := kwargs["name"].(string)
	cost := kwargs["cost"].(float64)
	modifier := kwargs["modifier"].(float64)
	duration := kwargs["duration"].(int)

	return func(options *models.AbOptions) {
		unit := options.Target
		if unit.Mp <= 0 {
			return
		}
		if options.Dmg <= 0 {
			return
		}

		dmg := modifier * options.Dmg
		mpCost := cost * dmg
		if unit.Mp < mpCost {
			mpCost = unit.Mp
			dmg = mpCost / cost
		}

		options.Dmg -= dmg
		if options.Dmg < 0 {
			options.Dmg = 0
		}
		unit.Mp -= mpCost

		unit.Affected[name] = &models.Effect{
			Owner:    unit.Player,
			Duration: duration,
		}
	}
}

var ManaBladeKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"modifier": 1.0,
	}
	kwargs = override(kw, kwargs)
	modifier := kwargs["modifier"].(float64)

	return func(options *models.AbOptions) {
		unit := options.UnitCurrent
		target := options.Target
		bonusAtk := modifier * unit.Mp
		options.Dmg += bonusAtk * status.DefCoef(unit, target)
		options.Atk[0] += bonusAtk
		options.Atk[1] += bonusAtk
	}
}

var TyrannicideKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"modifier": 2.0,
	}
	kwargs = override(kw, kwargs)
	modifier := kwargs["modifier"].(float64)

	return func(options *models.AbOptions) {
		target := options.Target
		if status.IsKing(target) {
			options.Dmg = options.Dmg * modifier
			options.Atk[0] = options.Atk[0] * modifier
			options.Atk[1] = options.Atk[1] * modifier
		}
	}
}

var GetSourceFromDmgKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"source": "hp",
		"coef":   0.0,
	}
	kwargs = override(kw, kwargs)
	source := kwargs["source"].(string)
	coef := kwargs["coef"].(float64)

	return func(options *models.AbOptions) {
		unit := options.Unit

		switch source {
		case "hp":
			unit.Hp += coef * options.Dmg
			if unit.Hp > unit.MaxHp {
				unit.Hp = unit.MaxHp
			}
		case "mp":
			unit.Mp += coef * options.Dmg
			if unit.Mp > unit.MaxMp {
				unit.Mp = unit.MaxMp
			}
		default:
			return
		}
	}
}

var RestoreSourceKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"castTime": 0.1,
		"source":   "hp",
		"rad":      6.0,
		"cost":     10.0,
		"value":    130.0,
	}
	kwargs = override(kw, kwargs)
	name := kwargs["name"].(string)
	rad := kwargs["rad"].(float64)
	source := kwargs["source"].(string)
	cost := kwargs["cost"].(float64)
	value := kwargs["value"].(float64)
	castTime := kwargs["castTime"].(float64)
	return func(options *models.AbOptions) {
		unit := options.Unit
		unitCurrent := options.UnitCurrent
		allies := options.Allies

		if unitCurrent.Mp < cost || status.IsCasting(unit) {
			return
		}

		for _, ally := range allies {
			if status.IsMarkedDead(ally) {
				continue
			}
			if (utils.XYRange(unit.Coords, ally.Coords) - ally.Size) > rad {
				continue
			}
			switch source {
			case "hp":
				//if ally.MaxHp-ally.Hp < value || (ally.MaxHp <= value && ally.Hp > 0.5*ally.MaxHp) {
				if ally.Hp == ally.MaxHp {
					continue
				}
			case "mp":
				//if ally.MaxMp-ally.Mp < value || (ally.MaxMp <= value && ally.Mp > 0.5*ally.MaxMp) {
				if ally.Mp == ally.MaxMp {
					continue
				}
			default:
				return
			}

			unit.Mp = unit.Mp - cost
			unit.Action = game.Cast
			unit.AspdTimer = castTime
			switch source {
			case "hp":
				ally.Hp = ally.Hp + value
				if ally.Hp > ally.MaxHp {
					ally.Hp = ally.MaxHp
				}
			case "mp":
				ally.Mp = ally.Mp + value
				if ally.Mp > ally.MaxMp {
					ally.Mp = ally.MaxMp
				}
			default:
				return
			}
			unit.Dir = ally.Coords
			ally.Affected[name] = &models.Effect{
				Owner:    unit.Player,
				Duration: 1,
			}
			fmt.Println(name, ally.Name, ally.Id, "for", value)
			return
		}
	}
}

var AddBuffKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"castTime": 0.1,
		"rad":      4.0,
		"cost":     10.0,
		"duration": 10.0,
	}
	kwargs = override(kw, kwargs)
	name := kwargs["name"].(string)
	rad := kwargs["rad"].(float64)
	cost := kwargs["cost"].(float64)
	duration := kwargs["duration"].(int)
	castTime := kwargs["castTime"].(float64)
	return func(options *models.AbOptions) {
		unit := options.Unit
		unitCurrent := options.UnitCurrent
		allies := options.Allies

		if unitCurrent.Mp < cost || status.IsCasting(unit) {
			return
		}

		for _, ally := range allies {
			if status.IsMarkedDead(ally) {
				continue
			}
			if (utils.XYRange(unit.Coords, ally.Coords) - ally.Size) > rad {
				continue
			}
			if _, ok := ally.Affected[name]; ok {
				continue
			}

			unit.Mp = unit.Mp - cost
			unit.Action = game.Cast
			unit.AspdTimer = castTime
			ally.Affected[name] = &models.Effect{
				Owner:    unit.Player,
				Duration: duration,
			}

			unit.Dir = ally.Coords

			//fmt.Println(name, ally.Name, ally.Id, "for", value)
			return
		}
	}
}

var AddSelfBuffKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"castTime": 0.1,
		"cost":     30.0,
		"duration": 15.0,
	}
	kwargs = override(kw, kwargs)
	cost := kwargs["cost"].(float64)
	name := kwargs["name"].(string)
	duration := kwargs["duration"].(int)
	castTime := kwargs["castTime"].(float64)
	return func(options *models.AbOptions) {
		unit := options.Unit
		unitCurrent := options.UnitCurrent

		if unitCurrent.Mp < cost || status.IsCasting(unit) {
			return
		}

		if _, ok := unit.Affected[name]; ok {
			return
		}

		unit.Mp = unit.Mp - cost
		unit.Action = game.Cast
		unit.AspdTimer = castTime
		unit.Affected[name] = &models.Effect{
			Owner:    unit.Player,
			Duration: duration,
		}
		fmt.Println(name, "add", unit.Name, unit.Id, "=> for time: ", duration)
	}
}

var SpawnUnitKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"unitName": "Orc Bruiser",
	}
	kwargs = override(kw, kwargs)
	unitName := kwargs["unitName"].(string)

	return func(options *models.AbOptions) {
		unit := options.Unit
		state := options.State
		stateCurrent := options.StateCurrent

		unitTemp, exists := units.UNITS_BY_NAME[unitName]
		if !exists {
			return
		}

		unitTemp.Id = utils.UniqId()
		unitTemp.Slot = unit.Slot
		unitTemp.Player = unit.Player
		unitTemp.Waypoint = unit.Waypoint
		unitTemp.Action = game.Idle
		unitTemp.Coords = unit.Coords
		unitTemp.Affected = make(map[string]*models.Effect)
		if !status.IsWave(unit) {
			state.PlayerUnits[unit.Player] = append(state.PlayerUnits[unit.Player], &unitTemp)
			stateCurrent.PlayerUnits[unit.Player] = append(stateCurrent.PlayerUnits[unit.Player], &unitTemp)
		} else if !status.IsInKingSlot(unit) {
			state.WaveUnits[unit.Slot] = append(state.WaveUnits[unit.Slot], &unitTemp)
			stateCurrent.WaveUnits[unit.Slot] = append(stateCurrent.WaveUnits[unit.Slot], &unitTemp)
		} else {
			for slot, _ := range state.WaveUnits {
				state.WaveUnits[slot] = append(state.WaveUnits[slot], &unitTemp)
				stateCurrent.WaveUnits[unit.Slot] = append(stateCurrent.WaveUnits[unit.Slot], &unitTemp)
				break
			}
		}
		fmt.Println(unit.Name, unit.Id, "spawned", unitTemp.Name, unitTemp.Id)
	}
}

// APPLIERS

var ApplyModifierKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"attribute": "atk",
		"modifier":  0.0,
	}
	kwargs = override(kw, kwargs)
	attribute := kwargs["attribute"].(string)
	modifier := kwargs["modifier"].(float64)

	return func(options *models.AbOptions) {
		unit := options.Unit

		switch attribute {
		case "atk":
			unit.Atk[0] = unit.Atk[0] * (1 + modifier)
			unit.Atk[1] = unit.Atk[1] * (1 + modifier)
		case "maxhp":
			unit.MaxHp = unit.MaxHp * (1 + modifier)
		case "def":
			unit.Def = unit.Def * (1 + modifier)
		case "aspd":
			unit.Aspd = unit.Aspd * (1 + modifier)
		case "movespeed":
			unit.Movespeed = unit.Movespeed * (1 + modifier)
		default:
			return
		}
		//fmt.Println("Modified ", unit.Name, unit.Id, " attr: ", kwargs["attribute"])
	}
}

var ApplyValueModifierKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"attribute": "atk",
		"value":     0.0,
	}
	kwargs = override(kw, kwargs)
	attribute := kwargs["attribute"].(string)
	value := kwargs["value"].(float64)

	return func(options *models.AbOptions) {
		unit := options.Unit

		switch attribute {
		case "atk":
			unit.Atk[0] = unit.Atk[0] + value
			unit.Atk[1] = unit.Atk[1] + value
		case "maxhp":
			unit.MaxHp = unit.MaxHp + value
		case "def":
			unit.Def = unit.Def + value
		case "aspd":
			unit.Aspd = unit.Aspd + value
		case "movespeed":
			unit.Movespeed = unit.Movespeed + value
		default:
			return
		}
		//fmt.Println("Modified ", unit.Name, unit.Id, " attr: ", kwargs["attribute"])
	}
}

var ChangeSourceByValueKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"name":   "Breath of life",
		"source": "hp",
		"value":  3.0,
		"step":   1.0,
	}
	kwargs = override(kw, kwargs)
	name := kwargs["name"].(string)
	source := kwargs["source"].(string)
	value := kwargs["value"].(float64)
	step := kwargs["step"].(int)

	return func(options *models.AbOptions) {
		unit := options.Unit
		//fmt.Println("Unit", unit.Name, "::", unit.Id, name, unit.Affected[name].Duration, step, unit.Affected[name].Duration%step)
		if unit.Affected[name].Duration%step != 0 {
			return
		}
		switch source {
		case "hp":
			unit.Hp = unit.Hp + value
			if unit.Hp > unit.MaxHp {
				unit.Hp = unit.MaxHp
			}
		case "mp":
			unit.Mp = unit.Mp + value
			if unit.Mp > unit.MaxMp {
				unit.Mp = unit.MaxMp
			}
			if unit.Mp < 0 {
				unit.Mp = 0
			}
		default:
			return
		}
		status.SetKiller(unit.Affected[name].Owner, unit)
	}
}

var ChangeSourceByCoefKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"name":   "Breath of life",
		"source": "hp",
		"coef":   0.01,
		"step":   1.0,
	}
	kwargs = override(kw, kwargs)
	name := kwargs["name"].(string)
	source := kwargs["source"].(string)
	coef := kwargs["coef"].(float64)
	step := kwargs["step"].(int)

	return func(options *models.AbOptions) {
		unit := options.Unit
		//fmt.Println("Unit", unit.Name, "::", unit.Id, name, unit.Affected[name].Duration, step, unit.Affected[name].Duration%step)
		if unit.Affected[name].Duration%step != 0 {
			return
		}
		switch source {
		case "hp":
			unit.Hp += coef * unit.MaxHp
			if unit.Hp > unit.MaxHp {
				unit.Hp = unit.MaxHp
			}
		case "mp":
			unit.Mp += coef * unit.MaxMp
			if unit.Mp > unit.MaxMp {
				unit.Mp = unit.MaxMp
			}
			if unit.Mp < 0 {
				unit.Mp = 0
			}
		default:
			return
		}
		status.SetKiller(unit.Affected[name].Owner, unit)
	}
}

var ApplySpeedModifiersKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"aspdModifier": -0.30,
		"mspdModifier": 0.20,
	}
	kwargs = override(kw, kwargs)
	aspdModifier := kwargs["aspdModifier"].(float64)
	mspdModifier := kwargs["mspdModifier"].(float64)

	return func(options *models.AbOptions) {
		unit := options.Unit

		unit.Aspd = unit.Aspd * (1 + aspdModifier)
		unit.Movespeed = unit.Movespeed * (1 + mspdModifier)
	}
}

var ApplyPoisonKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"step":         1.0,
		"dmg":          20.0,
		"aspdModifier": 0.20,
		"mspdModifier": -0.20,
	}
	kwargs = override(kw, kwargs)
	name := kwargs["name"].(string)
	step := kwargs["step"].(int)
	dmg := kwargs["dmg"].(float64)
	aspdModifier := kwargs["aspdModifier"].(float64)
	mspdModifier := kwargs["mspdModifier"].(float64)

	return func(options *models.AbOptions) {
		unit := options.Unit

		unit.Aspd = unit.Aspd * (1 + aspdModifier)
		unit.Movespeed = unit.Movespeed * (1 + mspdModifier)

		if unit.Affected[name].Duration%step != 0 {
			return
		}

		unit.Hp = unit.Hp - dmg
		status.SetKiller(unit.Affected[name].Owner, unit)
	}
}

var AddAbilityKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"abType":  "",
		"ability": "",
	}
	kwargs = override(kw, kwargs)
	abType := kwargs["abType"].(string)
	ability := kwargs["ability"].(string)

	return func(options *models.AbOptions) {
		unit := options.Unit

		if _, ok := unit.Abilities[abType]; ok {
			if !utils.ContainsString(unit.Abilities[abType], ability) {
				unit.Abilities[abType] = append(unit.Abilities[abType], ability)
			}
			return
		}

		unit.Abilities[abType] = []string{ability}
	}
}

var ApplyHymnOfBlazeKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"atkModifier":  0.03,
		"aspdModifier": -0.05,
	}
	kwargs = override(kw, kwargs)
	atkModifier := kwargs["atkModifier"].(float64)
	aspdModifier := kwargs["aspdModifier"].(float64)

	return func(options *models.AbOptions) {
		unit := options.Unit

		unit.Atk[0] = unit.Atk[0] * (1 + atkModifier)
		unit.Atk[1] = unit.Atk[1] * (1 + atkModifier)
		unit.Aspd = unit.Aspd * (1 + aspdModifier)
		//fmt.Println("Blazed ", unit.Name, unit.Id)
	}
}

var ApplyHymnOfPleadingKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"defValue":     1.0,
		"mspdModifier": 0.1,
	}
	kwargs = override(kw, kwargs)
	defValue := kwargs["defValue"].(float64)
	mspdModifier := kwargs["mspdModifier"].(float64)

	return func(options *models.AbOptions) {
		unit := options.Unit

		unit.Def = unit.Def + defValue
		unit.Movespeed = unit.Movespeed * (1 + mspdModifier)
		//fmt.Println("Pleaded ", unit.Name, unit.Id)
	}
}

var ApplyHymnOfMagnificenceKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"maxHpModifier": 0.1,
		"mpRegModifier": 0.1,
	}
	kwargs = override(kw, kwargs)
	maxHpModifier := kwargs["maxHpModifier"].(float64)
	mpRegModifier := kwargs["mpRegModifier"].(float64)

	return func(options *models.AbOptions) {
		unit := options.Unit

		unit.MaxHp = unit.MaxHp * (1 + maxHpModifier)
		unit.MpReg = unit.MpReg * (1 + mpRegModifier)
		//fmt.Println("Magnified ", unit.Name, unit.Id)
	}
}

var KingsBladeKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{}
	kwargs = override(kw, kwargs)

	return func(options *models.AbOptions) {
		unit := options.Unit
		target := options.Target
		state := options.State
		enemies := options.Enemies
		if !status.IsKing(unit) {
			return
		}
		atkGrade := state.KingAttrs[unit.Player%game.KingCoef].AtkGrade.CurrGrade
		rad := 0.1 * float64(atkGrade/10+1)
		coef := 0.01 * float64(atkGrade)
		for _, enemy := range enemies {
			if target.Id == enemy.Id || status.IsMarkedDead(enemy) {
				continue
			}
			if (utils.XYRange(target.Coords, enemy.Coords) - enemy.Size) > rad {
				continue
			}
			enemy.Hp -= options.Dmg * coef
			status.SetKiller(unit.Player, enemy)
		}
	}
}

var KingsBloodKw = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	var kw = map[string]interface{}{
		"name": "King's blood",
		"step": 1.0,
	}
	kwargs = override(kw, kwargs)
	name := kwargs["name"].(string)
	step := kwargs["step"].(int)

	return func(options *models.AbOptions) {
		unit := options.Unit
		state := options.State
		if unit.Affected[name].Owner/game.KingCoef != 1 {
			return
		}
		if unit.Affected[name].Duration%step != 0 {
			return
		}
		value := -float64(state.KingAttrs[unit.Affected[name].Owner%game.KingCoef].HpRegGrade.CurrGrade)
		unit.Hp += value
		if unit.Hp > unit.MaxHp {
			unit.Hp = unit.MaxHp
		}
		status.SetKiller(unit.Affected[name].Owner, unit)
	}
}

// ABILITIES

var ABILITIES = map[string]*models.Ability{
	"Bash": &models.Ability{
		AddMethod:   game.OnHit,
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: AddProcStatusKw,
			Options: map[string]interface{}{
				"chance":   0.2,
				"duration": 1.0,
				"status":   "stun",
				"dmg":      0.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"Blazing crush": &models.Ability{
		AddMethod:   game.OnHit,
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: AddProcStatusKw,
			Options: map[string]interface{}{
				"chance":   0.25,
				"duration": 1.0,
				"status":   "stun",
				"dmg":      120.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"Booming charges": &models.Ability{
		AddMethod:   game.OnHit,
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: BouncingAttackKw,
			Options: map[string]interface{}{
				"limit": 2,
				"coef":  0.75,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"Bouncing sparks": &models.Ability{
		AddMethod:   game.OnHit,
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: BouncingAttackKw,
			Options: map[string]interface{}{
				"limit": 1,
				"coef":  0.75,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"Breath of life": &models.Ability{
		AddMethod:   game.BuffAura,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AddStepModifierKw,
			Options: map[string]interface{}{
				"rad":      4.0,
				"duration": 1.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ChangeSourceByValueKw,
			Options: map[string]interface{}{
				"source": "hp",
				"value":  3.0,
				"step":   1.0,
			},
		},
	},
	"Burning embers": &models.Ability{
		AddMethod:   game.Active,
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: SpreadMissilesKw,
			Options: map[string]interface{}{
				"cost":       3.0,
				"limit":      3,
				"rad":        3.5,
				"projectile": "burning_ember",
				"projspeed":  3.0,
				"projtraj":   "",
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplyDamageKw,
			Options: map[string]interface{}{
				"dmg": 10.0,
			},
		},
	},
	"Commander aura": &models.Ability{
		AddMethod:   game.BuffAura,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AddModifierKw,
			Options: map[string]interface{}{
				"rad":      4.0,
				"duration": 2.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplyModifierKw,
			Options: map[string]interface{}{
				"attribute": "atk",
				"modifier":  0.04, //0.04
			},
		},
	},
	"Corruption": &models.Ability{
		AddMethod:   game.OnHit,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AddOnHitEffectKw,
			Options: map[string]interface{}{
				"duration": 4.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplyValueModifierKw,
			Options: map[string]interface{}{
				"attribute": "def",
				"value":     -4.0,
			},
		},
	},
	"Demonic might": &models.Ability{
		AddMethod:   game.BuffAura,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AddModifierKw,
			Options: map[string]interface{}{
				"rad":      4.0,
				"duration": 2.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplyModifierKw,
			Options: map[string]interface{}{
				"attribute": "atk",
				"modifier":  0.10,
			},
		},
	},
	"Demonic swiftness": &models.Ability{
		AddMethod:   game.BuffAura,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AddModifierKw,
			Options: map[string]interface{}{
				"rad":      4.0,
				"duration": 2.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplyModifierKw,
			Options: map[string]interface{}{
				"attribute": "movespeed",
				"modifier":  0.15,
			},
		},
	},
	"Divine hymn": &models.Ability{
		AddMethod:   game.Assist,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AddDivineHymnKw,
			Options: map[string]interface{}{
				"castTime": 0.8,
				"rad":      3.0,
				"cost":     15.0,
				"duration": 10.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"Hymn of blaze": &models.Ability{
		AddMethod:   game.Assist,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplyHymnOfBlazeKw,
			Options: map[string]interface{}{
				"atkModifier":  0.03,
				"aspdModifier": -0.05,
			},
		},
	},
	"Hymn of pleading": &models.Ability{
		AddMethod:   game.Assist,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplyHymnOfPleadingKw,
			Options: map[string]interface{}{
				"defValue":     1.0,
				"mspdModifier": 0.1,
			},
		},
	},
	"Hymn of magnificence": &models.Ability{
		AddMethod:   game.Assist,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplyHymnOfMagnificenceKw,
			Options: map[string]interface{}{
				"maxHpModifier": 0.1,
				"mpRegModifier": 0.1,
			},
		},
	},
	"Feast": &models.Ability{
		AddMethod:   game.BuffAura,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AddMeleeModifierKw,
			Options: map[string]interface{}{
				"rad":      4.0,
				"duration": 2.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: AddAbilityKw,
			Options: map[string]interface{}{
				"abType":  game.OnHit,
				"ability": "_(OnHit)Feast",
			},
		},
	},
	"_(OnHit)Feast": &models.Ability{
		AddMethod:   game.OnHit,
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: GetSourceFromDmgKw,
			Options: map[string]interface{}{
				"source": "hp",
				"coef":   0.17,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"Greater heal": &models.Ability{
		AddMethod:   game.Assist,
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: RestoreSourceKw,
			Options: map[string]interface{}{
				"castTime": 0.5,
				"source":   "hp",
				"rad":      6.0,
				"cost":     10.0,
				"value":    300.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"Great shield": &models.Ability{
		AddMethod:   game.Defensive,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AtkTypeModifierKw,
			Options: map[string]interface{}{
				"coef":    -0.35,
				"atkType": game.Piercing,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"Heal": &models.Ability{
		AddMethod:   game.Assist,
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: RestoreSourceKw,
			Options: map[string]interface{}{
				"castTime": 0.5,
				"source":   "hp",
				"rad":      6.0,
				"cost":     10.0,
				"value":    130.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"Heavy plate": &models.Ability{
		AddMethod:   game.Defensive,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: ProcDmgModifierKw,
			Options: map[string]interface{}{
				"chance": 0.3,
				"value":  -50.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"Illusive barrier": &models.Ability{
		AddMethod:   game.Defensive,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: ManaShieldKw,
			Options: map[string]interface{}{
				"cost":     0.25,
				"modifier": 0.9,
				"duration": 1.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"Illusive impulse": &models.Ability{
		AddMethod:   game.Offensive,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: ManaBladeKw,
			Options: map[string]interface{}{
				"modifier": 0.33,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"Mana tides": &models.Ability{
		AddMethod:   game.BuffAura,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AddStepModifierKw,
			Options: map[string]interface{}{
				"rad":      2.0,
				"duration": 1.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ChangeSourceByCoefKw,
			Options: map[string]interface{}{
				"source": "mp",
				"coef":   0.05,
				"step":   1.0,
			},
		},
	},
	"On spices": &models.Ability{
		AddMethod:   game.Assist,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AddBuffKw,
			Options: map[string]interface{}{
				"castTime": 0.2,
				"rad":      4.0,
				"cost":     10.0,
				"duration": 10.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplySpeedModifiersKw,
			Options: map[string]interface{}{
				"aspdModifier": -0.20,
				"mspdModifier": 0.20,
			},
		},
	},
	"Piercing howl": &models.Ability{
		AddMethod:   game.Aura,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AddModifierKw,
			Options: map[string]interface{}{
				"rad":      4.0,
				"duration": 2.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplySpeedModifiersKw,
			Options: map[string]interface{}{
				"aspdModifier": 0.20,
				"mspdModifier": -0.20,
			},
		},
	},
	"Slowing arrows": &models.Ability{
		AddMethod:   game.OnHit,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AddOnHitEffectKw,
			Options: map[string]interface{}{
				"duration": 5.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplySpeedModifiersKw,
			Options: map[string]interface{}{
				"aspdModifier": 0.50,
				"mspdModifier": -0.50,
			},
		},
	},
	"Solar spike": &models.Ability{
		AddMethod:   game.OnHit,
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: AddSolarSpikeKw,
			Options: map[string]interface{}{
				"rad":        2.0,
				"limit":      2,
				"projectile": "solar spike",
				"projspeed":  3.0,
				"projtraj":   "pohui",
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplyDamageKw,
			Options: map[string]interface{}{
				"dmg": 150.0,
			},
		},
	},
	"Spell blast": &models.Ability{
		AddMethod:   game.AttackReplacer,
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: CastOnAttackKw,
			Options: map[string]interface{}{
				"projectile": "spell_blast",
				"cost":       10.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplyDamageKw,
			Options: map[string]interface{}{
				"dmg": 60.0,
			},
		},
	},
	"Steel fur": &models.Ability{
		AddMethod:   game.BuffAura,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AddModifierKw,
			Options: map[string]interface{}{
				"rad":      4.0,
				"duration": 2.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplyValueModifierKw,
			Options: map[string]interface{}{
				"attribute": "def",
				"value":     3.0,
			},
		},
	},
	"Strange poison": &models.Ability{
		AddMethod:   game.OnHit,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AddOnHitStepEffectKw,
			Options: map[string]interface{}{
				"duration": 10.0,
				"step":     1.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplyPoisonKw,
			Options: map[string]interface{}{
				"step":         1.0,
				"dmg":          20.0,
				"aspdModifier": 0.20,
				"mspdModifier": -0.20,
			},
		},
	},
	"Strange regeneration": &models.Ability{
		AddMethod:   game.BuffAura,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AddStepModifierKw,
			Options: map[string]interface{}{
				"rad":      4.0,
				"duration": 1.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ChangeSourceByValueKw,
			Options: map[string]interface{}{
				"source": "hp",
				"value":  10.0,
				"step":   1.0,
			},
		},
	},
	"Tectonic flame": &models.Ability{
		AddMethod:   game.Aura,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AddStepModifierKw,
			Options: map[string]interface{}{
				"rad":      4.0,
				"duration": 1.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ChangeSourceByValueKw,
			Options: map[string]interface{}{
				"source": "hp",
				"value":  -10.0,
				"step":   1.0,
			},
		},
	},
	"Tyrannicide": &models.Ability{
		AddMethod:   game.Offensive,
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: TyrannicideKw,
			Options: map[string]interface{}{
				"modifier": 2.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"Welding arc": &models.Ability{
		AddMethod:   game.React,
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: ReflectMeleeKw,
			Options: map[string]interface{}{
				"coef": 0.16,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"Rocket barrage": &models.Ability{
		AddMethod:   game.Active,
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: AddRocketBarrageKw,
			Options: map[string]interface{}{
				"castTime":   0.5,
				"abRange":    2.75,
				"cost":       10.0,
				"projectile": "rocket_barrage",
				"projspeed":  3.0,
				"projtraj":   "pohui",
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplyRocketBarrageKw,
			Options: map[string]interface{}{
				"rad":      0.85,
				"dmg":      10.0,
				"duration": 1.0,
			},
		},
	},
	"Blessed quiver": &models.Ability{
		AddMethod:   game.AttackReplacer,
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: MultiAttackKw,
			Options: map[string]interface{}{
				"limit": 5,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"Splinter shot": &models.Ability{
		AddMethod:   game.OnHit,
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: SplinterAttackKw,
			Options: map[string]interface{}{
				"rad":   0.575,
				"coef":  0.25,
				"limit": 2,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"Scope": &models.Ability{
		AddMethod:   game.BuffAura,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AddRangedModifierKw,
			Options: map[string]interface{}{
				"rad":      3.0,
				"duration": 2.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplyModifierKw,
			Options: map[string]interface{}{
				"attribute": "atk",
				"modifier":  0.07,
			},
		},
	},
	"Scope 2.0": &models.Ability{
		AddMethod:   game.BuffAura,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AddRangedModifierKw,
			Options: map[string]interface{}{
				"rad":      3.0,
				"duration": 2.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplyModifierKw,
			Options: map[string]interface{}{
				"attribute": "atk",
				"modifier":  0.14,
			},
		},
	},
	"Overclocking": &models.Ability{
		AddMethod:   game.Assist,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AddSelfBuffKw,
			Options: map[string]interface{}{
				"castTime": 0.5,
				"cost":     30.0,
				"duration": 15.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplySpeedModifiersKw,
			Options: map[string]interface{}{
				"aspdModifier": -0.30,
				"mspdModifier": 0.20,
			},
		},
	},
	"Orc driver": &models.Ability{
		AddMethod:   game.OnDeath,
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: SpawnUnitKw,
			Options: map[string]interface{}{
				"unitName": "Orc Bruiser",
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"Auto-armored orc driver": &models.Ability{
		AddMethod:   game.OnDeath,
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: SpawnUnitKw,
			Options: map[string]interface{}{
				"unitName": "Auto-armored Orc",
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"stun": &models.Ability{
		AddMethod:   "status",
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"Final charge": &models.Ability{
		AddMethod:   game.Aura,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: ApplyValueModifierKw,
			Options: map[string]interface{}{
				"attribute": "def",
				"value":     2.0,
			},
		},
	},
	"Split splash": &models.Ability{
		AddMethod:   game.OnHit,
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: BouncingAttackKw,
			Options: map[string]interface{}{
				"limit": 1,
				"coef":  0.75,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"King's blade": &models.Ability{
		AddMethod:   game.OnHit,
		ApplyMethod: game.Instant,
		AddLogic: &models.Logic{
			HandlerKw: KingsBladeKw,
			Options:   map[string]interface{}{},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: Empty,
			Options:   map[string]interface{}{},
		},
	},
	"King's blood": &models.Ability{
		AddMethod:   game.Aura,
		ApplyMethod: game.Modify,
		AddLogic: &models.Logic{
			HandlerKw: AddStepModifierKw,
			Options: map[string]interface{}{
				"rad":      4.0,
				"duration": 1.0,
			},
		},
		ApplyLogic: &models.Logic{
			HandlerKw: KingsBloodKw,
			Options: map[string]interface{}{
				"step": 1.0,
			},
		},
	},
}

func init() {
	for name, _ := range ABILITIES {
		add_logic := ABILITIES[name].AddLogic
		add_logic.Options["name"] = name
		if _, ok := add_logic.Options["duration"]; ok {
			duration := add_logic.Options["duration"].(float64)
			add_logic.Options["duration"] = int(duration / game.WAR_STEP_DELAY)
		}
		if _, ok := add_logic.Options["step"]; ok {
			step := add_logic.Options["step"].(float64)
			add_logic.Options["step"] = int(step / game.WAR_STEP_DELAY)
		}
		if _, ok := add_logic.Options["castTime"]; ok {
			castTime := add_logic.Options["castTime"].(float64)
			add_logic.Options["castTime"] = castTime / game.WAR_STEP_DELAY
		}
		add_logic.Handler = add_logic.HandlerKw(add_logic.Options)

		apply_logic := ABILITIES[name].ApplyLogic
		apply_logic.Options["name"] = name
		if _, ok := apply_logic.Options["duration"]; ok {
			duration := apply_logic.Options["duration"].(float64)
			apply_logic.Options["duration"] = int(duration / game.WAR_STEP_DELAY)
		}
		if _, ok := apply_logic.Options["step"]; ok {
			step := apply_logic.Options["step"].(float64)
			apply_logic.Options["step"] = int(step / game.WAR_STEP_DELAY)
		}
		apply_logic.Handler = apply_logic.HandlerKw(apply_logic.Options)
	}
}

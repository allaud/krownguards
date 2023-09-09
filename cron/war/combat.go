package war

import (
	"fmt"
	"ws/abilities"
	"ws/game"
	"ws/models"
	//"ws/utils"
	"ws/utils/status"
)

func CombatStep(stateTemp, state *models.State, links *models.Links, friendly bool) {
	allUnitsTemp := []*models.Unit{}
	allUnits := []*models.Unit{}

	if ArenaPhase(state) {
		allUnitsTemp = GetUnitsBySlots(stateTemp, ArenaSlots(friendly), true)
		allUnits = GetUnitsBySlots(state, ArenaSlots(friendly), true)
	} else {
		allUnitsTemp = GetUnitsBySlots(stateTemp, []int{1, 2, 3, 4, 5, 6, 7, 8}, friendly)
		allUnits = GetUnitsBySlots(state, []int{1, 2, 3, 4, 5, 6, 7, 8}, friendly)
	}
	for uIndex, unitTemp := range allUnitsTemp {
		options := &models.AbOptions{}
		unit := allUnits[uIndex]
		alliesTemp := GetSideUnits(unitTemp, stateTemp, friendly)
		allies := GetSideUnits(unitTemp, state, friendly)
		enemiesTemp := GetSideUnits(unitTemp, stateTemp, !friendly)
		enemies := GetSideUnits(unitTemp, state, !friendly)

		if unit.Id != unitTemp.Id {
			fmt.Println(unit.Name, unit.Id, "is not", unitTemp.Name, unitTemp.Id, "!")
		}
		if len(allies) != len(alliesTemp) {
			fmt.Println(len(allies), "!=", len(alliesTemp))
		}
		if len(enemies) != len(enemiesTemp) {
			fmt.Println(len(enemies), "!=", len(enemiesTemp))
		}

		// trigger death of unit
		if !status.IsAlive(unitTemp) {
			TriggerDeath(stateTemp, state, unitTemp, unit, options)
			continue
		}

		// Proceed spell casting
		if status.IsCasting(unit) {
			unit.AspdTimer -= 1
			if unit.AspdTimer <= 0 {
				unit.Action = game.Idle
			}
		}

		// Refresh animation
		if IsControlled(unitTemp) {
			unit.Action = game.Idle
		}

		// applying buff auras
		if _, ok := unitTemp.Abilities[game.BuffAura]; ok {
			for aIndex, _ := range alliesTemp {
				ally := allies[aIndex]
				options.Unit = unitTemp
				options.Target = ally
				abilities.AddByType(unitTemp, abilities.ABILITIES, options, game.BuffAura)
			}
		}

		// apply assist skills
		if !IsControlled(unitTemp) && !status.IsCasting(unit) {
			options.Unit = unit
			options.UnitCurrent = unitTemp
			options.Allies = allies
			abilities.AddByType(unitTemp, abilities.ABILITIES, options, game.Assist)
		}

		// if unit isn't attacking set attackspeed delay
		if !IsAttacking(unitTemp) && !status.IsCasting(unit) {
			SetDelay(unit, HalfAspd(unit)-1)
		}

		// searching enemy target to attack
		targetTemp := &models.Unit{}
		targetTemp = nil
		target := &models.Unit{}
		target = nil
		for eIndex, enemyTemp := range enemiesTemp {
			enemy := enemies[eIndex]
			if !status.IsAlive(enemyTemp) {
				continue
			}
			if !IsVisible(unitTemp, enemyTemp) {
				continue
			}
			if (targetTemp == nil) || IsCloser(unitTemp, enemyTemp, targetTemp) {
				targetTemp = enemyTemp
				target = enemy
			}
			// applying debuff auras
			options.Unit = unitTemp
			options.Target = enemy
			abilities.AddByType(unitTemp, abilities.ABILITIES, options, game.Aura)
		}

		if !IsControlled(unitTemp) && !status.IsCasting(unit) {
			if targetTemp != nil {
				if target.Id != targetTemp.Id {
					fmt.Println(target.Name, target.Id, "is not", targetTemp.Name, targetTemp.Id, "!")
				}
				// apply active skills
				options.Unit = unit
				options.UnitCurrent = unitTemp
				options.Target = target
				options.Enemies = enemies
				options.State = state
				options.Links = links
				abilities.AddByType(unitTemp, abilities.ABILITIES, options, game.Active)
				// if unit isn't casting then try to attack
				if !status.IsCasting(unit) {
					if status.CanAttack(unitTemp, targetTemp) {
						options.Enemies = enemies
						Attack(state, links, unitTemp, target, unit, options)
					} else {
						MoveToCoords(unit, allies, targetTemp.Coords)
					}
				}
			} else {
				MoveToWaypoint(unit, allies, state)
			}
		}

		// reset affected attrs
		RefreshAttributes(state, unit)
		// apply affected effects lol
		options.Unit = unit
		options.State = state
		abilities.ApplyEffects(abilities.ABILITIES, options)
		RegenSources(unit)
	}
}

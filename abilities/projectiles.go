package abilities

import (
	"ws/game"
	"ws/models"
	"ws/utils"
	"ws/utils/status"
)

var ProjId = func() func() int {
	var id = 0

	return func() int {
		id = id + 1
		return id
	}
}()

func ResetProjId() {
	ProjId = func() func() int {
		var id = 0

		return func() int {
			id = id + 1
			return id
		}
	}()
}

func AttackHit(unit, target *models.Unit, options *models.AbOptions) {
	// get result dmg from attacking unit
	dmg := utils.Uniform(options.Atk[0], options.Atk[1]) * status.DefCoef(unit, target)
	options.Dmg = dmg
	// use all offensive abilities
	AddByType(unit, ABILITIES, options, game.Offensive)
	// use all target defensive abilities
	AddByType(target, ABILITIES, options, game.Defensive)
	target.Hp = target.Hp - options.Dmg
	// apply onHit effects
	AddByType(unit, ABILITIES, options, game.OnHit)
	// apply react effects from target
	AddByType(target, ABILITIES, options, game.React)
	status.SetKiller(unit.Player, target)
}

func Step(state *models.State, links *models.Links) {
	index := len(state.Projectiles) - 1
	for index >= 0 {
		if links.ProjLinks[index].Triggered {
			state.Projectiles[index] = nil
			links.ProjLinks[index] = nil
			state.Projectiles = append(state.Projectiles[:index], state.Projectiles[index+1:]...)
			links.ProjLinks = append(links.ProjLinks[:index], links.ProjLinks[index+1:]...)
		}
		index = index - 1
	}
	for pIndex, projectile := range state.Projectiles {
		Move(projectile, links.ProjLinks[pIndex])
		if IsHit(projectile, links.ProjLinks[pIndex]) {
			Trigger(state, links, pIndex)
		}
	}
}

func Create(state *models.State, links *models.Links, unit, target *models.Unit, options *models.AbOptions, ability string) {
	//startCoords := utils.GetVectorCoords(unit.Coords, target.Coords, unit.Size, 0)
	startCoords := unit.Coords
	// if proj start coords redefined
	if options.ProjCoords != [2]float64{0, 0} {
		startCoords = options.ProjCoords
	}
	// autoattack projectile
	prName := unit.Projectile
	prSpeed := unit.ProjSpeed
	if prName == game.Melee {
		prSpeed = utils.XYRange(startCoords, target.Coords) - target.Size
	}
	prTraj := unit.ProjTraj
	// ability projectile
	if ability != game.Attack {
		prName = options.Projectile
		prSpeed = options.ProjSpeed
		prTraj = options.ProjTraj
	}
	// id only for visible ranged projectiles
	prId := 0
	if prName != game.Melee && prName != game.Invisible {
		prId = ProjId()
	}
	projTemp := models.Projectile{
		Id:         prId,
		Name:       prName,
		Movespeed:  prSpeed,
		Trajectory: prTraj,
		Coords:     startCoords,
		Dir:        target.Coords,
		Ability:    ability,
	}
	projLinkTemp := models.ProjLink{
		Unit:    unit,
		Target:  target,
		Options: options,
	}

	state.Projectiles = append(state.Projectiles, &projTemp)
	links.ProjLinks = append(links.ProjLinks, &projLinkTemp)
}

func Move(projectile *models.Projectile, projLink *models.ProjLink) {
	// refresh direction
	projectile.Dir = projLink.Target.Coords
	coords := projectile.Coords
	dir := projectile.Dir
	step := projectile.Movespeed
	// placing new projectile on start position
	if coords == projLink.Unit.Coords {
		projectile.Coords = utils.GetVectorCoords(projLink.Unit.Coords, projLink.Target.Coords, projLink.Unit.Size+projLink.Unit.ProjStart[0], projLink.Unit.ProjStart[1])
		return
	}
	projectile.Coords = utils.GetVectorCoords(coords, dir, step, 0)
}

func Trigger(state *models.State, links *models.Links, index int) {
	unit := links.ProjLinks[index].Unit
	target := links.ProjLinks[index].Target
	options := links.ProjLinks[index].Options
	abName := state.Projectiles[index].Ability
	// check if it's autoattack projectile
	if abName == game.Attack {
		AttackHit(unit, target, options)
	} else {
		ApplyAbility(abName, options)
	}
	// check projectile triggered
	links.ProjLinks[index].Triggered = true
}

func IsHit(projectile *models.Projectile, projLink *models.ProjLink) bool {
	return utils.XYRange(projectile.Coords, projectile.Dir) <= projLink.Target.Size
}

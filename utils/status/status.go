package status

import (
	"math"
	"ws/game"
	"ws/models"
	"ws/utils"
)

func IsWave(unit *models.Unit) bool {
	return (unit.Tier == game.Wave) || (unit.Tier == game.Summon)
}

func IsSummon(unit *models.Unit) bool {
	return unit.Tier == game.Summon
}

func IsKing(unit *models.Unit) bool {
	return unit.Tier == game.King
}

func IsAlive(unit *models.Unit) bool {
	return unit.Hp > 0
}

func IsMarkedDead(unit *models.Unit) bool {
	return unit.Action == game.Dead
}

func IsRanged(unit *models.Unit) bool {
	return unit.Projectile != game.Melee
}

func IsCasting(unit *models.Unit) bool {
	return unit.Action == game.Cast
}

func IsInKingSlot(unit *models.Unit) bool {
	return unit.Slot/game.KingCoef == 1
}

func CanAttack(unit, target *models.Unit) bool {
	cRange := utils.XYRange(unit.Coords, target.Coords)
	return unit.AtkRange >= cRange-unit.Size-target.Size
}

func SetKiller(killer int, target *models.Unit) {
	if IsAlive(target) || !IsWave(target) || target.Killer != 0 {
		return
	}
	target.Killer = killer
}

func DefCoef(unit, target *models.Unit) float64 {
	typeCoef := game.AtkDefMultipliers[unit.AtkType][target.DefType]
	defCoef := 1 - math.Copysign(1-math.Abs(math.Pow(1.09, -math.Abs(target.Def))), target.Def)
	return typeCoef * defCoef
}

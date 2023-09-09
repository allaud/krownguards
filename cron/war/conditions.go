package war

import (
	"ws/game"
	"ws/models"
	"ws/utils"
	"ws/utils/status"
)

func IsWarFinished(state *models.State) bool {
	westDead := state.WestAlive == 0
	eastDead := state.EastAlive == 0
	return westDead && eastDead
}

func IsSideDefeated(state *models.State, side string) bool {
	switch side {
	case game.West:
		return state.EastAlive == 0
	case game.East:
		return state.WestAlive == 0
	default:
		return IsWarFinished(state)
	}
}

func ArenaPhase(state *models.State) bool {
	return state.Phase == game.Arena
}

func IsInArena(unit *models.Unit) bool {
	return unit.Slot/game.ArenaCoef == 1
}

func IsAttacking(unit *models.Unit) bool {
	return unit.Action == game.Attack || unit.Action == game.Wait
}

func IsVisible(unit, target *models.Unit) bool {
	targetInSlot := utils.IsInSlot(game.Slots, unit.Slot, target.Coords)
	unitInSlot := utils.IsInSlot(game.Slots, unit.Slot, unit.Coords)
	/*LOS := game.LineOfSight
	if status.IsKing(unit) {
		LOS = game.KingLOS
	}
	if IsInArena(unit) {
		LOS = game.ArenaLOS
	}*/
	inLineOfSight := !status.IsKing(unit) || (utils.XYRange(unit.Coords, target.Coords)-target.Size) < game.KingLOS
	return targetInSlot && unitInSlot && inLineOfSight
}

func IsCloser(unit, enemy, target *models.Unit) bool {
	newRange := utils.XYRange(unit.Coords, enemy.Coords)
	shortestRange := utils.XYRange(unit.Coords, target.Coords)
	return newRange < shortestRange
}

func IsControlled(unit *models.Unit) bool {
	for _, cc := range game.CrowdControl {
		for effect, _ := range unit.Affected {
			if effect == cc {
				return true
			}
		}
	}
	return false
}

func HalfAspd(unit *models.Unit) float64 {
	return float64(int((unit.Aspd + 1) / 2))
}

func SetDelay(unit *models.Unit, delay float64) {
	unit.AspdTimer = delay
}

func SetKingSlot(unit *models.Unit) {
	unit.Slot = game.KingCoef + unit.Slot/game.SideCoef
	unit.Waypoint = 1
}

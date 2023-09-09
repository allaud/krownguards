package war

import (
	"ws/game"
	"ws/models"
	"ws/utils"
	"ws/utils/status"
)

func GetSideUnits(unit *models.Unit, state *models.State, friendly bool) []*models.Unit {
	if ArenaPhase(state) {
		return GetUnitsBySlots(state, ArenaSlots(friendly), true)
	}

	if status.IsInKingSlot(unit) {
		return GetUnitsBySlots(state, SideSlots(unit.Slot), friendly)
	}

	if !status.IsWave(unit) {
		switch unit.Waypoint {
		case 1:
			assSlot, exists := utils.AssocSlot(unit.Slot)
			if exists {
				return GetUnitsBySlots(state, []int{unit.Slot, assSlot}, friendly)
			}
			return GetUnitsBySlots(state, []int{unit.Slot}, friendly)
		default:
			return GetUnitsBySlots(state, SideSlots(unit.Slot), friendly)
		}
	}

	switch unit.Waypoint {
	case 1:
		return GetUnitsBySlots(state, []int{unit.Slot}, friendly)
	case 2, 3:
		assSlot, exists := utils.AssocSlot(unit.Slot)
		if exists {
			return GetUnitsBySlots(state, []int{unit.Slot, assSlot}, friendly)
		}
		return GetUnitsBySlots(state, []int{unit.Slot}, friendly)
	default:
		return GetUnitsBySlots(state, SideSlots(unit.Slot), friendly)
	}
}

func GetUnitsBySlots(state *models.State, slots []int, friendly bool) []*models.Unit {
	unitsBySlots := []*models.Unit{}
	sideUnits := map[int][]*models.Unit{}
	if friendly {
		if !ArenaPhase(state) {
			for _, slot := range slots {
				if utils.Side(slot) == game.West {
					unitsBySlots = append(unitsBySlots, state.Kings[0])
					break
				}
			}
			for _, slot := range slots {
				if utils.Side(slot) == game.East {
					unitsBySlots = append(unitsBySlots, state.Kings[1])
					break
				}
			}
		}
		sideUnits = state.PlayerUnits
	} else {
		sideUnits = state.WaveUnits
	}
	for _, slot := range slots {
		unitsTemp := sideUnits[slot]
		unitsBySlots = append(unitsBySlots, unitsTemp...)
	}
	return unitsBySlots
}

func SideSlots(slot int) []int {
	if utils.Side(slot) == game.West {
		return []int{1, 2, 3, 4}
	}
	return []int{5, 6, 7, 8}
}

func ArenaSlots(side bool) []int {
	if side {
		return SideSlots(1)
	}
	return SideSlots(5)
}

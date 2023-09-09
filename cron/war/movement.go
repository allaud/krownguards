package war

import (
	"math"
	//"math/rand"
	"ws/datatypes"
	"ws/game"
	"ws/models"
	"ws/utils"
	"ws/utils/status"
)

func GetWaypoints(unit *models.Unit) datatypes.SlotMap {
	if status.IsWave(unit) {
		waveWaypoints := *game.WaveWaypoints
		return waveWaypoints[unit.Slot]
	}
	if status.IsKing(unit) {
		kingWaypoints := *game.KingWaypoints
		return kingWaypoints[unit.Slot]
	}
	unitWaypoints := *game.UnitWaypoints
	return unitWaypoints[unit.Slot]
}

func InWaypoint(unit *models.Unit, waypoint int) bool {
	x, y := unit.Coords[0], unit.Coords[1]
	unitWaypoint := GetWaypoints(unit)[waypoint]
	lX := x >= unitWaypoint.XRange[0]
	rX := x <= unitWaypoint.XRange[1]
	bY := y >= unitWaypoint.YRange[0]
	tY := y <= unitWaypoint.YRange[1]
	return lX && rX && bY && tY
}

func CloseIn(value float64, bounds [2]float64) float64 {
	if value < bounds[0] {
		value = bounds[0]
	}
	if value > bounds[1] {
		value = bounds[1]
	}
	return value
}

func CollisionExists(newCoords [2]float64, unit *models.Unit, others []*models.Unit) ([2]float64, float64, bool) {
	for _, other := range others {
		if unit.Id == other.Id || status.IsMarkedDead(other) {
			continue
		}
		if utils.XYRange(newCoords, other.Coords) < (unit.Size + other.Size) {
			return other.Coords, other.Size, true
		}
	}
	return [2]float64{0, 0}, 0, false
}

func ContainCoords(ways [][2]float64, coords [2]float64) bool {
	for _, way := range ways {
		if way == coords {
			return true
		}
	}
	return false
}

func MoveToCoords(unit *models.Unit, others []*models.Unit, coords [2]float64) {
	step := unit.Movespeed
	precStep := 2 * unit.Size
	wayCoords := utils.GetVectorCoords(unit.Coords, coords, precStep, 0)
	ways := [][2]float64{wayCoords}
	index := 0
	for index < len(ways) {
		coll, csize, exists := CollisionExists(ways[index], unit, others)
		inWall, xCoef, yCoef := utils.InWall(unit.Size, ways[index])
		if !exists && !inWall {
			/*if utils.XYRange(unit.Coords, coords) < utils.XYRange(ways[index], coords) {
				unit.Action = game.Wait
				return
			}*/
			if index == 0 {
				ways[index] = utils.GetVectorCoords(unit.Coords, coords, step, 0)
			}
			unit.Coords = ways[index]
			unit.Action = game.Move
			return
		}
		x, y := unit.Coords[0], unit.Coords[1]
		if utils.XYRange(unit.TargetCoords, coords) > utils.XYRange(unit.Coords, unit.TargetCoords) || unit.TargetSign == 0 {
			invert := math.Copysign(1, math.Abs(y-coords[1])-math.Abs(x-coords[0]))
			sign := invert * math.Copysign(1, x-coords[0]) * math.Copysign(1, y-coords[1])
			unit.TargetCoords = coords
			unit.TargetSign = sign
		}
		if inWall {
			unit.TargetSign = unit.TargetSign * -1
			if coords[0] != x {
				xCoef = coords[0] - x
			}
			wayCoords = [2]float64{x + math.Copysign(step, xCoef), y}
			if !ContainCoords(ways, wayCoords) {
				ways = append(ways, wayCoords)
			}
			if coords[1] != y {
				yCoef = coords[1] - y
			}
			wayCoords = [2]float64{x, y + math.Copysign(step, yCoef)}
		} else {
			// get collision round angle
			//angle := 0.5 * math.Pi
			sumSize := unit.Size + csize + 0.001 // DO NOT LEAN ON THE COLLISION
			rad := utils.XYRange(unit.Coords, coll)
			angle := utils.AngleABySides(sumSize, step, rad)
			if index == 0 {
				angle = utils.AngleABySides(sumSize, precStep, rad)
			}
			// is units stacked in each other - NEVER CALLED
			if rad < sumSize {
				angle = math.Pi
			}
			// get target deflection angle
			/*ctRange := utils.XYRange(coll, coords)
			utRange := utils.XYRange(unit.Coords, coords)
			defAngle := utils.AngleABySides(utRange, ctRange, rad)
			// if direction is near straight randomize round direction
			if math.Abs(defAngle) > math.Pi*0.99 {
				signs := [2]float64{-1, 1}
				unit.TargetSign = signs[rand.Intn(len(signs))]
			}*/
			wayCoords = utils.GetVectorCoords(unit.Coords, coll, step, unit.TargetSign*angle)
		}

		if !ContainCoords(ways, wayCoords) {
			ways = append(ways, wayCoords)
		}
		index += 1
	}
	unit.Action = game.Wait
}

func TeleportToCoords(unit *models.Unit, others []*models.Unit, xBounds, yBounds [2]float64) {
	lX := xBounds[0] + unit.Size
	rX := xBounds[1] - unit.Size
	bY := yBounds[0] + unit.Size
	tY := yBounds[1] - unit.Size
	ways := [][2]float64{[2]float64{utils.Uniform(lX, rX), utils.Uniform(bY, tY)}}
	index := 0
	for index < len(ways) {
		_, _, exists := CollisionExists(ways[index], unit, others)
		inWall, _, _ := utils.InWall(unit.Size, ways[index])
		if !exists && !inWall {
			unit.Coords = ways[index]
			unit.Action = game.Teleport
			return
		}
		wayCoords := [2]float64{utils.Uniform(lX, rX), utils.Uniform(bY, tY)}

		if !ContainCoords(ways, wayCoords) && len(ways) < len(others) {
			ways = append(ways, wayCoords)
		}

		index = index + 1
	}
	unit.Action = game.Wait
}

func MoveToWaypoint(unit *models.Unit, others []*models.Unit, state *models.State) {
	maxWaypoint := 1
	for wpoint, _ := range GetWaypoints(unit) {
		if wpoint > maxWaypoint {
			maxWaypoint = wpoint
		}
	}

	if InWaypoint(unit, maxWaypoint) && (status.IsInKingSlot(unit) || IsInArena(unit)) {
		unit.Action = game.Idle
		return
	}

	waypoint := GetWaypoints(unit)[unit.Waypoint]
	xBounds := waypoint.XRange
	yBounds := waypoint.YRange
	x := CloseIn(unit.Coords[0], xBounds)
	y := CloseIn(unit.Coords[1], yBounds)

	if !IsInArena(unit) && !status.IsWave(unit) && (unit.Waypoint == maxWaypoint) && (!status.IsInKingSlot(unit)) {
		// others => GetUnitsBySlots(state, SideSlots(unit.Slot), true)
		maxCount := (yBounds[1] - yBounds[0]) * (xBounds[1] - xBounds[0])
		teleported := GetUnitsBySlots(state, SideSlots(unit.Slot), true)
		shrinkCoef := math.Sqrt(float64(len(teleported)) / maxCount)
		tpXBounds := xBounds
		tpYBounds := yBounds
		if shrinkCoef < 0.1 {
			shrinkCoef = 0.1
		}
		if shrinkCoef < 1 {
			tpXBounds[0] = xBounds[0] + (0.5-0.5*shrinkCoef)*(xBounds[1]-xBounds[0])
			tpXBounds[1] = xBounds[0] + (0.5+0.5*shrinkCoef)*(xBounds[1]-xBounds[0])
			//tpYBounds[1] = yBounds[0] + shrinkCoef*(yBounds[1]-yBounds[0])
			tpYBounds[0] = yBounds[1] - shrinkCoef*(yBounds[1]-yBounds[0])
		}
		TeleportToCoords(unit, teleported, tpXBounds, tpYBounds)
	} else {
		MoveToCoords(unit, others, [2]float64{x, y})
	}
	if unit.Waypoint != maxWaypoint {
		if InWaypoint(unit, unit.Waypoint) {
			unit.Waypoint = unit.Waypoint + 1
		} else if !utils.IsInSlot(game.Slots, unit.Slot, unit.Coords) {
			// check if unit advanced to next WP without checking previous
			nextWaypoint := GetWaypoints(unit)[unit.Waypoint+1]
			nextXBounds := nextWaypoint.XRange
			nextYBounds := nextWaypoint.YRange
			nextX := CloseIn(unit.Coords[0], nextXBounds)
			nextY := CloseIn(unit.Coords[1], nextYBounds)
			WPRange := utils.XYRange(unit.Coords, [2]float64{x, y})
			nextWpRange := utils.XYRange(unit.Coords, [2]float64{nextX, nextY})
			//wpDist := utils.XYRange([2]float64{x, y}, [2]float64{nextX, nextY})
			//angle := utils.AngleABySides(wpDist, WPRange, nextWpRange)
			if WPRange > nextWpRange {
				unit.Waypoint = unit.Waypoint + 1
			}
		}
	}
	if IsInArena(unit) {
		return
	}
	if (!status.IsInKingSlot(unit)) && utils.IsInSlot(game.Slots, (game.KingCoef+unit.Slot/game.SideCoef), unit.Coords) {
		if status.IsWave(unit) {
			unit.Affected["Final charge"] = &models.Effect{
				Duration: int(120 / game.WAR_STEP_DELAY),
			}
			GetScore(state, unit, unit.Slot, false)
			DampBounty(state, unit)
		}
		SetKingSlot(unit)
	}
}

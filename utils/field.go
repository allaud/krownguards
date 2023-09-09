package utils

import (
	"math"
	"sort"
	"ws/datatypes"
	"ws/game"
	"ws/models"
)

func XYRange(coords1, coords2 [2]float64) float64 {
	return math.Sqrt(math.Pow(coords2[0]-coords1[0], 2) + math.Pow(coords2[1]-coords1[1], 2))
}

func AngleABySides(a, b, c float64) float64 {
	return math.Acos((math.Pow(b, 2) + math.Pow(c, 2) - math.Pow(a, 2)) / (2 * b * c))
}

func GetVectorCoords(coords1, coords2 [2]float64, length, delta float64) [2]float64 {
	x1, y1 := coords1[0], coords1[1]
	x2, y2 := coords2[0], coords2[1]
	rad := XYRange(coords1, coords2)
	if (rad < length || rad == 0) && delta == 0 {
		return coords2
	}
	angle := math.Copysign(1, y2-y1) * math.Acos((x2-x1)/rad)
	x := x1 + length*math.Cos(angle+delta)
	y := y1 + length*math.Sin(angle+delta)
	return [2]float64{x, y}
}

func IsInSlot(slots *datatypes.SlotMap, index int, coords [2]float64) bool {
	slotsValue := *slots
	InXMin := coords[0] >= slotsValue[index].XRange[0]
	InXMax := coords[0] <= slotsValue[index].XRange[1]
	InYMin := coords[1] >= slotsValue[index].YRange[0]
	InYMax := coords[1] <= slotsValue[index].YRange[1]
	return InXMin && InXMax && InYMin && InYMax
}

func IsInBuildSlot(slots *datatypes.SlotMap, index int, coords [2]float64) bool {
	CenterIn := IsInSlot(slots, index, coords)
	TopIn := IsInSlot(slots, index, [2]float64{coords[0] + 0.5, coords[1] + 0.5})
	BotIn := IsInSlot(slots, index, [2]float64{coords[0] - 0.5, coords[1] - 0.5})
	return CenterIn && TopIn && BotIn
}

func GetGridCoords(coords [2]float64) [2]float64 {
	x := float64(int(coords[0]/0.5))/2 + float64(int(math.Mod(coords[0], 0.5)/0.25))*0.5
	y := float64(int(coords[1]/0.5))/2 + float64(int(math.Mod(coords[1], 0.5)/0.25))*0.5
	return [2]float64{x, y}
}

func GetRandomCoords(size float64, slots *datatypes.SlotMap, index int) [2]float64 {
	slotsValue := *slots
	xMin := slotsValue[index].XRange[0] + size
	xMax := slotsValue[index].XRange[1] - size
	x := Uniform(xMin, xMax)
	yMin := slotsValue[index].YRange[0] + size
	yMax := slotsValue[index].YRange[1] - size
	y := Uniform(yMin, yMax)
	return [2]float64{x, y}
}

func RandomSpawn(unit *models.Unit, wave []*models.Unit, slots *datatypes.SlotMap) {
	randCoords := GetRandomCoords(unit.Size, slots, unit.Slot)
	attempts := 1
	index := 0
	for index < attempts {
		for _, waveUnit := range wave {
			if XYRange(randCoords, waveUnit.Coords) < (unit.Size + waveUnit.Size) {
				randCoords = GetRandomCoords(unit.Size, slots, unit.Slot)
				//if attempts < 50000 { //MAGIC NUMBER
				attempts = attempts + 1
				//}
				break
			}
		}
		index = index + 1
	}
	unit.Coords = randCoords
}

func AssocSlot(index int) (int, bool) {
	slot, ok := game.Mapping[index]
	if !ok {
		return 0, false
	}
	x, y := slot[0], slot[1]
	for sIndex, coords := range game.Mapping {
		if x == coords[0] && y == 3-coords[1] {
			return sIndex, true
		}
	}
	return 0, false
}

func Side(slot int) string {
	if (slot < game.SideCoef) || (slot == game.KingCoef) || (slot == game.ArenaCoef) {
		return game.West
	}
	return game.East
}

func InWall(size float64, coords [2]float64) (bool, float64, float64) {
	lX := coords[0] - size
	rX := coords[0] + size
	bY := coords[1] - size
	tY := coords[1] + size
	for _, wall := range game.Walls {
		inX := (rX > wall.XRange[0]) && (lX < wall.XRange[1])
		inY := (tY > wall.YRange[1]) || (bY < wall.YRange[0])
		if inX && inY {
			xCoef := math.Copysign(0.5, coords[0]-wall.XRange[1]) + math.Copysign(0.5, coords[0]-wall.XRange[0])
			yCoef := math.Copysign(0.5, wall.YRange[1]-tY) + math.Copysign(0.5, wall.YRange[0]-bY)
			return true, xCoef, yCoef
		}
	}

	for _, wall := range game.DiagonalWalls {
		inDX := (rX > wall.XRange[0]) && (lX < wall.XRange[1])
		yTop := wall.YRange[0]*coords[0] + wall.YRange[1]
		yBot := wall.YRange[0]*coords[0] + wall.YRange[1] - float64(game.SlotOverlap)*game.SqY
		inDY := (tY > yTop) || (bY < yBot)
		if inDX && inDY {
			xCoef := math.Copysign(0.5, coords[0]-wall.XRange[1]) + math.Copysign(0.5, coords[0]-wall.XRange[0])
			yCoef := math.Copysign(0.5, yTop-tY) + math.Copysign(0.5, yBot-bY)
			return true, xCoef, yCoef
		}
	}

	return false, 0, 0
}

func ArenaSpawn(state *models.State, unit *models.Unit, buildSlots, spawnSlots *datatypes.SlotMap) [2]float64 {
	// get player arena spawn slot
	spawnSlotsValue := *spawnSlots
	spawnSlot := spawnSlotsValue[game.ArenaCoef+unit.Slot/game.SideCoef]
	spawnCenter := spawnSlot.XRange[0] + 0.5*(spawnSlot.XRange[1]-spawnSlot.XRange[0])
	buildSlotsValue := *buildSlots
	buildSlot := buildSlotsValue[unit.Slot]
	halfSize := 0.5 * (buildSlot.XRange[1] - buildSlot.XRange[0])
	side := []int{}
	for slot, _ := range state.Players {
		if Side(unit.Slot) == Side(slot) {
			side = append(side, slot)
		}
	}
	sort.Ints(side)
	max := len(side)
	plIndex := 0
	for index, slot := range side {
		if slot == unit.Slot {
			plIndex = index + 1
		}
	}
	slotCenter := spawnCenter + halfSize*float64(2*plIndex-max-1)
	slot := datatypes.Size{
		XRange: [2]float64{
			slotCenter - halfSize,
			slotCenter + halfSize,
		},
		YRange: spawnSlot.YRange,
	}

	// get arena spawn coords
	x := slot.XRange[0] + unit.Coords[0] - buildSlot.XRange[0]
	y := slot.YRange[0] + unit.Coords[1] - buildSlot.YRange[0]
	return [2]float64{x, y}
}

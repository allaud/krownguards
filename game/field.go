package game

import (
	//"fmt"
	"math"
	"ws/datatypes"
)

var FieldSize = [2]float64{198, 56}

// set parameters
var BorderSize = [2]float64{1, 1}
var SlotSize = [2]float64{10, 23}
var BuildSlotSize = [2]float64{9, 15}
var SpawnSlotSize = [2]float64{9, 4}
var SlotOverlap float64 = 10
var GreenSize = [2]float64{6.5, 18}
var DecorSize = [2]float64{12.5, 24}
var KingSize = [2]float64{9, 21}
var GapSize = [2]float64{58, 56}
var ArenaSize = [2]float64{38, 52}
var ArenaSpawnSize = [2]float64{BuildSlotSize[1], 4 * BuildSlotSize[0]}
var IncomeBoxSize = [2]float64{9, 8}

var ParamFieldSize = [2]float64{
	6*BorderSize[0] + 4*SlotSize[0] + 4*GreenSize[0] + 4*DecorSize[0] + 2*KingSize[0] + GapSize[0],
	2*SlotSize[1] + SlotOverlap,
}

var SqX, SqY = FieldSize[0] / ParamFieldSize[0], FieldSize[1] / ParamFieldSize[1]

var MapCenter = [2]float64{ParamFieldSize[0] * SqX / 2, ParamFieldSize[1] * SqY / 2}

var Mapping = map[int][2]int{
	1: {1, 2},
	2: {2, 2},
	3: {1, 1},
	4: {2, 1},
	5: {3, 2},
	6: {4, 2},
	7: {3, 1},
	8: {4, 1},
}

var LineOfSight float64 = 9
var KingLOS float64 = 7.5
var ArenaLOS float64 = 40

var Walls = []datatypes.Size{
	// West left wall
	{
		XRange: [2]float64{
			0 * SqX,
			(1.5 * BorderSize[0]) * SqX,
		},
		YRange: [2]float64{
			(ParamFieldSize[1] / 2) * SqY,
			(ParamFieldSize[1] / 2) * SqY,
		},
	},
	// West King left wall
	{
		XRange: [2]float64{
			(1.5*BorderSize[0] + SlotSize[0] + GreenSize[0]) * SqX,
			(1.5*BorderSize[0] + SlotSize[0] + GreenSize[0] + DecorSize[0]) * SqX,
		},
		YRange: [2]float64{
			(ParamFieldSize[1] - 2*BorderSize[1] - DecorSize[1] - SlotOverlap) * SqY,
			(ParamFieldSize[1] - 2*BorderSize[1] - DecorSize[1]) * SqY,
		},
	},
	// West King slot
	{
		XRange: [2]float64{
			(1.5*BorderSize[0] + SlotSize[0] + GreenSize[0] + DecorSize[0]) * SqX,
			(1.5*BorderSize[0] + SlotSize[0] + GreenSize[0] + DecorSize[0] + KingSize[0]) * SqX,
		},
		YRange: [2]float64{
			(2*BorderSize[1] + GreenSize[1]) * SqY,
			(ParamFieldSize[1] - 1.5*BorderSize[1]) * SqY,
		},
	},
	// West King right wall
	{
		XRange: [2]float64{
			(1.5*BorderSize[0] + SlotSize[0] + GreenSize[0] + DecorSize[0] + KingSize[0]) * SqX,
			(1.5*BorderSize[0] + SlotSize[0] + GreenSize[0] + 2*DecorSize[0] + KingSize[0]) * SqX,
		},
		YRange: [2]float64{
			(ParamFieldSize[1] - 2*BorderSize[1] - DecorSize[1] - SlotOverlap) * SqY,
			(ParamFieldSize[1] - 2*BorderSize[1] - DecorSize[1]) * SqY,
		},
	},
	// West right wall
	{
		XRange: [2]float64{
			(1.5*BorderSize[0] + 2*SlotSize[0] + 2*GreenSize[0] + 2*DecorSize[0] + KingSize[0]) * SqX,
			(3*BorderSize[0] + 2*SlotSize[0] + 2*GreenSize[0] + 2*DecorSize[0] + KingSize[0]) * SqX,
		},
		YRange: [2]float64{
			(ParamFieldSize[1] / 2) * SqY,
			(ParamFieldSize[1] / 2) * SqY,
		},
	},
	// Arena left wall
	{
		XRange: [2]float64{
			(ParamFieldSize[0]/2 - ArenaSize[0]/2 - 2*BorderSize[0]) * SqX,
			(ParamFieldSize[0]/2 - ArenaSize[0]/2) * SqX,
		},
		YRange: [2]float64{
			(ParamFieldSize[1] / 2) * SqY,
			(ParamFieldSize[1] / 2) * SqY,
		},
	},
	// Arena
	{
		XRange: [2]float64{
			(ParamFieldSize[0]/2 - ArenaSize[0]/2) * SqX,
			(ParamFieldSize[0]/2 + ArenaSize[0]/2) * SqX,
		},
		YRange: [2]float64{
			(2 * BorderSize[1]) * SqY,
			(ParamFieldSize[1] - 2*BorderSize[1]) * SqY,
		},
	},
	// Arena right wall
	{
		XRange: [2]float64{
			(ParamFieldSize[0]/2 + ArenaSize[0]/2) * SqX,
			(ParamFieldSize[0]/2 + ArenaSize[0]/2 + 2*BorderSize[0]) * SqX,
		},
		YRange: [2]float64{
			(ParamFieldSize[1] / 2) * SqY,
			(ParamFieldSize[1] / 2) * SqY,
		},
	},
	// East left wall
	{
		XRange: [2]float64{
			(3*BorderSize[0] + 2*SlotSize[0] + 2*GreenSize[0] + 2*DecorSize[0] + KingSize[0] + GapSize[0]) * SqX,
			(4.5*BorderSize[0] + 2*SlotSize[0] + 2*GreenSize[0] + 2*DecorSize[0] + KingSize[0] + GapSize[0]) * SqX,
		},
		YRange: [2]float64{
			(ParamFieldSize[1] / 2) * SqY,
			(ParamFieldSize[1] / 2) * SqY,
		},
	},
	// East King left wall
	{
		XRange: [2]float64{
			(4.5*BorderSize[0] + 3*SlotSize[0] + 3*GreenSize[0] + 2*DecorSize[0] + KingSize[0] + GapSize[0]) * SqX,
			(4.5*BorderSize[0] + 3*SlotSize[0] + 3*GreenSize[0] + 3*DecorSize[0] + KingSize[0] + GapSize[0]) * SqX,
		},
		YRange: [2]float64{
			(ParamFieldSize[1] - 2*BorderSize[1] - DecorSize[1] - SlotOverlap) * SqY,
			(ParamFieldSize[1] - 2*BorderSize[1] - DecorSize[1]) * SqY,
		},
	},
	// East King slot
	{
		XRange: [2]float64{
			(4.5*BorderSize[0] + 3*SlotSize[0] + 3*GreenSize[0] + 3*DecorSize[0] + KingSize[0] + GapSize[0]) * SqX,
			(4.5*BorderSize[0] + 3*SlotSize[0] + 3*GreenSize[0] + 3*DecorSize[0] + 2*KingSize[0] + GapSize[0]) * SqX,
		},
		YRange: [2]float64{
			(2*BorderSize[1] + GreenSize[1]) * SqY,
			(ParamFieldSize[1] - 1.5*BorderSize[1]) * SqY,
		},
	},
	// East King right wall
	{
		XRange: [2]float64{
			(4.5*BorderSize[0] + 3*SlotSize[0] + 3*GreenSize[0] + 3*DecorSize[0] + 2*KingSize[0] + GapSize[0]) * SqX,
			(4.5*BorderSize[0] + 3*SlotSize[0] + 3*GreenSize[0] + 4*DecorSize[0] + 2*KingSize[0] + GapSize[0]) * SqX,
		},
		YRange: [2]float64{
			(ParamFieldSize[1] - 2*BorderSize[1] - DecorSize[1] - SlotOverlap) * SqY,
			(ParamFieldSize[1] - 2*BorderSize[1] - DecorSize[1]) * SqY,
		},
	},
	// West right wall
	{
		XRange: [2]float64{
			(ParamFieldSize[0] - 1.5*BorderSize[0]) * SqX,
			(ParamFieldSize[0]) * SqX,
		},
		YRange: [2]float64{
			(ParamFieldSize[1] / 2) * SqY,
			(ParamFieldSize[1] / 2) * SqY,
		},
	},
}

var DiagonalWalls = []datatypes.Size{
	{
		XRange: [2]float64{
			(1.5*BorderSize[0] + SlotSize[0]) * SqX,
			(1.5*BorderSize[0] + SlotSize[0] + GreenSize[0]) * SqX,
		},
		YRange: [2]float64{
			-0.5,
			(SlotSize[1]+SlotOverlap)*SqY + 0.5*(2*BorderSize[0]+SlotSize[0])*SqX,
		},
	},
	{
		XRange: [2]float64{
			(1.5*BorderSize[0] + SlotSize[0] + GreenSize[0] + 2*DecorSize[0] + KingSize[0]) * SqX,
			(1.5*BorderSize[0] + SlotSize[0] + 2*GreenSize[0] + 2*DecorSize[0] + KingSize[0]) * SqX},
		YRange: [2]float64{
			0.5,
			(SlotSize[1]+SlotOverlap)*SqY - 0.5*(1*BorderSize[0]+SlotSize[0]+2*GreenSize[0]+2*DecorSize[0]+KingSize[0])*SqX,
		},
	},
	{
		XRange: [2]float64{
			(4.5*BorderSize[0] + 3*SlotSize[0] + 2*GreenSize[0] + 2*DecorSize[0] + KingSize[0] + GapSize[0]) * SqX,
			(4.5*BorderSize[0] + 3*SlotSize[0] + 3*GreenSize[0] + 2*DecorSize[0] + KingSize[0] + GapSize[0]) * SqX,
		},
		YRange: [2]float64{
			-0.5,
			(SlotSize[1]+SlotOverlap)*SqY + 0.5*(5*BorderSize[0]+3*SlotSize[0]+2*GreenSize[0]+2*DecorSize[0]+KingSize[0]+GapSize[0])*SqX,
		},
	},
	{
		XRange: [2]float64{
			(4.5*BorderSize[0] + 3*SlotSize[0] + 3*GreenSize[0] + 4*DecorSize[0] + 2*KingSize[0] + GapSize[0]) * SqX,
			(4.5*BorderSize[0] + 3*SlotSize[0] + 4*GreenSize[0] + 4*DecorSize[0] + 2*KingSize[0] + GapSize[0]) * SqX,
		},
		YRange: [2]float64{
			0.5,
			(SlotSize[1]+SlotOverlap)*SqY - 0.5*(4*BorderSize[0]+3*SlotSize[0]+4*GreenSize[0]+4*DecorSize[0]+2*KingSize[0]+GapSize[0])*SqX,
		},
	},
}

var slotsValue = datatypes.SlotMap{}
var buildSlotsValue = datatypes.SlotMap{}
var spawnSlotsValue = datatypes.SlotMap{}
var unitWaypointsValue = datatypes.WaypointMap{}
var waveWaypointsValue = datatypes.WaypointMap{}
var kingWaypointsValue = datatypes.WaypointMap{}

var Slots = &slotsValue
var BuildSlots = &buildSlotsValue
var SpawnSlots = &spawnSlotsValue
var UnitWaypoints = &unitWaypointsValue
var WaveWaypoints = &waveWaypointsValue
var KingWaypoints = &kingWaypointsValue

var KingCoords = [2][2]float64{
	{
		(1.5*BorderSize[0] + SlotSize[0] + GreenSize[0] + DecorSize[0] + KingSize[0]/2) * SqX,
		(ParamFieldSize[1] - 4*BorderSize[1]) * SqY,
	},
	{
		(ParamFieldSize[0] - 1.5*BorderSize[0] - SlotSize[0] - GreenSize[0] - DecorSize[0] - KingSize[0]/2) * SqX,
		(ParamFieldSize[1] - 4*BorderSize[1]) * SqY,
	},
}

var KingDir = [2][2]float64{
	{
		KingCoords[0][0],
		MapCenter[1],
	},
	{
		KingCoords[1][0],
		MapCenter[1],
	},
}

var IncomeBox = &datatypes.SlotMap{
	0: datatypes.Size{
		XRange: [2]float64{
			(1.5*BorderSize[0] + SlotSize[0] + GreenSize[0] + DecorSize[0]) * SqX,
			(1.5*BorderSize[0] + SlotSize[0] + GreenSize[0] + DecorSize[0] + IncomeBoxSize[0]) * SqX,
		},
		YRange: [2]float64{
			(GreenSize[1] - IncomeBoxSize[1]) * SqY,
			(GreenSize[1]) * SqY,
		},
	},
	1: datatypes.Size{
		XRange: [2]float64{
			(ParamFieldSize[0] - 1.5*BorderSize[0] - SlotSize[0] - GreenSize[0] - DecorSize[0] - IncomeBoxSize[0]) * SqX,
			(ParamFieldSize[0] - 1.5*BorderSize[0] - SlotSize[0] - GreenSize[0] - DecorSize[0]) * SqX,
		},
		YRange: [2]float64{
			(GreenSize[1] - IncomeBoxSize[1]) * SqY,
			(GreenSize[1]) * SqY,
		},
	},
}

func init() {
	// init slots
	for slot, maps := range Mapping {
		x, y := maps[0], maps[1]
		// init Slot
		slotsValue[slot] = datatypes.Size{
			XRange: [2]float64{
				(1.5*BorderSize[0] + SlotSize[0]*float64(x-1) + (2*GreenSize[0]+2*DecorSize[0]+KingSize[0])*float64(x/2) + (3*BorderSize[0]+GapSize[0])*float64(x/3)) * SqX,
				(1.5*BorderSize[0] + SlotSize[0]*float64(x) + (2*GreenSize[0]+2*DecorSize[0]+KingSize[0])*float64(x/2) + (3*BorderSize[0]+GapSize[0])*float64(x/3)) * SqX,
			},
			YRange: [2]float64{
				(float64(y-1) * SlotSize[1]) * SqY,
				(SlotSize[1] + SlotOverlap + float64(y-1)*SlotSize[1]) * SqY,
			},
		}
		// init Build Slot
		buildSlotsValue[slot] = datatypes.Size{
			XRange: [2]float64{
				(2*BorderSize[0] + BuildSlotSize[0]*float64(x-1) + (BorderSize[0]+2*GreenSize[0]+2*DecorSize[0]+KingSize[0])*float64(x/2) + (4*BorderSize[0]+GapSize[0])*float64(x/3)) * SqX,
				(2*BorderSize[0] + BuildSlotSize[0]*float64(x) + (BorderSize[0]+2*GreenSize[0]+2*DecorSize[0]+KingSize[0])*float64(x/2) + (4*BorderSize[0]+GapSize[0])*float64(x/3)) * SqX,
			},
			YRange: [2]float64{
				(SlotSize[1] - BuildSlotSize[1] + float64(y-1)*(BuildSlotSize[1]+SlotOverlap)) * SqY,
				(SlotSize[1] + float64(y-1)*(BuildSlotSize[1]+SlotOverlap)) * SqY,
			},
		}
		// init Spawn Slot
		spawnSlotsValue[slot] = datatypes.Size{
			XRange: [2]float64{
				(2*BorderSize[0] + SpawnSlotSize[0]*float64(x-1) + (BorderSize[0]+2*GreenSize[0]+2*DecorSize[0]+KingSize[0])*float64(x/2) + (4*BorderSize[0]+GapSize[0])*float64(x/3)) * SqX,
				(2*BorderSize[0] + SpawnSlotSize[0]*float64(x) + (BorderSize[0]+2*GreenSize[0]+2*DecorSize[0]+KingSize[0])*float64(x/2) + (4*BorderSize[0]+GapSize[0])*float64(x/3)) * SqX,
			},
			YRange: [2]float64{
				(2*BorderSize[1] + float64(y-1)*(SpawnSlotSize[1]+2*BuildSlotSize[1]+4*BorderSize[1]+SlotOverlap)) * SqY,
				(2*BorderSize[1] + SpawnSlotSize[1] + float64(y-1)*(SpawnSlotSize[1]+2*BuildSlotSize[1]+4*BorderSize[1]+SlotOverlap)) * SqY,
			},
		}
	}
	//WEST KING SLOT
	slotsValue[10] = datatypes.Size{
		XRange: [2]float64{
			(1.5*BorderSize[0] + SlotSize[0] + GreenSize[0] + DecorSize[0]) * SqX,
			(1.5*BorderSize[0] + SlotSize[0] + GreenSize[0] + DecorSize[0] + KingSize[0]) * SqX,
		},
		YRange: [2]float64{
			(ParamFieldSize[1] - 1.5*BorderSize[1] - KingSize[1]) * SqY,
			(ParamFieldSize[1] - 1.5*BorderSize[1]) * SqY,
		},
	}

	//EAST KING SLOT
	slotsValue[11] = datatypes.Size{
		XRange: [2]float64{
			(ParamFieldSize[0] - 1.5*BorderSize[0] - SlotSize[0] - GreenSize[0] - DecorSize[0] - KingSize[0]) * SqX,
			(ParamFieldSize[0] - 1.5*BorderSize[0] - SlotSize[0] - GreenSize[0] - DecorSize[0]) * SqX,
		},
		YRange: [2]float64{
			(ParamFieldSize[1] - 1.5*BorderSize[1] - KingSize[1]) * SqY,
			(ParamFieldSize[1] - 1.5*BorderSize[1]) * SqY,
		},
	}

	//ARENA SLOT
	slotsValue[20] = datatypes.Size{
		XRange: [2]float64{
			((ParamFieldSize[0] - ArenaSize[0]) / 2) * SqX,
			((ParamFieldSize[0] + ArenaSize[0]) / 2) * SqX,
		},
		YRange: [2]float64{
			((ParamFieldSize[1] - ArenaSize[1]) / 2) * SqY,
			((ParamFieldSize[1] + ArenaSize[1]) / 2) * SqY,
		},
	}

	slotsValue[21] = slotsValue[20]

	//spawnSlotsValue[10] = slotsValue[10]
	//spawnSlotsValue[11] = slotsValue[11]

	//WEST ARENA SPAWN SLOT
	spawnSlotsValue[20] = datatypes.Size{
		XRange: [2]float64{
			((ParamFieldSize[0] - ArenaSize[0]/2 - ArenaSpawnSize[0]) / 2) * SqX,
			((ParamFieldSize[0] - ArenaSize[0]/2 + ArenaSpawnSize[0]) / 2) * SqX,
		},
		YRange: [2]float64{
			((ParamFieldSize[1] - ArenaSpawnSize[1]) / 2) * SqY,
			((ParamFieldSize[1] + ArenaSpawnSize[1]) / 2) * SqY,
		},
	}
	//EAST ARENA SPAWN SLOT
	spawnSlotsValue[21] = datatypes.Size{
		XRange: [2]float64{
			((ParamFieldSize[0] + ArenaSize[0]/2 - ArenaSpawnSize[0]) / 2) * SqX,
			((ParamFieldSize[0] + ArenaSize[0]/2 + ArenaSpawnSize[0]) / 2) * SqX,
		},
		YRange: [2]float64{
			((ParamFieldSize[1] - ArenaSpawnSize[1]) / 2) * SqY,
			((ParamFieldSize[1] + ArenaSpawnSize[1]) / 2) * SqY,
		},
	}
	// init Unit Waypoints
	for slot, maps := range Mapping {
		x, y := spawnSlotsValue[slot].XRange, spawnSlotsValue[slot].YRange
		kX, kY := slotsValue[KingCoef+slot/SideCoef].XRange, slotsValue[KingCoef+slot/SideCoef].YRange
		unitWaypointsValue[slot] = datatypes.SlotMap{
			1: datatypes.Size{
				XRange: x,
				YRange: [2]float64{
					y[0] + 0.5*(y[1]-y[0])*float64(maps[1]-1),
					y[1] + 0.5*(y[1]-y[0])*float64(maps[1]-2),
				},
			},
			2: datatypes.Size{
				XRange: kX,
				YRange: [2]float64{
					kY[0] + (1.0/21.0)*(kY[1]-kY[0]),
					kY[0] + (11.0/21.0)*(kY[1]-kY[0]),
				},
			},
		}
	}
	unitWaypointsValue[10] = datatypes.SlotMap{
		1: datatypes.Size{
			XRange: slotsValue[10].XRange,
			YRange: [2]float64{
				slotsValue[10].YRange[0],
				slotsValue[10].YRange[0] + (11.0/KingSize[1])*(slotsValue[10].YRange[1]-slotsValue[10].YRange[0]),
			},
		},
	}
	unitWaypointsValue[11] = datatypes.SlotMap{
		1: datatypes.Size{
			XRange: slotsValue[11].XRange,
			YRange: [2]float64{
				slotsValue[11].YRange[0],
				slotsValue[11].YRange[0] + (11.0/KingSize[1])*(slotsValue[11].YRange[1]-slotsValue[11].YRange[0]),
			},
		},
	}
	unitWaypointsValue[20] = datatypes.SlotMap{
		1: datatypes.Size{
			XRange: [2]float64{
				slotsValue[20].XRange[0] + 0.45*(slotsValue[20].XRange[1]-slotsValue[20].XRange[0]),
				slotsValue[20].XRange[0] + 0.55*(slotsValue[20].XRange[1]-slotsValue[20].XRange[0]),
			},
			YRange: slotsValue[20].YRange,
		},
	}
	unitWaypointsValue[21] = unitWaypointsValue[20]
	// init Wave Waypoints
	for slot, maps := range Mapping {
		x := slotsValue[slot].XRange
		kX, kY := slotsValue[KingCoef+slot/SideCoef].XRange, slotsValue[KingCoef+slot/SideCoef].YRange
		waveWaypointsValue[slot] = datatypes.SlotMap{
			1: datatypes.Size{
				XRange: [2]float64{
					x[0] + ((math.Pow(-1, float64(1-maps[0]%2))*0.1 + 0.1) * (x[1] - x[0])),
					x[0] + ((math.Pow(-1, float64(1-maps[0]%2))*0.1 + 0.9) * (x[1] - x[0])),
				},
				YRange: [2]float64{
					(ParamFieldSize[1]/2 - SlotOverlap/4 - SlotOverlap/2*float64(maps[1]/2)) * SqY,
					(ParamFieldSize[1]/2 + SlotOverlap/4 + SlotOverlap/2*float64(maps[1]%2)) * SqY,
				},
			},
			2: datatypes.Size{
				XRange: [2]float64{
					kX[0] - DecorSize[0]*SqX,
					kX[1] + DecorSize[0]*SqX,
				},
				YRange: [2]float64{
					(2*BorderSize[1] + GreenSize[1]) * SqY,
					(2*BorderSize[1] + GreenSize[1] + 0.6*SlotOverlap) * SqY,
				},
			},
			3: datatypes.Size{
				XRange: [2]float64{
					kX[0] + 0.45*(kX[1]-kX[0]),
					kX[1] - 0.45*(kX[1]-kX[0]),
				},
				YRange: [2]float64{
					(2*BorderSize[1] + GreenSize[1]) * SqY,
					(2*BorderSize[1] + GreenSize[1] + 2*SlotOverlap) * SqY,
				},
			},
			4: datatypes.Size{
				XRange: kX,
				YRange: [2]float64{
					0.5 * (kY[0] + kY[1]),
					kY[1],
				},
			},
		}
	}
	waveWaypointsValue[10] = datatypes.SlotMap{
		1: datatypes.Size{
			XRange: slotsValue[10].XRange,
			YRange: [2]float64{
				KingCoords[0][1] - 0.1*(slotsValue[10].YRange[1]-slotsValue[10].YRange[0]),
				KingCoords[0][1] + 0.1*(slotsValue[10].YRange[1]-slotsValue[10].YRange[0]),
			},
		},
	}
	waveWaypointsValue[11] = datatypes.SlotMap{
		1: datatypes.Size{
			XRange: slotsValue[11].XRange,
			YRange: [2]float64{
				KingCoords[1][1] - 0.1*(slotsValue[11].YRange[1]-slotsValue[11].YRange[0]),
				KingCoords[1][1] + 0.1*(slotsValue[11].YRange[1]-slotsValue[11].YRange[0]),
			},
		},
	}
	waveWaypointsValue[20] = unitWaypointsValue[20]
	waveWaypointsValue[21] = unitWaypointsValue[20]
	// King Waypoints
	kingWaypointsValue[10] = datatypes.SlotMap{
		1: datatypes.Size{
			XRange: [2]float64{
				KingCoords[0][0] - 0.1*(slotsValue[10].XRange[1]-slotsValue[10].XRange[0]),
				KingCoords[0][0] + 0.1*(slotsValue[10].XRange[1]-slotsValue[10].XRange[0]),
			},
			YRange: [2]float64{
				KingCoords[0][1] - 0.1*(slotsValue[10].YRange[1]-slotsValue[10].YRange[0]),
				KingCoords[0][1] + 0.1*(slotsValue[10].YRange[1]-slotsValue[10].YRange[0]),
			},
		},
	}
	kingWaypointsValue[11] = datatypes.SlotMap{
		1: datatypes.Size{
			XRange: [2]float64{
				KingCoords[1][0] - 0.1*(slotsValue[11].XRange[1]-slotsValue[11].XRange[0]),
				KingCoords[1][0] + 0.1*(slotsValue[11].XRange[1]-slotsValue[11].XRange[0]),
			},
			YRange: [2]float64{
				KingCoords[1][1] - 0.1*(slotsValue[11].YRange[1]-slotsValue[11].YRange[0]),
				KingCoords[1][1] + 0.1*(slotsValue[11].YRange[1]-slotsValue[11].YRange[0]),
			},
		},
	}
	//buildSlotsValue[1] = slotsValue[11]
	/*fmt.Println("FS", ParamFieldSize)
	fmt.Println("KingsCoords", KingCoords)
	fmt.Println("Slot 1", slotsValue[1])
	fmt.Println("Slot 2", slotsValue[2])
	fmt.Println("Slot 3", slotsValue[3])
	fmt.Println("Slot 4", slotsValue[4])
	fmt.Println("Slot 5", slotsValue[5])
	fmt.Println("Slot 6", slotsValue[6])
	fmt.Println("Slot 7", slotsValue[7])
	fmt.Println("Slot 8", slotsValue[8])
	fmt.Println("Slot 10", slotsValue[10])
	fmt.Println("Slot 10", slotsValue[11])
	fmt.Println("Slot 1", buildSlotsValue[1])
	fmt.Println("Slot 2", buildSlotsValue[2])
	fmt.Println("Slot 3", buildSlotsValue[3])
	fmt.Println("Slot 4", buildSlotsValue[4])
	fmt.Println("Slot 5", buildSlotsValue[5])
	fmt.Println("Slot 6", buildSlotsValue[6])
	fmt.Println("Slot 7", buildSlotsValue[7])
	fmt.Println("Slot 8", buildSlotsValue[8])
	fmt.Println("Slot 10", slotsValue[10])
	incboxval := *IncomeBox
	fmt.Println("IncBox0", incboxval[0])*/
}

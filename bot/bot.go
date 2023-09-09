package bot

import (
	//"fmt"
	"math/rand"
	"strconv"
	"strings"
	"ws/datatypes"
	"ws/game"
	"ws/game/units"
	"ws/json"
	"ws/models"
	"ws/websocket"
)

var BOTS []int
var COORDS datatypes.SlotMap
var Answer func(*websocket.Conn, []byte, *models.State, *models.Links, datatypes.ChannelsMap)
var kingAttrs = []string{"hpreg", "atk"}

func buildingStep(slot int, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	player := state.Players[slot]
	playerUnits := state.PlayerUnits[slot]

	if state.Servertime < 2 {
		return
	}
	if float64(player.Value) >= 1.1*units.Waves[state.Wave].RecValue {
		return
	}
	// Upgrade existing units
	for _, unit := range playerUnits {
		if len(unit.Upgrades) > 0 {
			for _, upgradeName := range unit.Upgrades {
				upgradeProto := units.UNITS_BY_NAME[upgradeName]
				if player.Gold >= upgradeProto.Price && (player.Food[1]-player.Food[0]) >= (upgradeProto.Food-unit.Food) {
					upgrade(upgradeName, unit.Id, slot, state, links, channels)
					return
				}
			}
		}
	}
	// Build units
	for index, _ := range player.AvailableUnits {
		unitProto := units.UNITS_BY_NAME[player.AvailableUnits[len(player.AvailableUnits)-index-1]]
		if player.Gold >= unitProto.Price && player.Food[1]-player.Food[0] >= unitProto.Food {
			build(unitProto.Name, slot, state, links, channels)
			return
		}
	}
}

func step(slot int, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	player := state.Players[slot]

	if state.Servertime < 2 {
		return
	}

	// upgrade farm if food cap is near
	if player.Food[1]-player.Food[0] < 3 && player.Gold >= game.FOOD_GOLD && player.Stone >= game.FOOD_STONE {
		upgradeFarm(slot, state, links, channels)
		return
	}

	if player.Stone > game.FOOD_STONE+game.KING_GRADE_PRICE {
		upgradeKing(slot, state, links, channels, kingAttrs[rand.Intn(len(kingAttrs))])
		return
	}
	if float64(player.Value) < 1.3*units.Waves[state.Wave].RecValue {
		return
	}
	if player.Gold >= player.StoneGrade.Price[0] && player.Stone >= player.StoneGrade.Price[1] {
		upgradeStone(slot, state, links, channels)
		return
	}
}

func upgrade(name string, id int, slot int, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	message, _ := json.Marshal(map[string]interface{}{
		"type":    "upgrade_unit!",
		"uid":     "ankxel",
		"slot":    slot,
		"id":      id,
		"upgrade": name,
	})
	Answer(nil, message, state, links, channels)
}

func build(name string, slot int, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	coords := COORDS[slot]
	unitProto := units.UNITS_BY_NAME[name]
	yCoords := coords.YRange
	if (coords.YRange[0] > game.MapCenter[1] && unitProto.Projectile == game.Melee) || (coords.YRange[0] < game.MapCenter[1] && unitProto.Projectile != game.Melee) {
		yCoords[0] += 0.5 * (coords.YRange[1] - coords.YRange[0])
	} else {
		yCoords[1] -= 0.5 * (coords.YRange[1] - coords.YRange[0])
	}
	if state.Players[slot].Guild == game.Recruits {
		yCoords = coords.YRange
	}
	x, y := Uniform(coords.XRange[0], coords.XRange[1]), Uniform(yCoords[0], yCoords[1])
	message, _ := json.Marshal(map[string]interface{}{
		"type":   "build!",
		"uid":    "ankxel",
		"slot":   slot,
		"coords": []float64{x, y},
		"name":   name,
	})
	Answer(nil, message, state, links, channels)
}

func upgradeFarm(slot int, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	message, _ := json.Marshal(map[string]interface{}{
		"type": "upgrade_farm!",
		"uid":  "ankxel",
		"slot": slot,
	})
	Answer(nil, message, state, links, channels)
}

func upgradeKing(slot int, state *models.State, links *models.Links, channels datatypes.ChannelsMap, attr string) {
	message, _ := json.Marshal(map[string]interface{}{
		"type":      "upgrade_king!",
		"uid":       "ankxel",
		"attribute": attr,
		"slot":      slot,
	})
	Answer(nil, message, state, links, channels)
}

func upgradeStone(slot int, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	message, _ := json.Marshal(map[string]interface{}{
		"type": "upgrade_stone!",
		"uid":  "ankxel",
		"slot": slot,
	})
	Answer(nil, message, state, links, channels)
}

func JoinBots(slots *string, state *models.State, links *models.Links, channels datatypes.ChannelsMap, onMessage func(*websocket.Conn, []byte, *models.State, *models.Links, datatypes.ChannelsMap)) {
	if *slots == "" {
		return
	}
	Answer = onMessage

	COORDS = *game.BuildSlots
	// = map[int][][]float64{
	//	5: [][]int{[]int{130, 139}, []int{33, 48}},
	//	6: [][]int{[]int{187, 196}, []int{33, 48}},
	//	3: [][]int{[]int{2, 11}, []int{8, 23}},
	//	2: [][]int{[]int{59, 68}, []int{33, 48}},
	//	8: [][]int{[]int{187, 196}, []int{8, 23}},
	//	1: [][]int{[]int{2, 11}, []int{33, 48}},
	//	4: [][]int{[]int{59, 68}, []int{8, 23}},
	//	7: [][]int{[]int{130, 139}, []int{8, 23}},
	//}
	ids := strings.Split(*slots, ",")
	for _, slot := range ids {
		id, _ := strconv.Atoi(slot)
		BOTS = append(BOTS, id)

		message, _ := json.Marshal(map[string]interface{}{
			"type":     "connect!",
			"username": "bot" + slot,
			"slot":     id,
		})
		Answer(nil, message, state, links, channels)
		message, _ = json.Marshal(map[string]interface{}{
			"type":  "pick_guild!",
			"guild": game.Random,
			"slot":  id,
			"uid":   "ankxel",
		})
		Answer(nil, message, state, links, channels)
	}
}

func BotsStep(state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	for _, slot := range BOTS {
		if state.Phase == game.Building {
			buildingStep(slot, state, links, channels)
		}
		step(slot, state, links, channels)
	}
}

func Uniform(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

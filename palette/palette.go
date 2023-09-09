package palette

import (
	"sync"
	"ws/datatypes"
	"ws/models"
	"ws/utils"
	"ws/websocket"
)

var PALETTE = map[string]datatypes.HandlerFunc{
	"connect!":              connect,
	"pause!":                authorize(pause),
	"disconnect!":           authorize(disconnect),
	"pick_guild!":           authorize(pickGuild),
	"reset_guild!":          authorize(resetGuild),
	"build!":                authorize(build),
	"upgrade_unit!":         authorize(upgradeUnit),
	"sell_unit!":            authorize(sellUnit),
	"upgrade_farm!":         authorize(upgradeFarm),
	"cancel_farm_upgrade!":  authorize(cancelFarmUpgrade),
	"upgrade_stone!":        authorize(upgradeStone),
	"cancel_stone_upgrade!": authorize(cancelStoneUpgrade),
	"send_income!":          authorize(sendIncome),
	"upgrade_king!":         authorize(upgradeKing),
	"leave!":                authorize(leave),
	"state":                 state,
	"add_resources!":        authorize(addResources),
	"change_wave!":          authorize(changeWave),
	"change_king_hp!":       authorize(changeKingHp),
}

var DEFAULT_ACTION = PALETTE["state"]

func getOrDefault(data datatypes.HandlerHash, key string, defaultVal datatypes.HandlerFunc) datatypes.HandlerFunc {
	if val, ok := data[key]; ok {
		return val
	} else {
		return defaultVal
	}
}

var locks = map[string]*sync.Mutex{}

func Answer(ws *websocket.Conn, message []byte, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	defer utils.CatchAny()
	message_json := utils.Recieve(message)
	messageType := message_json["type"].(string)

	locks[messageType].Lock()
	action := getOrDefault(PALETTE, messageType, DEFAULT_ACTION)
	action(ws, message_json, state, links, channels)
	locks[messageType].Unlock()

}

func init() {
	for key, _ := range PALETTE {
		locks[key] = &sync.Mutex{}
	}
}

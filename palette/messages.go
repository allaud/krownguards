package palette

import (
	"fmt"
	"github.com/satori/go.uuid"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
	"ws/cron"
	"ws/datatypes"
	"ws/flow"
	"ws/game"
	"ws/game/tooltips"
	"ws/game/units"
	"ws/models"
	"ws/utils"
	"ws/websocket"
)

func authorize(fn func(*websocket.Conn, datatypes.MessageData, *models.State, *models.Links, datatypes.ChannelsMap)) func(*websocket.Conn, datatypes.MessageData, *models.State, *models.Links, datatypes.ChannelsMap) {
	return func(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
		uid, ok := data["uid"].(string)
		if !ok {
			return
		}
		slotFl, ok := data["slot"].(float64)
		if !ok {
			return
		}
		slot := int(slotFl)

		if channel, ok := channels[slot]; ok {
			if uid != channel.Uid && uid != "ankxel" {
				fmt.Println("WRONG UID!")
				return
			}
			fn(ws, data, state, links, channels)
		}
	}
}

func pause(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	slotFl, ok := data["slot"].(float64)
	if !ok {
		return
	}
	slot := int(slotFl)
	player, ok := state.Players[slot]
	if !ok {
		return
	}
	if state.Pause == 0 {
		if player.PauseCap > 0 {
			state.Pause = slot
			utils.SendNotification(state, channels, player.Name+" paused the game!", 0)
			utils.SendAll(channels, &datatypes.MessageData{
				"type":  "state",
				"state": state,
			}, 0)
			flow.Messagebox <- datatypes.Signal{
				Type:     "pause",
				Duration: 1 * time.Second,
				Fn:       func() { cron.PauseTimer(state, links, channels) },
			}
		}
		return
	}
	pauser, ok := state.Players[state.Pause]
	if !ok {
		state.Pause = slot
	}
	if state.Pause == slot {
		state.Pause = 0
	} else {
		if pauser.PauseCap <= 0 {
			state.Pause = 0
		}
	}
	if state.Pause == 0 {
		utils.SendNotification(state, channels, player.Name+" resumed the game!", 0)
		utils.SendAll(channels, &datatypes.MessageData{
			"type":  "state",
			"state": state,
		}, 0)
	}
}

func connect(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	var uid string
	username, ok := data["username"].(string)
	if !ok {
		return
	}
	slotFl, ok := data["slot"].(float64)
	if !ok {
		return
	}
	slot := int(slotFl)
	if (slot < 1) || (slot > 8) {
		return
	}

	if !utils.ContainsString(state.Usernames, username) && !utils.ContainsString(game.BotNames, username) {
		return
	}
	for plSlot, player := range state.Players {
		if username == player.Name && slot != plSlot {
			fmt.Println("FALSE: this username already in other slot")
			return
		}
	}
	if _, contains := state.Players[slot]; contains {
		if username != state.Players[slot].Name {
			fmt.Println("FALSE: trying to connect in existing slot")
			return
		}
		// reconnect block
		channels[slot].Channel = ws
		channels[slot].Channel.Username = username
		channels[slot].Mutex = &sync.Mutex{}
		uid = channels[slot].Uid
	} else {
		// new connection block. check if "start"
		// make a copy of Summons
		playerSummons := map[string]*models.Summon{}
		for k, v := range units.Summons[1] {
			tempV := *v
			playerSummons[k] = &tempV
		}
		// add new player to game state
		state.Players[slot] = &models.Player{
			Gold:  game.START_GOLD,
			Stone: game.START_STONE,
			StoneGrade: models.StoneUpgrade{
				Price:     [2]int{game.STONE_GRADE[1][0], game.STONE_GRADE[1][1]},
				CurrGrade: 1,
				MaxGrade:  len(game.STONE_GRADE),
				StoneInc:  game.STONE_GRADE[0][2],
				Cooldown:  game.STONE_GRADE_CD,
			},
			Food: game.START_FOOD,
			FoodGrade: models.FoodUpgrade{
				Cooldown: game.FOOD_GRADE_CD,
			},
			Name:        username,
			GuildsList:  append(units.GUILDS, game.Random),
			Resets:      game.GUILD_RESET_CAP,
			ResetPrice:  game.GUILD_RESET_PRICE[0],
			SummonUnits: playerSummons,
			PauseCap:    game.PLAYER_PAUSE_CAP,
		}

		uid = uuid.NewV4().String()
		if ws != nil {
			ws.Username = username
		}
		channels[slot] = &datatypes.Connection{
			Channel: ws,
			Mutex:   &sync.Mutex{},
			Uid:     uid,
		}
		state.WaveUnits[slot] = []*models.Unit{}
		state.PlayerUnits[slot] = []*models.Unit{}
	}

	utils.SendNotification(state, channels, username+" connected", 0)
	utils.Send(ws, &datatypes.MessageData{
		"type": "connect!",
		"uid":  uid,
		"gamedata": map[string]interface{}{
			"units":         units.UNITS_BY_NAME,
			"summons":       units.UNITS_BY_TIER[game.Summon],
			"waves":         units.Waves,
			"slots":         *game.Slots,
			"spawnslots":    *game.SpawnSlots,
			"buildslots":    *game.BuildSlots,
			"unitwaypoints": *game.UnitWaypoints,
			"wavewaypoints": *game.WaveWaypoints,
			"typetooltips":  tooltips.TypeTooltips,
			"abtooltips":    tooltips.AbTooltips,
			"aftooltips":    tooltips.AfTooltips,
			"guildtooltips": tooltips.GuildTooltips,
		},
		"state": state,
	}, channels[slot].Mutex, false)
	utils.SendAll(channels, &datatypes.MessageData{
		"type":  "state",
		"state": state,
	}, slot)

	if len(state.Players) == state.WaitFor {
		state.WaitFor = 0
		cron.Run(state, links, channels)
	}
}

func disconnect(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	username, ok := data["username"].(string)
	if !ok {
		return
	}
	utils.SendNotification(state, channels, username+" disconnected", 0)
}

func leave(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	utils.SendAll(channels, &datatypes.MessageData{
		"type":  "leave!",
		"state": state,
	}, 0)
	slotFl, ok := data["slot"].(float64)
	if !ok {
		return
	}
	slot := int(slotFl)

	username := state.Players[slot].Name
	index := utils.GetIndex(state.Usernames, username)
	// move leaver from usernames to leaver names
	if index != -1 {
		state.Leavers = append(state.Leavers, username)
		state.Usernames = append(state.Usernames[:index], state.Usernames[index+1:]...)
	}
	// delete player entities
	delete(channels, slot)
	if state.Phase == game.Building {
		delete(state.Players, slot)
		delete(state.PlayerUnits, slot)
	}
	ws.Close()
	resp, err := http.Get(state.LobbyAddr + "/leavecallback?username=" + username)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func pickGuild(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	if state.Pause != 0 {
		return
	}
	slotFl, ok := data["slot"].(float64)
	if !ok {
		return
	}
	slot := int(slotFl)
	player, ok := state.Players[slot]
	if !ok {
		return
	}
	if player.Guild != "" {
		return
	}
	guild, ok := data["guild"].(string)
	if !ok {
		return
	}
	if !utils.ContainsString(player.GuildsList, guild) {
		return
	}
	if guild == game.Random {
		guild = units.GUILDS[rand.Intn(len(units.GUILDS))]
	}
	player.Guild = guild
	player.AvailableUnits = units.UNITS_BY_GUILD[guild]
	if guild == game.Draft || guild == game.DynamicDraft {
		for tier := 1; tier <= 6; tier++ {
			unit := units.RandomByTier(tier)
			player.AvailableUnits = append(player.AvailableUnits, unit)
		}
	}
	player.GuildsList = []string{}

	utils.SendAll(channels, &datatypes.MessageData{
		"type":  "pick_guild!",
		"state": state,
	}, 0)
}

func resetGuild(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	if state.Pause != 0 {
		return
	}
	slotFl, ok := data["slot"].(float64)
	if !ok {
		return
	}
	slot := int(slotFl)
	player, ok := state.Players[slot]
	if !ok {
		return
	}
	if player.Guild == "" || len(player.GuildsList) != 0 {
		return
	}
	if player.Resets == 0 {
		return
	}
	if player.Gold < player.ResetPrice[0] || player.Stone < player.ResetPrice[1] {
		return
	}

	player.Resets = player.Resets - 1
	player.Gold = player.Gold - player.ResetPrice[0]
	player.Stone = player.Stone - player.ResetPrice[1]
	if game.GUILD_RESET_CAP-player.Resets < len(game.GUILD_RESET_PRICE) {
		player.ResetPrice = game.GUILD_RESET_PRICE[game.GUILD_RESET_CAP-player.Resets]
	}
	player.Guild = ""
	player.AvailableUnits = []string{}

	randGuild := ""
	for len(player.GuildsList) < game.GUILDS_LIST_CAP {
		randGuild = units.GUILDS[rand.Intn(len(units.GUILDS))]
		if !utils.ContainsString(player.GuildsList, randGuild) {
			player.GuildsList = append(player.GuildsList, randGuild)
		}
	}

	utils.SendAll(channels, &datatypes.MessageData{
		"type":  "pick_guild!",
		"state": state,
	}, 0)
}

func build(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	if state.Pause != 0 {
		return
	}
	if state.Phase != game.Building {
		return
	}
	slotFl, ok := data["slot"].(float64)
	if !ok {
		return
	}
	slot := int(slotFl)
	player, ok := state.Players[slot]
	if !ok {
		return
	}
	coordsI, ok := data["coords"].([]interface{})
	if !ok {
		return
	}
	coordX, ok := coordsI[0].(float64)
	if !ok {
		return
	}
	coordY, ok := coordsI[1].(float64)
	if !ok {
		return
	}
	coords := [2]float64{coordX, coordY}

	gridCoords := utils.GetGridCoords(coords)

	if player.Guild == "" {
		return
	}
	if !utils.IsInBuildSlot(game.BuildSlots, slot, gridCoords) {
		return
	}
	name, ok := data["name"].(string)
	if !ok {
		return
	}
	if !utils.ContainsString(player.AvailableUnits, name) {
		return
	}
	unit, exists := units.UNITS_BY_NAME[name]
	if !exists {
		return
	}
	if player.Gold < unit.Price {
		return
	}
	if (player.Food[1] - player.Food[0]) < unit.Food {
		return
	}
	for _, pUnit := range state.PlayerUnits[slot] {
		if utils.XYRange(gridCoords, pUnit.Coords) < 1 {
			return
		}
	}
	/*
		playersToCheck := []int{slot}
		assSlot, exists := utils.AssocSlot(slot)
		if exists && state.Players[assSlot] != nil {
			playersToCheck = append(playersToCheck, assSlot)
		}
		for _, owner := range playersToCheck {
			for _, ownersUnit := range state.PlayerUnits[owner] {
				if utils.XYRange(gridCoords, ownersUnit.Coords) < 1 {
					return
				}
			}
		}*/
	player.Gold = player.Gold - unit.Price
	player.Food[0] = player.Food[0] + unit.Food
	if player.Guild == game.Recruits {
		tier, _ := strconv.Atoi(name[5:])
		unit = units.UNITS_BY_NAME[units.RandomByTier(tier)]
	}
	player.Value += unit.Price
	unit.Id = utils.UniqId()
	unit.Slot = slot
	unit.Player = slot
	unit.Waypoint = 1
	unit.Action = game.Idle
	if player.Guild == game.Recruits || player.Guild == game.DynamicDraft {
		utils.DampPrice(&unit)
		//unit.SellPrice += unit.Price / 2
		//unit.Price = 0
	}
	if player.Guild == game.DynamicDraft {
		player.AvailableUnits[unit.Tier-1] = units.RandomByTier(unit.Tier)
	}
	unit.Coords = gridCoords
	unit.Affected = make(map[string]*models.Effect)
	unit.Abilities = utils.CopyAbilities(unit.Abilities)
	state.PlayerUnits[slot] = append(state.PlayerUnits[slot], &unit)
	utils.SendAll(channels, &datatypes.MessageData{
		"type":  "build!",
		"state": state,
	}, 0)
}

func upgradeUnit(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	if state.Pause != 0 {
		return
	}
	if state.Phase != game.Building {
		return
	}
	slotFl, ok := data["slot"].(float64)
	if !ok {
		return
	}
	slot := int(slotFl)
	player, ok := state.Players[slot]
	if !ok {
		return
	}
	idFl, ok := data["id"].(float64)
	if !ok {
		return
	}
	id := int(idFl)
	for uIndex, unit := range state.PlayerUnits[slot] {
		if unit.Id != id {
			continue
		}
		upgrade, ok := data["upgrade"].(string)
		if !ok {
			return
		}
		if !utils.ContainsString(unit.Upgrades, upgrade) {
			return
		}
		unitUpgrade, exists := units.UNITS_BY_NAME[upgrade]
		if !exists {
			return
		}
		if player.Gold < unitUpgrade.Price {
			return
		}
		if (player.Food[1] - player.Food[0]) < (unitUpgrade.Food - unit.Food) {
			return
		}
		player.Gold = player.Gold - unitUpgrade.Price
		player.Food[0] = player.Food[0] + (unitUpgrade.Food - unit.Food)
		player.Value = player.Value + unitUpgrade.Price
		unitUpgrade.Price += unit.Price
		unitUpgrade.SellPrice = unit.SellPrice
		unitUpgrade.Id = utils.UniqId()
		unitUpgrade.Slot = unit.Slot
		unitUpgrade.Player = unit.Player
		unitUpgrade.Waypoint = 1
		unitUpgrade.Action = game.Idle
		unitUpgrade.Coords = unit.Coords
		unitUpgrade.Abilities = utils.CopyAbilities(unitUpgrade.Abilities)
		unitUpgrade.Affected = make(map[string]*models.Effect)
		state.PlayerUnits[slot][uIndex] = &unitUpgrade
		utils.SendAll(channels, &datatypes.MessageData{
			"type":  "upgrade_unit!",
			"state": state,
		}, 0)
		return
	}
}

func sellUnit(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	if state.Pause != 0 {
		return
	}
	if state.Phase != game.Building {
		return
	}
	slotFl, ok := data["slot"].(float64)
	if !ok {
		return
	}
	slot := int(slotFl)
	player, ok := state.Players[slot]
	if !ok {
		return
	}
	idFl, ok := data["id"].(float64)
	if !ok {
		return
	}
	id := int(idFl)
	for uIndex, unit := range state.PlayerUnits[slot] {
		if unit.Id != id {
			continue
		}
		player.Gold += unit.Price + unit.SellPrice
		player.Food[0] = player.Food[0] - unit.Food
		player.Value -= unit.Price + unit.SellPrice*2
		state.PlayerUnits[slot][uIndex] = nil
		state.PlayerUnits[slot] = append(state.PlayerUnits[slot][:uIndex], state.PlayerUnits[slot][uIndex+1:]...)
		utils.SendAll(channels, &datatypes.MessageData{
			"type":  "sell_unit!",
			"state": state,
		}, 0)
		return
	}
}

func sendIncome(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	if state.Pause != 0 {
		return
	}
	slotFl, ok := data["slot"].(float64)
	if !ok {
		return
	}
	slot := int(slotFl)
	player, ok := state.Players[slot]
	if !ok {
		return
	}
	name, ok := data["name"].(string)
	if !ok {
		return
	}
	summon, ok := player.SummonUnits[name]
	if !ok {
		return
	}
	unit, exists := units.UNITS_BY_NAME[name]
	if !exists {
		return
	}
	if player.Stone < unit.Price {
		return
	}
	if summon.Current < 1 {
		return
	}
	player.Stone = player.Stone - unit.Price
	player.Income = player.Income + unit.Bounty
	summon.Current = summon.Current - 1
	unit.Id = utils.UniqId()
	unit.Slot = slot
	unit.Player = slot
	unit.Waypoint = 1
	unit.Action = game.Idle
	unit.Coords = utils.GetRandomCoords(unit.Size, game.IncomeBox, slot/game.SideCoef)
	unit.Abilities = utils.CopyAbilities(unit.Abilities)
	unit.Affected = make(map[string]*models.Effect)
	player.IncomeUnits = append(player.IncomeUnits, &unit)
	//fmt.Println("Player", player.Name, "summoned:", player.IncomeUnits[len(player.IncomeUnits)-1].Name, "::", player.IncomeUnits[len(player.IncomeUnits)-1].Id)
	//fmt.Println("Now player", player.Name, "have:", player.IncomeUnits)
	utils.SendAll(channels, &datatypes.MessageData{
		"type":  "send_income!",
		"state": state,
	}, 0)
}

func upgradeFarm(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	if state.Pause != 0 {
		return
	}
	slotFl, ok := data["slot"].(float64)
	if !ok {
		return
	}
	slot := int(slotFl)
	player, ok := state.Players[slot]
	if !ok {
		return
	}
	if player.FoodGrade.InProgress {
		return
	}
	if (player.Gold < game.FOOD_GOLD) || (player.Stone < game.FOOD_STONE) {
		return
	}
	player.Gold = player.Gold - game.FOOD_GOLD
	player.Stone = player.Stone - game.FOOD_STONE
	player.FoodGrade.InProgress = true
	utils.SendAll(channels, &datatypes.MessageData{
		"type":  "upgrade_farm!",
		"state": state,
	}, 0)
}

func cancelFarmUpgrade(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	if state.Pause != 0 {
		return
	}
	slotFl, ok := data["slot"].(float64)
	if !ok {
		return
	}
	slot := int(slotFl)
	player, ok := state.Players[slot]
	if !ok {
		return
	}
	if !player.FoodGrade.InProgress {
		return
	}
	player.Gold = player.Gold + game.FOOD_GOLD
	player.Stone = player.Stone + game.FOOD_STONE
	player.FoodGrade.InProgress = false
	player.FoodGrade.Timer = 0
	utils.SendAll(channels, &datatypes.MessageData{
		"type":  "cancel_farm_upgrade!",
		"state": state,
	}, 0)
}

func upgradeStone(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	if state.Pause != 0 {
		return
	}
	slotFl, ok := data["slot"].(float64)
	if !ok {
		return
	}
	slot := int(slotFl)
	player, ok := state.Players[slot]
	if !ok {
		return
	}
	if player.StoneGrade.InProgress {
		return
	}
	if player.StoneGrade.CurrGrade >= player.StoneGrade.MaxGrade {
		return
	}
	if player.Gold < player.StoneGrade.Price[0] {
		return
	}
	if player.Stone < player.StoneGrade.Price[1] {
		return
	}
	player.Gold = player.Gold - player.StoneGrade.Price[0]
	player.Stone = player.Stone - player.StoneGrade.Price[1]
	player.StoneGrade.InProgress = true
	utils.SendAll(channels, &datatypes.MessageData{
		"type":  "upgrade_stone!",
		"state": state,
	}, 0)
}

func cancelStoneUpgrade(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	if state.Pause != 0 {
		return
	}
	slotFl, ok := data["slot"].(float64)
	if !ok {
		return
	}
	slot := int(slotFl)
	player, ok := state.Players[slot]
	if !ok {
		return
	}
	if !player.StoneGrade.InProgress {
		return
	}
	player.Gold = player.Gold + player.StoneGrade.Price[0]
	player.Stone = player.Stone + player.StoneGrade.Price[1]
	player.StoneGrade.InProgress = false
	player.StoneGrade.Timer = 0
	utils.SendAll(channels, &datatypes.MessageData{
		"type":  "cancel_stone_upgrade!",
		"state": state,
	}, 0)
}

func upgradeKing(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	if state.Pause != 0 {
		return
	}
	slotFl, ok := data["slot"].(float64)
	if !ok {
		return
	}
	slot := int(slotFl)
	player, ok := state.Players[slot]
	if !ok {
		return
	}
	attr, ok := data["attribute"].(string)
	if !ok {
		return
	}
	king := state.Kings[slot/game.SideCoef]
	kingAttr := state.KingAttrs[slot/game.SideCoef]
	var attrGrade models.Upgrade
	switch attr {
	case "atk":
		attrGrade = kingAttr.AtkGrade
	case "maxhp":
		attrGrade = kingAttr.MaxHpGrade
	case "hpreg":
		attrGrade = kingAttr.HpRegGrade
	default:
		return
	}
	if attrGrade.CurrGrade >= attrGrade.MaxGrade {
		return
	}
	if player.Stone < attrGrade.Price {
		return
	}
	player.Stone = player.Stone - attrGrade.Price
	player.Income = player.Income + game.KING_GRADE_INCOME
	attrName := ""
	attrLvl := ""
	switch attr {
	case "atk":
		king.Atk[0] += kingAttr.AtkGrade.GradeStep
		king.Atk[1] += kingAttr.AtkGrade.GradeStep
		kingAttr.AtkGrade.CurrGrade += 1
		attrName = "attack"
		attrLvl = strconv.Itoa(kingAttr.AtkGrade.CurrGrade) + "/" + strconv.Itoa(kingAttr.AtkGrade.MaxGrade)
	case "maxhp":
		king.MaxHp += kingAttr.MaxHpGrade.GradeStep
		king.Hp += kingAttr.MaxHpGrade.GradeStep
		king.Def -= float64(kingAttr.MaxHpGrade.CurrGrade / 4)
		kingAttr.MaxHpGrade.CurrGrade += 1
		king.Def += float64(kingAttr.MaxHpGrade.CurrGrade / 4)
		attrName = "max health"
		attrLvl = strconv.Itoa(kingAttr.MaxHpGrade.CurrGrade) + "/" + strconv.Itoa(kingAttr.MaxHpGrade.MaxGrade)
	case "hpreg":
		king.HpReg += kingAttr.HpRegGrade.GradeStep
		kingAttr.HpRegGrade.CurrGrade += 1
		attrName = "health regeneration"
		attrLvl = strconv.Itoa(kingAttr.HpRegGrade.CurrGrade) + "/" + strconv.Itoa(kingAttr.HpRegGrade.MaxGrade)
	default:
		return
	}
	utils.SendNotification(state, channels, player.Name+
		" upgraded "+king.Name+"'s "+attrName+" up to "+attrLvl, 0)
	utils.SendAll(channels, &datatypes.MessageData{
		"type":  "upgrade_king!",
		"state": state,
	}, 0)
}

func state(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	//utils.Send(ws, datatypes.MessageData{
	//	"type":  "state",
	//	"state": state,
	//})
}

func addResources(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	slot, ok := data["slot"].(float64)
	if !ok {
		return
	}
	gold, ok := data["gold"].(float64)
	if !ok {
		return
	}
	stone, ok := data["stone"].(float64)
	if !ok {
		return
	}
	food, ok := data["food"].(float64)
	if !ok {
		return
	}
	player, ok := state.Players[int(slot)]
	if !ok {
		return
	}
	player.Gold += int(gold)
	player.Stone += int(stone)
	player.Food[1] += int(food)
	utils.SendNotification(state, channels, "TEST: Added "+strconv.FormatFloat(gold, 'f', 0, 64)+
		" gold, "+strconv.FormatFloat(stone, 'f', 0, 64)+" stone, "+strconv.FormatFloat(food, 'f', 0, 64)+
		" food to player in slot "+strconv.FormatFloat(slot, 'f', 0, 64), 0)
	utils.SendAll(channels, &datatypes.MessageData{
		"type":  "state",
		"state": state,
	}, 0)
}

func changeWave(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	wave, ok := data["wave"].(float64)
	if !ok {
		return
	}
	state.Wave = int(wave)
	utils.SendNotification(state, channels, "TEST: Wave changed to "+strconv.FormatFloat(wave, 'f', 0, 64), 0)
	utils.SendAll(channels, &datatypes.MessageData{
		"type":  "state",
		"state": state,
	}, 0)
}

func changeKingHp(ws *websocket.Conn, data datatypes.MessageData, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	kingId, ok := data["king"].(float64)
	if !ok {
		return
	}
	hp, ok := data["hp"].(float64)
	if !ok {
		return
	}
	king := state.Kings[int(kingId)]
	king.Hp += hp
	if king.Hp > king.MaxHp {
		king.Hp = king.MaxHp
	}
	utils.SendNotification(state, channels, "TEST: "+king.Name+" hp changed by "+strconv.FormatFloat(hp, 'f', 0, 64), 0)
	utils.SendAll(channels, &datatypes.MessageData{
		"type":  "state",
		"state": state,
	}, 0)
}

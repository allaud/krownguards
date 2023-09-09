package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"ws/bot"
	"ws/cron"
	"ws/datatypes"
	"ws/flow"
	"ws/game"
	"ws/game/units"
	"ws/json"
	"ws/models"
	"ws/palette"
	"ws/utils"
	"ws/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var AppState = &models.State{
	Servertime:    0,
	WaveCountdown: game.PREPARATION_TIME,
	Phase:         game.Start,
	Wave:          1,
	ArenaWaves:    game.ARENA_EVERY_WAR,
	Usernames:     make([]string, 0),
	Leavers:       make([]string, 0),
	Notifications: make([]string, 0),
	Players:       make(map[int]*models.Player),
	PlayerUnits:   make(map[int][]*models.Unit),
	WaveUnits:     make(map[int][]*models.Unit),
}

var AppLinks = &models.Links{
	ProjLinks: []*models.ProjLink{},
}

var Channels = datatypes.ChannelsMap{}

func InitKings(state *models.State) {
	westKing := units.King
	westKing.Name = "West " + westKing.Name
	westKing.Coords = game.KingCoords[0]
	westKing.Dir = game.KingDir[0]
	westKing.Id = utils.UniqId()
	westKing.Action = game.Idle
	westKing.Slot = 10
	westKing.Player = 10
	westKing.Affected = make(map[string]*models.Effect)
	westKingAttrs := units.KingAttr
	eastKing := units.King
	eastKing.Name = "East " + eastKing.Name
	eastKing.Coords = game.KingCoords[1]
	eastKing.Dir = game.KingDir[1]
	eastKing.Id = utils.UniqId()
	eastKing.Action = game.Idle
	eastKing.Slot = 11
	eastKing.Player = 11
	eastKing.Affected = make(map[string]*models.Effect)
	eastKingAttrs := units.KingAttr
	state.Kings = [2]*models.Unit{&westKing, &eastKing}
	state.KingAttrs = [2]*models.KingUpgrades{&westKingAttrs, &eastKingAttrs}
}

func onMessage(ws *websocket.Conn, message []byte, state *models.State, links *models.Links, channels datatypes.ChannelsMap) {
	flow.Messagebox <- datatypes.Signal{
		Type:     "message",
		Duration: 0,
		Fn:       func() { palette.Answer(ws, message, state, links, channels) },
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil, true)
	defer func() {
		fmt.Println("Defered close websocket")
		json_data, _ := json.Marshal(map[string]interface{}{
			"slot":     1,
			"uid":      "ankxel",
			"username": ws.Username,
			"type":     "disconnect!",
		})
		onMessage(nil, json_data, AppState, AppLinks, Channels)
		ws.Username = "offline"
		ws.Close()
	}()
	if err != nil {
		return
	}
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Connection closed")
			return
		}
		onMessage(ws, p, AppState, AppLinks, Channels)
	}
}

func stateHandler(w http.ResponseWriter, r *http.Request) {
	json_data, _ := json.Marshal(AppState)
	fmt.Println("===> /state")
	fmt.Fprintf(w, string(json_data))
}

func channelsHandler(w http.ResponseWriter, r *http.Request) {
	json_data, _ := json.Marshal(Channels)
	fmt.Fprintf(w, string(json_data))
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	utils.EndGame(AppState)
}

func pauseHandler(w http.ResponseWriter, r *http.Request) {
	json_data, _ := json.Marshal(map[string]interface{}{
		"slot": 1,
		"uid":  "ankxel",
		"type": "pause!",
	})
	onMessage(nil, json_data, AppState, AppLinks, Channels)
}

func leaveHandler(w http.ResponseWriter, r *http.Request) {
	slot, _ := strconv.Atoi(r.FormValue("slot"))
	json_data, _ := json.Marshal(map[string]interface{}{
		"slot": slot,
		"uid":  "ankxel",
		"type": "leave!",
	})
	onMessage(nil, json_data, AppState, AppLinks, Channels)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	slot, ok := strconv.Atoi(r.FormValue("slot"))
	if ok == nil {
		gold, ok := strconv.Atoi(r.FormValue("gold"))
		if !(ok == nil) {
			gold = 0
		}
		stone, ok := strconv.Atoi(r.FormValue("stone"))
		if !(ok == nil) {
			stone = 0
		}
		food, ok := strconv.Atoi(r.FormValue("food"))
		if !(ok == nil) {
			food = 0
		}
		resources_data, _ := json.Marshal(map[string]interface{}{
			"slot":  slot,
			"gold":  gold,
			"stone": stone,
			"food":  food,
			"uid":   "ankxel",
			"type":  "add_resources!",
		})
		onMessage(nil, resources_data, AppState, AppLinks, Channels)
	}

	wave, ok := strconv.Atoi(r.FormValue("wave"))
	if ok == nil {
		wave_data, _ := json.Marshal(map[string]interface{}{
			"slot": 1,
			"wave": wave,
			"uid":  "ankxel",
			"type": "change_wave!",
		})
		onMessage(nil, wave_data, AppState, AppLinks, Channels)
	}

	king, ok := strconv.Atoi(r.FormValue("king"))
	if ok == nil {
		if king != 0 && king != 1 {
			king = 0
		}
		hp, ok := strconv.Atoi(r.FormValue("hp"))
		if !(ok == nil) {
			hp = 0
		}
		king_data, _ := json.Marshal(map[string]interface{}{
			"slot": 1,
			"king": king,
			"hp":   hp,
			"uid":  "ankxel",
			"type": "change_king_hp!",
		})
		onMessage(nil, king_data, AppState, AppLinks, Channels)
	}
}

func main() {
	id := flag.String("id", "", "Game unique ID in lobby sessions")
	port := flag.String("port", "8888", "Port to listen on")
	waitFor := flag.Int("wait_for", 1, "Wait to connect N users")
	botSlots := flag.String("bots", "", "Bot slots (e.g. --bots=5,6)")
	lobbyAddr := flag.String("lobby", "http://krownguards.com:8899", "Lobby address")
	usernames := flag.String("usernames", "", "Usernames of expected players")
	flag.Parse()

	AppState.Id = *id
	AppState.WaitFor = *waitFor
	AppState.Port = *port
	AppState.LobbyAddr = *lobbyAddr
	AppState.Usernames = strings.Split(*usernames, ",")
	fmt.Println("WAIT FOR", AppState.WaitFor)

	rand.Seed(time.Now().UTC().UnixNano())
	InitKings(AppState)

	go flow.Run(AppState)
	bot.JoinBots(botSlots, AppState, AppLinks, Channels, onMessage)

	fmt.Println("Server started on " + *port + "...")
	cron.ConnectDelay(AppState, AppLinks, Channels)
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/state", stateHandler)
	http.HandleFunc("/channels", channelsHandler)
	http.HandleFunc("/shutdown", shutdownHandler)
	http.HandleFunc("/pause", pauseHandler)
	http.HandleFunc("/leave", leaveHandler)
	http.HandleFunc("/test", testHandler)

	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		panic("Error: " + err.Error())
	}
	fmt.Println("FINISH")
}

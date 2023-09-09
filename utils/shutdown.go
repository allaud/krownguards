package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"ws/game"
	"ws/json"
	"ws/models"
)

func winner(state *models.State) int {
	westDead := state.Kings[0].Action == game.Dead
	eastDead := state.Kings[1].Action == game.Dead
	if westDead && eastDead {
		return -1
	}

	if westDead {
		return 1
	}

	if eastDead {
		return 0
	}

	return -1
}

func apiStats(state *models.State) map[string]interface{} {
	game := map[string]interface{}{
		"servertime": state.Servertime,
		"wave":       state.Wave,
		"port":       state.Port,
		"state":      state.Phase,
		"winner":     winner(state),
	}
	players := map[string]interface{}{}
	for slot, player := range state.Players {
		players[player.Name] = map[string]interface{}{
			"slot":   slot,
			"gold":   player.Gold,
			"stone":  player.Stone,
			"income": player.Income,
			"value":  player.Value,
			"leaked": player.Leaked,
			"score":  player.Score,
		}
	}

	return map[string]interface{}{
		"game":    game,
		"players": players,
	}
}

func EndGame(state *models.State) {
	names := []string{}
	for _, player := range state.Players {
		names = append(names, player.Name)
	}

	jsonStr, _ := json.Marshal(map[string]interface{}{
		"id":     state.Id,
		"port":   state.Port,
		"result": winner(state),
		"users":  append(state.Usernames, state.Leavers...),
	})

	//API callback
	apiJson, _ := json.Marshal(apiStats(state))
	apireq, _ := http.NewRequest("POST", "http://krownguards.com/api/game", bytes.NewBuffer(apiJson))
	apiclient := &http.Client{}
	apiresp, err := apiclient.Do(apireq)
	if err != nil {
		panic(err)
	}
	defer apiresp.Body.Close()

	//Shutdown callback
	fmt.Println("Shutdown callback: ", state.LobbyAddr)
	req, err := http.NewRequest("POST", state.LobbyAddr+"/shutdown", bytes.NewBuffer(jsonStr))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	os.Exit(0)
}

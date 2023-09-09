package utils

import (
	"fmt"
	"sync"
	"ws/datatypes"
	"ws/game"
	"ws/json"
	"ws/models"
	"ws/websocket"
)

var offlineSent = 0

func Recieve(message []byte) datatypes.MessageData {
	var message_json datatypes.MessageData

	if err := json.Unmarshal(message, &message_json); err != nil {
		fmt.Println("error unmarshall!", message)
		return message_json
	}
	fmt.Println(">", message_json)
	return message_json
}

func Send(ws *websocket.Conn, data *datatypes.MessageData, mutex *sync.Mutex, processed bool) {
	if ws == nil {
		return
	}

	if !processed {
		data = preprocess(data, SIMPLE_TRANSFORMATIONS)
	}

	json_data, err := json.Marshal(*data)
	if err != nil {
		fmt.Println("Error marshalling json", err)
		fmt.Println("DATA: ", data, *data)
		return
	}

	//fmt.Println("< ", string(json_data)[:20])
	mutex.Lock()
	defer mutex.Unlock()
	//TODO Unlock here!!!
	err = ws.WriteMessage(websocket.TextMessage, json_data)
	if err != nil {
		fmt.Println("Error writing: ", err)
		return
	}
}

func SendNotification(state *models.State, channels datatypes.ChannelsMap, text string, avoid int) {
	state.Notifications = append(state.Notifications, text)
}

func SendAll(channels datatypes.ChannelsMap, data *datatypes.MessageData, avoid int) {
	state := (*data)["state"].(*models.State)

	//lazy add notifications
	if avoid == 0 && len(state.Notifications) > 0 {
		(*data)["notifications"] = state.Notifications
		state.Notifications = []string{}
	}

	data = preprocess(data, FULL_TRANSFORMATIONS)

	offlineUsers := 0
	for slot, channel := range channels {
		if slot == avoid {
			continue
		}

		if channel.Channel != nil {
			if channel.Channel.Username == "offline" {
				offlineUsers = offlineUsers + 1
			}
		} else {
			offlineUsers = offlineUsers + 1
		}
		if offlineUsers == len(channels) {
			if state.Phase == game.War {
				offlineSent += 1
			} else {
				offlineSent += 10
			}
			if offlineSent > 500 {
				EndGame(state)
			}

		}
		go Send(channel.Channel, data, channel.Mutex, true)
	}
}

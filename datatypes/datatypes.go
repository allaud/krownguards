package datatypes

import (
	"sync"
	"time"
	"ws/models"
	"ws/websocket"
)

type Connection struct {
	Mutex   *sync.Mutex
	Channel *websocket.Conn
	Uid     string
}

type Signal struct {
	Type     string
	Duration time.Duration
	Fn       func()
}

type MessageData map[string]interface{}

type HandlerFunc func(*websocket.Conn, MessageData, *models.State, *models.Links, ChannelsMap)

type HandlerHash map[string]HandlerFunc

type ChannelsMap map[int]*Connection

type Size struct {
	XRange [2]float64
	YRange [2]float64
}

type SlotMap map[int]Size

type WaypointMap map[int]SlotMap

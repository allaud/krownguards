package models

import (
	"fmt"
	"time"
)

func Trace(s string) (string, time.Time) {
	return s, time.Now()
}

func Un(s string, startTime time.Time) {
	endTime := time.Now()
	fmt.Println("[benchmark] #"+s, ": ", endTime.Sub(startTime))
}

type State struct {
	Id            string
	Pause         int
	Servertime    int
	WaveCountdown int //Timer before wave's start
	WaitFor       int //Timer before game start
	Phase         string
	Wave          int
	ArenaWaves    []int //Arena battles start before this wave
	WestAlive     int   //Number of enemies alive on west side
	EastAlive     int   //Number of enemies alive on east side
	Port          string
	LobbyAddr     string
	Usernames     []string         // Usernames of players bound to this game instance
	Players       map[int]*Player  // Players {slot: {player}}
	PlayerUnits   map[int][]*Unit  // Player's units
	Leavers       []string         // Players who left game permanently
	WaveUnits     map[int][]*Unit  // Wave units + summons
	Kings         [2]*Unit         // West King + East King
	KingAttrs     [2]*KingUpgrades // Kings upgrades
	Projectiles   []*Projectile    // Projectiles
	Notifications []string
}

// part of non-sent data
type Links struct {
	ProjLinks []*ProjLink // non-sent projectiles data
}

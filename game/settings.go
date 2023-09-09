package game

var TIMER_DELAY = 1

var CONNECT_DELAY = 60
var PREPARATION_TIME = 80
var SWITCH_WAR_DELAY = 40

var COUNTDOWN = 60

var WAR_STEP_DELAY = 0.1

var PLAYER_PAUSE_CAP = 60

var MAX_CHAT_MESSAGES_COUNT = 5

var END_GAME_SHUTDOWN_TIMEOUT = 5

var ARENA_EVERY_WAR = []int{6, 11, 16, 21, 26, 31}

var GUILD_RESET_CAP = 5
var GUILD_RESET_PRICE = [][2]int{{0, 0}, {100, 100}, {200, 200}, {300, 300}, {400, 400}}
var GUILDS_LIST_CAP = 2

var STONE_INCOME_STEP = 5
var STONE_GRADE_CD = 20
var STONE_GRADE = [][3]int{
	// {gold, stone, stone per tick}
	{40, 10, 1},
	{40, 10, 2},
	{40, 10, 3},
	{40, 10, 4},
	{40, 10, 5},
	{40, 10, 6},
	{40, 10, 7},
	{40, 40, 10},
	{55, 55, 14},
	{70, 70, 18},
	{85, 85, 21},
	{100, 100, 24},
	{115, 115, 28},
	{130, 130, 32},
	{150, 150, 35},
	{150, 150, 38},
	{150, 150, 42},
	{150, 150, 46},
	{150, 150, 49},
	{150, 150, 52},
	{150, 150, 56},
	{150, 150, 60},
}

var START_GOLD = 200
var START_STONE = 60
var START_FOOD = [2]int{0, 15}

var KING_GRADE_PRICE = 80
var KING_GRADE_INCOME = 3

var FOOD_GRADE_STEP = 8
var FOOD_GRADE_CD = 15
var FOOD_GOLD = 24
var FOOD_STONE = 80

var UnitScore float64 = 10

var CrowdControl = []string{"stun"}

var AtkDefMultipliers = map[string]map[string]float64{
	Piercing: {
		Unarmored: 1.0,
		Light:     1.3,
		Medium:    0.9,
		Heavy:     0.8,
		Fortified: 0.7,
	},
	Normal: {
		Unarmored: 1.0,
		Light:     0.9,
		Medium:    1.2,
		Heavy:     0.9,
		Fortified: 0.8,
	},
	Magic: {
		Unarmored: 1.0,
		Light:     1.1,
		Medium:    0.8,
		Heavy:     1.2,
		Fortified: 0.7,
	},
	Siege: {
		Unarmored: 1.0,
		Light:     0.9,
		Medium:    0.9,
		Heavy:     0.9,
		Fortified: 1.25,
	},
	Chaos: {
		Unarmored: 1.0,
		Light:     1.0,
		Medium:    1.0,
		Heavy:     1.0,
		Fortified: 1.0,
	},
}

var BotNames = []string{
	"bot1",
	"bot2",
	"bot3",
	"bot4",
	"bot5",
	"bot6",
	"bot7",
	"bot8",
}

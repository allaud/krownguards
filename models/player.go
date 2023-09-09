package models

type StoneUpgrade struct {
	StoneInc   int    // Stone value per tick
	MaxGrade   int    // Max step of upgrade
	CurrGrade  int    // Current step of upgrade
	InProgress bool   // flag for upgrading progress
	Timer      int    // how much time's passed, sec
	Cooldown   int    // time need for upgrade
	Price      [2]int // [Gold, Stone] cost of upgrade
}

type FoodUpgrade struct {
	InProgress bool // flag for upgrading progress
	Timer      int  // how much time's passed, sec
	Cooldown   int  // time need for upgrade
}

type Player struct {
	Name           string             // Player's name
	Gold           int                // Money
	Stone          int                // Special money
	Food           [2]int             // Food [reserved, all]
	Income         int                // Players bonus gold every wave
	StoneGrade     StoneUpgrade       // Stone upgrade info
	FoodGrade      FoodUpgrade        // Food upgrade info
	Guild          string             // player's curent guild
	GuildsList     []string           // guilds player can pick
	Resets         int                // amount of guild resets
	ResetPrice     [2]int             // current guild reset price
	AvailableUnits []string           // Units for defence
	SummonUnits    map[string]*Summon // Units for attack
	IncomeUnits    []*Unit            // Already summoned units
	Value          int                // Total player's units value
	Leaked         int                // Number of units player leaked
	Score          float64            // Spherical "coolness" of player
	PauseCap       int                // Pause cap for player
}

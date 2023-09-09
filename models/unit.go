package models

type Effect struct {
	Owner    int // Player who applied effect
	Duration int // Effect time counter
}

type Unit struct {
	Id           int                 // Unit unique identifier
	Name         string              // Unit's name
	Projectile   string              // Projectile type
	ProjSpeed    float64             // Projectile speed
	ProjTraj     string              // Projectile trajectory type
	ProjStart    [2]float64          // Projectile start offset [delta range, delta angle]
	Player       int                 // Unit's owner
	Atk          [2]float64          // [MinAttack, MaxAttack]
	AtkRange     float64             // Range unit can attack on
	Aspd         float64             // Attack speed (attack every n tacts)
	AspdTimer    float64             // Aspd tact counter
	AtkType      string              // Type of unit's attack
	Def          float64             // Unit's armor value
	DefType      string              // Type of unit's armor
	Movespeed    float64             // square/sec
	MaxHp        float64             // Maximum health points
	Hp           float64             // Current health points
	HpReg        float64             // Health point per second
	MaxMp        float64             // Maximum mana points
	Mp           float64             // Current mana points
	MpReg        float64             // Mana points regeneration per tact
	Price        int                 // Unit price in gold or stone
	Food         int                 // Amount of food unit reserve
	Bounty       int                 // Bouty for killing this unit
	Size         float64             // Unit collider radius
	Guild        string              // Unit's guild name
	Tier         int                 // Unit's subclass
	Coords       [2]float64          // [x coord, y coord]
	Dir          [2]float64          // [direction x coord, direction y coord]
	TargetCoords [2]float64          // [target x coord, target y coord]
	TargetSign   float64             // Turn direction
	Slot         int                 // Unit's current game slot
	Waypoint     int                 // Number of WP in current slot
	SellPrice    int                 // Sellprice damped part
	Action       string              // Unit's current action
	Killer       int                 // Slot of player who got the bounty
	Upgrades     []string            // Unit's upgrades
	Abilities    map[string][]string // AbilityType: [ability, ability]
	Affected     map[string]*Effect  // Effects on unit
}

type Wave struct {
	Unit     string  // unit wave consists of
	Count    int     // wave size
	Bounty   int     // bounty for finishing wave
	RecValue float64 // value recommended for succesful defense
}

type Summon struct {
	Current  int
	Cap      int
	Timer    int
	Cooldown int
}

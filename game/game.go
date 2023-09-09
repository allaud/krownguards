package game

const (
	// PHASES
	Start    = "start"
	Building = "building"
	War      = "war"
	Arena    = "arena"
	GameOver = "game_over"
	// SIDES
	West = "west"
	East = "east"
	// SPECIAL COEF
	SideCoef  = 5
	KingCoef  = 10
	ArenaCoef = 20
	// SPECIAL TIERS
	Wave   = 0
	Summon = -1
	King   = 99
	// SPECIAL GUILDS
	Draft        = "Draft"
	DynamicDraft = "Dynamic Draft"
	Recruits     = "Recruits"
	Random       = "Random"
	// ABILITY ADD TYPES
	AttackReplacer = "attackReplacer"
	OnHit          = "onHit"
	Active         = "active"
	React          = "react"
	Assist         = "assist"
	Aura           = "aura"
	BuffAura       = "buffAura"
	OnDeath        = "onDeath"
	Defensive      = "defensive"
	Offensive      = "offensive"
	// ABILITY APPLY TYPES
	Instant = "instant"
	Modify  = "modify"
	// ACTIONS
	Idle     = "idle"
	Move     = "move"
	Wait     = "wait"
	Attack   = "attack"
	Dead     = "dead"
	Cast     = "cast"
	Teleport = "teleport"
	// ATTACK TYPES
	Piercing = "piercing"
	Normal   = "normal"
	Magic    = "magic"
	Siege    = "siege"
	Chaos    = "chaos"
	// ARMOR TYPES
	Unarmored = "unarmored"
	Light     = "light"
	Medium    = "medium"
	Heavy     = "heavy"
	Fortified = "fortified"
	// PROJECTILE SPECIAL TYPES
	Melee     = "melee"
	Invisible = "invisible"
)

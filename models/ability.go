package models

type Ability struct {
	AddMethod   string
	ApplyMethod string
	AddLogic    *Logic
	ApplyLogic  *Logic
}

type Logic struct {
	HandlerKw func(map[string]interface{}) func(options *AbOptions)
	Handler   func(*AbOptions)
	Options   map[string]interface{}
}

type AbOptions struct {
	Unit         *Unit
	UnitCurrent  *Unit
	Target       *Unit
	Allies       []*Unit
	Enemies      []*Unit
	State        *State
	StateCurrent *State
	Links        *Links
	Hit          bool
	Limit        int
	Atk          [2]float64
	Dmg          float64
	Projectile   string
	ProjSpeed    float64
	ProjTraj     string
	ProjCoords   [2]float64
}

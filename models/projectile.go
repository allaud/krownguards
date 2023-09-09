package models

type Projectile struct {
	Id         int        // Unique identifier
	Name       string     // Projectile's name
	Movespeed  float64    // square/sec
	Trajectory string     // Trajectory type
	Coords     [2]float64 // [x coord, y coord]
	Dir        [2]float64 // [direction x coord, direction y coord]
	Ability    string     // Name of ability'd started projectile
}

type ProjLink struct {
	Unit      *Unit      // Creator of projectile
	Target    *Unit      // Target of projectile
	Options   *AbOptions // Options
	Triggered bool       // True when projectile triggered
}

package units

import (
	"ws/game"
	"ws/models"
)

var King = models.Unit{
	Name:       "King",
	Size:       1.5,
	Tier:       game.King,
	Atk:        [2]float64{55, 65},
	AtkType:    game.Chaos,
	AtkRange:   7.5,
	Projectile: "king_missile",
	ProjSpeed:  15,
	ProjTraj:   "linear",
	MaxHp:      5000,
	Hp:         5000,
	MaxMp:      0,
	Mp:         0,
	MpReg:      0,
	HpReg:      1,
	Aspd:       float64(0.5 / game.WAR_STEP_DELAY),
	Def:        5,
	DefType:    game.Unarmored,
	Movespeed:  1 * game.WAR_STEP_DELAY,
	Waypoint:   1,
	Abilities: map[string][]string{
		game.Aura:  {"King's blood"},
		game.OnHit: {"King's blade"},
	},
}

var KingAttr = models.KingUpgrades{
	AtkGrade: models.Upgrade{
		Price:     game.KING_GRADE_PRICE,
		CurrGrade: 0,
		MaxGrade:  30,
		GradeStep: 10,
	},
	HpRegGrade: models.Upgrade{
		Price:     game.KING_GRADE_PRICE,
		CurrGrade: 0,
		MaxGrade:  20,
		GradeStep: 2,
	},
	MaxHpGrade: models.Upgrade{
		Price:     game.KING_GRADE_PRICE,
		CurrGrade: 0,
		MaxGrade:  20,
		GradeStep: 450,
	},
}

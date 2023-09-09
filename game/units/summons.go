package units

import (
	"ws/models"
)

var Summons = map[int]map[string]*models.Summon{
	1: {
		"Militia": {
			Current:  3,
			Cap:      3,
			Cooldown: 10,
		},
		"Marksman": {
			Current:  2,
			Cap:      2,
			Cooldown: 15,
		},
		"Swordsman": {
			Current:  2,
			Cap:      2,
			Cooldown: 15,
		},
		"Rogue": {
			Current:  2,
			Cap:      2,
			Cooldown: 20,
		},
		"Lizard": {
			Current:  2,
			Cap:      2,
			Cooldown: 25,
		},
		"Hermit": {
			Current:  2,
			Cap:      2,
			Cooldown: 25,
		},
		"Spearman": {
			Current:  2,
			Cap:      2,
			Cooldown: 30,
		},
		"Commander": {
			Current:  2,
			Cap:      2,
			Cooldown: 30,
		},
		"Furry Bear": {
			Current:  2,
			Cap:      2,
			Cooldown: 40,
		},
		"Flying Spearman": {
			Current:  2,
			Cap:      2,
			Cooldown: 40,
		},
		"Catapult": {
			Current:  2,
			Cap:      2,
			Cooldown: 45,
		},
		"Plate Runner": {
			Current:  2,
			Cap:      2,
			Cooldown: 45,
		},
		"Not-so-furry Wolf": {
			Cooldown: 50,
		},
		"Slow Troll": {
			Cooldown: 50,
		},
		"Warlock": {
			Cooldown: 55,
		},
		"Extremely Furry Bear": {
			Cooldown: 55,
		},
		"Ballista": {
			Cooldown: 60,
		},
		"Spice Merchant": {
			Cooldown: 60,
		},
		"Flesheater": {
			Cooldown: 60,
		},
		"Monstrous Sparrow": {
			Cooldown: 70,
		},
		"Giant War Bear": {
			Cooldown: 90,
		},
		"Something Strange": {
			Cooldown: 90,
		},
		"Moving Mountain": {
			Cooldown: 120,
		},
		"Dimon": {
			Cooldown: 120,
		},
	},
	11: {
		"Not-so-furry Wolf": {
			Current:  2,
			Cap:      2,
			Cooldown: 50,
		},
		"Slow Troll": {
			Current:  2,
			Cap:      2,
			Cooldown: 50,
		},
		"Warlock": {
			Current:  2,
			Cap:      2,
			Cooldown: 55,
		},
		"Extremely Furry Bear": {
			Current:  2,
			Cap:      2,
			Cooldown: 55,
		},
		"Ballista": {
			Current:  2,
			Cap:      2,
			Cooldown: 60,
		},
		"Spice Merchant": {
			Current:  2,
			Cap:      2,
			Cooldown: 60,
		},
		"Flesheater": {
			Current:  2,
			Cap:      2,
			Cooldown: 60,
		},
		"Monstrous Sparrow": {
			Current:  2,
			Cap:      2,
			Cooldown: 70,
		},
	},
	15: {
		"Giant War Bear": {
			Current:  1,
			Cap:      1,
			Cooldown: 90,
		},
		"Something Strange": {
			Current:  1,
			Cap:      1,
			Cooldown: 90,
		},
		"Moving Mountain": {
			Current:  1,
			Cap:      1,
			Cooldown: 120,
		},
		"Dimon": {
			Current:  1,
			Cap:      1,
			Cooldown: 120,
		},
	},
}

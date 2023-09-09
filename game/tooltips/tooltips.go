package tooltips

import (
	//"fmt"
	"math"
	"strconv"
	"strings"
	"ws/abilities"
	"ws/game"
	"ws/utils"
)

var Tc = 1 / game.WAR_STEP_DELAY // Time coef
var Rc = 10.0                    // Range coef

func fToS(number float64) string {
	number = float64(int(100*number/1))/100 + float64(int(math.Mod(1000*number, 10)/5))*0.01
	if math.Mod(number, 1) != 0 {
		return strconv.FormatFloat(number, 'f', 1, 64)
	}
	return strconv.FormatFloat(number, 'f', 0, 64)
}

func ripVal(name, field string) string {
	var value interface{}
	ability, ok := abilities.ABILITIES[name]
	if !ok {
		return ""
	}
	addValue, ok := ability.AddLogic.Options[field]
	if ok {
		value = addValue
	}
	applyValue, ok := ability.ApplyLogic.Options[field]
	if ok {
		value = applyValue
	}
	textVal, ok := value.(string)
	if ok {
		return textVal
	}
	intVal, ok := value.(int)
	if ok {
		if intVal < 0 {
			intVal = -intVal
		}
		if !utils.ContainsString([]string{"duration", "step"}, field) {
			return strconv.Itoa(intVal)
		}
		floatVal := float64(intVal) / Tc
		return fToS(floatVal)
	}
	floatVal, ok := value.(float64)
	if ok {
		floatVal = math.Abs(floatVal)
		if strings.Contains(field, "odifier") {
			field = "modifier"
		}
		if utils.ContainsString([]string{"modifier", "coef", "chance"}, field) {
			floatVal = 100 * floatVal
		}
		if utils.ContainsString([]string{"rad", "abRange"}, field) {
			floatVal = floatVal * Rc
		}
		return fToS(floatVal)
	}
	return ""
}

var GuildTooltips = map[string]string{
	game.Draft:    "The guild will be randomly formed from other guilds' units",
	game.Recruits: "Every time you build unit from some tier you will get random unit of that tier. Sell price of that unit always 50%",
	game.Random:   "You will get random guild",
}

// TypeTooltips = Type : Tooltip
var TypeTooltips = map[string]string{}

// AbTooltips = Ability : [Type | Range | Cost | Description]
var AbTooltips = map[string][4]string{
	"Auto-armored orc driver": {
		"Passive",
		"",
		"",
		ripVal("Auto-armored orc driver", "unitName") + " leaves vehicle after destruction",
	},
	"Bash": {
		"Passive",
		"",
		"",
		"Has " + ripVal("Bash", "chance") + "% chance to apply " +
			ripVal("Bash", "status") + " on enemy for " + ripVal("Bash", "duration") +
			" sec with each attack",
	},
	"Blazing crush": {
		"Passive",
		"",
		"",
		"Has " + ripVal("Blazing crush", "chance") + "% chance to apply " +
			ripVal("Blazing crush", "status") + " on enemy for " +
			ripVal("Blazing crush", "duration") + " sec and deal " +
			ripVal("Blazing crush", "dmg") + " extra damage with each attack",
	},
	"Blessed quiver": {
		"Passive",
		"",
		"",
		"Hits up to " + ripVal("Blessed quiver", "limit") + " enemies with each attack",
	},
	"Booming charges": {
		"Passive",
		"",
		"",
		"Hits " + ripVal("Booming charges", "limit") + " additional enemies with " +
			ripVal("Booming charges", "coef") + "% of initial damage",
	},
	"Bouncing sparks": {
		"Passive",
		"",
		"",
		"Hits " + ripVal("Bouncing sparks", "limit") + " additional enemy with " +
			ripVal("Bouncing sparks", "coef") + "% of initial damage",
	},
	"Breath of life": {
		"Aura",
		ripVal("Breath of life", "rad") + " ft",
		"",
		"Regenerates " + ripVal("Breath of life", "value") + " " +
			ripVal("Breath of life", "source") + " to nearby allies every second",
	},
	"Burning embers": {
		"Offensive",
		ripVal("Burning embers", "rad") + " ft",
		ripVal("Burning embers", "cost") + " mp",
		"Launch " + ripVal("Burning embers", "limit") + " burning embers to nearby enemies dealing " +
			ripVal("Rocket barrage", "dmg") + " damage",
	},
	"Commander aura": {
		"Aura",
		ripVal("Commander aura", "rad") + " ft",
		"",
		"Gives " + ripVal("Commander aura", "modifier") + "% bonus attack to nearby allies",
	},
	"Corruption": {
		"Passive",
		"",
		"",
		"Reduces enemy armor by " + ripVal("Corruption", "value") + " for " +
			ripVal("Corruption", "duration") + " sec with each attack",
	},
	"Demonic might": {
		"Aura",
		ripVal("Demonic might", "rad") + " ft",
		"",
		"Increases attack of nearby allies by " + ripVal("Demonic might", "modifier") + "%",
	},
	"Demonic swiftness": {
		"Aura",
		ripVal("Demonic swiftness", "rad") + " ft",
		"",
		"Increases movement speed of nearby allies by " + ripVal("Demonic swiftness", "modifier") + "%",
	},
	"Divine hymn": {
		"Assist",
		ripVal("Divine hymn", "rad") + " ft",
		ripVal("Divine hymn", "cost") + " mp",
		"Buffs nearby allies with one of the Divine hymns for " + ripVal("Divine hymn", "duration") + " sec",
	},
	"Feast": {
		"Aura",
		ripVal("Feast", "rad") + " ft",
		"",
		"Nearby allies restore " + ripVal("_(OnHit)Feast", "coef") + "% melee damage done as health",
	},
	"Great shield": {
		"Passive",
		"",
		"",
		"Takes " + ripVal("Great shield", "coef") + "% less damage from " +
			ripVal("Great shield", "atkType") + " attacks",
	},
	"Greater heal": {
		"Assist",
		ripVal("Greater heal", "rad") + " ft",
		ripVal("Greater heal", "cost") + " mp",
		"Restore " + ripVal("Greater heal", "value") + " " + ripVal("Greater heal", "source") + " to an ally",
	},
	"Heal": {
		"Assist",
		ripVal("Heal", "rad") + " ft",
		ripVal("Heal", "cost") + " mp",
		"Restore " + ripVal("Heal", "value") + " " + ripVal("Heal", "source") + " to an ally",
	},
	"Heavy plate": {
		"Passive",
		"",
		"",
		"Has " + ripVal("Heavy plate", "chance") + "% chance to block " +
			ripVal("Heavy plate", "value") + " damage when attacked",
	},
	"King's blade": {
		"Passive",
		"",
		"",
		"Deals splash damage based on king's current attack upgrade level to nearby enemies with each attack",
	},
	"King's blood": {
		"Aura",
		ripVal("King's blood", "rad") + " ft",
		"",
		"Deals damage based on king's current regeneration upgrade level to nearby enemies every second",
	},
	"Illusive barrier": {
		"Passive",
		"",
		/*ripVal("Illusive barrier", "cost") + */ "0.25 mp/dmg",
		"Absorb " + ripVal("Illusive barrier", "modifier") + "% of incoming damage in the cost of mana",
	},
	"Illusive impulse": {
		"Passive",
		"",
		"",
		"Add " + ripVal("Illusive impulse", "modifier") + "% of current mana as bonus attack",
	},
	"Mana tides": {
		"Aura",
		ripVal("Mana tides", "rad") + " ft",
		"",
		"Regenerates " + ripVal("Mana tides", "coef") + "% of maximum mana to nearby allies every second",
	},
	"On spices": {
		"Assist",
		ripVal("On spices", "rad") + " ft",
		ripVal("On spices", "cost") + " mp",
		"Increase ally's attack speed by " + ripVal("On spices", "aspdModifier") +
			"% and movement speed by " + ripVal("On spices", "mspdModifier") + "% for " +
			ripVal("On spices", "duration") + " sec",
	},
	"Orc driver": {
		"Passive",
		"",
		"",
		ripVal("Orc driver", "unitName") + " leaves vehicle after destruction",
	},
	"Overclocking": {
		"Selfbuff",
		"",
		ripVal("Overclocking", "cost") + " mp",
		"Increases unit's attack speed by " + ripVal("Overclocking", "aspdModifier") +
			"% and movement speed by " + ripVal("Overclocking", "mspdModifier") + "% for " +
			ripVal("Overclocking", "duration") + " sec",
	},
	"Piercing howl": {
		"Aura",
		ripVal("Piercing howl", "rad") + " ft",
		"",
		"Reduces nearby enemies attack speed by " + ripVal("Piercing howl", "aspdModifier") +
			"% and movement speed by " + ripVal("Piercing howl", "mspdModifier") + "%",
	},
	"Rocket barrage": {
		"Offensive",
		ripVal("Rocket barrage", "abRange") + " ft (" + ripVal("Rocket barrage", "rad") + " ft AoE)",
		ripVal("Rocket barrage", "cost") + " mp",
		"Launch missle to the enemy that explodes dealing " + ripVal("Rocket barrage", "dmg") +
			" damage and apply stun for " + ripVal("Rocket barrage", "duration") + " second to nearby enemies",
	},
	"Scope": {
		"Aura",
		ripVal("Scope", "rad") + " ft",
		"",
		"Increases attack of nearby ranged allies by " + ripVal("Scope", "modifier") + "%",
	},
	"Scope 2.0": {
		"Aura",
		ripVal("Scope 2.0", "rad") + " ft",
		"",
		"Increases attack of nearby ranged allies by " + ripVal("Scope 2.0", "modifier") + "%",
	},
	"Slowing arrows": {
		"Passive",
		"",
		"",
		"Reduces enemy attack speed by " + ripVal("Slowing arrows", "aspdModifier") +
			"% and movement speed by " + ripVal("Slowing arrows", "mspdModifier") + "% for " +
			ripVal("Slowing arrows", "duration") + " sec with each attack",
	},
	"Solar spike": {
		"Passive",
		ripVal("Solar spike", "rad") + " ft",
		"",
		"Hits up to " + ripVal("Solar spike", "limit") + " additional enemies for " +
			ripVal("Solar spike", "dmg") + " damage with each attack",
	},
	"Spell blast": {
		"Offensive",
		"",
		ripVal("Spell blast", "cost") + " mp",
		"Adds " + ripVal("Spell blast", "dmg") + "magic damage to the next attack",
	},
	"Splinter shot": {
		"Passive",
		ripVal("Splinter shot", "rad") + " ft",
		"",
		"Hits up to " + ripVal("Splinter shot", "limit") + " additional enemies for " +
			ripVal("Splinter shot", "coef") + "% of damage with each attack",
	},
	"Steel fur": {
		"Aura",
		ripVal("Steel fur", "rad") + " ft",
		"",
		"Increases armor of nearby allies by " + ripVal("Steel fur", "value"),
	},
	"Strange poison": {
		"Passive",
		"",
		"",
		"Attacks apply poison that reduces enemy's attack speed by " + ripVal("Strange poison", "aspdModifier") +
			"% and movement speed by " + ripVal("Strange poison", "mspdModifier") +
			"% and deals " + ripVal("Strange poison", "dmg") + " damage every second for " +
			ripVal("Strange poison", "duration") + " sec",
	},
	"Strange regeneration": {
		"Aura",
		ripVal("Strange regeneration", "rad") + " ft",
		"",
		"Regenerates " + ripVal("Strange regeneration", "value") +
			" " + ripVal("Strange regeneration", "source") + " to nearby allies every second",
	},
	"Tectonic flame": {
		"Aura",
		ripVal("Tectonic flame", "rad") + " ft",
		"",
		"Deals " + ripVal("Tectonic flame", "value") + " damage to nearby enemies every second",
	},
	"Tyrannicide": {
		"Passive",
		"",
		"",
		"Deals " + ripVal("Tyrannicide", "modifier") + "% damage to the King",
	},
	"Welding arc": {
		"Passive",
		"",
		"",
		"Reflects " + ripVal("Welding arc", "coef") + "% melee damage back to the attacker",
	},
}

// AfTooltips = Effect : [Type | Description]
var AfTooltips = map[string][2]string{
	"stun": {
		"Debuff",
		"Stunned. Can't do anything",
	},
	"Breath of life": {
		"Buff",
		"Regenerates " + ripVal("Breath of life", "value") + " " +
			ripVal("Breath of life", "source") + " every second",
	},
	"Commander aura": {
		"Buff",
		"Attack increased by " + ripVal("Commander aura", "modifier") + "%",
	},
	"Corruption": {
		"Debuff",
		"Armor reduced by " + ripVal("Corruption", "value"),
	},
	"Demonic might": {
		"Buff",
		"Attack increased by " + ripVal("Demonic might", "modifier") + "%",
	},
	"Demonic swiftness": {
		"Buff",
		"Movement speed increased by " + ripVal("Demonic swiftness", "modifier") + "%",
	},
	"Feast": {
		"Buff",
		"Restores " + ripVal("_(OnHit)Feast", "coef") + "% melee damage done as health",
	},
	"Final charge": {
		"Buff",
		"This unit broke through first line of defense! Armor increased by " + ripVal("Steel fur", "value"),
	},
	"Hymn of blaze": {
		"Buff",
		"Attack increased by " + ripVal("Hymn of blaze", "atkModifier") +
			"% and attack speed by " + ripVal("Hymn of blaze", "aspdModifier") + "%",
	},
	"Hymn of magnificence": {
		"Buff",
		"Max health increased by " + ripVal("Hymn of magnificence", "maxHpModifier") +
			"% and mana regeneration by " + ripVal("Hymn of magnificence", "mpRegModifier") + "%",
	},
	"Hymn of pleading": {
		"Buff",
		"Armor increased by " + ripVal("Hymn of pleading", "defValue") +
			" and movement speed by " + ripVal("Hymn of pleading", "mspdModifier") + "%",
	},
	"Illusive barrier": {
		"",
		"Absorb " + ripVal("Illusive barrier", "modifier") + "% of incoming damage in the cost of mana",
	},
	"Mana tides": {
		"Buff",
		"Regenerates " + ripVal("Mana tides", "coef") + "% of maximum mana every second",
	},
	"On spices": {
		"Buff",
		"Attack speed increased by " + ripVal("On spices", "aspdModifier") +
			"% and movement speed by " + ripVal("On spices", "mspdModifier") + "%",
	},
	"Overclocking": {
		"Buff",
		"Attack speed increased by " + ripVal("Overclocking", "aspdModifier") +
			"% and movement speed by " + ripVal("Overclocking", "mspdModifier") + "%",
	},
	"Piercing howl": {
		"Debuff",
		"Attack speed decreased by " + ripVal("Piercing howl", "aspdModifier") +
			"% and movement speed by " + ripVal("Piercing howl", "mspdModifier") + "%",
	},
	"Rocket barrage": {
		"Debuff",
		"You are target of a rocket! Boom!",
	},
	"Scope": {
		"Buff",
		"Attack increased by " + ripVal("Scope", "modifier") + "%",
	},
	"Scope 2.0": {
		"Buff",
		"Attack increased by " + ripVal("Scope 2.0", "modifier") + "%",
	},
	"Slowing arrows": {
		"Debuff",
		"Attack speed decreased by " + ripVal("Slowing arrows", "aspdModifier") +
			"% and movement speed by " + ripVal("Slowing arrows", "mspdModifier") + "%",
	},
	"Steel fur": {
		"Buff",
		"Armor increased by " + ripVal("Steel fur", "value"),
	},
	"Strange poison": {
		"Debuff",
		"Attack speed decreased by " + ripVal("Strange poison", "aspdModifier") +
			"% and movement speed by " + ripVal("Strange poison", "mspdModifier") +
			"%. Takes " + ripVal("Strange poison", "dmg") + " damage every second",
	},
	"Strange regeneration": {
		"Buff",
		"Regenerate " + ripVal("Strange regeneration", "value") +
			" " + ripVal("Strange regeneration", "source") + " every second",
	},
	"Tectonic flame": {
		"Debuff",
		"Takes " + ripVal("Tectonic flame", "value") + " damage every second",
	},
}

func init() {
	// Generate tooltips for all armor and attack types
	for atkType, defTypes := range game.AtkDefMultipliers {
		if atkType == game.Chaos {
			TypeTooltips[atkType] = "Deals 100% damage to any type of armor"
			continue
		}
		for defType, coef := range defTypes {
			if defType == game.Unarmored {
				TypeTooltips[defType] = "Takes 100% damage from any type of attack"
				continue
			}
			if _, exists := TypeTooltips[atkType]; !exists {
				TypeTooltips[atkType] = "Deals " + fToS(coef*100) + "% damage to " + defType + " armor"
			} else {
				TypeTooltips[atkType] += "\nDeals " + fToS(coef*100) + "% damage to " + defType + " armor"
			}
			if _, exists := TypeTooltips[defType]; !exists {
				TypeTooltips[defType] = "Takes " + fToS(coef*100) + "% damage from " + atkType + " attacks"
			} else {
				TypeTooltips[defType] += "\nTakes " + fToS(coef*100) + "% damage from " + atkType + " attacks"
			}
		}
	}

	/*for k, v := range AbTooltips {
		fmt.Println(k, v)
	}
	for k, v := range AfTooltips {
		fmt.Println(k, v)
	}
	for k, v := range TypeTooltips {
		fmt.Println(k, v)
		fmt.Println()
	}*/
}

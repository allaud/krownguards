package abilities

import (
	"ws/models"
)

func override(s1, s2 map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, _ := range s1 {
		res[k] = s1[k]
	}
	for k, _ := range s2 {
		res[k] = s2[k]
	}
	return res
}

func AddByType(actor *models.Unit, abilities map[string]*models.Ability, options *models.AbOptions, abType string) {
	//unit := options[actor].(*models.Unit)
	for _, abName := range actor.Abilities[abType] {
		ability := abilities[abName].AddLogic.Handler
		ability(options)
	}
}

func ApplyEffects(abilities map[string]*models.Ability, options *models.AbOptions) {
	//unit := options["unit"].(*models.Unit)
	for eName, effect := range options.Unit.Affected {
		if effect.Duration == 0 {
			delete(options.Unit.Affected, eName)
			continue
		}
		ability := abilities[eName].ApplyLogic.Handler
		ability(options)
		options.Unit.Affected[eName].Duration -= 1
	}
}

func ApplyAbility(abName string, options *models.AbOptions) {
	ability := ABILITIES[abName].ApplyLogic.Handler
	ability(options)
}

func RefreshStatus(target, unit *models.Unit, status string, duration int) {
	if _, ok := target.Affected[status]; !ok {
		target.Affected[status] = &models.Effect{
			Owner:    unit.Player,
			Duration: duration,
		}
		return
	}
	if duration > target.Affected[status].Duration {
		target.Affected[status].Duration = duration
	}
	target.Affected[status].Owner = unit.Player
}

var Empty = func(kwargs map[string]interface{}) func(options *models.AbOptions) {
	return func(options *models.AbOptions) {
	}
}

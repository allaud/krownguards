package utils

import (
	"math/rand"
	"ws/models"
)

func MapKeys(dict map[string]models.Player) []string {
	keys := []string{}
	for k := range dict {
		keys = append(keys, k)
	}
	return keys
}

func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GetIndex(s []string, e string) int {
	for k, v := range s {
		if v == e {
			return k
		}
	}
	return -1
}

func Uniform(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func CopyAbilities(original map[string][]string) map[string][]string {
	temp := map[string][]string{}
	for k, v := range original {
		temp[k] = v
	}
	return temp
}

func DampPrice(unit *models.Unit) {
	unit.SellPrice += unit.Price / 2
	unit.Price = 0
}

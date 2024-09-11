package main

import (
	"fmt"
	"mhn/armors"
	"slices"
	"strings"
)

func main() {
	armorSetList := armors.ReadArmorCollection()

	exclude := []string{"Diablos", "Kushala", "Kaiser", "Rath Soul"}
	filteredArmorSetList := excludeSets(exclude, armorSetList)

	foundSets := armors.FindArmorSets(filteredArmorSetList, []armors.Skill{
		{Name: "Offensive Guard", Level: 2},
		{Name: "Guard", Level: 1},
		{Name: "Artillery", Level: 3},
		{Name: "Focus", Level: 2},
	})

	printSets(foundSets)
	fmt.Println("Found", len(foundSets))
}

func excludeSets(sets []string, armorSetList [][]armors.ArmorPiece) [][]armors.ArmorPiece {
	filteredArmorSetList := [][]armors.ArmorPiece{}
	for _, set := range armorSetList {
		filteredArmorSetList = append(filteredArmorSetList, slices.DeleteFunc(set, func(x armors.ArmorPiece) bool {
			for _, ex := range sets {
				if strings.Contains(x.Name, ex) {
					return true
				}
			}
			return false
		}))
	}

	return filteredArmorSetList
}

func printSets(sets []armors.GradedArmorSet) {
	for _, set := range sets {
		for _, piece := range []armors.GradedArmorPiece{set.Head, set.Chest, set.Arms, set.Waist, set.Legs} {
			fmt.Println(piece.Name, piece.Grade, piece.EnabledSkills)
		}
		fmt.Println()
	}
}

package main

import (
	"fmt"
	"mhn/armors"
	"time"
)

func main() {
	armorSetList := armors.ReadArmorCollection()
	start := time.Now()

	foundSets := armors.FindArmorSets(armorSetList, []armors.Skill{
		{Name: "Offensive Guard", Level: 1},
		{Name: "Guard", Level: 1},
		{Name: "Attack Boost", Level: 1},
	})

	elapsed := time.Since(start)
	printSets(foundSets)
	fmt.Println("Found", len(foundSets), elapsed)
}

func printSets(sets []armors.GradedArmorSet) {
	for _, set := range sets {
		for _, piece := range []armors.GradedArmorPiece{set.Head, set.Chest, set.Arms, set.Waist, set.Legs} {
			fmt.Println(piece.Name, piece.Grade, piece.EnabledSkills)
		}
		fmt.Println()
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type ArmorSets struct {
	ArmorSets []ArmorSet `json:"armor_sets"`
}

type ArmorSet struct {
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Skills []Skill `json:"skills"`
}

type Skill struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
	Grade int    `json:"grade"`
}

var armor_collection ArmorSets

func main() {
	readArmorCollection()
	findArmorSets([]Skill{
		{
			Name:  "Focus",
			Level: 1,
		},
		{
			Name:  "Attack Boost",
			Level: 1,
		},
	})
}

func readArmorCollection() {
	bytes, _ := os.ReadFile("assets/armor_collection.json")
	json.Unmarshal(bytes, &armor_collection)
}

func findArmorSets(skills []Skill) {
	for _, skill := range skills {
		fmt.Println(skill)
	}
}

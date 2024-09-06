package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type SkillName string

const (
	Focus       SkillName = "Focus"
	AttackBoost SkillName = "Attack Boost"
)

type ArmorType string

const (
	HEAD  ArmorType = "Head"
	CHEST ArmorType = "Chest"
	ARMS  ArmorType = "Arms"
	WAIST ArmorType = "Waist"
	LEGS  ArmorType = "Legs"
)

type ArmorSets struct {
	ArmorSets []ArmorPiece `json:"armor_sets"`
}

type ArmorPiece struct {
	Name   string    `json:"name"`
	Type   ArmorType `json:"type"`
	Skills []Skill   `json:"skills"`
}

type Skill struct {
	Name  SkillName `json:"name"`
	Level int       `json:"level"`
	Grade int       `json:"grade"`
}

type ArmorSet struct {
	Head  ArmorPiece
	Chest ArmorPiece
	Arms  ArmorPiece
	Waist ArmorPiece
	Legs  ArmorPiece
}

var armorCollection ArmorSets

// Armor collection organized into their types
var headSets = []ArmorPiece{}
var chestSets = []ArmorPiece{}
var armSets = []ArmorPiece{}
var waistSets = []ArmorPiece{}
var legSets = []ArmorPiece{}

func main() {
	readArmorCollection()
}

func readArmorCollection() {
	bytes, error := os.ReadFile("assets/armor_collection.json")
	json.Unmarshal(bytes, &armorCollection)

	if error != nil {
		fmt.Println("Error reading in armor_collection.json")
		os.Exit(1)
	}

	// Organize armor pieces by type
	for _, armor := range armorCollection.ArmorSets {
		switch armor.Type {
		case HEAD:
			headSets = append(headSets, armor)
		case CHEST:
			chestSets = append(chestSets, armor)
		case ARMS:
			armSets = append(armSets, armor)
		case WAIST:
			waistSets = append(waistSets, armor)
		case LEGS:
			legSets = append(legSets, armor)
		}
	}
}

func search() {

	/*
		Check if the current set contains all the skills requested
		If it does contain all the skills, add it to the list of valid sets and return
		If it does not contain all the skills, return

		Make a valid choice based on which armor type we're currently looking for
		Recursively search with the current set
		Undo the previous choice

		Once we're exhausted all options, return and complete the search
	*/

	var currSet = ArmorSet{}
	var armorPieces = []*ArmorPiece{
		&currSet.Head, &currSet.Chest, &currSet.Arms, &currSet.Waist, &currSet.Legs,
	}
	*armorPieces[0] = ArmorPiece{}
}

func findArmorSets(skills []Skill) []ArmorSet {
	var validSets = []ArmorSet{}

	search()

	return validSets
}

func isValidSet(armorSet ArmorSet, skills []Skill) bool {
	return true
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
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

type GradedArmorPiece struct {
	OriginalPiece ArmorPiece
	Grade         int
	EnabledSkills []Skill
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

type GradedArmorSet struct {
	Head  GradedArmorPiece
	Chest GradedArmorPiece
	Arms  GradedArmorPiece
	Waist GradedArmorPiece
	Legs  GradedArmorPiece
}

var armorSetList [][]ArmorPiece

func main() {
	readArmorCollection()
	var foundSets = findArmorSets([]Skill{
		{Name: "Attack Boost", Level: 6},
		// {Name: "Focus", Level: 1},
		// {Name: "Offensive Guard", Level: 1},
		// {Name: "Burst", Level: 1},
		// {Name: "Quick Work", Level: 1},
	})

	fmt.Println(len(foundSets))
	for _, set := range foundSets {
		for _, piece := range []GradedArmorPiece{set.Head, set.Chest, set.Arms, set.Waist, set.Legs} {
			fmt.Println(piece.OriginalPiece.Name, piece.Grade, piece.EnabledSkills)
		}
		fmt.Println("===================================================")
	}
}

func modify(numbers []int) {
	numbers = append(numbers, 1)
}

func readArmorCollection() {
	bytes, error := os.ReadFile("assets/armor_collection.json")
	var armorSets ArmorSets
	json.Unmarshal(bytes, &armorSets)

	if error != nil {
		fmt.Println("Error reading in armor_collection.json")
		os.Exit(1)
	}
	var headSets = []ArmorPiece{}
	var chestSets = []ArmorPiece{}
	var armSets = []ArmorPiece{}
	var waistSets = []ArmorPiece{}
	var legSets = []ArmorPiece{}

	// Organize armor pieces by type
	for _, armor := range armorSets.ArmorSets {
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

	armorSetList = [][]ArmorPiece{
		headSets, chestSets, armSets, waistSets, legSets,
	}
}

func search(skillReqs []Skill, validSets *[]GradedArmorSet, currPieces []*GradedArmorPiece, currPieceIdx int, armorTypeIdx int) {
	// If we're searching at the 5th index, we're done for now and can check if currPieces is valid to append
	if currPieceIdx == 5 {
		// Check if currPieces satifies the skill requirements
		var set = GradedArmorSet{
			Head: *currPieces[0], Chest: *currPieces[1], Arms: *currPieces[2], Waist: *currPieces[3], Legs: *currPieces[4],
		}
		if isValidSet(set, skillReqs) {
			*validSets = append(*validSets, set)
		}
		return
	}

	// Preprocessing skill requirements for faster lookup in isValidPiece
	var skillReqMap = make(map[SkillName]bool)
	for _, skill := range skillReqs {
		skillReqMap[skill.Name] = true
	}

	var foundValidPiece = false
	for _, potentialPiece := range armorSetList[armorTypeIdx] {
		for _, gradedArmorPiece := range expandArmorPieceByGrade(potentialPiece) {
			if isValidPiece(gradedArmorPiece, skillReqMap) {
				foundValidPiece = true
				// Choose potential piece and keep track of the previous piece for later
				var previousPiece = *currPieces[currPieceIdx]
				*currPieces[currPieceIdx] = gradedArmorPiece
				// Recursively search for the next armor type
				search(skillReqs, validSets, currPieces, currPieceIdx+1, armorTypeIdx+1)
				// Undo last choice
				*currPieces[currPieceIdx] = previousPiece
			}
		}
	}
	// If no valid piece is found, continue with the next armor type
	if !foundValidPiece {
		search(skillReqs, validSets, currPieces, currPieceIdx+1, armorTypeIdx+1)
	}
}

func findArmorSets(skills []Skill) []GradedArmorSet {
	var validSets = []GradedArmorSet{}
	var initialSet = GradedArmorSet{}
	var armorPieces = []*GradedArmorPiece{
		&initialSet.Head, &initialSet.Chest, &initialSet.Arms, &initialSet.Waist, &initialSet.Legs,
	}

	search(skills, &validSets, armorPieces, 0, 0)

	return validSets
}

func isValidSet(armorSet GradedArmorSet, requiredSkills []Skill) bool {
	var summedSkillReqs = make(map[SkillName]int)
	// Summing up the levels of all required skills
	for _, skill := range requiredSkills {
		summedSkillReqs[skill.Name] += skill.Level
	}

	// Do the same for the armorSet
	var summedCurrSkills = make(map[SkillName]int)
	for _, piece := range []GradedArmorPiece{armorSet.Head, armorSet.Chest, armorSet.Arms, armorSet.Waist, armorSet.Legs} {
		for _, skill := range piece.EnabledSkills {
			summedCurrSkills[skill.Name] += skill.Level
		}
	}

	// Now validate requested skills against target armor set
	for skillName, reqLevel := range summedSkillReqs {
		level, exists := summedCurrSkills[skillName]
		if !exists {
			return false
		}
		if level < reqLevel {
			return false
		}
	}

	return true
}

func isValidPiece(gradedArmorPiece GradedArmorPiece, skillReqs map[SkillName]bool) bool {
	// An armor piece is valid if it has a skill that exists in skillReqs
	// This can just be checked by looking for a matching Skill.name
	for _, skill := range gradedArmorPiece.EnabledSkills {
		if _, exists := skillReqs[skill.Name]; exists {
			return true
		}
	}

	return false
}

func expandArmorPieceByGrade(armorPiece ArmorPiece) []GradedArmorPiece {
	var gradedArmorPieces = []GradedArmorPiece{}

	// Determine what gradeSkillMap are needed for skills on this piece
	var gradeSkillMap = make(map[int]map[SkillName]Skill)
	for _, piece := range armorPiece.Skills {
		gradeSkillMap[piece.Grade] = make(map[SkillName]Skill)
	}
	// Then determine what skills are enabled for each grade
	for grade := range gradeSkillMap {
		for _, skill := range armorPiece.Skills {
			gradeSkillMap[grade] = append(gradeSkillMap[grade], skill)
		}
		gradedArmorPieces = append(gradedArmorPieces, GradedArmorPiece{
			OriginalPiece: armorPiece,
			Grade:         grade,
			EnabledSkills: gradeSkillMap[grade],
		})
	}

	return gradedArmorPieces
}

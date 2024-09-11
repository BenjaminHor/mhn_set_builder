package armors

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
)

type SkillName string

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
	ArmorPiece
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

func ReadArmorCollection() [][]ArmorPiece {
	bytes, error := os.ReadFile("assets/armor_collection.json")
	var armorSets ArmorSets
	json.Unmarshal(bytes, &armorSets)

	if error != nil {
		fmt.Println("Error reading in armor_collection.json")
		os.Exit(1)
	}
	headSets := []ArmorPiece{}
	chestSets := []ArmorPiece{}
	armSets := []ArmorPiece{}
	waistSets := []ArmorPiece{}
	legSets := []ArmorPiece{}

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

	return [][]ArmorPiece{
		headSets, chestSets, armSets, waistSets, legSets,
	}
}

func FindArmorSets(armorSetList [][]ArmorPiece, skills []Skill) []GradedArmorSet {
	validSets := []GradedArmorSet{}
	initialSet := GradedArmorSet{}
	armorPieces := []*GradedArmorPiece{
		&initialSet.Head, &initialSet.Chest, &initialSet.Arms, &initialSet.Waist, &initialSet.Legs,
	}

	search(armorSetList, skills, &validSets, armorPieces, 0)

	return validSets
}

func search(armorSetList [][]ArmorPiece, skillReqs []Skill, validSets *[]GradedArmorSet, currPieces []*GradedArmorPiece, currPieceIdx int) {
	// Check if currPieces satisfies the skill requirements
	currGradedSet := GradedArmorSet{
		Head: *currPieces[0], Chest: *currPieces[1], Arms: *currPieces[2], Waist: *currPieces[3], Legs: *currPieces[4],
	}
	if isValidSet(currGradedSet, skillReqs) {
		*validSets = append(*validSets, currGradedSet)
		return
	}

	// We've reached the end without finding a valid set, return
	if currPieceIdx >= 5 {
		return
	}

	// Preprocessing skill requirements for faster lookup in isValidPiece
	skillReqMap := make(map[SkillName]bool)
	for _, skill := range skillReqs {
		skillReqMap[skill.Name] = true
	}

	for _, potentialPiece := range armorSetList[currPieceIdx] {
		for _, gradedArmorPiece := range expandArmorPieceByGrade(potentialPiece) {
			if isValidPiece(gradedArmorPiece, skillReqMap) {
				// Choose potential piece and keep track of the previous piece for later
				previousPiece := *currPieces[currPieceIdx]
				*currPieces[currPieceIdx] = gradedArmorPiece
				// Recursively search for the next armor type
				search(armorSetList, skillReqs, validSets, currPieces, currPieceIdx+1)
				// Undo last choice
				*currPieces[currPieceIdx] = previousPiece
			}
		}
	}
	// Continue searching in next armor slot
	search(armorSetList, skillReqs, validSets, currPieces, currPieceIdx+1)
}

func isValidSet(armorSet GradedArmorSet, requiredSkills []Skill) bool {
	summedSkillReqs := make(map[SkillName]int)
	// Summing up the levels of all required skills
	for _, skill := range requiredSkills {
		summedSkillReqs[skill.Name] += skill.Level
	}

	// Do the same for the armorSet
	summedCurrSkills := make(map[SkillName]int)
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
	gradedArmorPieces := []GradedArmorPiece{}

	// Determine what grades are needed for skills on this piece
	gradeSkillMap := make(map[int]bool)
	for _, piece := range armorPiece.Skills {
		gradeSkillMap[piece.Grade] = true
	}
	// Then determine what skills are enabled for each grade
	for currGrade := range gradeSkillMap {
		enabledSkills := make(map[SkillName]Skill)
		for _, skill := range armorPiece.Skills {
			if skill.Grade <= currGrade {
				if _, exists := enabledSkills[skill.Name]; !exists {
					enabledSkills[skill.Name] = skill
					continue
				}

				if skill.Level >= enabledSkills[skill.Name].Level {
					enabledSkills[skill.Name] = skill
				}
			}
		}

		// Create gradedArmorPiece
		enabledSkillsSlice := []Skill{}
		for skill := range maps.Values(enabledSkills) {
			enabledSkillsSlice = append(enabledSkillsSlice, skill)
		}
		gradedArmorPiece := GradedArmorPiece{
			ArmorPiece:    armorPiece,
			Grade:         currGrade,
			EnabledSkills: enabledSkillsSlice,
		}
		gradedArmorPieces = append(gradedArmorPieces, gradedArmorPiece)
	}

	return gradedArmorPieces
}

package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/JohnG-Dev/army_builder_api/internal/models"
	"gopkg.in/yaml.v3"
)

var StatMapping = map[string]string{
	"Move": "move", "M": "move",
	"Health": "health", "W": "health",
	"Save": "save", "Sv": "save",
	"Control": "control", "OC": "control",
	"Toughness": "toughness", "T": "toughness",
	"Leadership": "leadership", "Ld": "leadership",
	"Bravery": "leadership", "Ward": "ward",
}

type Converter struct{}

func main() {
	rawDir := "./data/raw"
	outDir := "./data/factions"
	files, err := os.ReadDir(rawDir)
	if err != nil {
		log.Fatalf("failed to read: %v", err)
	}
	cv := &Converter{}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".cat") {
			continue
		}
		filePath := filepath.Join(rawDir, file.Name())
		xmlData, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", file.Name(), err)
			continue
		}
		var catalogue Catalogue
		err = xml.Unmarshal(xmlData, &catalogue)
		if err != nil {
			fmt.Printf("Unmarshal error in %s: %v\n", file.Name(), err)
			continue
		}

		fmt.Printf("File: %s\n", catalogue.Name)
		fmt.Printf(" - Selection: %d, Links: %d, Shared: %d\n",
			len(catalogue.SelectionEntries),
			len(catalogue.EntryLinks),
			len(catalogue.SharedEntries))

		var allUnits []SelectionEntry
		allUnits = cv.findUnits(catalogue.SelectionEntries, allUnits)
		allUnits = cv.findUnits(catalogue.EntryLinks, allUnits)
		allUnits = cv.findUnits(catalogue.SharedEntries, allUnits)
		allUnits = cv.findUnits(catalogue.SharedGroups, allUnits)

		seed := models.SeedData{
			GameName: "Age of Sigmar",
			Factions: []models.FactionSeed{{
				Name:   catalogue.Name,
				Source: "Battlescribe Data",
				Units:  []models.UnitSeed{},
			}},
		}

		for _, entry := range allUnits {
			unit := cv.transformUnit(entry)
			seed.Factions[0].Units = append(seed.Factions[0].Units, unit)
		}

		outPath := filepath.Join(outDir, strings.ToLower(strings.ReplaceAll(catalogue.Name, " ", "_"))+".yaml")
		yamlData, err := yaml.Marshal(seed)
		if err != nil {
			fmt.Printf("Error marshaling YAML for %s: %v\n", catalogue.Name, err)
			continue
		}
		err = os.WriteFile(outPath, yamlData, 0o644)
		if err != nil {
			fmt.Printf("Failed to write YAML for %s: %v\n", catalogue.Name, err)
			continue
		}
		fmt.Printf("Converted %s (%d units)\n", catalogue.Name, len(seed.Factions[0].Units))
	}
}

func (c *Converter) transformUnit(entry SelectionEntry) models.UnitSeed {
	unit := models.UnitSeed{
		Name:            entry.Name,
		AdditionalStats: make(map[string]string),
	}

	allProfiles := entry.Profiles
	allProfiles = append(allProfiles, c.collectProfiles(entry)...)

	for _, p := range allProfiles {
		if strings.Contains(p.TypeName, "Unit") {
			c.mapStats(p.Characteristics, &unit)
		} else if strings.Contains(p.TypeName, "Weapon") {
			unit.Weapons = append(unit.Weapons, c.mapWeapon(p))
		}
	}

	for _, cat := range entry.Categories {
		unit.Keywords = append(unit.Keywords, strings.ToUpper(cat.Name))
	}
	return unit
}

func (c *Converter) mapStats(chars []Characteristic, unit *models.UnitSeed) {
	for _, char := range chars {
		if target, ok := StatMapping[char.Name]; ok {
			switch target {
			case "move":
				unit.Move = char.Value
			case "health":
				unit.Health = char.Value
			case "save":
				unit.Save = char.Value
			case "ward":
				unit.Ward = char.Value
			case "control":
				unit.Control = char.Value
			case "toughness":
				unit.Toughness = char.Value
			case "leadership":
				unit.Leadership = char.Value
			}
		} else {
			unit.AdditionalStats[char.Name] = char.Value
		}
	}
}

func (c *Converter) mapWeapon(p Profile) models.WeaponSeed {
	w := models.WeaponSeed{Name: p.Name}
	for _, char := range p.Characteristics {
		switch char.Name {
		case "Range":
			w.Range = char.Value
		case "Attacks", "A":
			w.Attacks = char.Value
		case "To Hit", "WS":
			w.ToHit = char.Value
		case "To Wound", "S":
			w.ToWound = char.Value
		case "Rend", "AP":
			w.Rend = char.Value
		case "Damage", "D":
			w.Damage = char.Value
		}
	}
	return w
}

func (c *Converter) isUnit(entry SelectionEntry) bool {
	for _, p := range entry.Profiles {
		if strings.Contains(p.TypeName, "Unit") || strings.Contains(p.TypeName, "Models") {
			return true
		}
	}

	if len(entry.Categories) > 1 {
		return true
	}

	return false
}

func (c *Converter) findUnits(entries []SelectionEntry, found []SelectionEntry) []SelectionEntry {
	for _, entry := range entries {
		if c.isUnit(entry) {
			found = append(found, entry)
		}

		found = c.findUnits(entry.ChildEntries, found)
		found = c.findUnits(entry.LinkEntries, found)
		found = c.findUnits(entry.GroupEntries, found)
	}
	return found
}

func (c *Converter) collectProfiles(entry SelectionEntry) []Profile {
	var found []Profile

	containers := [][]SelectionEntry{entry.ChildEntries, entry.LinkEntries, entry.GroupEntries}
	for _, container := range containers {
		for _, child := range container {
			found = append(found, child.Profiles...)
			found = append(found, c.collectProfiles(child)...)
		}
	}
	return found
}

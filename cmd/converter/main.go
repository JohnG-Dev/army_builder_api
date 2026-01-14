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

type Converter struct {
	MasterEntries map[string]SelectionEntry
}

func main() {
	rawDir := "./data/raw"
	outDir := "./data/factions"

	allFiles, err := os.ReadDir(rawDir)
	if err != nil {
		log.Fatalf("failed to read directory: %v", err)
	}

	cv := &Converter{
		MasterEntries: make(map[string]SelectionEntry),
	}

	for _, file := range allFiles {
		name := file.Name()
		if !strings.HasSuffix(name, ".cat") && !strings.HasSuffix(name, ".gst") {
			continue
		}

		filePath := filepath.Join(rawDir, name)
		xmlData, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading %s for indexing: %v\n", name, err)
			continue
		}

		var catalogue Catalogue
		err = xml.Unmarshal(xmlData, &catalogue)
		if err != nil {
			fmt.Printf("Unmarshal error in %s during indexing: %v\n", name, err)
			continue
		}

		cv.indexCatalogue(catalogue)
	}

	fmt.Printf("Indexing Complete. Master Map contains %d entries\n", len(cv.MasterEntries))

	for _, file := range allFiles {
		name := file.Name()
		isLibrary := strings.Contains(name, "Library")
		isGst := strings.HasSuffix(name, ".gst")
		isNotCat := !strings.HasSuffix(name, ".cat")

		if isLibrary || isGst || isNotCat {
			continue
		}

		filePath := filepath.Join(rawDir, name)
		xmlData, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading %s for conversion: %v\n", name, err)
			continue
		}

		var catalogue Catalogue
		err = xml.Unmarshal(xmlData, &catalogue)
		if err != nil {
			fmt.Printf("Unmarshal error in %s during conversion: %v\n", name, err)
			continue
		}

		fmt.Printf("Converting File: %s\n", catalogue.Name)

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
		fmt.Printf("Successfully converted %s (%d units)\n", catalogue.Name, len(seed.Factions[0].Units))
	}
}

func (c *Converter) transformUnit(entry SelectionEntry) models.UnitSeed {
	unit := models.UnitSeed{
		Name:            entry.Name,
		AdditionalStats: make(map[string]string),
	}

	for _, cost := range entry.Costs {
		if cost.Name == "pts" {
			var p float64
			_, err := fmt.Sscanf(cost.Value, "%d", &unit.Points)
			if err != nil {
				fmt.Printf("Warning: could not parse points '%s' for %s: %v\n", cost.Value, entry.Name, err)
				continue
			}
			unit.Points = int(p)
		}
	}

	for _, cat := range entry.CategoryLinks {
		unit.Keywords = append(unit.Keywords, strings.ToUpper(cat.Name))
	}

	allProfiles := c.collectProfiles(entry)

	for _, p := range allProfiles {
		if strings.Contains(p.TypeName, "Unit") || strings.Contains(p.TypeName, "Model") {
			c.mapStats(p.Characteristics, &unit)
		} else if strings.Contains(p.TypeName, "Weapon") {
			unit.Weapons = append(unit.Weapons, c.mapWeapon(p))
		}
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
	for _, cost := range entry.Costs {
		if cost.Name == "pts" && cost.Value != "0" && cost.Value != "" {
			return true
		}
	}

	for _, p := range entry.Profiles {
		if strings.Contains(p.TypeName, "Unit") || strings.Contains(p.TypeName, "Models") {
			return true
		}
	}

	for _, cat := range entry.CategoryLinks {
		upperCat := strings.ToUpper(cat.Name)
		if strings.Contains(upperCat, "INFANTRY") || strings.Contains(upperCat, "HERO") {
			return true
		}
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
		found = c.findUnits(entry.SelectionEntryGroups, found)
	}
	return found
}

func (c *Converter) collectProfiles(entry SelectionEntry) []Profile {
	var found []Profile

	found = append(found, entry.Profiles...)

	if entry.TargetID != "" {
		if target, ok := c.MasterEntries[entry.TargetID]; ok {
			targetProfiles := c.collectProfiles(target)
			found = append(found, targetProfiles...)
		}
	}

	containers := [][]SelectionEntry{
		entry.ChildEntries,
		entry.LinkEntries,
		entry.SelectionEntryGroups,
	}
	for _, container := range containers {
		for _, child := range container {
			childProfiles := c.collectProfiles(child)
			found = append(found, childProfiles...)
		}
	}
	return found
}

func (c *Converter) indexEntry(entry SelectionEntry) {
	if entry.ID != "" {
		c.MasterEntries[entry.ID] = entry
	}

	for _, child := range entry.ChildEntries {
		c.indexEntry(child)
	}
	for _, link := range entry.LinkEntries {
		c.indexEntry(link)
	}
	for _, group := range entry.SelectionEntryGroups {
		c.indexEntry(group)
	}
}

func (c *Converter) indexCatalogue(cat Catalogue) {
	for _, entry := range cat.SelectionEntries {
		c.indexEntry(entry)
	}
	for _, link := range cat.EntryLinks {
		c.indexEntry(link)
	}
	for _, shared := range cat.SharedEntries {
		c.indexEntry(shared)
	}
	for _, group := range cat.SharedGroups {
		c.indexEntry(group)
	}
}

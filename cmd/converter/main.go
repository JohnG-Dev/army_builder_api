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
		if err := xml.Unmarshal(xmlData, &catalogue); err != nil {
			fmt.Printf("Unmarshal error in %s: %v\n", file.Name(), err)
			continue
		}
		seed := models.SeedData{
			GameName: "Age of Sigmar",
			Factions: []models.FactionSeed{{
				Name:   catalogue.Name,
				Source: "Battlescribe Data",
				Units:  []models.UnitSeed{},
			}},
		}
		for _, entry := range catalogue.SelectionEntries {
			if entry.Type == "model" || entry.Type == "unit" {
				unit := cv.transformUnit(entry)
				seed.Factions[0].Units = append(seed.Factions[0].Units, unit)
			}
		}
		outPath := filepath.Join(outDir, strings.ToLower(strings.ReplaceAll(catalogue.Name, " ", "_"))+".yaml")
		yamlData, _ := yaml.Marshal(seed)
		os.WriteFile(outPath, yamlData, 0o644)
		fmt.Printf("Converted %s (%d units)\n", catalogue.Name, len(seed.Factions[0].Units))
	}
}

func (c *Converter) transformUnit(entry SelectionEntry) models.UnitSeed {
	unit := models.UnitSeed{Name: entry.Name, AdditionalStats: make(map[string]string)}
	for _, p := range entry.Profiles {
		if p.TypeName == "Unit" {
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

package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/JohnG-Dev/army_builder_api/internal/models"
)

var StatMapping = map[string]string{
	"Move":       "move",
	"M":          "move",
	"Health":     "health",
	"W":          "health",
	"Save":       "save",
	"Sv":         "save",
	"Control":    "control",
	"OC":         "control",
	"Toughness":  "toughness",
	"T":          "toughness",
	"Leadership": "leadership",
	"Ld":         "leadership",
	"Bravery":    "leadership",
	"Ward":       "ward",
	"Inv":        "invuln",
}

func main() {
	rawDir := "./data/raw"
	files, err := os.ReadDir(rawDir)
	if err != nil {
		log.Fatalf("failed to read raw directory: %v", err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".cat") {
			continue
		}

		filePath := filepath.Join(rawDir, file.Name())

		xmlData, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("error reading file %s: %v\n", file.Name(), err)
		}

		var catalogue Catalogue
		err = xml.Unmarshal(xmlData, &catalogue)
		if err != nil {
			fmt.Printf("failed to unmasrhal %s: %v", file.Name(), err)
			continue
		}
		fmt.Printf("Successfully unmarshaled: %s\n", catalogue.Name)
	}
	for _, entry := range catalogue.SelectionEntries {
		switch entry.Type {
		case "model", "unit":
			unit := sr.processUnit(entry)
			faction.Units = append(faction.Units, unit)
		case "upgrade":

			enh := sr.processEnhancement(entry)
			faction.Enhancements = append(faction.Enhancements, enh)
		case "selection":

			bf := sr.processBattleFormation(entry)
			faction.BattleFormations = append(faction.BattleFormations, bf)
		}
	}
}

func getStat(characteristics []Characteristic, statName string) string {
	for _, c := range characteristics {
		if c.Name == statName {
			return c.Value
		}
	}
	return ""
}

func (c *Converter) extractUnitStats(profile Profile, unit *models.UnitSeed) {
	if unit.AdditionalStats == nil {
		unit.AdditionalStats = make(map[string]string)
	}

	for _, char := range profile.Characterisitcs {
		targetColumn, isKnown := StatMapping[char.Name]

		if isKnown {
			switch targetColumn {
			case "move":
				unit.Move = char.Value
			case "health":
				unit.Health = char.Value
			case "save":
				unit.Save = char.Value
			case "ward":
				unit.Ward = char.Value
			case "invuln":
				unit.Invuln = char.Value
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

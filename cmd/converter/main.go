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
	"Save": "save", "Sv": "save", "SV": "save",
	"Control": "control", "OC": "control",
	"Toughness": "toughness", "T": "toughness",
	"Leadership": "leadership", "Ld": "leadership",
	"Bravery": "leadership", "Ward": "ward",
	"Strength": "strength", "S": "strength",
}

var AllegianceMap = map[string]string{
	"STORMCAST": "ORDER", "CITIES": "ORDER", "DAUGHTERS": "ORDER", "FYRESLAYERS": "ORDER", "IDONETH": "ORDER", "KHARADRON": "ORDER", "LUMINETH": "ORDER", "SERAPHON": "ORDER", "SYLVANETH": "ORDER",
	"KHORNE": "CHAOS", "TZEENTCH": "CHAOS", "NURGLE": "CHAOS", "SLAANESH": "CHAOS", "SKAVEN": "CHAOS", "SLAVES": "CHAOS", "BEASTS": "CHAOS", "HELSMITHS": "CHAOS",
	"FLESH-EATER": "DEATH", "NIGHTHAUNT": "DEATH", "OSSIARCH": "DEATH", "SOULBLIGHT": "DEATH",
	"GLOOMSPITE": "DESTRUCTION", "IRONJAWZ": "DESTRUCTION", "KRULEBOYZ": "DESTRUCTION", "OGOR": "DESTRUCTION", "SONS": "DESTRUCTION", "ORRUK": "DESTRUCTION", "BONESPLITTERZ": "DESTRUCTION",
}

type Converter struct {
	MasterEntries map[string]SelectionEntry
	GameSystems   map[string]string
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
		GameSystems:   make(map[string]string),
	}

	for _, file := range allFiles {
		name := file.Name()
		if !strings.HasSuffix(name, ".cat") && !strings.HasSuffix(name, ".gst") {
			continue
		}

		filePath := filepath.Join(rawDir, name)
		xmlData, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		if strings.HasSuffix(name, ".gst") {
			var gs GameSystem
			if err := xml.Unmarshal(xmlData, &gs); err == nil {
				cv.GameSystems[gs.ID] = gs.Name
				cv.indexGameSystem(gs)
			}
		} else {
			var catalogue Catalogue
			if err := xml.Unmarshal(xmlData, &catalogue); err == nil {
				cv.indexCatalogue(catalogue)
			}
		}
	}

	fmt.Printf("Indexing Complete. Master Map contains %d entries\n", len(cv.MasterEntries))

	for _, file := range allFiles {
		name := file.Name()
		if strings.Contains(name, "Library") || !strings.HasSuffix(name, ".cat") {
			continue
		}

		filePath := filepath.Join(rawDir, name)
		xmlData, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		var catalogue Catalogue
		if err := xml.Unmarshal(xmlData, &catalogue); err != nil {
			continue
		}

		gameName := cv.GameSystems[catalogue.GameSystemID]
		if gameName == "" {
			gameName = "Unknown Game"
		}

		fmt.Printf("Converting File: %s (%s)\n", catalogue.Name, gameName)

		seed := models.SeedData{
			GameName: gameName,
			Factions: []models.FactionSeed{{
				Name:   catalogue.Name,
				Source: "Battlescribe Data",
				Units:  []models.UnitSeed{},
			}},
		}

		upperName := strings.ToUpper(catalogue.Name)
		for key, allegiance := range AllegianceMap {
			if strings.Contains(upperName, key) {
				seed.Factions[0].Allegiance = allegiance
				break
			}
		}

		var allUnits []SelectionEntry
		allUnits = cv.findUnits(catalogue.SelectionEntries, allUnits)
		allUnits = cv.findUnits(catalogue.EntryLinks, allUnits)
		allUnits = cv.findUnits(catalogue.SharedEntries, allUnits)
		allUnits = cv.findUnits(catalogue.SharedGroups, allUnits)

		for _, entry := range allUnits {
			unit := cv.transformUnit(entry, name)
			if seed.Factions[0].Allegiance != "" {
				unit.Keywords = append(unit.Keywords, seed.Factions[0].Allegiance)
			}
			seed.Factions[0].Units = append(seed.Factions[0].Units, unit)
		}

		gameSlug := strings.ToLower(strings.ReplaceAll(gameName, " ", "_"))
		factionSlug := strings.ToLower(strings.ReplaceAll(catalogue.Name, " ", "_"))

		finalOutDir := filepath.Join(outDir, gameSlug)
		os.MkdirAll(finalOutDir, 0o755)

		outPath := filepath.Join(finalOutDir, factionSlug+".yaml")
		yamlData, _ := yaml.Marshal(seed)
		os.WriteFile(outPath, yamlData, 0o644)
		fmt.Printf(" -> Saved to %s (%d units)\n", outPath, len(seed.Factions[0].Units))
	}
}

func (c *Converter) transformUnit(entry SelectionEntry, fileName string) models.UnitSeed {
	unit := models.UnitSeed{
		Name:            entry.Name,
		MatchedPlay:     true,
		Source:          "Battlescribe Data",
		AdditionalStats: make(map[string]string),
	}

	if strings.Contains(strings.ToUpper(fileName), "LEGENDS") {
		unit.MatchedPlay = false
	}

	unit.Points = c.parsePoints(entry.Costs)
	if unit.Points == 0 && entry.TargetID != "" {
		if target, ok := c.MasterEntries[entry.TargetID]; ok {
			unit.Points = c.parsePoints(target.Costs)
		}
	}

	allConstraints := c.collectConstraints(entry)
	c.processConstraints(allConstraints, &unit)

	// Reinforcement Check (AoS 4.0 specific)
	if c.canBeReinforced(entry) {
		unit.MaxUnitSize = unit.MinUnitSize * 2
	}

	c.collectKeywords(entry, &unit)

	allProfiles := c.collectProfiles(entry)
	for _, p := range allProfiles {
		if strings.Contains(p.TypeName, "Unit") || strings.Contains(p.TypeName, "Model") || strings.Contains(p.TypeName, "Stats") {
			c.mapStats(p.Characteristics, &unit)
		} else if strings.Contains(p.TypeName, "Weapon") {
			unit.Weapons = append(unit.Weapons, c.mapWeapon(p))
		}
	}

	return unit
}

func (c *Converter) processConstraints(constraints []Constraint, unit *models.UnitSeed) {
	for _, cons := range constraints {
		var val int
		fmt.Sscanf(cons.Value, "%d", &val)

		if cons.Type == "max" && val == 1 && cons.Scope == "roster" {
			unit.IsUnique = true
		}

		// 2. UNIT SIZE DETECTION
		if cons.Field == "selections" && (cons.Scope == "parent" || cons.Scope == "self") {
			if cons.Type == "min" && val > 0 {
				if val > unit.MinUnitSize {
					unit.MinUnitSize = val
				}
			}
			if cons.Type == "max" && val > 0 {
				if val > unit.MaxUnitSize {
					unit.MaxUnitSize = val
				}
			}
		}
	}
}

func (c *Converter) collectKeywords(entry SelectionEntry, unit *models.UnitSeed) {
	if entry.TargetID != "" {
		if target, ok := c.MasterEntries[entry.TargetID]; ok {
			c.collectKeywords(target, unit)
		}
	}

	for _, cat := range entry.CategoryLinks {
		keyword := strings.ToUpper(cat.Name)
		isDup := false
		for _, existing := range unit.Keywords {
			if existing == keyword {
				isDup = true
				break
			}
		}
		if !isDup {
			unit.Keywords = append(unit.Keywords, keyword)
		}

		if strings.Contains(keyword, "MANIFESTATION") || strings.Contains(keyword, "ENDLESS SPELL") {
			unit.IsManifestation = true
		}
		if keyword == "UNIQUE" {
			unit.IsUnique = true
		}
		if keyword == "LEGENDS" || keyword == "WARHAMMER LEGENDS" {
			unit.MatchedPlay = false
		}
	}
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
		case "Range", "Rng":
			w.Range = char.Value
		case "Attacks", "A", "Atk":
			w.Attacks = char.Value
		case "To Hit", "WS", "Hit":
			w.ToHit = char.Value
		case "To Wound", "S", "Wnd":
			w.ToWound = char.Value
		case "Rend", "AP", "Rnd":
			w.Rend = char.Value
		case "Damage", "D", "Dmg":
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
			found = append(found, c.collectProfiles(target)...)
		}
	}
	containers := [][]SelectionEntry{entry.ChildEntries, entry.LinkEntries, entry.SelectionEntryGroups}
	for _, container := range containers {
		for _, child := range container {
			found = append(found, c.collectProfiles(child)...)
		}
	}
	return found
}

func (c *Converter) collectConstraints(entry SelectionEntry) []Constraint {
	var found []Constraint
	found = append(found, entry.Constraints...)

	if entry.TargetID != "" {
		if target, ok := c.MasterEntries[entry.TargetID]; ok {
			found = append(found, c.collectConstraints(target)...)
		}
	}

	containers := [][]SelectionEntry{entry.ChildEntries, entry.LinkEntries, entry.SelectionEntryGroups}
	for _, container := range containers {
		for _, child := range container {
			found = append(found, c.collectConstraints(child)...)
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

func (c *Converter) indexGameSystem(gs GameSystem) {
	for _, entry := range gs.SharedEntries {
		c.indexEntry(entry)
	}
	for _, group := range gs.SharedGroups {
		c.indexEntry(group)
	}
}

func (c *Converter) indexCatalogue(cat Catalogue) {
	containers := [][]SelectionEntry{cat.SelectionEntries, cat.EntryLinks, cat.SharedEntries, cat.SharedGroups}
	for _, container := range containers {
		for _, entry := range container {
			c.indexEntry(entry)
		}
	}
}

func (c *Converter) parsePoints(costs []Cost) int {
	for _, cost := range costs {
		if cost.Name == "pts" {
			var p float64
			_, err := fmt.Sscanf(cost.Value, "%f", &p)
			if err == nil {
				return int(p)
			}
		}
	}
	return 0
}

func (c *Converter) canBeReinforced(entry SelectionEntry) bool {
	reinforcedID := "1b37-82b8-c062-eb82"

	for _, mod := range entry.Modifiers {
		for _, cond := range mod.Conditions {
			if cond.ChildID == reinforcedID {
				return true
			}
		}
		for _, rep := range mod.Repeats {
			if rep.ChildID == reinforcedID {
				return true
			}
		}
	}

	if entry.TargetID != "" {
		if target, ok := c.MasterEntries[entry.TargetID]; ok {
			if c.canBeReinforced(target) {
				return true
			}
		}
	}

	containers := [][]SelectionEntry{entry.ChildEntries, entry.LinkEntries, entry.SelectionEntryGroups}
	for _, container := range containers {
		for _, child := range container {
			if c.canBeReinforced(child) {
				return true
			}
		}
	}

	return false
}

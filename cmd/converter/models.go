package main

import "encoding/xml"

type Catalogue struct {
	XMLName          xml.Name         `xml:"catalogue"`
	Name             string           `xml:"name,attr"`
	ID               string           `xml:"id,attr"`           // Needed to index GST files
	GameSystemID     string           `xml:"gameSystemId,attr"` // Fixed case: gameSystemId
	SelectionEntries []SelectionEntry `xml:"selectionEntries>selectionEntry"`
	EntryLinks       []SelectionEntry `xml:"entryLinks>entryLink"`
	SharedEntries    []SelectionEntry `xml:"sharedSelectionEntries>selectionEntry"`
	SharedGroups     []SelectionEntry `xml:"sharedSelectionEntryGroups>selectionEntryGroup"`
	SharedProfiles   []Profile        `xml:"sharedProfiles>profile"`
	CatalogueLinks   []CatalogueLink  `xml:"catalogueLinks>catalogueLink"`
	CategoryLinks    []CategoryLink   `xml:"categoryLinks>categoryLink"`    // Added this
	CategoryEntries  []Category       `xml:"categoryEntries>categoryEntry"` // Fixed typo: CategoryEntries
}

type GameSystem struct {
	XMLName        xml.Name         `xml:"gameSystem"`
	ID             string           `xml:"id,attr"`
	Name           string           `xml:"name,attr"`
	SharedEntries  []SelectionEntry `xml:"sharedSelectionEntries>selectionEntry"`
	SharedGroups   []SelectionEntry `xml:"sharedSelectionEntryGroups>selectionEntryGroup"`
	SharedProfiles []Profile        `xml:"sharedProfiles>profile"`
}

type SelectionEntry struct {
	Name                 string           `xml:"name,attr"`
	Type                 string           `xml:"type,attr"`
	ID                   string           `xml:"id,attr"`
	TargetID             string           `xml:"targetId,attr"`
	Profiles             []Profile        `xml:"profiles>profile"`
	Costs                []Cost           `xml:"costs>cost"`
	CategoryLinks        []CategoryLink   `xml:"categoryLinks>categoryLink"`
	ChildEntries         []SelectionEntry `xml:"selectionEntries>selectionEntry"`
	LinkEntries          []SelectionEntry `xml:"entryLinks>entryLink"`
	SelectionEntryGroups []SelectionEntry `xml:"selectionEntryGroups>selectionEntryGroup"`
	Constraints          []Constraint     `xml:"constraints>constraint"`
	Modifiers            []Modifier       `xml:"modifiers>modifier"` // Added this
}

type Modifier struct {
	Type       string      `xml:"type,attr"`
	Field      string      `xml:"field,attr"`
	Value      string      `xml:"value,attr"`
	Conditions []Condition `xml:"conditions>condition"`
	Repeats    []Repeat    `xml:"repeats>repeat"`
}

type Condition struct {
	ChildID string `xml:"childId,attr"`
}

type Repeat struct {
	ChildID string `xml:"childId,attr"`
}

type Profile struct {
	Name            string           `xml:"name,attr"`
	ID              string           `xml:"id,attr"`
	TypeName        string           `xml:"typeName,attr"`
	Characteristics []Characteristic `xml:"characteristics>characteristic"`
}

type Characteristic struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type RuleEntry struct {
	Name        string `xml:"name,attr"`
	Description string `xml:"description"`
}

type Category struct {
	Name string `xml:"name,attr"`
}

type Cost struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type CategoryLink struct {
	Name    string `xml:"name,attr"`
	Primary bool   `xml:"primary,attr"`
}

type CatalogueLink struct {
	TargetID string `xml:"targetId,attr"`
	Name     string `xml:"name,attr"`
}

type Constraint struct {
	Type  string `xml:"type,attr"`
	Value string `xml:"value,attr"`
	Field string `xml:"field,attr"`
	Scope string `xml:"scope,attr"`
}

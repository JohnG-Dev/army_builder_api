package main

import "encoding/xml"

type Catalogue struct {
	XMLName          xml.Name         `xml:"catalogue"`
	Name             string           `xml:"name,attr"`
	SelectionEntries []SelectionEntry `xml:"selectionEntries>selectionEntry"`
	EntryLinks       []SelectionEntry `xml:"entryLinks>entryLink"`
	SharedEntries    []SelectionEntry `xml:"sharedSelectionEntries>selectionEntry"`
	SharedGroups     []SelectionEntry `xml:"sharedSelectionEntryGroups>selectionEntryGroup"`
}

type SelectionEntry struct {
	Name         string           `xml:"name,attr"`
	Type         string           `xml:"type,attr"`
	Profiles     []Profile        `xml:"profiles>profile"`
	Categories   []Category       `xml:"categoryLinks>categoryLink"`
	ChildEntries []SelectionEntry `xml:"selectionEntries>selectionEntry"`
	LinkEntries  []SelectionEntry `xml:"entryLinks>entryLink"`
	GroupEntries []SelectionEntry `xml:"selectionEntryGroups>selectionEntryGroup>selectionEntries>selectionEntry"`
}

type Profile struct {
	Name            string           `xml:"name,attr"`
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

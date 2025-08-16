package metadata

import (
	"encoding/xml"
	"io"
	"os"
)

type Model struct {
	XMLName      xml.Name     `xml:"Edmx"`
	DataServices DataServices `xml:"DataServices"`
}

type DataServices struct {
	Schema Schema `xml:"Schema"`
}

type Schema struct {
	Namespace       string          `xml:"Namespace,attr"`
	EntityTypes     []EntityType    `xml:"EntityType"`
	ComplexTypes    []ComplexType   `xml:"ComplexType"`
	EnumTypes       []EnumType      `xml:"EnumType"`
	EntityContainer EntityContainer `xml:"EntityContainer"`
}

type EntityContainer struct {
	Name       string      `xml:"Name,attr"`
	EntitySets []EntitySet `xml:"EntitySet"`
}

type EntitySet struct {
	Name       string `xml:"Name,attr"`
	EntityType string `xml:"EntityType,attr"`
}

type EntityType struct {
	Name                 string               `xml:"Name,attr"`
	Key                  Key                  `xml:"Key"`
	Properties           []Property           `xml:"Property"`
	NavigationProperties []NavigationProperty `xml:"NavigationProperty"`
}

type ComplexType struct {
	Name       string     `xml:"Name,attr"`
	Properties []Property `xml:"Property"`
}

type EnumType struct {
	Name    string       `xml:"Name,attr"`
	Members []EnumMember `xml:"Member"`
}

type EnumMember struct {
	Name  string `xml:"Name,attr"`
	Value string `xml:"Value,attr"`
}

type Key struct {
	PropertyRefs []PropertyRef `xml:"PropertyRef"`
}

type PropertyRef struct {
	Name string `xml:"Name,attr"`
}

type NavigationProperty struct {
	Name           string `xml:"Name,attr"`
	Type           string `xml:"Type,attr"`
	ContainsTarget string `xml:"ContainsTarget,attr"`
	Partner        string `xml:"Partner,attr"`
}

type Property struct {
	Name     string `xml:"Name,attr"`
	Type     string `xml:"Type,attr"`
	Nullable string `xml:"Nullable,attr"`
}

// Parse parses the metadata .xml file and returns a model
// representation.
func Parse(path string) (*Model, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var edmx Model
	if err := xml.Unmarshal(data, &edmx); err != nil {
		return nil, err
	}

	return &edmx, nil
}

package converter

// This file contains structs that model the XML Schema in Go.
// The model below was created by looking over http://www.w3.org/2001/XMLSchema.xsd

import (
	"encoding/xml"
)

type Schema *struct {
	Name         string
	Annotation   []Annotation  `xml:"annotation"`
	ComplexTypes []ComplexType `xml:"complexType"`
	SimpleTypes  []SimpleType  `xml:"simpleType"`
	Elements     []Element     `xml:"element"`
	Comment      string        `xml:",comment"`

	// Types used by the schema
	Aliases map[Alias]SimpleKind `xml:"-"`
}

type ComplexType struct {
	Comment        string         `xml:",comment"`
	Name           string         `xml:"name,attr"`
	Sequences      []Sequence     `xml:"sequence"`
	Choices        []Choice       `xml:"choice"`
	Attributes     []Attribute    `xml:"attribute"`
	Alias          Alias          `xml:"restriction"`
	ComplexContent ComplexContent `xml:"complexContent"`
	SimpleContent  SimpleContent  `xml:"simpleContent"`
}

type ComplexContent struct {
}

type SimpleContent struct {
}

type Extension struct {
}

type Element struct {
	XMLName   xml.Name
	Type      string `xml:"type,attr"`
	MinOccurs int    `xml:"minOccurs,attr"`
	MaxOccurs int    `xml:"maxOccurs,attr"`
	//ComplexTypes []ComplexType `xml:"complexType"`
}

type Choice struct {
	Elements  []Element `xml:"element"`
	MinOccurs int       `xml:"minOccurs,attr"`
	MaxOccurs int       `xml:"maxOccurs,attr"`
}

type Sequence struct {
	Elements []Element `xml:"element"`
}

type Annotation struct {
	Source        string `xml:"source"`
	Documentation string `xml:"documentation"`
}

type Attribute struct {
	XMLName xml.Name
	Type    string `xml:"type,attr"`
}

// The following types are ignored for now because they are too complex
// type Key struct {}
// type KeyRef struct {}
// type Selector struct {}
// type Field struct {}
// type Include struct {}
// type Redefine struct {}
// type Group struct {}
// type Import struct {}
// type Notation struct {}

type SimpleKind int

// These are a list of the types we may generate as we parse the input XSD.
// Clients that process our model only have to handle these types.
const (
	TYPE_STRING SimpleKind = iota
	TYPE_BOOLEAN
	TYPE_FLOAT
	TYPE_DOUBLE
	TYPE_DECIMAL
	TYPE_DURATION
	TYPE_DATETIME
	TYPE_TIME
	TYPE_DATE
	TYPE_G_YEAR_MONTH
	TYPE_G_YEAR
	TYPE_G_MONTH_DAY
	TYPE_G_DAY
	TYPE_G_MONTH
	TYPE_HEX_BINARY
	TYPE_BASE_64_BINARY
	TYPE_ANY_URI
	TYPE_Q_NAME
	TYPE_NORMALIZED_STRING
	// Skipping some elements here, that seem to be only useful for
	// modelling xsd itself.
	TYPE_TOKEN
	TYPE_INTEGER
	TYPE_INT
	TYPE_SHORT
	TYPE_BYTE
	TYPE_U_INTEGER
	TYPE_U_INT
	TYPE_U_SHORT
	TYPE_U_BYTE
	TYPE_POSITIVE_INTEGER
)

type Alias struct {
	Base string `xml:"base,attr"`
}

type SimpleType struct {
	Name        string       `xml:"name,attr"`
	Alias       Alias        `xml:"restriction"`
	Annotations []Annotation `xml:"annotation"`
	Sequences   []Sequence   `xml:"sequence"`
	Attributes  []Attribute  `xml:"attribute"`
}

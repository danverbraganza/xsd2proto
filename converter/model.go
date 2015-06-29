package converter

// This file contains structs that model the XML Schema in Go.
// The model below was created by looking over http://www.w3.org/2001/XMLSchema.xsd

import (
	"encoding/xml"
	"fmt"
//	"log"
)

type Schema struct {
	Name         string
	Annotation   []Annotation  `xml:"annotation"`
	Attributes   []Attribute   `xml:"attribute"`
	ComplexTypes []ComplexType `xml:"complexType"`
	SimpleTypes  []SimpleType  `xml:"simpleType"`
	Elements     []Element     `xml:"element"`
	Comment      string        `xml:",comment"`

	// Types used by the schema
	AttrRefs    map[string]Attribute `xml:"-"`
	TypeRefs    map[string]Type      `xml:"-"`
	ElementRefs map[string]Element   `xml:"-"`
}

// Must be a pointer to either a simple or a complex type, or nil.
type Type interface{
	GetName() string
	DeRef(s Schema) Type
}

type ComplexType struct {
	Comment        string         `xml:",comment"`
	Name           *string        `xml:"name,attr"`
	Ref            *string        `xml:"ref,attr"`
	Sequence       *Sequence      `xml:"sequence"`
	Choice         Choice         `xml:"choice"`
	Attributes     []Attribute    `xml:"attribute"`
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
	XMLName     xml.Name
	Name        *string      `xml:"name,attr"`
	Type        *string      `xml:"type,attr"`
	Ref         *string      `xml:"ref,attr"`
	MinOccurs   int          `xml:"minOccurs,attr"`
	MaxOccurs   int          `xml:"maxOccurs,attr"`
	ComplexType *ComplexType `xml:"complexType"`
	SimpleType  *SimpleType  `xml:"simpleType"`
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
	Name    *string `xml:"name,attr"`
	Ref     *string `xml:"ref,attr"`
	Type    *string `xml:"type,attr"`
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
	TYPE_INVALID SimpleKind = iota
	TYPE_STRING
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
)

type Restriction struct {
	Ref *string `xml:"base,attr"`
}

// Handle the ref types by making them referencable.
type SimpleType struct {
	Name        *string      `xml:"name,attr"`
	Ref         *string      `xml:"ref,attr"`
	Restriction       Restriction        `xml:"restriction"`
	Annotations []Annotation `xml:"annotation"`
	Sequences   []Sequence   `xml:"sequence"`
	Attributes  []Attribute  `xml:"attribute"`
}

func (s *Schema) BuildRefs() {

	s.ElementRefs = map[string]Element{}
	for _, e := range s.Elements {
		if e.Name != nil {
			s.ElementRefs[*e.Name] = e
		}
		fmt.Println(s.ElementRefs)
	}

	s.AttrRefs = map[string]Attribute{}
	for _, a := range s.Attributes {
		if a.Name != nil {
			s.AttrRefs[*a.Name] = a
		}
	}

	s.TypeRefs = map[string]Type{}
	for _, c := range s.ComplexTypes {
		if c.Name != nil {
			s.TypeRefs[*c.Name] = c
		}
	}

	for _, t := range s.SimpleTypes {
		if t.Name != nil {
			s.TypeRefs[*t.Name] = t
		}
	}
}

func (a Attribute) DeRef(s Schema) Attribute {
	if a.Ref == nil {
		return a
	}
	return s.AttrRefs[*a.Ref]
}

func (e Element) DeRef(s Schema) Element {
	if e.Ref == nil {
		return e
	}
	if e, ok := s.ElementRefs[*e.Ref]; ok {
		return e
	}
	fmt.Println("Element refs", s.ElementRefs)
	return e
}


func (r Restriction) DeRef(s Schema) Type {
	if r.Ref == nil {
		return r
	}
	return s.TypeRefs[*r.Ref]
}

func (r Restriction) GetName() string {
	return *r.Ref
}


var XmlTypes = map[string]SimpleKind {
	"xs:string": TYPE_STRING,
	"xs:boolean": TYPE_INT,
	"xs:float": TYPE_FLOAT,
	"xs:double": TYPE_DOUBLE,
	"xs:decimal": TYPE_DECIMAL,
	"xs:duration": TYPE_DURATION,
	"xs:dateTime": TYPE_DATETIME,
	"xs:time": TYPE_TIME,
	"xs:date": TYPE_DATE,
	"xs:gYearMonth": TYPE_G_YEAR_MONTH,
	"xs:gYear": TYPE_G_YEAR,
	"xs:gMonthDay": TYPE_G_MONTH_DAY,
	"xs:gDay": TYPE_G_DAY,
	"xs:gMonth": TYPE_G_MONTH,
	"xs:hexBinary": TYPE_HEX_BINARY,
	"xs:base64Binary": TYPE_BASE_64_BINARY,
	"xs:anyURI": TYPE_ANY_URI,
	"xs:QName": TYPE_Q_NAME,
	"xs:normalizedString": TYPE_STRING,
	"xs:token": TYPE_TOKEN,
	"xs:language": TYPE_STRING,
	"xs:IDREFS": TYPE_STRING,
	"xs:ENTITIES": TYPE_STRING,
	"xs:NMTOKEN": TYPE_STRING,
	"xs:NCName": TYPE_STRING,
	"xs:ID": TYPE_STRING,
	"xs:IDREF": TYPE_STRING,
	"xs:ENTITY": TYPE_STRING,
	"xs:integer": TYPE_INTEGER,
	"xs:nonPositiveInteger": TYPE_INTEGER,
	"xs:negativeInteger": TYPE_INTEGER,
	"xs:long": TYPE_INTEGER,
	"xs:int": TYPE_INT,
	"xs:short": TYPE_SHORT,
	"xs:byte": TYPE_BYTE,
	"xs:nonNegativeInteger": TYPE_U_INTEGER,
	"xs:unsignedLong": TYPE_U_INTEGER,
	"xs:unsignedInt": TYPE_U_INTEGER,
	"xs:unsignedShort": TYPE_U_INTEGER,
	"xs:unsignedByte": TYPE_U_INTEGER,
	"xs:positiveInteger": TYPE_U_INTEGER,
}

func (a Attribute) Kind() SimpleKind {
	return XmlTypes[*a.Type]
}

// Requires access to the Schema, because we may need to walk a tree to find the
// base type in the xs: namespace.
func (t SimpleType) Kind(s Schema) SimpleKind {
	if a, ok := XmlTypes[*t.Name]; ok {
		return a
	}
	// Couldn't find it, let's find the base type.
	return t.Restriction.DeRef(s).(SimpleType).Kind(s)
}



func (t SimpleType) DeRef(s Schema) Type {
	if t.Ref == nil {
		return t
	}
	return s.TypeRefs[*t.Ref]
}

func (t SimpleType) GetName() string {
	return *t.Name
}


func (t ComplexType) DeRef(s Schema) Type {
	if t.Ref == nil {
		return t
	}
	return s.TypeRefs[*t.Ref]
}

func (t ComplexType) GetName() string {
	if t.Name != nil {
		return *t.Name
	}
	return ""

}


func (e Element) Child() Type {
	switch {
	case e.ComplexType != nil:
		return e.ComplexType
	case e.SimpleType != nil:
		return e.SimpleType
	default:
		return nil

	}
}

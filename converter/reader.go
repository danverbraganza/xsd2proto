package converter;
// This file contains utilities to read xml schema definitions into a Go struct.
// Based on http://www.w3.org/2001/XMLSchema.xsd

import (
	"encoding/xml"
	"io"
)



func ReadXsd(reader io.Reader, name string) (s Schema, e error) {
	d := xml.NewDecoder(reader)
 	e = d.Decode(&s)
	s.Name = name
	s.BuildRefs()
	return
}

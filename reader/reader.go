package reader

// This file contains utilities to read xml schema definitions into a Go struct.
// Based on http://www.w3.org/2001/XMLSchema.xsd

import (
	"encoding/xml"
	"io"

	"github.com/danverbraganza/xsd2proto/model"
)

func ReadXsd(reader io.Reader, name string) (s model.Schema, e error) {
	d := xml.NewDecoder(reader)
	e = d.Decode(&s)
	s.Name = name
	s.BuildRefs()
	return
}

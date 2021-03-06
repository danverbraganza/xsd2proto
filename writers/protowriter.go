package writers

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/danverbraganza/varcaser/varcaser"

	"github.com/danverbraganza/xsd2proto/model"
)

type Type struct {
	Name    string
	Imports []string
	Fields  []Field
}

type Field struct {
	Name       string
	TypeRef    string               // If complex type: contains the type reference.
	SimpleKind model.SimpleKind // If simple type: contains the number.
	Repeated   bool
	Required   bool
}

type ProtoBuilder struct {
	Schema      model.Schema
	TypesByName map[string]*Type
	AllTypes    []*Type
}

var caseConverter = varcaser.Caser{
	From: varcaser.LowerCamelCase,
	To:   varcaser.LowerSnakeCase,
}

// AllTypes inspects a Schema and returns a slice of all the types that need to
// be constructed in it.
func FromSchema(s model.Schema) (p ProtoBuilder) {

	p.Schema = s

	p.TypesByName = map[string]*Type{}

	for _, c := range s.ComplexTypes {
		p.LoadComplexType(c)
	}

	for _, e := range s.Elements {
		p.LoadElement(e)
	}
	return p
}

// LoadElement is called to load an element. It turns an Element into a Type of
// ProtoBuilder. It also returns a Field, because Elements are
// eponymous fields of their superiors.
func (p *ProtoBuilder) LoadElement(e model.Element) Field {
	e = e.DeRef(p.Schema)
	fmt.Println(e.String())

	// A type was provided.
	if e.Type != nil {
		t := p.Schema.TypeRefs[*e.Type]
		fmt.Println(t)
		switch c := t.(type) {
		case model.ComplexType:
			if c.Name == nil {
				c.Name = e.Name
			}
			p.LoadComplexType(c)
			return Field{
				Name:    *e.Name,
				TypeRef: *c.Name,
			}
		case model.SimpleType:
			return Field{
				Name:       *e.Name,
				SimpleKind: c.Kind(p.Schema),
			}
		default: // Didn't find it, it's probably an xs: type
			return Field{
				Name:       *e.Name,
				SimpleKind: model.XmlTypes[*e.Type],
			}
		}
	}

	if child := e.Child(); child != nil {
		child = child.DeRef(p.Schema)
		switch c := child.(type) {
		case model.ComplexType:
			if c.Name == nil {
				c.Name = e.Name
			}
			p.LoadComplexType(c)
			return Field{
				Name:     *e.Name,
				TypeRef:  *e.Name,
				Repeated: false,
				Required: false,
			}
		case model.SimpleType:
			return Field{
				Name:       child.GetName(),
				SimpleKind: c.Kind(p.Schema),
				Repeated:   false,
				Required:   false,
			}
		}

	}

	log.Fatalf(e.String())
	return Field{Name: *e.Name, SimpleKind: model.TYPE_STRING}
}

// We are parsing a complex type, and we turn it into an entry in our map.
func (p *ProtoBuilder) LoadComplexType(c model.ComplexType) {
	if c.Name == nil {
		return
	}

	name := *c.Name

	if _, ok := p.TypesByName[name]; ok {
		// Already seen this type, it's being added.
		return
	}

	t := &Type{name, nil, nil}
	p.TypesByName[*c.Name] = t // Put this in the map so we know we've seen it if we recurse again.

	for _, a := range c.Attributes {
		a = a.DeRef(p.Schema)
		fmt.Println(a)
		if a.Type == nil {
			t.AddField(Field{Name: *a.Name, SimpleKind: model.TYPE_STRING})
		} else {
			t.AddField(Field{Name: *a.Name, SimpleKind: a.Kind()})
		}
	}

	if c.Sequence != nil {
		for _, e := range c.Sequence.Elements {
			field := p.LoadElement(e)
			t.AddField(field)
			if field.TypeRef != "" {
				// Was a complex type: import
				t.AddImport(field.TypeRef)
			}
		}
	}

	p.AllTypes = append(p.AllTypes, t)
}

func (t *Type) AddField(f Field) {
	t.Fields = append(t.Fields, f)
}

func (t *Type) AddImport(name string) {
	t.Imports = append(t.Imports, name)
}

func (t Type) ToProtoFile(dirname string, s model.Schema) error {
	f, err := os.Create(path.Join(dirname, t.FileName()))
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(`syntax = "proto2";`)
	f.WriteString("\n")

	f.WriteString("package ")
	f.WriteString(s.Name)
	f.WriteString(";\n\n")


	fmt.Println(t)

	for _, impor := range t.Imports {
		f.WriteString(`import "` + impor + `.proto";`)
		f.WriteString("\n")
	}
	f.WriteString("\n")

	f.WriteString("message ")
	f.WriteString(t.Name)
	f.WriteString(" {\n")

	for i, field := range t.Fields {
		f.WriteString("    ")
		f.WriteString(field.String())
		f.WriteString(strconv.Itoa(i + 1))
		f.WriteString(";\n")
	}

	f.WriteString("} \n")

	return nil
}

var ToProtoTypes = map[model.SimpleKind]string{
	model.TYPE_INVALID:           "",
	model.TYPE_STRING:            "string",
	model.TYPE_BOOLEAN:           "bool",
	model.TYPE_FLOAT:             "float",
	model.TYPE_DOUBLE:            "double",
	model.TYPE_DECIMAL:           "fixed64",
	model.TYPE_ANY_URI:           "string",
	model.TYPE_Q_NAME:            "string",
	model.TYPE_NORMALIZED_STRING: "string",
	model.TYPE_TOKEN:             "string",
	model.TYPE_INTEGER:           "int64",
	model.TYPE_INT:               "int64",
	model.TYPE_SHORT:             "int32",
	model.TYPE_BYTE:              "int32",
	model.TYPE_U_INTEGER:         "uint64",
	model.TYPE_U_INT:             "uint64",
	model.TYPE_U_SHORT:           "uint32",
	model.TYPE_U_BYTE:            "uint32",
}

func (f Field) String() string {
	b := &bytes.Buffer{}
	if f.Repeated {
		b.WriteString("repeated ")
	} else if f.Required {
		b.WriteString("required ")
	} else {
		b.WriteString("optional ")
	}
	if f.SimpleKind != model.TYPE_INVALID {
		b.WriteString(ToProtoTypes[f.SimpleKind])
	} else {
		b.WriteString(f.TypeRef)
	}
	b.WriteString(" ")
	b.WriteString(caseConverter.String(f.Name))
	b.WriteString(" = ")
	return b.String()
}

func (t Type) FileName() string {
	return t.Name + ".proto"

}

func DumpToProto(dirname string, s model.Schema) error {
	fmt.Println("Preparing to dump to protobuf in folder", dirname)

	if err := os.MkdirAll(dirname, os.ModeDir|os.ModePerm); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	for _, t := range FromSchema(s).AllTypes {
		err := t.ToProtoFile(dirname, s)
		if err != nil {
			return err
		}
	}

	return nil
}

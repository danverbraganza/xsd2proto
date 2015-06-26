package main;

import (
	"fmt"
	"os"

	"github.com/danverbraganza/xsd2proto/converter"
)


func main() {
	s := "XMLSchema.xsd"
	f, _ := os.Open("samples/" + s)

	a, err := converter.ReadXsd(f)
	a.Name = s

	fmt.Printf("%+v, %v", a, err)
}

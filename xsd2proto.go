package main;

import (
	"fmt"
	"os"

	"github.com/danverbraganza/xsd2proto/converter"
	"github.com/danverbraganza/xsd2proto/writers"
)


func main() {
	f, _ := os.Open("samples/LibraryBooks.xsd")

	a, err := converter.ReadXsd(f, "LibraryBooks")

	fmt.Printf("%+v, %v\n", a, err)

	writers.DumpToProto("temp", a)


}

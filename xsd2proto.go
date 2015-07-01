package main;

import (
	"fmt"
	"os"

	"github.com/danverbraganza/xsd2proto/reader"
	"github.com/danverbraganza/xsd2proto/writers"
)


func main() {
	f, _ := os.Open("samples/LibraryBooks.xsd")

	a, err := reader.ReadXsd(f, "LibraryBooks")

	fmt.Printf("%+v, %v\n", a, err)

	writers.DumpToProto("temp", a)
}

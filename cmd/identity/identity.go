package main

import "code.google.com/p/gogramps/xml"
import "flag"
import "fmt"
import "os"

var inFilename = flag.String("in", "", "The name of the gramps file to read")
var outFilename = flag.String("out", "", "The name of the gramps file to write")

func main() {
	f, err := os.Open(*inFilename)
	if err != nil {
		fmt.Println("Could not read file: ", err)
		return
	}
	db, err := xml.Parse(f)
	if err != nil {
		fmt.Println("Could not parse XML: ", err)
		return
	}
	if err = db.Serialize(*outFilename); err != nil {
		fmt.Println(err)
		return
	}
}

func init() {
	flag.Parse()
}

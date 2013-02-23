package main

import "flag"
import  "code.google.com/p/gogramps/xml"
import "fmt"

var inFilename = flag.String("in", "", "The name of the gramps file to read")
var outFilename = flag.String("out", "", "The name of the gramps file to write")

func main() {
  db, err := xml.Parse(*inFilename)
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

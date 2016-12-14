package main

import "fmt"
import "flag"

var flagVerbose bool

func main() {
	flag.BoolVar(&flagVerbose, "v", false, "Prints detailed operations")
	flag.Parse()


	fmt.Println("Verbose:", flagVerbose)
	


}
package main

import "fmt"
import "flag"
import "path/filepath"

var flagVerbose bool

func walktree(droot string) (int, error) {

	fmt.Println("Processing", droot)
	//err := filepath.Walk(droot, visit)
  //fmt.Printf("filepath.Walk() returned %v\n", err)

	return 1, nil
}


func main() {
	flag.BoolVar(&flagVerbose, "v", false, "Prints detailed operations")

	flag.Parse()

	items := []string{"."}

	fmt.Println("Verbose:", flagVerbose)
	fmt.Println("Args", flag.Args())

	if flag.NArg() > 0 {
		items = flag.Args()
	}

	for _,i := range items {
		topitems, err := filepath.Glob(i)
		if err != nil { fmt.Println(error(err))	}
		for _,j := range topitems {
			count, err := walktree(j)
			if err != nil { fmt.Println(error(err))	}
			fmt.Printf("Processed %d items in dir[%s]\n", count, i)
		}

	}

}



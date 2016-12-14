package main

import "fmt"
import "log"
import "flag"
import "path/filepath"
import "io/ioutil"

var flagVerbose bool
var flagDotted bool

func walktree(droot string) (int, error) {

	droot, err := filepath.Abs(droot)
	if err != nil { log.Fatal(err) }

	count := 0
	fmt.Println("Processing", droot)
	flist, err := ioutil.ReadDir(droot)
	if err != nil { log.Fatal(err) }

  for _, i := range flist {

  	if i.Name()[0] == '.' {
  		if !flagDotted {
  			fmt.Println("Skipping hidden", i.Name())
  			continue
  		}
  	}
  	count++

  	fmt.Println(i.Name(), i.Mode().IsDir())

  	if i.Mode().IsDir() {
  		dapath, err := filepath.Abs(filepath.Join(droot, i.Name()))
  		count2, err := walktree(dapath)
  		if err != nil { log.Fatal(err) }
  		count += count2
  	}
  }

	return count, nil
}


func main() {
	flag.BoolVar(&flagVerbose, "v", false, "Prints detailed operations")
	flag.BoolVar(&flagDotted, "d", false, "Follow hidden dot directorys")
	flag.Parse()

	items := []string{"."}  // default arguments to use if omitted

	//fmt.Println("Verbose:", flagVerbose)
	//fmt.Println("Args", flag.Args())

	if flag.NArg() > 0 {
		items = flag.Args()
	}

	for _, i := range items {
		topitems, err := filepath.Glob(i)
		if err != nil { log.Fatal(err) }
		for _, j := range topitems {
			count, err := walktree(j)
			if err != nil { log.Fatal(err)	}
			fmt.Printf("Processed %d items in dir[%s]\n", count, i)
		}

	}

}



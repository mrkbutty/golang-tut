/*
treetime will set directory timestamps to match most recent of contents.

*/

package main

import "fmt"
import "os"
import "log"
import "flag"
import "path/filepath"
import "io/ioutil"
import "time"

var flagVerbose bool
var flagDotted bool
var flagIgnoreDir bool

func walktree(droot string) (int, error) {

	droot, err := filepath.Abs(droot)
	if err != nil { log.Fatal(err) }
	drootinfo, err := os.Stat(droot)
	if err != nil { log.Fatal(err) }	

	count := 0
	maxtime := time.Time{}


	if flagVerbose { fmt.Printf("Processing %s\n", droot) }
	flist, err := ioutil.ReadDir(droot)
	if err != nil { log.Fatal(err) }

	//first time through we check for directories recursively call this function 
  for _, i := range flist {

  	if i.Name()[0] == '.' {
  		if !flagDotted {
  			if flagVerbose { fmt.Println("Skipping hidden", i.Name()) }
  			continue
  		}
  	}
  	count++

  	//fmt.Println(i.Name(), i.Mode().IsDir())

  	if i.Mode().IsDir() {
  		dapath, err := filepath.Abs(filepath.Join(droot, i.Name()))
  		count2, err := walktree(dapath)
  		if err != nil { log.Fatal(err) }
  		count += count2
  	}
  }

  // Must read again in case lower directories have updated time
  if !flagIgnoreDir {
		flist, err = ioutil.ReadDir(droot)  
		if err != nil { log.Fatal(err) }

  }

	//now get most updated modification time
  for _, i := range flist {
  	if flagIgnoreDir && i.Mode().IsDir() { continue }
  	if i.ModTime().After(maxtime) { maxtime = i.ModTime()}
  }

  if flagVerbose {
  	fmt.Printf("[%s]\t%v\t-->\t%v\n", filepath.Base(droot), 
	  	drootinfo.ModTime().Format(time.UnixDate), 
	  	maxtime.Format(time.UnixDate))
  }


	return count, nil
}


func main() {
	flag.BoolVar(&flagVerbose, "v", false, "Prints detailed operations")
	flag.BoolVar(&flagDotted, "d", false, "Follow hidden dot directorys")
	flag.BoolVar(&flagIgnoreDir, "i", false, "Ignore directory timestamps")
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



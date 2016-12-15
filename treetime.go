/*
treetime will set directory timestamps to match most recent of contents below.

In default mode it will recursively travel down the named directorys looking at
modification times setting the parent directory to the most recent.  This includes file and directory timestamps unless changed with "-i".

Usage: treetime <directory>...

	<directory> = list of directories to process.

  -d    Follow hidden dot directorys
  -i    Ignore directory timestamps
  -q    No output apart from errors
  -t    Test only do not change
  -v    Prints detailed operations

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
var flagQuiet bool
var flagDotted bool
var flagIgnoreDir bool
var flagTestOnly bool
var dateFormat=time.RFC1123

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
	checked := 0
  for _, i := range flist {
  	if flagIgnoreDir && i.Mode().IsDir() { continue }
  	if i.ModTime().After(maxtime) { maxtime = i.ModTime()}
  	checked++
  }

  if checked > 0 { // does not change empty directories
  	if maxtime.IsZero() { 
  		fmt.Printf("Warning: Only zero time found: %s", droot)
  	} else {
  		if flagVerbose {
		  	fmt.Printf("[%s]  %v  -->  %v\n", filepath.Base(droot), 
			  	drootinfo.ModTime().Format(dateFormat), 
			  	maxtime.Format(dateFormat))
		  }
		  if !flagTestOnly {
			  err = os.Chtimes(droot, time.Now(), maxtime)
			  if err != nil { log.Fatal(err) }		  	
		  }

		}
  }


	return count, nil
}


func main() {
	flag.BoolVar(&flagVerbose, "v", false, "Prints detailed operations")
	flag.BoolVar(&flagQuiet, "q", false, "No output apart from errors")
	flag.BoolVar(&flagDotted, "d", false, "Follow hidden dot directorys")
	flag.BoolVar(&flagIgnoreDir, "i", false, "Ignore directory timestamps")
	flag.BoolVar(&flagTestOnly, "t", false, "Test only do not change")
	flag.Parse()

	//items := []string{"."}  // default arguments to use if omitted

	if flag.NArg() == 0 {
		flag.PrintDefaults()
		return
	}

	start := time.Now()
	total := 0
	for _, i := range flag.Args() {
		topitems, err := filepath.Glob(i)
		if err != nil { log.Fatal(err) }
		for _, j := range topitems {
			count, err := walktree(j)
			if err != nil { log.Fatal(err)	}
			total += count
			if flagVerbose {
				fmt.Printf("Completed %d items in [%s]\n" , count, i)
			}
		}
	}
	if !flagQuiet {
		elapsed := time.Since(start)
		fmt.Printf("Processed %d total items in %v\n" , total, elapsed)
	}
}



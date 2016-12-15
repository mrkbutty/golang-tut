# treetime v0.1

My first ever Go program in an attempt to learn about the language.

Summary
-------
treetime will set directory timestamps to match most recent of contents below.

In default mode it will recursively travel down the named directorys looking at
modification times setting the parent directory to the most recent.  This includes file and directory timestamps unless changed with "-i".


**Usage:** 
```
treetime [directory] ...

	directory = defaults to "."

  -d    Follow hidden dot directorys
  -i    Ignore directory timestamps
  -q    No output apart from errors
  -t    Test only do not change
  -v    Prints detailed operations
```
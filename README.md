# GoVal

Small program to count stuff in your project, like number of packages, files, functions (internal and exported), number of functions with and without docs.

Some of this stuff can be excluded in the config

```
Hello GoVal

Packages    Files   Functions   Internal    Exported    No docs   With docs	
5           6       10          4           6           10        0		

```

With more details (this should be added as flag, not hardcoded)
```
Hello GoVal


Package		File		        Function          Ln. Start       Ln. End     Lines		
main		main.go		        main              23              35          13		
main		main.go		        parseDir          37              116         80		


Package		File                    Function          Ln. Start	  Ln. End     Lines		
util		util/error.go		ShowError         8               12          5		


.
.
.


Packages    Files   Functions   Internal    Exported    No docs   With docs	
5           6       10          4           6           10        0		

```

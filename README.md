# Go-Make 
Go-Make is a small, make-like, build tool inspired by make, designed to run on Windows. 

It has the basic features of make. 

# Usage 
To use download [dist/go-make.exe](dist/go-make.exe) and save it anywhere on your computer, then add the path you saved it in to your `PATH` environmental variable. 

# Basic Features 
Your make file can have: 
* Multiple sections, specified in any order that can call eachother 
* Global scoped and local scoped variables 
* Comments (no inline comments!)

Comming Soon: 
* Setting variables in command line arguments 
* Inline comments 



# The Makefile 
Below is an example of a simple go-make build script 
```makefile 
# This is a comment 
GLOBAL = I am a global variable 


section1: 
	echo $(GLOBAL)
	
all: 
	GLOBAL = I am a local variable overriding a global variable 
	echo $(GLOBAL)
	
	# Execute commands in the other section 
	section1 
```
Save this as `makefile` in your directory.

To execute this, simply use 
`go-make` in the command prompt in your working directory 


You will see the following output: 
```
echo I am a local variable overriding a global variable
        I am a local variable overriding a global variable

echo I am a global variable
        I am a global variable

Make finished
```

## Output formatting 
The output for go-make is formatted in the following way: 
```
Command 1 
	Command 1 Output 
.
.
.
Command n 
	Command n Output 

Make Finished 
```


## Command Line Parameters 
Go-make currently supports two command line arguments: 
```
-f filename # This specifies the filename of the makefile you want to run (default: makefile)

-s section # This specifies the section name to run (default: all)
```

# peepingJim

A multi-core multi-threaded take on peepingTom from https://bitbucket.org/LaNMaSteR53/peepingtom

I was having issues with the XML parsing so I wrote this instead

## Pre-package

You will need phantomjs and capture.js in the same directory as peepingJim

## Build Executable

* Linux Install
  * make clean && make linux
* OSX Install
  * make clean && make osx

## Examples

### Running from executable

```
peepingJim 3.1.1 by James Cook <@_jbcook>
Usage:
  -dir string
    	dir of xml files
  -list string
    	file that contains a list of URLs
  -output string
    	where to write folder
  -threads int
    	Number of Threads to use (default 1)
  -timeout int
    	time out in seconds (default 8)
  -url string
    	single URL to scan
  -verbose
    	Verbose
  -version
    	Print version
  -xml string
    	xml file to parse
```
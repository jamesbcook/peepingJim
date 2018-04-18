# peepingJim [![Docker Build Status](https://img.shields.io/docker/build/b00stfr3ak/peepingjim.svg)](https://hub.docker.com/r/b00stfr3ak/peepingjim/)

A multi-core multi-threaded take on peepingTom from https://bitbucket.org/LaNMaSteR53/peepingtom

I was having issues with the XML parsing so I wrote this instead

## Pre-package

You will need chrome installed to take screenshots

## Build Executable

* Linux Install
  * make linux
* OSX Install
  * make osx

## Examples

### Running from executable

```
peepingJim 4.0.0 by James Cook <@_jbcook>
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
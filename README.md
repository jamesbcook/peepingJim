#### peepingJim


A multi-core multi-threaded take on peepingTom from https://bitbucket.org/LaNMaSteR53/peepingtom

I was having issues with the XML parsing so I wrote this instead


## Pre-package

You will need phantomjs in the same directory as peepingJim

## Install NMAP Parsing Lib

    go get github.com/lair-framework/go-nmap

## Build Executable

   go build -o peepingJim peepingJim.go

## Examples

# Running from executable

      ./peepingJim -h
      Usage of ./peepingJim:
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
        -xml string
              xml file to parse

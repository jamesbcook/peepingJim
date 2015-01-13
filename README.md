#### peepingJim


A multi-core multi-threaded take on peepingTom from https://bitbucket.org/LaNMaSteR53/peepingtom

I was having issues with the XML parsing so I wrote this instead


## Pre-package

You will need phantomjs in the same directory as peepingJim

## Install NMAP Parsing Lib

    go get github.com/b00stfr3ak/nmap

## Build Executable

   go build -o peepingJim peepingJim.go

## Examples


# Running from source

    ➜  peepingJim  go run peepingJim.go -help
       Usage of peepingJim:
         -cores=1: Number of Cores to use
         -dir="": dir of xml files
         -list="": file that contains a list of URLs
         -output="": where to write folder
         -threads=1: Number of Threads to use
         -timeout=8: time out in seconds
         -url="": single URL to scan
         -verbose=0: Verbose level 0,1,2
         -xml="": xml file to parse

# Running from executable

    ➜  peepingJim  ./peepingJim.go -help
       Usage of peepingJim:
         -cores=1: Number of Cores to use
         -dir="": dir of xml files
         -list="": file that contains a list of URLs
         -output="": where to write folder
         -threads=1: Number of Threads to use
         -timeout=8: time out in seconds
         -url="": single URL to scan
         -verbose=0: Verbose level 0,1,2
         -xml="": xml file to parse


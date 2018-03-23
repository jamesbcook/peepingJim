package peepingJim

import (
	"log"
	"os"
	"sync"
)

var (
	requiredFiles = [2]string{"phantomjs", "capture.js"}
	//Version of the package
	Version string
	//Author of the package
	Author string
)

//App settings
type App struct {
	InputType string
	Threads   int
}

//Client info needed for the worker
type Client struct {
	Output    string
	TimeOut   int
	PhantomJS Opts
	Sync      sync.RWMutex
	Verbose   bool
}

func init() {
	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			log.Fatal(file, " was not found in this directory")
		}
	}
}

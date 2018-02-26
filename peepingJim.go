package peepingJim

import (
	"log"
	"os"
	"sync"
)

const (
	version = "3.0.0"
	author  = "James Cook <@_jbcook>"
)

var (
	requiredFiles = [2]string{"phantomjs", "capture.js"}
)

//App settings
type App struct {
	InputType string
	Threads   int
	Verbose   int
}

//Client info needed for the worker
type Client struct {
	Output    string
	TimeOut   int
	PhantomJS Opts
	Sync      sync.RWMutex
}

func init() {
	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			log.Fatal(file, " was not found in this directory")
		}
	}
}

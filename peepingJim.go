package peepingJim

import (
	"sync"
)

var (
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
	Output  string
	TimeOut int
	Sync    sync.RWMutex
	Verbose bool
	Chrome
}

func init() {
	LocateChrome()
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/jamesbcook/peepingJim"
)

const (
	version = "3.1.0"
	author  = "James Cook <@_jbcook>"
)

//flagOpts hold all the possible options a user could pass at the cli
type flagOpts struct {
	url     string
	dir     string
	xml     string
	list    string
	output  string
	threads int
	timeout int
	verbose int
}

//flags is a function that builds the flagOpts struct
func flags() *flagOpts {
	xmlOpt := flag.String("xml", "", "xml file to parse")
	listOpt := flag.String("list", "", "file that contains a list of URLs")
	dirOpt := flag.String("dir", "", "dir of xml files")
	urlOpt := flag.String("url", "", "single URL to scan")
	threadOpt := flag.Int("threads", 1, "Number of Threads to use")
	outputOpt := flag.String("output", "", "where to write folder")
	timeoutOpt := flag.Int("timeout", 8, "time out in seconds")
	verboseOpt := flag.Int("verbose", 0, "Verbose level 0,1,2")
	flag.Parse()
	return &flagOpts{url: *urlOpt, dir: *dirOpt, xml: *xmlOpt, list: *listOpt,
		output: *outputOpt, threads: *threadOpt, timeout: *timeoutOpt,
		verbose: *verboseOpt}
}

func main() {
	//Gather all the cli arguments
	options := flags()
	var dstPath string
	//Creating Directory to store all output from phantom and curl
	if options.output != "" {
		if _, err := os.Stat(options.output); err == nil {
			log.Fatal(options.output + " already exists")
		} else {
			dstPath = options.output
		}
	} else {
		dstPath = "peepingJim_" + time.Now().Format("2006_01_02_15_04_05")
	}
	var targets []string
	if options.xml != "" {
		targets = peepingJim.GetTargets(peepingJim.InputType(peepingJim.XML), options.xml)
	} else if options.list != "" {
		targets = peepingJim.GetTargets(peepingJim.InputType(peepingJim.List), options.list)
	} else if options.dir != "" {
		targets = peepingJim.GetTargets(peepingJim.InputType(peepingJim.Dir), options.dir)
	} else if options.url != "" {
		targets = peepingJim.GetTargets(peepingJim.InputType(peepingJim.Plane), options.url)
	} else {
		log.Fatal("Need an input source")
	}
	app := peepingJim.App{}
	client := peepingJim.Client{}
	app.Threads = options.threads
	client.Output = dstPath
	client.TimeOut = options.timeout
	client.PhantomJS = peepingJim.RunPhantom()
	os.Mkdir(dstPath, 0755)
	//Making a list of targets to scan
	db := []map[string]string{}
	//Report name
	report := "peepingJim.html"
	outFile := fmt.Sprintf("%s/%s", dstPath, report)
	fmt.Printf("Loading %d targets\n", len(targets))
	// capture ctrl+c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("captured %v, stopping scanner and exiting...", sig)
			peepingJim.BuildReport(db, outFile)
			os.Exit(1)
		}
	}()
	queue := make(chan string)
	//spawn workers
	for i := 0; i <= app.Threads; i++ {
		go client.Worker(queue, &db)
	}
	//make work
	for _, target := range targets {
		queue <- target
	}
	//fill queue with finished work
	for n := 0; n <= app.Threads; n++ {
		queue <- ""
	}
	//Bulding the final html file
	peepingJim.BuildReport(db, outFile)
}

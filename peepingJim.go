package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/b00stfr3ak/nmap"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

//Enumerate Ports to see if they are valid or not and if they are HTTP or HTTPS
func enumPort(port *nmap.Port) (string, string) {
	httpPorts := []string{"80", "81", "8000", "8080", "8081", "8082"}
	httpsPorts := []string{"443", "8443"}
	// TODO
	// also check for http or https
	if port.State.Status == "open" {
		for _, value := range httpsPorts {
			if value == port.ID {
				return "https", port.ID
			}
		}
		for _, value := range httpPorts {
			if value == port.ID {
				return "http", port.ID
			}
		}
	}
	return "", ""
}

//parseNmap takes an array of structs from the imported nmap lib and
//builds a list of targets
func parseNmap(res []*nmap.ReportHost) []string {
	targets := []string{}
	for _, host := range res {
		for _, port := range host.Ports.Port {
			proto, portID := enumPort(&port)
			if portID != "" {
				var hostName string
				if host.Host.Names.Name == "" {
					hostName = host.Address.Addr
				} else {
					hostName = host.Host.Names.Name
				}
				url := fmt.Sprintf("%s://%s:%s", proto, hostName, portID)
				targets = append(targets, url)
			}
		}

	}
	return targets
}

//getTargets takes the pointer to the flagOpts struct and either
//makes a target list off one url, a list of URL's from a file,
//from an xml file or a dir of xml files
func getTargets(options *flagOpts) []string {
	var targets []string
	if options.url != "" {
		targets = append(targets, options.url)
	} else if options.xml != "" {
		res, _ := nmap.Parse(options.xml)
		targets = parseNmap(res)
	} else if options.list != "" {
		file, err := os.Open(options.list)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			targets = append(targets, scanner.Text())
		}
	} else if options.dir != "" {
		files, _ := filepath.Glob(options.dir + "/*.xml")
		for _, file := range files {
			res, _ := nmap.Parse(file)
			targets = append(targets, parseNmap(res)...)
		}
	} else {
		log.Fatal("No Targets were given")
	}
	return targets
}

//runPhantom sets up runCommand to run the phantom binary with all the options
func runPhantom(url, imgPath string, timeout int) string {
	phantomCMD := fmt.Sprintf("--ignore-ssl-errors=yes capture.js %s %s %d", url, imgPath, timeout*1000)
	opts := strings.Fields(phantomCMD)
	return runCommand("./phantomjs", opts)
}

//getHeader sets up runCommand to run the phantom binary with all the options
func getHeader(url, srcpath string, timeout int, c chan string) {
	curlOpts := fmt.Sprintf("-sLkD - %s -o %s --max-time %d", url, srcpath, timeout)
	opts := strings.Fields(curlOpts)
	c <- runCommand("curl", opts)
}

//runCommand takes a binary and it's ops and runs them
func runCommand(bin string, opts []string) string {
	cmd := exec.Command(bin, opts...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	return out.String()
}

//buildReport takes a hashmap and builds an html file that will be written
//to the file system
func buildReport(db []map[string]string, outFile string) {
	var liveMarkup string
	for _, d := range db {
		liveMarkup += fmt.Sprintf(`<tr><td class='img'><a href='%s' target='_blank'><img src='%s' onerror="this.parentNode.parentNode.innerHTML='No image available.';" /></a></td><td class='head'><a href='%s' target='_blank'>%s</a> (<a href='%s' target='_blank'>source</a>)<br /><pre>%s</pre></td></tr>`, d["imgPath"], d["imgPath"], d["url"], d["url"], d["srcPath"], d["headers"])
	}
	htmlBody := fmt.Sprintf(`<!doctype html>
<head>
<style>
table, td, th {border: 1px solid black;border-collapse: collapse;padding: 5px;font-size: .9em;font-family: tahoma;}
table {width: 100%%;table-layout: fixed;min-width: 1000px;}
td.img {width: 40%%;}
img {width: 100%%;}
td.head {vertical-align: top;word-wrap: break-word;}
</style>
</head>
<body>
<table>
%s
</table>
</body>
</html>`, liveMarkup)
	file, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString(htmlBody)
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
}

//Making a regex to later remove :// and : from a URL
var reg = regexp.MustCompile("(://)|(:)")

func worker(id int, queue chan string, options *flagOpts, dstPath string, db *[]map[string]string) {
	for {
		target := <-queue
		if target == "" {
			break
		}
		fmt.Printf("Scanning %s\n", target)
		//Cleaning URL so we can write to a file
		targetFixed := reg.ReplaceAllString(target, "")
		imgName := fmt.Sprintf("%s.png", targetFixed)
		srcName := fmt.Sprintf("%s.txt", targetFixed)
		imgPath := fmt.Sprintf("%s/%s", dstPath, imgName)
		srcPath := fmt.Sprintf("%s/%s", dstPath, srcName)
		//Making a channel to store curl output to
		c := make(chan string)
		go getHeader(target, srcPath, options.timeout, c)
		runPhantom(target, imgPath, options.timeout)
		//Writing output to a hash map and appending it to an array
		targetData := make(map[string]string)
		targetData["url"] = target
		targetData["imgPath"] = imgName
		targetData["srcPath"] = srcName
		targetData["headers"] = <-c
		*db = append(*db, targetData)
	}
}

//flagOpts hold all the possible options a user could pass at the cli
type flagOpts struct {
	url     string
	dir     string
	xml     string
	list    string
	output  string
	threads int
	timeout int
	cores   int
	verbose int
}

//flags is a function that builds the flagOpts struct
func flags() *flagOpts {
	xmlOpt := flag.String("xml", "", "xml file to parse")
	listOpt := flag.String("list", "", "file that contains a list of URLs")
	dirOpt := flag.String("dir", "", "dir of xml files")
	urlOpt := flag.String("url", "", "single URL to scan")
	coreOpt := flag.Int("cores", 1, "Number of Cores to use")
	threadOpt := flag.Int("threads", 1, "Number of Threads to use")
	outputOpt := flag.String("output", "", "where to write folder")
	timeoutOpt := flag.Int("timeout", 8, "time out in seconds")
	verboseOpt := flag.Int("verbose", 0, "Verbose level 0,1,2")
	flag.Parse()
	return &flagOpts{url: *urlOpt, dir: *dirOpt, xml: *xmlOpt, list: *listOpt,
		output: *outputOpt, cores: *coreOpt, threads: *threadOpt, timeout: *timeoutOpt,
		verbose: *verboseOpt}
}

//coreCheck takes a number and checks if you have that many cores or not.
//If you do it sets the max procs to that number
func coreCheck(cores int) {
	if cores > runtime.NumCPU() || cores <= 0 {
		log.Fatal(`You don't have that many cores... you can use up to `,
			runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(cores)
		if verbose == 1 {
			fmt.Println("Using " + string(runtime.GOMAXPROCS(cores)) + " cores")
		}
	}
}

var verbose int

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
		dstPath = "peepingGO_" + time.Now().Format("2006_01_02_15_04_05")
	}
	targets := getTargets(options)
	os.Mkdir(dstPath, 0755)
	//Checking the number of cores the user wants to use
	coreCheck(options.cores)
	verbose = options.verbose
	//Making a list of targets to scan
	db := []map[string]string{}
	//Report name
	report := "peepingGo.html"
	outFile := fmt.Sprintf("%s/%s", dstPath, report)
	fmt.Printf("Loading %d targets\n", len(targets))
	// capture ctrl+c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("captured %v, stopping scanner and exiting...", sig)
			buildReport(db, outFile)
			os.Exit(1)
		}
	}()
	threads := options.threads
	queue := make(chan string)
	//spawn workers
	for i := 0; i <= threads; i++ {
		go worker(i, queue, options, dstPath, &db)
	}
	//make work
	for _, target := range targets {
		queue <- target
	}
	//fill queue with finished work
	for n := 0; n <= threads; n++ {
		queue <- ""
	}
	//Bulding the final html file
	buildReport(db, outFile)
}

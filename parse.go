package peepingJim

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	nmap "github.com/lair-framework/go-nmap"
)

//Making a regex to later remove :// and : from a URL
var (
	reg     = regexp.MustCompile("(://)|(:)")
	httpReg = regexp.MustCompile("http")
)

//parseNmap takes an array of structs from the imported nmap lib and
//builds a list of targets
func parseNmap(res *nmap.NmapRun) []string {
	targets := []string{}
	var serviceName string
	for _, host := range res.Hosts {
		for _, port := range host.Ports {
			if port.State.State == "open" && httpReg.MatchString(port.Service.Name) {
				switch port.Service.Name {
				case "http":
					serviceName = "http"
				case "https":
					serviceName = "https"
				case "http-alt":
					serviceName = "http"
				case "https-alt":
					serviceName = "https"
				case "http-proxy":
					serviceName = "http"
				case "wbem-http":
					serviceName = "http"
				case "wbem-https":
					serviceName = "https"
				case "radan-http":
					serviceName = "http"
				}
				url := fmt.Sprintf("%s://%s:%d", serviceName, host.Addresses[0].Addr, port.PortId)
				targets = append(targets, url)
			}
		}
	}
	return targets
}

//InputType type match
type InputType func(string) []string

//XML parsing function
func XML(s string) []string {
	data, err := ioutil.ReadFile(s)
	if err != nil {
		log.Fatal("Couldn't Read File", err.Error())
	}
	res, _ := nmap.Parse(data)
	return parseNmap(res)
}

//List parsing function
func List(s string) []string {
	var targets []string
	file, err := os.Open(s)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		targets = append(targets, scanner.Text())
	}
	return targets
}

//Dir parsing function
func Dir(s string) []string {
	var targets []string
	files, _ := filepath.Glob(s + "/*.xml")
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal("Couldn't Read File", err.Error())
		}
		res, _ := nmap.Parse(data)
		targets = append(targets, parseNmap(res)...)
	}
	return targets
}

//Plane parsing function
func Plane(s string) []string {
	var targets []string
	targets = append(targets, s)
	return targets
}

//GetTargets takes the pointer to the flagOpts struct and either
//makes a target list off one url, a list of URL's from a file,
//from an xml file or a dir of xml files
func GetTargets(kind InputType, source string) []string {
	return kind(source)
}

package peepingJim

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const (
	phantomJS  = "./phantomjs"
	phantomCMD = "--ignore-ssl-errors=yes capture.js %s %s %d"
)

//Opts is a function that can pass options to a format string
//to be used on the cli
type Opts func(opts interface{}) string

type phantomStruct struct {
	url, imgPath string
	timeOut      int
}

func phantomOpts(opts interface{}) string {
	return fmt.Sprintf(phantomCMD, opts.(phantomStruct).url,
		opts.(phantomStruct).imgPath, opts.(phantomStruct).timeOut)
}

//RunPhantom returns a function that can be used to to pass cli options
func RunPhantom() Opts {
	return runCommand(phantomJS, Opts(phantomOpts))
}

//runCommand takes a binary and it's ops and runs them
func runCommand(bin string, fs func(interface{}) string) func(options interface{}) string {
	return func(options interface{}) string {
		allOpts := strings.Fields(fs(options))
		cmd := exec.Command(bin, allOpts...)
		var out, err bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &err
		cmd.Run()
		if err.Len() > 0 {
			log.Println(err.String())
		}
		return out.String()
	}
}

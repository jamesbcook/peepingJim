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

type phantomStruct struct {
	url, imgPath string
	timeOut      int
}

type opts func(opts interface{}) string
type fs func(opts interface{}) string

func phantomOpts(opts interface{}) string {
	return fmt.Sprintf(phantomCMD, opts.(phantomStruct).url,
		opts.(phantomStruct).imgPath, opts.(phantomStruct).timeOut)
}

func runPhantom() opts {
	return runCommand(phantomJS, fs(phantomOpts))
}

//runCommand takes a binary and it's ops and runs them
func runCommand(bin string, fs func(interface{}) string) func(opts interface{}) string {
	return func(opts interface{}) string {
		allOpts := strings.Fields(fs(opts))
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

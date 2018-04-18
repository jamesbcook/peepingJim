package peepingJim

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

//Worker thread performs all the scanning
func (client *Client) Worker(queue chan string, db *[]map[string]string) {
	for {
		target := <-queue
		if target == "" {
			break
		}
		if client.Verbose {
			fmt.Printf("Scanning %s\n", target)
		} else {
			fmt.Printf(".")
		}
		//Cleaning URL so we can write to a file
		targetFixed := reg.ReplaceAllString(target, "")
		targetFixed = strings.TrimSuffix(targetFixed, "/")
		u, err := url.Parse(target)
		if err != nil {
			log.Println(err)
			continue
		}
		imgName := fmt.Sprintf("%s.png", targetFixed)
		srcName := fmt.Sprintf("%s.txt", targetFixed)
		imgPath := fmt.Sprintf("%s/%s", client.Output, imgName)
		srcPath := fmt.Sprintf("%s/%s", client.Output, srcName)
		//Making a channel to store curl output to
		c := make(chan string)
		go getHeader(u, srcPath, client.TimeOut, c)
		if err := client.takeScreenshot(u, imgPath); err != nil {
			log.Println(err)
		}
		//Writing output to a hash map and appending it to an array
		targetData := make(map[string]string)
		targetData["url"] = target
		targetData["imgPath"] = imgName
		targetData["srcPath"] = srcName
		targetData["headers"] = <-c
		client.Sync.Lock()
		*db = append(*db, targetData)
		client.Sync.Unlock()
	}
}

package peepingJim

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

//getHeader returns the header and body of a page, while following any redirect
func getHeader(url, srcpath string, timeout int, c chan string) {
	var headers []string
	for {
		client := http.Client{
			Timeout: time.Duration(timeout) * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		resp, err := client.Get(url)
		if err != nil {
			log.Println(err)
			break
		}
		defer resp.Body.Close()
		if resp.StatusCode == 302 || resp.StatusCode == 301 {
			header, err := httputil.DumpResponse(resp, false)
			if err != nil {
				log.Println(err)
				break
			}
			headers = append(headers, string(header))
			var tmpURL string
			tmpURL = resp.Header.Get("Location")
			if tmpURL[0] == '/' {
				url = strings.TrimSuffix(url, "/") + tmpURL
			} else {
				url = tmpURL
			}
		} else {
			header, err := httputil.DumpResponse(resp, false)
			if err != nil {
				log.Println(err)
				break
			}
			headers = append(headers, string(header))
			var body bytes.Buffer
			body.ReadFrom(resp.Body)
			err = ioutil.WriteFile(srcpath, body.Bytes(), 0755)
			if err != nil {
				log.Println(err)
				break
			}
			break
		}
	}
	c <- strings.Join(headers, "\n")
}

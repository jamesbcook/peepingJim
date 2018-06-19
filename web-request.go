package peepingJim

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

const (
	redirectLimit = 3
)

//getHeader returns the header and body of a page, while following any redirect
func getHeader(u *url.URL, srcpath string, timeout int, c chan string) {
	var headers []string
	redirects := 0
	targetURL := u.String()
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
		req, err := http.NewRequest("GET", targetURL, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("User-Agent", UserAgent)
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			break
		}
		defer resp.Body.Close()
		if (resp.StatusCode == 302 || resp.StatusCode == 301) && (redirects < redirectLimit) {
			header, err := httputil.DumpResponse(resp, false)
			if err != nil {
				log.Println(err)
				break
			}
			headers = append(headers, string(header))
			var tmpURL string
			tmpURL = resp.Header.Get("Location")
			targetURL = strings.TrimSuffix(tmpURL, "/")
			redirects++
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

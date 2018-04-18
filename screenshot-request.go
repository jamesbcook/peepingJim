package peepingJim

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"time"
)

const (
	resolution    = "1200,800"
	chromeTimeout = 90
)

var (
	paths = []string{
		"/usr/bin/chromium",
		"/usr/bin/chromium-browser",
		"/usr/bin/google-chrome-stable",
		"/usr/bin/google-chrome",
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
		"C:/Program Files (x86)/Google/Chrome/Application/chrome.exe",
	}
	proxyPort int
)

//Chrome data struct
type Chrome struct {
	Path       string
	Resolution string
	UserAgent  string
}

//LocateChrome on the system
func LocateChrome() string {
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		return path
	}
	log.Fatal("Could not find chrome")
	return ""
}

func (client *Client) takeScreenshot(u *url.URL, output string) error {
	basicArguments := []string{
		"--headless", "--disable-gpu", "--hide-scrollbars",
		"--disable-crash-reporter",
		"--user-agent=" + UserAgent,
		"--window-size=" + resolution, "--screenshot=" + output,
	}
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("Error getting the current user object %v", err)
	}
	if currentUser.Uid == "0" {
		basicArguments = append(basicArguments, "--no-sandbox")
	}
	if u.Scheme == "https" {
		proxySetup(u)
		proxyURL, err := url.Parse("http://localhost:" + strconv.Itoa(proxyPort) + "/")
		if err != nil {
			log.Fatal(err)
		}
		proxyURL.Path = u.Path
		basicArguments = append(basicArguments, "--allow-insecure-localhost")
		basicArguments = append(basicArguments, proxyURL.String())
	} else {
		basicArguments = append(basicArguments, u.String())
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(chromeTimeout*time.Second))
	defer cancel()

	cmd := exec.CommandContext(ctx, client.Chrome.Path, basicArguments...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("[Error] starting the chrome command %v", err)
	}
	if err := cmd.Wait(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("[Error] Context time out")
		}
	}
	return nil
}

func proxySetup(u *url.URL) {
	u.Path = "/"
	rp := httputil.NewSingleHostReverseProxy(u)
	rp.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	proxyPort = listener.Addr().(*net.TCPAddr).Port
	go func() {
		httpServer := http.NewServeMux()
		httpServer.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			r.Host = u.Host
			rp.ServeHTTP(w, r)
		})
		if err := http.Serve(listener, httpServer); err != nil {
			log.Printf("Error serving reverse proxy %v", err)
			return
		}
	}()
}

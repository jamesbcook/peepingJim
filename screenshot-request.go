package peepingJim

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"time"
)

const (
	resolution    = "1200,800"
	userAgent     = "peeingJim/4.0.0"
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

func (client *Client) takeScreenshot(url, output string) error {
	basicArguments := []string{
		"--headless", "--disable-gpu", "--hide-scrollbars",
		"--disable-crash-reporter",
		"--user-agent=" + userAgent,
		"--window-size=" + resolution, "--screenshot=" + output,
	}
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("Error getting the current user object %v", err)
	}
	if currentUser.Uid == "0" {
		basicArguments = append(basicArguments, "--no-sandbox")
	}
	basicArguments = append(basicArguments, url)
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

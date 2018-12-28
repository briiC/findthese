package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
)

var appname = "findthese"
var version = "v0.1"

// flags
var argSourcePath string
var argEndpoint string
var argMethod = "HEAD"                                                              // assigned default value
var argOutput = "./findthese.report"                                                // assigned default value
var argDelay = 250                                                                  // assigned default value
var argDepth = 0                                                                    // assigned default value
var argSkip = []string{"jquery", "css", "img", "images", "i18n", "po"}              // assigned default value
var argSkipExts = []string{".png", ".jpeg", "jpg", "Gif", ".CSS", ".less", ".sass"} // assigned default value
var argSkipCodes = []string{"404"}                                                  // assigned default value
var argSkipSizes = []string{}                                                       // assigned default value
var argDirOnly = false                                                              // assigned default value

// asterisk "*" replaced by filename
// if no asterisk found treat as suffix
var argBackups = []string{"~", ".swp", ".swo", ".tmp", ".dmp", ".TMP", ".bkp", ".backup", ".bak", ".old", "_*", "~*"} // assigned default value

// Walk mode. Before real check/fetch count ETA
const walkModeCount = 0
const walkModeProcess = 1

var walkMode = walkModeCount
var dirItemCount = 0
var totalScanCount = 0

func main() {
	parseArgs()

	// TODO: Count items in source path folder and calc ~ETA
	walkMode = walkModeCount
	filepath.Walk(argSourcePath, localFileVisit)
	// durETA := time.Duration(totalScanCount*(argDelay+200)) * time.Millisecond
	printUsedArgs()

	// Walk local source directory
	log.Printf("(START) -- (%d items + %d mutations)", dirItemCount, totalScanCount)
	walkMode = walkModeProcess
	fmt.Println(strings.Repeat("-", 80)) // cleans \r
	if err := filepath.Walk(argSourcePath, localFileVisit); err != nil {
		fmt.Printf("ERR: Local directory: %v\n", err)
	}
	fmt.Println(strings.Repeat("-", 80)) // cleans \r
	log.Printf("(END)")

}

// callback
func localFileVisit(fpath string, f os.FileInfo, err error) error {
	fpath = strings.TrimPrefix(fpath, argSourcePath) // without local directory path
	depth := strings.Count(fpath, "/") + 1

	if fpath == "" {
		return nil
	}

	//  skip file if allowed to scan only directories
	if argDirOnly && !f.IsDir() {
		return nil
	}

	// Skip predefined dirs
	if f.IsDir() {

		// Skip by allowed depth
		if argDepth > 0 && depth > argDepth {
			return filepath.SkipDir
		}

		if inSlice(f.Name(), []string{".", "..", ".hg", ".git"}) {
			// fmt.Printf("-- SKIP ALWAYS [%s] --", f.Name())
			return filepath.SkipDir
		}
	}

	// Skip by name
	if inSlice(f.Name(), argSkip) {
		if f.IsDir() {
			return filepath.SkipDir // to skip whole tree
		}
		return nil // skip one item
	}

	// Skip by file extension
	if !f.IsDir() {
		ext := strings.ToLower(filepath.Ext(fpath))
		if inSlice(ext, argSkipExts) {
			// fmt.Printf("-- SKIP [%s] --", ext)
			return nil
		}
	}

	// counting mode
	if walkMode == walkModeCount {
		dirItemCount++
		totalScanCount += len(argBackups) - 1
		return nil
	}

	// generate mutations fpath list based on given fpath
	var fpaths []string
	fpaths = filePathMutations(fpath, argBackups)

	// Loop throw all fpath versions
	cleanupLen := 0 // cleaning current line with previous line length
	for _, fpath := range fpaths {
		fullURL := argEndpoint + fpath
		fname := filepath.Base(fpath)

		// Delay after basic checks and right before call
		time.Sleep(time.Duration(argDelay) * time.Millisecond)

		// Fetch
		resp, err := fetchURL(argMethod, fullURL)
		if err != nil {
			color.Red("ERR: %v", err)
			fmt.Println()
			continue
		}

		sCode := fmt.Sprintf("%d", resp.StatusCode)

		// try to read real body length if -1 found
		if resp.ContentLength == -1 {
			buf, _ := ioutil.ReadAll(resp.Body)
			resp.ContentLength = int64(len(buf))
		}
		sLength := fmt.Sprintf("%d", resp.ContentLength)

		sMore := "" // add at the end of line
		switch {

		case inSlice(sCode, argSkipCodes) || inSlice(sLength, argSkipSizes):
			// do not print out
			fmt.Printf("\r")
			fmt.Printf(strings.Repeat(" ", cleanupLen)) // cleaning
			fmt.Printf("\r")
			sLine := fmt.Sprintf("-> %s%s \tCODE:%s SIZE:%s ", color.MagentaString(argEndpoint), fpath, sCode, sLength)
			cleanupLen = len(sLine)
			fmt.Printf(sLine)
			// fmt.Printf(strings.Repeat(" ", 40)) // cleaning
			fmt.Printf("\r")
			continue

		case sCode == "200":
			sCode = color.GreenString(sCode)
			sMore += color.GreenString(fullURL)

		case sCode[:1] == "3": // 3xx codes
			sCode = color.CyanString(sCode)
			sMore += color.CyanString(fullURL)

		case sCode[:1] == "4": // 4xx codes
			sCode = color.RedString(sCode)
			sMore += color.RedString(fullURL)

		case sCode[:1] == "5": // 5xx codes
			sCode = color.BlueString(sCode)
			sMore += color.BlueString(fullURL)
		}

		fmt.Printf("depth=%d %20s | %-7s ", depth, fname, argMethod)
		fmt.Printf("CODE:%-4s SIZE:%-10s %-10s", sCode, sLength, sMore)
		fmt.Println()

	}

	return nil
}

// Fetches url content to dataTarget
func fetchURL(method, URL string) (*http.Response, error) {
	client := requestClient(URL)

	// Request
	req, _ := http.NewRequest(method, URL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.75 Safari/537.36")

	// Make request
	resp, reqErr := client.Do(req)
	if reqErr != nil {
		log.Printf("ERROR: [FETCH] %s -- %v", URL, reqErr)
		return nil, reqErr
	}

	return resp, nil
}

// Common request http client for data fetch
func requestClient(URL string) *http.Client {
	u, _ := url.Parse(URL)

	tr := &http.Transport{}

	if u.Scheme == "https" {
		tr.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	client := &http.Client{
		Transport: tr,
		// Timeout:   time.Second * 5,
		// CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// 	return http.ErrNoLocation
		// },
	}

	return client
}

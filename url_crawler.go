
package main

import "fmt"
import "io/ioutil"
import "net/http"
import "regexp"

import log "github.com/dmuth/google-go-log4go"


/**
* Our response object.
*/
type Response struct {
	//
	// The URL we just crawled
	//
	Url string
	//
	// HTTP code
	//
	Code int
	//
	// The actual page content.
	//
	Body string
}

//
// Keep track of if we crawled hosts with specific URLs
//
var hostsCrawled map [string]map[string]bool


/**
* Spin up 1 or more goroutines to do crawling.
*
* @param {int} num_instances
* @returm {chan string, chan Response} Our channel to read URLs from, 
*	our channel to write responses to.
*/
func NewUrlCrawler(NumInstances uint) (in chan string, out chan Response) {

	hostsCrawled = make(map[string]map[string]bool)

	//
	// I haven't yet decided if I want a buffer for 
	//InBufferSize := 1000
	InBufferSize := 0
	//OutBufferSize := 1000
	OutBufferSize := 0
	in = make(chan string, InBufferSize)
	out = make(chan Response, OutBufferSize)

	for i:=uint(0); i< NumInstances; i++ {
		log.Infof("Spun up crawler instance #%d", (i+1))
		go crawl(in, out)
	}

	return in, out

} // End of NewUrlCrawler()


/**
* This is run as a goroutine which is responsible for doing the crawling and 
* returning the results.
*
* @param {chan string} in Our channel to read URLs to crawl from
* @param {chan Response} out Responses will be written on this channel
*
* @return {Response} A response consisting of our code and body
*/
func crawl(in chan string, out chan Response) {

	for {

		log.Debug("About to ingest a URL...")
		url := <-in

		//
		// Grab our URL parts
		//
		regex, _ := regexp.Compile("(https?://[^/]+)(.*)")
		results := regex.FindStringSubmatch(url)
		Host := results[1]
		Uri := results[2]

		//
		// Create our host entry if we don't already have it.
		//
		if _, ok := hostsCrawled[Host]; !ok {
			hostsCrawled[Host] = make(map[string]bool)
		}

		//
		// If this is our first time here, cool. Otherwise, skip.
		//
		if _, ok := hostsCrawled[Host][Uri]; !ok {
			hostsCrawled[Host][Uri] = true
		} else {
			log.Warnf("We've already been to '%s', skipping!", url)
			continue
		}

		log.Infof("About to crawl '%s'...", url)
		out <-httpGet(url)
		log.Infof("Done crawling '%s'!", url)

	}

} // End of crawl()


/**
* Retrieve a URL via HTTP GET.
*
* @param {string} url The URL to retrieve.
* @return {Response} A response consisting of our code and body
*/
func httpGet(url string) (retval Response) {

	retval.Url = url

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Warnf("Error fetching %s: %s", url, err)
		retval.Body = fmt.Sprintf("%s", err)
		retval.Code = 0
		return(retval)
	}

	req.Header.Set("User-Agent", 
		"Dmuth's crawler. Please report bugs to me: http://www.dmuth.org/contact")

	resp, err := client.Do(req)
	if err != nil {
		log.Warnf("Error fetching %s: %s", url, err)
		retval.Body = fmt.Sprintf("%s", err)
		retval.Code = 0
		return(retval)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("Error fetching %s: %s", url, err)
		retval.Body = fmt.Sprintf("%s", err)
		retval.Code = 0
		return(retval)
	}

	retval.Body = fmt.Sprintf("%s", body)
	retval.Code = resp.StatusCode

	return(retval)

} // End of httpGet()



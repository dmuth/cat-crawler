
package main

import "fmt"
import "io/ioutil"
import "net/http"
import "regexp"
import "strings"

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

		url = filterUrl(url)

		if (beenHere(url)) {
			log.Debugf("We've already been to '%s', skipping!", url)
			continue
		}

		if (!sanityCheck(url)) {
			//
			// In the future, I might make the in channel take a data 
			// structure which includes the referrer so I can dig 
			// into bad URLs. With a backhoe.
			//
			log.Warnf("URL '%s' fails sanity check, skipping!", url)
			continue
		}

		log.Infof("About to crawl '%s'...", url)
		out <-httpGet(url)
		log.Infof("Done crawling '%s'!", url)

	}

} // End of crawl()


/**
* Filter meaningless things out of URLs. Like hashmarks.
*
* @param {string} url The URL
*
* @return {string} The filtered URL
*/
func filterUrl(url string) (string) {

	//
	// First, nuke hashmarks (thanks, Apple!)
	//
	regex, _ := regexp.Compile("([^#]+)#")
	results := regex.FindStringSubmatch(url)
	if (len(results) >= 2) {
		url = results[1]
	}

	//
	// Replace groups of 2 or more slashes with a single slash (thanks, log4j!)
	//
	regex, _ = regexp.Compile("[^:](/[/]+)")
	for {
		results = regex.FindStringSubmatch(url)
		if (len(results) < 2) {
			break
		}

		Dir := results[1]
		//url = regex.ReplaceAllString(url, "/")
		url = strings.Replace(url, Dir, "/", -1)

	}

	//
	// Now, remove references to parent directories, because that's just 
	// ASKING for path loops. (thanks, Apple!)
	//
	// Do this by looping as long as we have ".." present.
	//
	regex, _ = regexp.Compile("([^/]+/\\.\\./)")
	for  {
		results = regex.FindStringSubmatch(url)
		if (len(results) < 2) {
			break
		}

		Dir := results[1]
		url = strings.Replace(url, Dir, "", -1)

	}

	//
	// Replace paths of single dots
	//
	regex, _ = regexp.Compile("/\\./")
	url = regex.ReplaceAllString(url, "/")

	return(url)

} // End of filterUrl()


/**
* Have we already been to this URL?
*
* @param {string} url The URL we want to crawl
*
* @return {bool} True if we've crawled this URL before, false if we have not.
*/
func beenHere(url string) (retval bool) {

	retval = true

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
		retval = false
	}

	return retval

} // End of beenHere()


/**
* Check to see if this URL is sane.
*
* @return {bool} True if the URL looks okay, false otherwise.
*/
func sanityCheck(url string) (retval bool) {

	retval = true

	regex, _ := regexp.Compile(" ")
	result := regex.FindString(url)

	if (result != "") {
		retval = false
	}

	return(retval)

} // End of sanityCheck()


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



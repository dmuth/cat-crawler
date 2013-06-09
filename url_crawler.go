
package main

import "fmt"
import "io/ioutil"
import "net/http"

import log "github.com/dmuth/google-go-log4go"


/**
* Our response object.
*/
type Response struct {
	//
	// HTTP code
	//
	Code int
	//
	// The actual page content.
	//
	Data string
}


/**
* Spin up 1 or more goroutines to do crawling.
*
* @param {chan string} in Our channel to read URLs to crawl from
* @param {chan Response} out Responses will be written on this channel
* @param {int} num_instances
*/
func NewUrlCrawler(in chan string, out chan Response,
	num_instances uint) () {

	for i:=uint(0); i< num_instances; i++ {
		log.Infof("Spun up crawler instance #%d", (i+1))
		go crawl(in, out)
	}

} // End of NewUrlCrawler()


/**
*
*/
func crawl(in chan string, out chan Response) {

	url := <-in
	log.Infof("About to crawl '%s'...", url)

	out <-httpGet(url)
	log.Infof("Done crawling '%s'!", url)

} // End of crawl()


/**
* Retrieve a URL via HTTP GET.
*
* @param {string} url The URL to retrieve.
*/
func httpGet(url string) (retval Response) {

	resp, err := http.Get(url)
	if (err != nil) {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if (err != nil) {
		return
	}

	retval.Data = fmt.Sprintf("%s", body)
	retval.Code = resp.StatusCode

	return(retval)

} // End of httpGet()



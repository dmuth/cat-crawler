
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
	Body string
}


/**
* Spin up 1 or more goroutines to do crawling.
*
* @param {int} num_instances
*/
func NewUrlCrawler(NumInstances uint) (in chan string, out chan Response) {

	in = make(chan string, NumInstances)
	out = make(chan Response)

	for i:=uint(0); i< NumInstances; i++ {
		log.Infof("Spun up crawler instance #%d", (i+1))
		go crawl(in, out)
	}

	return in, out

} // End of NewUrlCrawler()


/**
* This is run as a goroutine which is responsible for doing the carwling and 
* returning the results.
*
* @param {chan string} in Our channel to read URLs to crawl from
* @param {chan Response} out Responses will be written on this channel
*
* @return {Response} A response consisting of our code and body
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
* @return {Response} A response consisting of our code and body
*/
func httpGet(url string) (retval Response) {

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



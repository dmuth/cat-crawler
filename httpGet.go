
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




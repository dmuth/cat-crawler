
package main

//import "fmt"
import "regexp"
import "os"

import log "github.com/dmuth/google-go-log4go"


//
// Keep track of if we crawled hosts with specific URLs
//
var hostsCrawledImages map [string]map[string]bool


/**
* Fire up 1 or more crawlers to start grabbing images.
*
* @param {chan} Image in Image data structures will be read from here.
* @param {uint} NumConnections How many go threads to fire up?
*
*/
func NewImageCrawler(in chan Image, NumConnections uint) {

	hostsCrawledImages = make(map[string]map[string]bool)
	for i:=0; i<int(NumConnections); i++ {
		go crawlImages(in)
	}

} // End of NewImageCrawler()


/**
* Continuously read images and crawl them.
*/
func crawlImages(in chan Image) {

	for {
		image := <-in

		// src, alt, title
		Url := image.src

		if (Url == "") {
			continue
		}

		if (imageBeenHere(Url)) {
			log.Debugf("crawlImages(): We've already been to '%s', skipping!", Url)
			continue
		}

		response := httpGet(Url)
		log.Infof("Response code %d on URL '%s'", response.Code, response.Url)

		filename := getFilenameFromUrl(Url)

		writeImage(filename, response.Body)

	}

} // End of crawlImages()


/**
* Have we already been to this image?
*
* @param {string} url The URL we want to crawl
*
* @return {bool} True if we've crawled this URL before, false if we have not.
*/
func imageBeenHere(url string) (retval bool) {

	retval = true

	//
	// Grab our URL parts
	//
	results := getUrlParts(url)
	if (len(results) < 5) {
		log.Warnf("imageBeenHere(): Unable to parse URL: '%s'", url)
		return(true)
	}
	Host := results[1]
	Uri := results[4]

	//
	// Create our host entry if we don't already have it.
	//
	if _, ok := hostsCrawledImages[Host]; !ok {
		hostsCrawledImages[Host] = make(map[string]bool)
	}

	//
	// If this is our first time here, cool. Otherwise, skip.
	//
	if _, ok := hostsCrawledImages[Host][Uri]; !ok {
		hostsCrawledImages[Host][Uri] = true
		retval = false
	}

	return retval

} // End of imageBeenHere()


/**
* Convert our URL into a filename
*/
func getFilenameFromUrl(Url string) (retval string) {

	retval = Url

	results := getUrlParts(Url)
	Host := results[3]
	Uri := results[4]

	regex, _ := regexp.Compile("/$")
	Uri = regex.ReplaceAllLiteralString(Uri, "")

	retval = Host + Uri

	//
	// Trim the filename if it's too long.
	// This isn't a perfect fix, but it'll work for now.
	//
	MaxLen := 80
	if (len(retval) > MaxLen) {
		retval = retval[:(MaxLen - 1)]
	}

	return(retval)

} // End of getFilenameFromUrl()


/**
* Write our image out to disk.
*
* @param {string} Filename The name of the file
*
* @param {string} Body The content of the image
*/
func writeImage(Filename string, Body string) {

	cwd, _ := os.Getwd()

	target := cwd + "/downloads/" + Filename

	//
	// Try making our directory
	// We want to panic if there are any issues since it could 
	// mean awful things like a full disk!
	//
	result := os.MkdirAll(target, 0750)
	if (result != nil) {
		log.Errorf("%s", result)
		panic(result)
	}

	//
	// 
	//
// TEST

} // End of writeImage()




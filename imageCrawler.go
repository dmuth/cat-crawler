package main

//import "fmt"
import "regexp"
import "strings"
import "os"

import log "github.com/dmuth/google-go-log4go"
import stats "github.com/dmuth/cat-crawler/stats"

//
// Keep track of if we crawled hosts with specific URLs
//
var hostsCrawledImages map[string]map[string]bool

/**
* Fire up 1 or more crawlers to start grabbing images.
*
* @param {config} Our configuration
* @param {chan} Image in Image data structures will be read from here.
* @param {uint} NumConnections How many go threads to fire up?
*
 */
func NewImageCrawler(config Config, in chan Image, NumConnections uint) {

	hostsCrawledImages = make(map[string]map[string]bool)
	for i := 0; i < int(NumConnections); i++ {
		go crawlImages(config, in)
	}

} // End of NewImageCrawler()

/**
* Continuously read images and crawl them.
*
* @param {config} Our configuration
*
 */
func crawlImages(config Config, in chan Image) {

	for {
		stats.IncrStat("go_image_crawler_waiting")
		image := <-in
		stats.DecrStat("go_image_crawler_waiting")
		stats.DecrStat("images_to_be_crawled")

		// src, alt, title
		Url := image.src

		if Url == "" {
			continue
		}

		//
		// If we've been here before, stop
		//
		if imageBeenHereUrl(Url) {
			log.Debugf("crawlImages(): We've already been to '%s', skipping!", Url)
			stats.IncrStat("images_skipped")
			continue
		}
		setImageBeenHereUrl(Url)

		match := false
		if strings.Contains(strings.ToLower(image.alt), config.SearchString) {
			match = true
			log.Debugf("Match found on ALT tag for URL '%s'!", Url)
		}
		if strings.Contains(strings.ToLower(image.title), config.SearchString) {
			match = true
			log.Infof("Match found on TITLE tag for URL '%s'!", Url)
		}

		if !match {
			log.Debugf("No match for %s found in alt and title tags for URL '%s', stopping!", Url)
			stats.IncrStat("images_not_matched")
			continue
		}

		log.Infof("Image: About to crawl '%s'...", Url)
		response := httpGet(Url)
		stats.IncrStat("images_crawled")
		log.Infof("Image: Response code %d on URL '%s'", response.Code, response.Url)

		//
		// If the content-type isn't an image, stop.
		//
		regex, _ := regexp.Compile("^image")
		results := regex.FindString(response.ContentType)
		if len(results) == 0 {
			log.Errorf("Skipping Content-Type of '%s', on URL '%s'",
				response.ContentType, response.Url)
			continue
		}

		filename := getFilenameFromUrl(Url)

		writeImage(filename, response.Body)

	}

} // End of crawlImages()

/**
* Wrapper for imageBeenHere() which takes a URL
 */
func imageBeenHereUrl(url string) bool {

	//
	// Grab our URL parts
	//
	results := getUrlParts(url)
	if len(results) < 5 {
		log.Warnf("imageBeenHereUrl(): Unable to parse URL: '%s'", url)
		return (true)
	}
	host := results[1]
	uri := results[4]

	if imageBeenHere(host, uri) {
		return (true)
	}

	//
	// Assume false
	//
	return (false)

}

/**
* Wrapper for setImageBeenHere() which takes a URL.
 */
func setImageBeenHereUrl(url string) {

	//
	// Grab our URL parts
	//
	results := getUrlParts(url)
	host := results[1]
	uri := results[4]

	setImageBeenHere(host, uri)

}

/**
* Make the deterination if we've been to this image before.
*
* @param {string} host The hostname
* @param {string} uri The URI
*
* @return {bool} True if we've crawled this image before, false otherwise.
 */
func imageBeenHere(host string, uri string) bool {

	//
	// Create our host entry if we don't already have it.
	//
	if _, ok := hostsCrawledImages[host]; !ok {
		hostsCrawledImages[host] = make(map[string]bool)
	}

	//
	// See if we've been here before.
	//
	_, ok := hostsCrawledImages[host][uri]
	if ok {
		return (true)
	} else {
		return (false)
	}

} // End of imageBeenHere()

/**
* We've been to this image before!
 */
func setImageBeenHere(host string, uri string) {

	//
	// Create our host entry if we don't already have it.
	//
	if _, ok := hostsCrawledImages[host]; !ok {
		hostsCrawledImages[host] = make(map[string]bool)
	}

	hostsCrawledImages[host][uri] = true

}

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
	if len(retval) > MaxLen {
		retval = retval[:(MaxLen - 1)]
	}

	return (retval)

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

	//
	// Create our target and nuke the filename from the end
	//
	target := cwd + "/cat-crawler-downloads/" + Filename
	regex, _ := regexp.Compile("/[^/]+$")
	dir := regex.ReplaceAllLiteralString(target, "")

	//
	// Try making our directory
	// We want to panic if there are any issues since it could
	// mean awful things like a full disk!
	//
	result := os.MkdirAll(dir, 0750)
	if result != nil {
		log.Errorf("Error creating directory: %s", result)
		panic(result)
	}

	//
	// Now write the file.
	//
	file, err := os.Create(target)
	if err != nil {
		log.Warnf("Error opening file: %s", err)

	} else {
		n, err := file.Write([]byte(Body))
		if err != nil {
			log.Warnf("Error writing file: '%s': %s", target, err)
		}
		log.Infof("%d bytes written to file '%s'", n, target)

		err = file.Close()
		if err != nil {
			log.Errorf("Error closing file: %s", err)
			panic(err)
		}

	}

} // End of writeImage()

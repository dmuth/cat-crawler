
package main

import log "github.com/dmuth/google-go-log4go"


/**
* Fire up 1 or more crawlers to start grabbing images.
*
* @param {chan} Image in Image data structures will be read from here.
* @param {uint} NumConnections How many go threads to fire up?
*
*/
func NewImageCrawler(in chan Image, NumConnections uint) {

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

		//
		// No src URL?  Stop.
		//
		if (image.src == "") {
			continue
		}

		// src, alt, title
		response := httpGet(image.src)
		log.Infof("Response code %d on URL '%s'", response.Code, response.Url)

	}

} // End of crawlImages()



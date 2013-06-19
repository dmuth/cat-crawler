
package main

import "os"
import "os/signal"
import "syscall"


import log "github.com/dmuth/google-go-log4go"


func main() {

	//
	// Parse our arguments and report them
	//
	config := ParseArgs()
	log.Infof("Config: %s", config)
	log.Infof("SeedURLs: %s", config.SeedUrls)
	if (len(config.AllowUrls) > 0 ) {
		log.Infof("Only allowing URLs starting with: %s", config.AllowUrls)
	}

	//
	// Catch our interrupt signal
	//
	go sigInt()

	NumConnections := config.NumConnections

	//
	// Start the crawler and seed it with our very first URL
	//
	UrlCrawlerIn, UrlCrawlerOut := NewUrlCrawler(uint(NumConnections), config.AllowUrls)

	//UrlCrawlerIn <- "http://localhost:8080/" // Debugging
	for _, value := range config.SeedUrls {
		UrlCrawlerIn <- value
	}

	//
	// Create our HTML parser
	//
	HtmlBodyIn, ImageCrawlerIn := NewHtml(UrlCrawlerIn)

	//
	// Start up our image crawler
	//
	NewImageCrawler(ImageCrawlerIn, NumConnections)

	for {
		//
		// Read a result from our crawler
		//
		Res := <-UrlCrawlerOut

		if (Res.Code != 200) {
			log.Debugf("Skipping non-2xx response of %d on URL '%s'", 
				Res.Code, Res.Url)
			continue
		}

		//
		// Pass it into the HTML parser.  It will in turn send any URLs
		// it finds into the URL Crawler and any images to the Image Crawler.
		//
		HtmlBodyIn <- []string { Res.Url, Res.Body, Res.ContentType }

	}

} // End of main()


/**
* Wait for ctrl-c to happen, then exit!
*/
func sigInt() {
        ch := make(chan os.Signal)
        signal.Notify(ch, syscall.SIGINT)
        <-ch
        log.Error("CTRL-C; exiting")
        os.Exit(0)
}




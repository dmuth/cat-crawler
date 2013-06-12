
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
	log.Infof(
		"SeedURL: %s",
		config.SeedUrl)

	//
	// Catch our interrupt signal
	//
	go sigInt()

	//
	// Start the crawler and seed it with our very first URL
	//
	NumConnections := config.NumConnections
	UrlCrawlerIn, UrlCrawlerOut := NewUrlCrawler(uint(NumConnections))
	UrlCrawlerIn <- config.SeedUrl
	//UrlCrawlerIn <- "http://localhost:8080/" // Debugging

	//
	// Create our HTML parser
	//
	HtmlBodyIn, ImageCrawlerIn := NewHtml(UrlCrawlerIn)
log.Infof("%s", ImageCrawlerIn)

	//
	// TODO: Start the image crawler here
	// NewImageCrawler(ImageCrawlerIn)

	for {
		//
		// Read a result from our crawler
		//
		Res := <-UrlCrawlerOut

		//
		// Pass it into the HTML parser.  It will in turn send any URLs
		// it finds into the URL Crawler and any images to the Image Crawler.
		//
		HtmlBodyIn <- []string { Res.Url, Res.Body }

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




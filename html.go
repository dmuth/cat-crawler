package main

//import "fmt"
import "regexp"

import log "github.com/dmuth/google-go-log4go"
import stats "github.com/dmuth/cat-crawler/stats"

/**
* Representation of our parsed Html
 */
type Image struct {
	html  string
	src   string
	alt   string
	title string
}
type HtmlParsed struct {
	links  []string
	images []Image
}

/**
* Set up our parser to run in the background.
*
* @param {chan string} UrlCrawlerIn URLs written to this will be sent
*	off to the URL crawler.
*
* @return {chan string} A channel which will be used to ingest HTML for parsing.
 */
func NewHtml(UrlCrawlerIn chan string) (
	chan []string, chan Image) {

	HtmlCrawlerIn := make(chan []string)

	BufferSize := 1000
	//BufferSize = 1 // Debugging
	ImageCrawlerIn := make(chan Image, BufferSize)

	go HtmlParseWorker(HtmlCrawlerIn, UrlCrawlerIn, ImageCrawlerIn)

	return HtmlCrawlerIn, ImageCrawlerIn

} // End of NewHtml()

/**
* This function is run as a goroutine and ingests Html to parse. It then
* sends off URLs and image URLs.
*
* @param {chan string} HtmlIn Incoming HTML
* @param {chan string} UrlCrawlerIn URLs written to this will be sent
*	off to the URL crawler.
* @param {chan string} ImageCrawlerOut URLs written to this will be sent
*	off to the image crawler.
*
 */
func HtmlParseWorker(HtmlIn chan []string, UrlCrawlerIn chan string,
	ImageCrawlerIn chan Image) {

	//
	// Loop through HTML and parse all the things.
	//
	for {
		in := <-HtmlIn
		BaseUrl := in[0]
		Html := in[1]
		Parsed := HtmlParseString(BaseUrl, Html)

		//
		// Put these into goroutines so that we can get back to parsing
		//
		go HtmlParseWorkerLinks(&Parsed, UrlCrawlerIn)
		go HtmlParseWorkerImages(&Parsed, ImageCrawlerIn)

	}

} // End of HtmlParseWorker()

/**
* Another goroutine that loops through our links and sends them off to UrlCrawler
*
* @param {HtmlParsed} Our parsed HTML elements.
* @param {chan string} The channel to send URLs to our URL crawler
*
 */
func HtmlParseWorkerLinks(Parsed *HtmlParsed, UrlCrawlerIn chan string) {

	for i := range Parsed.links {
		Row := Parsed.links[i]
		log.Debugf("Sending to UrlCrawler: %d: %s", i, Row)
stats.IncrStat("urls_to_be_crawled")
		UrlCrawlerIn <- Row
	}

} // End of HtmlParseWorkerLinks()

/**
* Another goroutine that loops through our images and sends them off to
* the ImageCrawler.
*
* @param {HtmlParsed} Our parsed HTML elements.
* @param {ImageCrawlerIn} The channel to send images to our image crawler
*
 */
func HtmlParseWorkerImages(Parsed *HtmlParsed, ImageCrawlerIn chan Image) {

	for i := range Parsed.images {
		Row := Parsed.images[i]
		log.Debugf("Sending to ImageCrawler: %d: %s", i, Row)
stats.IncrStat("images_to_be_crawled")
		ImageCrawlerIn <- Row
	}

} // End of HtmlParseWorkerImages()

/**
* Parse an HTML response and get links and images.
*
* @param {string} BaseUrl The URL of the page we got these links from
* @param {string} Body The body of the page
*
* @return {HtmlParsed} Structure of links and images
 */
func HtmlParseString(BaseUrl string, Body string) (retval HtmlParsed) {

	//
	// Break up our base URL into a host and URI
	//
	regex, _ := regexp.Compile("(https?://[^/]+)(.*)")
	results := regex.FindStringSubmatch(BaseUrl)
	BaseUrlHost := results[1]
	BaseUrlUri := results[2]

	retval.links = HtmlParseLinks(BaseUrlHost, BaseUrlUri, Body)
	retval.images = HtmlParseImages(BaseUrlHost, BaseUrlUri, Body)

	return (retval)

} // End of HtmlParseString()

/**
* Grab our links out of the body, and fully qualify them.
*
* @param {string} BaseUrlHost The http:// and hostname part of our base URL
* @param {string} BaseUrlUri Our base URI
* @param {string} Body The body of the webpage
*
* @return {[]string} Array of links
 */
func HtmlParseLinks(BaseUrlHost string, BaseUrlUri string, Body string) (retval []string) {

	//
	// Get all of our links
	//
	regex, _ := regexp.Compile("(?s)" +
		"href=\"" +
		"(" +
		"(https?://([^/]+))?" +
		"([^\"]+)" +
		")\"")
	results := regex.FindAllStringSubmatch(Body, -1)

	for i := range results {

		result := results[i]

		HostAndMethod := result[2]
		Uri := result[4]

		//
		// If a host and method is specified, just glue them back together.
		//
		Url := ""
		if HostAndMethod != "" {
			Url = HostAndMethod + Uri

		} else {
			//
			// Otherwise, it's on the same host. Determine if
			// it's a relative or absolute link.
			//
			FirstChar := string(Uri[0])
			if FirstChar == "/" {
				Url = BaseUrlHost + Uri
			} else {
				Url = BaseUrlHost + BaseUrlUri + "/" + Uri
			}

		}

		//fmt.Println("FINAL URL", Url)

		retval = append(retval, Url)

	}

	return (retval)

} // End of HtmlParseLinks()

/**
* Grab image links out of the body, and fully qualify them.
*
* @param {string} BaseUrlHost The http:// and hostname part of our base URL
* @param {string} BaseUrlUri Our base URI
* @param {string} Body The body of the webpage
*
* @return {[]Image} Array of images
 */
func HtmlParseImages(BaseUrlHost string, BaseUrlUri string, Body string) (retval []Image) {

	retval = htmlParseImageTags(Body)

	for i := range retval {
		htmlParseSrc(BaseUrlHost, BaseUrlUri, &retval[i])
		if retval[i].src != "" {
			htmlParseAlt(&retval[i])
			htmlParseTitle(&retval[i])
		}
	}

	return (retval)

} // End of HtmlParseImages()

/**
* Grab our image tags out of the body.
*
* @param {string} Body The HTML body
*
* @return {[]Image} Array of Image elements
 */
func htmlParseImageTags(Body string) (retval []Image) {

	regex, _ := regexp.Compile("(?s)" +
		"<img[^>]+>")
	results := regex.FindAllStringSubmatch(Body, -1)

	for i := range results {
		image := Image{results[i][0], "", "", ""}
		retval = append(retval, image)
	}

	return (retval)

} // End of htmlParseImageTags()

/**
* Parse the src tag out of our image.
*
* @param {string} BaseUrlHost The http:// and hostname part of our base URL
* @param {string} BaseUrlUri Our base URI
* @param {*Image} Pointer to our image structure
 */
func htmlParseSrc(BaseUrlHost string, BaseUrlUri string, image *Image) {

	regex, _ := regexp.Compile("(?s)" +
		"<img[^>]+src=\"" +
		"(" +
		"(https?://([^/]+))?" +
		"([^\"]+)" +
		")\"")
	result := regex.FindStringSubmatch(image.html)

	//
	// Bail out if we have no source
	//
	if len(result) == 0 {
		return
	}

	HostAndMethod := result[2]
	Uri := result[4]

	//
	// If a host and method is specified, just glue them back together.
	//
	Url := ""
	if HostAndMethod != "" {
		Url = HostAndMethod + Uri

	} else {
		//
		// Otherwise, it's on the same host. Determine if
		// it's a relative or absolute link.
		//
		FirstChar := string(Uri[0])
		if FirstChar == "/" {
			Url = BaseUrlHost + Uri
		} else {
			Url = BaseUrlHost + BaseUrlUri + "/" + Uri
		}

	}

	image.src = Url

} // End of htmlParseSrc()

/**
* Parse the alt tag out of our image.
 */
func htmlParseAlt(image *Image) {

	regex, _ := regexp.Compile("(?s)" +
		"<img[^>]+alt=\"([^\"]+)\"")
	result := regex.FindStringSubmatch(image.html)
	if len(result) > 1 {
		image.alt = result[1]
	}

} // End of htmlParseAlt()

/**
* Parse the title tag out of our image.
 */
func htmlParseTitle(image *Image) {

	regex, _ := regexp.Compile("(?s)" +
		"<img[^>]+title=\"([^\"]+)\"")
	result := regex.FindStringSubmatch(image.html)
	if len(result) > 1 {
		image.title = result[1]
	}

} // End of htmlParseTitle()

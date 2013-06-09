
package main

//import "fmt"
import "regexp"


/**
* Representation of our parsed Html
*/
type Image struct {
	html string
	src string
	alt string
	title string
}
type HtmlParsed struct {
	links []string
	images []Image
}


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

	return(retval)

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
		"("+
		"(https?://([^/]+))?" +
		"([^\"]+)" +
		")\"")
	results := regex.FindAllStringSubmatch(Body, -1)

	for i:= range results {

		result := results[i]

		HostAndMethod := result[2]
		Uri := result[4]

		//
		// If a host and method is specified, just glue them back together.
		//
		Url := ""
		if (HostAndMethod != "") {
			Url = HostAndMethod + Uri

		} else {
			//
			// Otherwise, it's on the same host. Determine if 
			// it's a relative or absolute link.
			//
			FirstChar := string(Uri[0])
			if (FirstChar == "/") {
				Url = BaseUrlHost + Uri
			} else {
				Url = BaseUrlHost + BaseUrlUri + "/" + Uri
			}

		}

		//fmt.Println("FINAL URL", Url)

		retval = append(retval, Url)

	}

	return(retval)

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

	for i:= range retval {
		htmlParseSrc(BaseUrlHost, BaseUrlUri, &retval[i])
		htmlParseAlt(&retval[i])
		htmlParseTitle(&retval[i])
	}

	return(retval)

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

	for i:= range results {
		image := Image{ results[i][0], "", "", "" }
		retval = append(retval, image)
	}

	return(retval)

} // End of htmlParseImagesTags()


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
		"("+
		"(https?://([^/]+))?" +
		"([^\"]+)" +
		")\"")
	result := regex.FindStringSubmatch(image.html)

	HostAndMethod := result[2]
	Uri := result[4]

	//
	// If a host and method is specified, just glue them back together.
	//
	Url := ""
	if (HostAndMethod != "") {
		Url = HostAndMethod + Uri

	} else {
		//
		// Otherwise, it's on the same host. Determine if 
		// it's a relative or absolute link.
		//
		FirstChar := string(Uri[0])
		if (FirstChar == "/") {
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
	if (len(result) > 1) {
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
	if (len(result) > 1) {
		image.title = result[1]
	}

} // End of htmlParseTitle()





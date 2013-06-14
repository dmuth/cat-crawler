
package main

import "flag"
//import "fmt"
import "regexp"
import "strings"
import "os"

import log "github.com/dmuth/google-go-log4go"


//
// Configuration for what was passed in on the command line.
//
type Config struct {
	SeedUrls []string
	AllowUrls []string
	NumConnections uint
}


/**
* Parse our command line arguments.
* @return {config} Our configuration info
*/
func ParseArgs() (retval Config) {

	retval = Config{[]string{}, []string{}, 1}

	hostnames := flag.String("seed-url",
		"http://www.cnn.com/",
		"URL to start with.")
	allowUrls := flag.String("allow-urls",
		"", "Url base names to crawl. " +
		"If specified, this basically acts like a whitelist. " + 
		"This may be a comma-delimited list. " + 
		"Examples: http://cnn.com/, http://www.apple.com/store")
	flag.UintVar(&retval.NumConnections, "num-connections",
		1, "How many concurrent outbound connections?")
	h := flag.Bool("h", false, "To get this help")
	help := flag.Bool("help", false, "To get this help")
	debug_level := flag.String("debug-level", "info", "Set the debug level")

	flag.Parse()

	log.SetLevelString(*debug_level)
	log.Error("Debug level: " + *debug_level)

	if (*h || *help) {
		flag.PrintDefaults()
		os.Exit(1)
	}

	retval.SeedUrls = SplitHostnames(*hostnames)
	retval.AllowUrls = SplitHostnames(*allowUrls)

	return(retval)

} // End of ParseArgs()


/**
* Take a comma-delimited string of hostnames and turn it into an array of URLs.
*
* @param {string} Input The comma-delimited string
*
* @return {[]string} Array of URLs
*/
func SplitHostnames(Input string) (retval []string) {

	Results := strings.Split(Input, ",")

	for _, value := range Results {

		if (value != "") {
			pattern := "^http(s)?://"
			match, _ := regexp.MatchString(pattern, value)
			if (!match) {
				value = "http://" + value
			}

		}

		retval = append(retval, value)

	}

	return(retval)

} // End of SplitHostnames()




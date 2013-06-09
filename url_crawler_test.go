
package main

//import "fmt"
import "regexp"
import "testing"

//import log "github.com/dmuth/google-go-log4go"
import server "github.com/dmuth/procedural-webserver"


func Test(t *testing.T) {

	//log.SetLevelString("info")

	//
	// Start up our server
	//
	port := 8080
	server_obj := server.NewServer(port, 5, 20, 5, 20, "test_seed")
	go server_obj.Start()

	in := make(chan string)
	out := make(chan Response)

	NewUrlCrawler(in, out, 10)

	in <- "http://localhost:8080/test2"
	result := <-out

	if (result.Code != 200) {
		t.Errorf("Code %d does not match 200!", result.Code)
	}

	in <- "http://localhost:8080/test2?code=404"
	result = <-out

	if (result.Code != 404) {
		t.Errorf("Code %d does not match 404!", result.Code)
	}

	//
	// Try a bad port
	//
	in <- "http://localhost:12345/test2?code=404"
	result = <-out

	if (result.Code != 0) {
		t.Errorf("Code %d does not match 0!", result.Code)
	}

	pattern := "connection refused"
	match, _ := regexp.MatchString(pattern, result.Body)
	if (!match) {
		t.Errorf("Could not find pattern '%s' in result '%s'", pattern, result)
	}

	//in <- "http://www.cnn.com/robots.txt"
	//in <- "http://localhost:8080/test2?delay=1s"
	//in <- "http://httpbin.org/headers"

	server_obj.Stop()

} // End of Test()



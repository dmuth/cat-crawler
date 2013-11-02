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

	in, out := NewUrlCrawler(10, []string{})

	url := "http://localhost:8080/test2"
	in <- url
	result := <-out

	if result.Url != url {
		t.Errorf("URL '%s' does not match '%s'!", result.Url, url)
	}

	if result.Code != 200 {
		t.Errorf("Code %d does not match 200!", result.Code)
	}

	in <- "http://localhost:8080/test2?code=404"
	result = <-out

	if result.Code != 404 {
		t.Errorf("Code %d does not match 404!", result.Code)
	}

	//
	// Try a bad port
	//
	in <- "http://localhost:12345/test2?code=404"
	result = <-out

	if result.Code != 0 {
		t.Errorf("Code %d does not match 0!", result.Code)
	}

	pattern := "connection refused"
	match, _ := regexp.MatchString(pattern, result.Body)
	if !match {
		t.Errorf("Could not find pattern '%s' in result '%s'", pattern, result)
	}

	//in <- "http://www.cnn.com/robots.txt"
	//in <- "http://localhost:8080/test2?delay=1s"
	//in <- "http://httpbin.org/headers"

	server_obj.Stop()

} // End of Test()

func TestFilterUrl(t *testing.T) {

	Urls := []string{
		"http://www.apple.com/",
		"http://www.apple.com/#",
		"http://www.apple.com/#foobar",
		"http://www.apple.com/what#foobar",
		"http://www.apple.com/CSS/ie7.css/../Resources/foobar",
		"http://www.apple.com/CSS/ie7.css/../Resources/../foobar",
		"http://www.apple.com/CSS/ie7.css/../Resources/../foobar/baz",
		"http://logging.apache.org/log4j/1.2/./css/print.css",
		"http://logging.apache.org/log4j/1.2///////css/print.css",
		"http://logging.apache.org/log4j/1.2///..///..///./css/print.css",
		"http:/www.flickr.com/photos/dmuth/6071648896/",
		"https:/www.flickr.com/photos/dmuth/6071648896/",
	}
	Expected := []string{
		"http://www.apple.com/",
		"http://www.apple.com/",
		"http://www.apple.com/",
		"http://www.apple.com/what",
		"http://www.apple.com/CSS/Resources/foobar",
		"http://www.apple.com/CSS/foobar",
		"http://www.apple.com/CSS/foobar/baz",
		"http://logging.apache.org/log4j/1.2/css/print.css",
		"http://logging.apache.org/log4j/1.2/css/print.css",
		"http://logging.apache.org/css/print.css",
		"http://www.flickr.com/photos/dmuth/6071648896/",
		"https://www.flickr.com/photos/dmuth/6071648896/",
	}

	for i := range Urls {
		Url := filterUrl(Urls[i])
		if Url != Expected[i] {
			t.Errorf("Filtered URL '%s' does not match expected URL '%s'!",
				Url, Expected[i])
		}
	}

} // End of TestFilterUrl()

func TestIsUrlAllowed(t *testing.T) {

	_, _ = NewUrlCrawler(10, []string{
		"http://foo/",
		"https://bar/baz",
	})
	Urls := []string{
		"http://google.com/",
		"http://foo",
		"http://foo/",
		"http://foo/bar",
		"http://bar",
		"http://bar/baz",
		"https://bar/baz",
		"https://bar/baz/",
	}
	Expected := []bool{
		false,
		false,
		true,
		true,
		false,
		false,
		true,
		true,
	}

	for key, value := range Urls {
		result := isUrlAllowed(value)
		if result != Expected[key] {
			t.Errorf("For URL '%s', expected %s and got %s",
				value, Expected[key], result)
		}
	}

	_, _ = NewUrlCrawler(10, []string{})
	for key, value := range Urls {
		result := isUrlAllowed(value)
		if result != true {
			t.Errorf("For URL '%s', expected %s and got %s",
				value, Expected[key], result)
		}
	}

} // End of TestIsUrlAllowed()

package main

//import "fmt"
import "testing"

func TestGetFilenameFromUrl(t *testing.T) {

	Urls := []string{
		"http://www.apple.com/image.png",
		"http://www.apple.com/image",
		"http://www.apple.com/CSS/Resources/foobar.png",
		"http://www.apple.com/CSS/foobar",
		"http://logging.apache.org/log4j/1.2/css/print.css",
		"http://logging.apache.org/css/print.css",
		"http://www.flickr.com/photos/dmuth/6071648896",
		"http://www.flickr.com/photos/dmuth/6071648896/",
		"https://www.flickr.com/photos/dmuth/6071648896/",
		"https://www.flickr.com/photos/dmuth/6071648896/" +
			"1234567890" +
			"1234567890" +
			"1234567890" +
			"1234567890" +
			"1234567890" +
			"1234567890",
	}

	Expected := []string{
		"www.apple.com/image.png",
		"www.apple.com/image",
		"www.apple.com/CSS/Resources/foobar.png",
		"www.apple.com/CSS/foobar",
		"logging.apache.org/log4j/1.2/css/print.css",
		"logging.apache.org/css/print.css",
		"www.flickr.com/photos/dmuth/6071648896",
		"www.flickr.com/photos/dmuth/6071648896",
		"www.flickr.com/photos/dmuth/6071648896",
		"www.flickr.com/photos/dmuth/6071648896/1234567890123456789012345678901234567890",
	}

	for key, value := range Urls {
		Result := getFilenameFromUrl(value)
		if Result != Expected[key] {
			t.Errorf("Filename '%s' != expected value '%s",
				Result, Expected[key])
		}
	}

} // End of TestGetFilenameFromUrl()

func TestImageBeenHereUrl(t *testing.T) {

	hostsCrawledImages = make(map[string]map[string]bool)

	url := "http://www.cnn.com/"
	result := imageBeenHereUrl(url)
	if result != false {
		t.Errorf("Expected result to be false for %s, got %s", url, result)
	}

	setImageBeenHereUrl(url)
	result = imageBeenHereUrl(url)
	if result != true {
		t.Errorf("Expected result to be true for %s, got %s", url, result)
	}

	url = "http://www.cnn.com/foobar"
	result = imageBeenHereUrl(url)
	if result != false {
		t.Errorf("Expected result to be false for %s, got %s", url, result)
	}

	setImageBeenHereUrl(url)
	result = imageBeenHereUrl(url)
	if result != true {
		t.Errorf("Expected result to be true for %s, got %s", url, result)
	}

}

func TestImageBeenHere(t *testing.T) {

	hostsCrawledImages = make(map[string]map[string]bool)

	result := imageBeenHere("test", "test2")
	if result != false {
		t.Errorf("Expected result to be false, got %s", result)
	}

	setImageBeenHere("test", "test2")

	result = imageBeenHere("test", "test2")
	if result != true {
		t.Errorf("Expected result to be true, got %s", result)
	}

	result = imageBeenHere("test", "test")
	if result != false {
		t.Errorf("Expected result to be false, got %s", result)
	}

}

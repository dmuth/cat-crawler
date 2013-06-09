
package main

//import "fmt"
import "testing"

//import log "github.com/dmuth/google-go-log4go"


func TestHtml(t *testing.T) {

	//log.SetLevelString("info")

	HtmlString := "<a href=\"foobar1\">foobar1 content</a>" +
		"<a href=\"/foobar2\">foobar2 content</a>" +
		"<a href=\"http://localhost/foobar3\">foobar3 content</a>" +
		"<a href=\"https://localhost/foobar4\">foobar4 content</a>\n" +
		"<a href=\"http://localhost:8080/foobar5\">foobar5 content</a>\n" +
		"<a href=\"https://localhost:8080/foobar6\">foobar6 content</a>\n" +
		"<img src=\"foobar1.png\" alt=\"foobar1 alt tag\">" +
		"<img src=\"/foobar2.png\" alt=\"foobar2 alt tag\" />" +
		"<img src=\"http://localhost/foobar3.png\" alt=\"foobar3 alt tag\" />" +
		"<img src=\"https://localhost/foobar4.png\" title=\"foobar4 title\" />" +
		"<img src=\"http://localhost:8080/foobar5.png\" alt=\"foobar5 alt tag\" />" +
		"<img src=\"https://localhost:8080/foobar6.png\" alt=\"foobar6 alt tag\" />" +
		""

	Results := HtmlParseString("http://www.cnn.com/world", HtmlString)

	ExpectedLinks := []string{
		"http://www.cnn.com/world/foobar1", 
		"http://www.cnn.com/foobar2", 
		"http://localhost/foobar3", 
		"https://localhost/foobar4", 
		"http://localhost:8080/foobar5", 
		"https://localhost:8080/foobar6", 
		}
	ExpectedImages := []string{
		"http://www.cnn.com/world/foobar1.png", 
		"http://www.cnn.com/foobar2.png", 
		"http://localhost/foobar3.png", 
		"https://localhost/foobar4.png", 
		"http://localhost:8080/foobar5.png", 
		"https://localhost:8080/foobar6.png", 
		}
	ExpectedAlt := []string{
		"foobar1 alt tag",
		"foobar2 alt tag",
		"foobar3 alt tag",
		"",
		"foobar5 alt tag",
		"foobar6 alt tag",
		}
	ExpectedTitles := []string{
		"",
		"",
		"",
		"foobar4 title",
		"",
		"",
		}

	for i := range Results.links {
		if (Results.links[i] != ExpectedLinks[i]) {
			t.Errorf("Result '%s' didn't match expected '%s'", Results.links[i], ExpectedLinks[i])
		}
	}

	for i := range Results.images {
		if (Results.images[i].src != ExpectedImages[i]) {
			t.Errorf("Images '%s' didn't match expected '%s'", Results.images[i].src, ExpectedImages[i])
		}
		if (Results.images[i].alt != ExpectedAlt[i]) {
			t.Errorf("Alt '%s' didn't match expected '%s'", Results.images[i].alt, ExpectedAlt[i])
		}
		if (Results.images[i].title != ExpectedTitles[i]) {
			t.Errorf("Title '%s' didn't match expected '%s'", Results.images[i].title, ExpectedTitles[i])
		}
	}


} // End of TestHtml()



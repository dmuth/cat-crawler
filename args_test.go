
package main

//import "fmt"
import "testing"

//import log "github.com/dmuth/google-go-log4go"


func TestSplitHostnames(t *testing.T) {

	//log.SetLevelString("info")

	Input := "test,test2,http://test3,https://test4/,test5:8080/,test6:8080/foobar"
	Output := SplitHostnames(Input)
	Expected := []string {
		"http://test",
		"http://test2",
		"http://test3",
		"https://test4/",
		"http://test5:8080/",
		"http://test6:8080/foobar",
	}

	for key, value := range Output {

		if (value != Expected[key]) {
			t.Errorf("Value '%s' doesn't match expected '%s'!", value, Expected[key])
		}

	}

	Input = "test"
	Output = SplitHostnames(Input)
	Expected = []string {
		"http://test",
	}

	for key, value := range Output {

		if (value != Expected[key]) {
			t.Errorf("Value '%s' doesn't match expected '%s'!", value, Expected[key])
		}

	}



} // End of TestSplitHostnames()




package main

//import "fmt"
import "testing"

func TestMain(t *testing.T) {

	key := "key1"
	IncrStatus(key, 1)

	result := Status(key)
	expected := 1
	if result != expected {
		t.Errorf("Value '%d' doesn't match expected '%d'!", result, expected)
	}

	DecrStatus(key, 3)
	result = Status(key)
	expected = -2
	if result != expected {
		t.Errorf("Value '%d' doesn't match expected '%d'!", result, expected)
	}

}


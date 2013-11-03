
package main

//import "fmt"
import "testing"

func TestMain(t *testing.T) {

	key := "key1"
	IncrStats(key, 1)

	result := Stats(key)
	expected := 1
	if result != expected {
		t.Errorf("Value '%d' doesn't match expected '%d'!", result, expected)
	}

	DecrStats(key, 3)
	result = Stats(key)
	expected = -2
	if result != expected {
		t.Errorf("Value '%d' doesn't match expected '%d'!", result, expected)
	}


	IncrStats("key2", 2)
	data := StatsAll()

	expected = -2
	if data[key] != expected {
		t.Errorf("Value '%d' doesn't match expected '%d'!", result, expected)
	}

	expected = 2
	if data["key2"] != expected {
		t.Errorf("Value '%d' doesn't match expected '%d'!", data["key2"], expected)
	}

}


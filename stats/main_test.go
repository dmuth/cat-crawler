
package stats

//import "fmt"
import "testing"

func TestMain(t *testing.T) {

	key := "key1"
	key2 := "key2"
	IncrStat(key)

	result := Stat(key)
	expected := 1
	if result != expected {
		t.Errorf("Value '%d' doesn't match expected '%d'!", result, expected)
	}

	DecrStat(key)
	result = Stat(key)
	expected = 0
	if result != expected {
		t.Errorf("Value '%d' doesn't match expected '%d'!", result, expected)
	}

	IncrStat(key2)
	data := StatAll()

	expected = 0
	if data[key] != expected {
		t.Errorf("Value '%d' doesn't match expected '%d'!", result, expected)
	}

	expected = 1
	if data[key2] != expected {
		t.Errorf("Value '%d' doesn't match expected '%d'!", data[key2], expected)
	}

	AddStat(key, 1)
	SubStat(key2, 3)

	data = StatAll()

	expected = 1
	if data[key] != expected {
		t.Errorf("Value '%d' doesn't match expected '%d'!", data[key], expected)
	}

	expected = -2
	if data[key2] != expected {
		t.Errorf("Value '%d' doesn't match expected '%d'!", data[key2], expected)
	}

}


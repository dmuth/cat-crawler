

package main

//import "fmt" // Debugging

//
// Our data structure holding key/value pairs
//
var data map[string]int
var beenhere bool

func initData(key string) {
	if !beenhere {
		data = make(map[string]int)
		beenhere = true
	}

	_, ok := data[key]
	if !ok {
		data[key] = 0
	}

}

/**
* Increment a key by a specified value.
*/
func IncrStats(key string, value int) {
	initData(key)
	data[key] += value
}

/**
* Decrement a key by a specified value.
*/
func DecrStats(key string, value int) {
	initData(key)
	data[key] -= value
}

/**
* Grab the value of a specific key.
*/
func Stats(key string) (int) {
	initData(key)
	return(data[key])
}


/**
* Get stats for all keys.
*/
func StatsAll() (map[string]int) {
	return(data)
}


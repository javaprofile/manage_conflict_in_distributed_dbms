package main

import (
	"fmt"
	"sync"
	"time"
)

type Version struct {
	value     int
	timestamp int64
}

type MVCCData struct {
	versions []Version
	mutex    sync.RWMutex
}

var database = map[string]*MVCCData{}
var dbMutex = sync.Mutex{}

func write(key string, value int, timestamp int64) {
	dbMutex.Lock()
	if _, ok := database[key]; !ok {
		database[key] = &MVCCData{}
	}
	dbMutex.Unlock()

	data := database[key]
	data.mutex.Lock()
	defer data.mutex.Unlock()
	data.versions = append(data.versions, Version{value, timestamp})
}

func read(key string, timestamp int64) (int, bool) {
	dbMutex.Lock()
	data, ok := database[key]
	dbMutex.Unlock()
	if !ok {
		return 0, false
	}
	data.mutex.RLock()
	defer data.mutex.RUnlock()
	var result Version
	found := false
	for _, v := range data.versions {
		if v.timestamp <= timestamp && (!found || v.timestamp > result.timestamp) {
			result = v
			found = true
		}
	}
	if found {
		return result.value, true
	}
	return 0, false
}

func main() {
	now := time.Now().UnixNano()
	write("x", 10, now)
	time.Sleep(time.Millisecond)
	write("x", 20, time.Now().UnixNano())

	val, ok := read("x", now)
	fmt.Println("Read at", now, "Value:", val, "Success:", ok)
}

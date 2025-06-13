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

var (
	database     = map[string]*MVCCData{}
	dbMutex      = sync.Mutex{}
	conflictLock = sync.Mutex{}
	conflicts    = 0
)

func write(key string, value int, timestamp int64) {
	dbMutex.Lock()
	if _, ok := database[key]; !ok {
		database[key] = &MVCCData{}
	}
	dbMutex.Unlock()

	data := database[key]
	data.mutex.Lock()
	defer data.mutex.Unlock()
	for _, v := range data.versions {
		if v.timestamp == timestamp {
			conflictLock.Lock()
			conflicts++
			conflictLock.Unlock()
			return
		}
	}
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
	t1 := time.Now().UnixNano()
	t2 := t1
	write("a", 1, t1)
	write("a", 2, t2)
	val, ok := read("a", t1)
	fmt.Println("Read:", val, "Success:", ok)
	fmt.Println("Conflicts:", conflicts)
}

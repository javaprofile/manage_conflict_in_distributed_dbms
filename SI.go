package main
import (
	"fmt"
	"sync"
	"time"
)
type DataItem struct {
	value     int
	timestamp int64
}
type Transaction struct {
	id        int
	startTime int64
	writeSet  map[string]int
}
var (
	data      = map[string]DataItem{}
	dataMutex = sync.RWMutex{}
	txID      = 0
)
func startTransaction() Transaction {
	txID++
	return Transaction{
		id:        txID,
		startTime: time.Now().UnixNano(),
		writeSet:  make(map[string]int),
	}
}
func read(tx Transaction, key string) (int, bool) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()
	item, exists := data[key]
	if !exists || item.timestamp > tx.startTime {
		return 0, false
	}
	return item.value, true
}
func write(tx *Transaction, key string, value int) {
	tx.writeSet[key] = value
}
func commit(tx Transaction) bool {
	dataMutex.Lock()
	defer dataMutex.Unlock()
	for key, item := range data {
		if item.timestamp > tx.startTime {
			if _, inWriteSet := tx.writeSet[key]; inWriteSet {
				return false
			}
		}
	}
	for key, value := range tx.writeSet {
		data[key] = DataItem{
			value:     value,
			timestamp: time.Now().UnixNano(),
		}
	}
	return true
}
func main() {
	tx1 := startTransaction()
	v, ok := read(tx1, "x")
	fmt.Println("Tx1 Read x:", v, ok)
	write(&tx1, "x", 10)
	success := commit(tx1)
	fmt.Println("Tx1 Commit:", success)
	tx2 := startTransaction()
	v, ok = read(tx2, "x")
	fmt.Println("Tx2 Read x:", v, ok)
	write(&tx2, "x", 20)
	success = commit(tx2)
	fmt.Println("Tx2 Commit:", success)
}

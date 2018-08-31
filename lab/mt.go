ppackage main

// handle all the imports
import (
	"fmt"
	//"os"
    "sync"
    //"sync/atomic"
    "time"
    //"reflect"
)

var mutex1 = &sync.Mutex{}
var mutex2 = &sync.Mutex{}

func main() {
    //var mutex1 = &sync.Mutex{}
    //var mutex2 = &sync.Mutex{}
    fmt.Println("Application started. starting thread 1+2")
    go thread1(mutex1)
    go thread2(mutex2)
    fmt.Println("Sleep")    
    time.Sleep(30000000000)
    fmt.Println("Starting thread 3+4")
    go thread3(mutex1)
    go thread4(mutex2)
    fmt.Println("Sleep")    
    time.Sleep(1000000000000)
}


func thread1 (mutex *sync.Mutex ) {
    fmt.Println("Thread 1 running")
    mutex.Lock()
    fmt.Println("Mutex 1 Locked")
    time.Sleep(90000000000)
    mutex.Unlock()
}

func thread2 (mutex *sync.Mutex ) {
    fmt.Println("Thread 2 running")
    mutex.Lock()
    fmt.Println("Mutex 2 Locked")
    time.Sleep(100000000000)
    mutex.Unlock()
}

func thread3 (mutex *sync.Mutex) {
    fmt.Println("Thread 3 running. Trying to unlock mutex 1")
    mutex.Lock()    
    fmt.Println("Thread 3: Mutex 1 Locked")
}

func thread4 (mutex *sync.Mutex) {
    fmt.Println("Thread 4 running. Trying to unlock mutex 2")
    mutex.Lock()
    fmt.Println("Thread 4: Mutex 2 locked")
}
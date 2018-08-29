package main

import (
    //"os"
    "fmt"
    "time"
    "storage"
)




func main() {
    fmt.Println("Test started")
    var i =  storage.CreateEntityIdent("test");
    fmt.Println("Test ended , newid: ",i)
    var entity storage.Entity
    entity.Type  = 1
    entity.Ident = storage.EntityRIdents["test"]
    entity.Context = "thats it"
    tmp := 1
    fmt.Println("Starting massinsert:")
    start := time.Now()
    for tmp < 10000000 {
         storage.CreateEntity(entity)
        tmp++
    }
    elapsed := time.Since(start)
    fmt.Println("Insert done in:",elapsed)
    
    
    
    
    
}
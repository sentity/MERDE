package main

import (
    //"os"
    //"errors"
    "fmt"
    "time"
    "storage"
)




func main() {
    fmt.Println("Test started")
    var i =  storage.CreateEntityIdent("test");
    fmt.Println("Created new ident , retrieved id: ",i)
    // ---------------------
    var entity storage.Entity
    entity.Ident   = storage.EntityRIdents["test"]
    entity.Context = "thats it"
    entity.Value   = "what a wonderfull world"
    // ---------------------
    fmt.Println("Starting mass testing:")
    tmp := 1
    //max := 1000000
    max := 10
    fmt.Println("Defined value for mass tests: ",max)
    start := time.Now()
    for tmp < max + 1 {
        entity.ID = tmp
        storage.CreateEntity(entity)
        tmp++
    }
    elapsed := time.Since(start)
    fmt.Println("Entity insert done in:",elapsed)
    // ---------------------
    start2  := time.Now()
    tmp2    := 1
    for tmp2 < max {
        storage.GetEntityByPath(1,tmp2)
        tmp2++
    }
    elapsed2 := time.Since(start2)
    fmt.Println("Entity read done in:",elapsed2)
    // ---------------------
    var relation storage.Relation
    relation.Context = "thats it"
    start3  := time.Now()
    var tmp3 = 1
    var tmp4 = 2
    for tmp4 < max {
        relation.SourceIdent  = 1
        relation.SourceID     = 1
        relation.TargetIdent  = 1
        relation.TargetID     = tmp4
        storage.CreateRelation(1,1,1,tmp4,relation)
        tmp3++
        tmp4++
    }
    elapsed3 := time.Since(start3)
    fmt.Println("Relation insert done in: ",elapsed3)
    // ---------------------
    start4   := time.Now()
    storage.GetRelationsBySourceIdentAndSourceId(1,1)
    elapsed4 := time.Since(start4)
    fmt.Println("Relation read (nax* out from 1:1) done in: ",elapsed4)
}
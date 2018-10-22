package main

import (
    //"os"
    //"errors"
    "time"
    //"storage"   
    "fmt"
    "mapper"
    "encoding/json"
    "connector"
)


var JsonMapThreadTest = make(map[string]bool)

func main() {
    connector.Listen()
    //tests()
}


func tests() {
    //fmt.Println("Test started")
    //i,_ :=  storage.CreateEntityIdent("test");
    //fmt.Println("Created new ident , retrieved id: ",i)
    // ---------------------
    //var entity storage.Entity
    //entity.Ident   = storage.EntityRIdents["test"]
    //entity.Context = "thats it"
    //entity.Value   = "what a wonderfull world"
    // ---------------------
    //fmt.Println("Starting mass testing:")
    //tmp := 1
    //max := 100000
    ////max := 10
    //fmt.Println("Defined value for mass tests: ", max)
    //start := time.Now()
    //for tmp < max + 1 {
    //    entity.ID = tmp
    //    _, err := storage.CreateEntity(entity)
    //    if err == nil {
    //        tmp++
    //    }
    //
    //}
    //elapsed := time.Since(start)
    //fmt.Println("Entity insert done in:", elapsed)
    // ---------------------
    //start2  := time.Now()
    //tmp2    := 1
    //for tmp2 < max {
    //    storage.GetEntityByPath(1,tmp2)
    //    tmp2++
    //}
    //elapsed2 := time.Since(start2)
    //fmt.Println("Entity read done in:",elapsed2)
    // ---------------------
    //var relation storage.Relation
    //relation.Context = "thats it"
    //start3  := time.Now()
    //var tmp3 = 1
    //var tmp4 = 2
    //for tmp4 < max {
    //    relation.SourceIdent  = 1
    //    relation.SourceID     = 1
    //    relation.TargetIdent  = 1
    //    relation.TargetID     = tmp4
    //    storage.CreateRelation(1,1,1,tmp4,relation)
    //    tmp3++
    //    tmp4++
    //}
    //elapsed3 := time.Since(start3)
    //fmt.Println("Relation insert done in: ",elapsed3)
    // ---------------------
    //start4   := time.Now()
    //storage.GetRelationsBySourceIdentAndSourceId(1,1)
    //elapsed4 := time.Since(start4)
    //fmt.Println("Relation read (nax* out from 1:1) done in: ",elapsed4)
    // ----------------------
    //fmt.Println("Testing json mapper - inserting")
    //var id = TestJsonMap()
    //fmt.Println("Testing json mapper - getting from ident ip: and id: ",id)
    //TestJsonGet(id)
    // ----------------------    
    //start5   := time.Now()
    //// testing mass insert
    //var i = 0
    //for i < 100000 {
    //    TestJsonMap()
    //    i++
    //}
    //elapsed5 := time.Since(start5)
    //fmt.Println("3 level entity mapped mass in : ",elapsed5)
    //start6   := time.Now()
    //for i < 100000 {
    //    TestJsonGet(i)
    //    i++
    //}
    //elapsed6 := time.Since(start6)
    //fmt.Println("3 level entity read mass in : ",elapsed6)
    // ----------------------    
    start7   := time.Now()
    max      := 1000000
    // testing mass insert
    JsonMapThreadTest["thread1"] = false
    JsonMapThreadTest["thread2"] = false
    go JsonMapThread1(max)
    //go JsonMapThread2(max)
    for JsonMapThreadTest["thread1"] != true && JsonMapThreadTest["thread2"] != true {
        time.Sleep(100000)
    }
    elapsed7 := time.Since(start7)
    fmt.Println("3 level entity mapped in 2 threads mass in : ",elapsed7)
    fmt.Println("Mass read :");
    start8   := time.Now()
    i := 0
    for i < max {
        mapper.GetEntityRecursive(mapper.HandleIdent("ip"),i)
        i++
    }
    elapsed8 := time.Since(start8)
    fmt.Println("3 level entity read in 1 thread mass amount : ",max , " in " , elapsed8)
}

func JsonMapThread1(max int)  {
    i := 0
    for i < max{
        TestJsonMap()
        fmt.Println("Thead 1 run: ",i);
        i++
    }
    fmt.Println("Thread 1 done. Wrote ", i , " 3level entities")
    JsonMapThreadTest["thread1"] = true
}

func JsonMapThread2(max int) {
    i := 0
    for i < max{
        TestJsonMap()
        fmt.Println("Thead 2 run: ",i);
        i++
    }
    fmt.Println("Thread 2 done. Wrote ", i , " 3level entities")
    JsonMapThreadTest["thread2"] = true
}


func TestJsonMap() (int){
    JsonByteArray := []byte(`{"Context":"asd","Ident":"ip","Value":"it works yippiyey","Properties":{"onekey":"onevalue","twokey":"twovalue"},"Children":{}}`)
    //JsonByteArray := []byte(`{"Context":"asd","Ident":"ip","Value":"it works yippiyey","Properties":{"onekey":"onevalue","twokey":"twovalue"},"Children":{"1":{"Context":"im the subobject","Ident":"Port","Value":"subobject kinda cooly","Properties":{"onekey":"udabedi","twokey":"dabedei"},"Children":{"2":{"Context":"im the third and best","Ident":"state","Value":"third depp so deep","Properties":{"onekey":"bass","twokey":"boom"},"Children":{}}}}}}`)
    var id,_ = mapper.MapJson(JsonByteArray)
    return id
}

func TestJsonGet(id int) () {
    ret,_:= mapper.GetEntityRecursive(mapper.HandleIdent("ip"),id)
    //mapper.GetEntityRecursive(mapper.HandleIdent("ip"),id)
    //fmt.Printf("%#v", ret)
    out, _    := json.MarshalIndent(ret, "", "  ")
    fmt.Print(string(out))  
}
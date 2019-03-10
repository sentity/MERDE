package main

import (
	"encoding/json"
	"fmt"
	"goer/libs/connector"
	"goer/libs/mapper"
	//	"goer/libs/query"
	"goer/libs/storage"
	//"errors"
	//"os"
	"time"
)

var JsonMapThreadTest = make(map[string]bool)

func main() {
	connector.Listen()
	//tests()
	//testQueryWhereParser()
	//mapGiantChain()
	//testQueryParser()
}

func testQueryParser() {
	//multiquery := "[{\"type\":\"find.entity\",\"ident\":\"ip\",\"where\":[\"property.type ='ipv6' && property.active='true'\",\"value='2.3.4.5'\"]},{\"type\":\"find.entity\",\"direction\":\"parent\",\"ident\":\"domain\",\"traverse\":3}]"
	//query.HandleQuery(multiquery)

}

func testQueryWhereParser() {
	testConditionString("property.test=='42'")
	testConditionString("property.test>'42'&&value%='12\\&\\&\\'\\'&&test==42")

}

func testConditionString(param string) {
	//fmt.Println("- - - - - - - - - - - - - - - - -")
	//	start := time.Now()
	//	max := int(100000)

	//for i := 0; i <= max; i++ {
	//	arrRet, _ := query.ParseConditions(param)
	//	if i == max {
	//	query.DebugPrint(arrRet)
	//	}
	//}
	//elapsed := time.Since(start)
	//fmt.Println("Parseind condition: ", param, " | ", max, " - Times | in: ", elapsed)
}

func mapGiantChain() {
	max := 6000000
	entityTypeId, _ := storage.CreateEntityIdent("test")
	// we create the most upper parent entity
	tmpEntity := storage.Entity{
		ID:      -1,
		Ident:   entityTypeId,
		Value:   "test",
		Context: "test",
	}
	start := time.Now()
	// create new entity
	rootId, _ := storage.CreateEntity(tmpEntity)
	parentId := rootId
	for i := 0; i < max; i++ {
		tmpEntity := storage.Entity{
			ID:      -1,
			Ident:   entityTypeId,
			Value:   "test",
			Context: "test",
		}
		// create new entity
		childId, _ := storage.CreateEntity(tmpEntity)
		tmpRelation := storage.Relation{
			SourceIdent: entityTypeId,
			SourceID:    parentId,
			TargetIdent: entityTypeId,
			TargetID:    childId,
		}
		storage.CreateRelation(entityTypeId, parentId, entityTypeId, childId, tmpRelation)
		parentId = childId
	}
	elapsed := time.Since(start)
	fmt.Println("Created ", max, " datasets long chain - most root id: ", rootId, " | time taken ", elapsed)
	//ret, _ := mapper.GetEntityRecursive(entityTypeId, rootId)

	//fmt.Println("Read out 500k datasets linked in a chain in ", elapsed, " seconds - testing return ")
	//testReturnTraverseDepth(ret, 0)
}

func testReturnTraverseDepth(data mapper.Entity, depth int) {
	if childData, ok := data.Children[1]; ok {
		depth = depth + 1
		testReturnTraverseDepth(childData, depth)
	} else {
		fmt.Println("Traversed through return object, depth found ", depth)
	}
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
	start7 := time.Now()
	max := 100000
	// testing mass insert
	JsonMapThreadTest["thread1"] = false
	JsonMapThreadTest["thread2"] = false
	go JsonMapThread1(max)
	go JsonMapThread2(max)
	for JsonMapThreadTest["thread1"] != true && JsonMapThreadTest["thread2"] != true {
		time.Sleep(100000)
	}
	elapsed7 := time.Since(start7)
	fmt.Println("3 level entity mapped in 2 threads mass in : ", elapsed7)
	fmt.Println("Mass read :")
	start8 := time.Now()
	i := 0
	for i < max {
		mapper.GetEntityRecursive(mapper.HandleIdent("ip"), i)
		i++
	}
	elapsed8 := time.Since(start8)
	fmt.Println("3 level entity read in 1 thread mass amount : ", max, " in ", elapsed8)
}

func JsonMapThread1(max int) {
	i := 0
	for i < max {
		TestJsonMap()
		//fmt.Println("Thead 1 run: ", i)
		i++
	}
	fmt.Println("Thread 1 done. Wrote ", i, " 3level entities")
	JsonMapThreadTest["thread1"] = true
}

func JsonMapThread2(max int) {
	i := 0
	for i < max {
		TestJsonMap()
		//fmt.Println("Thead 2 run: ", i)
		i++
	}
	fmt.Println("Thread 2 done. Wrote ", i, " 3level entities")
	JsonMapThreadTest["thread2"] = true
}

func TestJsonMap() int {
	//JsonByteArray := []byte(`{"Context":"asd","Ident":"ip","Value":"it works yippiyey","Properties":{"onekey":"onevalue","twokey":"twovalue"},"Children":{}}`)
	JsonByteArray := []byte(`{"Context":"asd","Ident":"ip","Value":"it works yippiyey","Properties":{"onekey":"onevalue","twokey":"twovalue"},"Children":{"1":{"Context":"im the subobject","Ident":"Port","Value":"subobject kinda cooly","Properties":{"onekey":"udabedi","twokey":"dabedei"},"Children":{"2":{"Context":"im the third and best","Ident":"state","Value":"third depp so deep","Properties":{"onekey":"bass","twokey":"boom"},"Children":{}}}}}}`)
	var id, _ = mapper.MapJson(JsonByteArray)
	return id
}

func TestJsonGet(id int) {
	ret, _ := mapper.GetEntityRecursive(mapper.HandleIdent("ip"), id)
	//mapper.GetEntityRecursive(mapper.HandleIdent("ip"),id)
	//fmt.Printf("%#v", ret)
	out, _ := json.MarshalIndent(ret, "", "  ")
	fmt.Print(string(out))
}

package storage

// handle all the imports
import (
	"fmt"
	"os"
    //"math/rand"
    "sync"
    "time"
) 

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// entity struct 
type Entity struct {
    ID         int
    Type       int
    Context    string
    Ident      int
    Properties map[string]string
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// entity storage map            [Ident] [Type]  [ID]
var EntityStorage      = make(map[int]map[int]map[int]Entity)

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// entity storage id max         [Ident] [Type]  [ID]
var EntityIDMax        = make(map[int]map[int]map[int]int)

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// entity storage mutex is a per ident mutex so write
// operations only block on ident + type
//                               [ident] [type]
var EntityStorageMutex = make(map[int]map[int]*sync.Mutex )

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// maps to translate Idents to their INT and reverse
var EntityIdents       = make(map[string]int)
var EntityRIdents      = make(map[int]string)
// and a fitting max ID
var EntityIdentIDMax   = int

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// entity ident mutex (for adding and deleting ident types)
var EntityIdentMutex   = &sync.Mutex{}


// - - - - - - - - - - - - - - - - - - - - - - - - - -
// relation struct 
type Relation struct {
    ID         int
    Type       int
    Context    string
    Source     string
    Target     string
    Properties map[string]string
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// relation storage map         [Type] [Source]   [Target]   [ID]
var RelationStorage  = make(map[int]map[string]map[string]map[int]Relation)

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// relation reverse storage map
// (for faster queries)        [Type]  [Target]   [source]   [ID] [path]
var RelationRStorage = make(map[int]map[string]map[string]map[int]string)

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// relation index max id       [Type]  [Source]   [Target]   [ID]
var RelationIDMax    = make(map[int]map[string]map[string]map[int]int)


// - - - - - - - - - - - - - - - - - - - - - - - - - -
// + + + + + + FUNCTIONS + + + + + + 
// - - - - - - - - - - - - - - - - - - - - - - - - - -

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// Create an entity ident
func CreateEntityIdent(name string) (int){
    // first of allw e lock
    EntityIdentMutex.Lock()
    // lets check if the ident allready exists
    // if it does we just return the ID
    if id, ok := EntityRIdents[name]; ok {
        // dont forget to unlock
        EntityIdentMutex.Unlock()
        return id
    }
    // ok entity doesnt exist yet, lets
    // upcount our ID Max and copy it
    // into another variable so we can be sure
    // between unlock of the ressource and return
    // it doesnt get upcounted
    EntityIdentIDMax++
    var newID             = EntityIdentIDMax
    // finally create the new ident in our
    // EntityIdents index and reverse index
    EntityIdents[newID]   = name
    EntityRIdents[name]   = newID
    // now we create a special mutex for
    // this ident
    EntityIdentsMutex[newID]      = &sync.Mutex{}
    // and create mutex for EntityStorage ident+type
    EntityStorageMutex[newID][1]  = &sync.Mutex{}
    EntityStorageMutex[newID][2]  = &sync.Mutex{}
    // finally set the maxID for the new
    // ident types
    EntityIDMax[newID][1] = 0
    EntityIDMax[newID][1] = 0
    // now we unlock the mutex
    // and return the new id
    EntityIdentMutex.Unlock()
    return
}

//func DeleteEntityIdent() {
//    
//}

func CreateEntity(entity Entity) (int){
    // first we lock fitting
    // entity storage mutex
    if _, ok := EntityIdents[entity.Ident]; ok {
        return -1
    }
    // upcount our ID Max and copy it
    // into another variable so we can be sure
    // between unlock of the ressource and return
    // it doesnt get upcounted
    // lets upcount the entity id max fitting to
    //         [ident]  and  [type]
    EntityIDMax[entity.Ident][entity.Type]++
    var newID = EntityIDMax[entity.Ident][entity.Type]
    // and tell the entity its own id
    entity.ID = newID
    // now we store the entity element
    // in the EntityStorage
    EntityStorage[entity.Ident][entity.Type][newId] = entity
    // since we now stored the entity we can unlock
    // the storage ressource and return the ID
    EntityStorageMutex[entity.Ident][entity.Type].Unlock()
    return newId
}

func GetEntityByPath(ident int, Type int, id int) (Entity){
    if id, ok := EntityRIdents[name]; ok {
        // dont forget to unlock
        EntityIdentMutex.Unlock()
        return id
    }
}

func GetEntityByIdentAndType() {
    
}

func GetEntityByIdent() {
    
}

func DeleteEntity() {
    
}

func UpdateEntity() {
    
}

func CreateRelation() {
    
}

func GetRelation() {
    
}

func DeleteRelation() {
    
}

func UpdateRelation() {
    
}


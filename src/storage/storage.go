package storage

// handle all the imports
import (
	//"fmt"
	//"os"
    //"math/rand"
    "sync"
    //"time"
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
// entity storage id max         [Ident] [Type]  
var EntityIDMax        = make(map[int]map[int]int)

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// entity storage mutex is a per ident mutex so write
// operations only block on ident + type
//                               [ident] [type]
var EntityStorageMutex = make(map[int]map[int]*sync.Mutex )

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// maps to translate Idents to their INT and reverse
var EntityIdents       = make(map[int]string)
var EntityRIdents      = make(map[string]int)
// and a fitting max ID
var EntityIdentIDMax  int = 0


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
    // and create mutex for EntityStorage ident+type
    var tmpMap                    = make(map[int]*sync.Mutex)
    EntityStorageMutex[newID]     = tmpMap
    EntityStorageMutex[newID][1]  = &sync.Mutex{}
    EntityStorageMutex[newID][2]  = &sync.Mutex{}
    // now we prepare the submaps in the entity
    // storage itseöf....
    var tmpMap2                   = make(map[int]map[int]Entity)
    EntityStorage[newID]          = tmpMap2
    var tmpMap3                   = make(map[int]Entity)
    EntityStorage[newID][1]       = tmpMap3
    EntityStorage[newID][2]       = tmpMap3
//    map[int]map[int]map[int]Entity
    // finally set the maxID for the new
    // ident types
    var tmpMap4                = make(map[int]int)
    EntityIDMax[newID]         = tmpMap4
    EntityIDMax[newID][1]      = 0
    EntityIDMax[newID][1]      = 0
    // now we unlock the mutex
    // and return the new id
    EntityIdentMutex.Unlock()
    return newID
}

//func DeleteEntityIdent() {
//    
//}

func CreateEntity(entity Entity) (int){
    // first we lock the entity ident mutex
    // to make sure while we check for the
    // existence it doesnt get deletet, this
    // may sound like a very rare upcoming case,
    //but better be safe than sorry
    EntityIdentMutex.Lock()
    // now 
    if _, ok := EntityIdents[entity.Ident]; !ok {
        // the ident doest exist, lets unlock
        // the ident mutex and return -1 for fail0r
        EntityIdentMutex.Unlock()
        return -1
    }
    // the ident seems to exist, now lets lock the
    // storage mutex before Unlocking the Entity
    // ident mutex to prevent the ident beeing
    // deleted before we start locking (small
    // timing still possible )
    EntityStorageMutex[entity.Ident][entity.Type].Lock()
    EntityIdentMutex.Unlock()
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
    EntityStorage[entity.Ident][entity.Type][newID] = entity
    // since we now stored the entity we can unlock
    // the storage ressource and return the ID
    EntityStorageMutex[entity.Ident][entity.Type].Unlock()
    return newID
}

func GetEntityByPath(ident int, Type int, id int) (Entity){
    if entity, ok := EntityStorage[ident][Type][id]; ok {
        // dont forget to unlock
        EntityIdentMutex.Unlock()
        return entity
    }
    var entity = Entity{}
    return entity
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


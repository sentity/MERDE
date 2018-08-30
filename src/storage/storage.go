package storage

// handle all the imports
import (
	//"fmt"
	//"os"
    //"math/rand"
    "sync"
    //"time"
    "errors"
    //"strings"
) 

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// entity struct 
type Entity struct {
    ID         int
    Ident      int
    Context    string
    Value      string
    Properties map[string]string
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// entity storage map            [Ident] [ID]
var EntityStorage      = make(map[int]map[int]Entity)

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// entity storage id max         [Ident]  
var EntityIDMax        = make(map[int]int)

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// entity storage mutex is a per ident mutex so write
// operations only block on ident + type
//                               [ident] 
var EntityStorageMutex = make(map[int]*sync.Mutex )

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
    Context    string
    Source     string
    Target     string
    Properties map[string]string
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// s prefix = source
// t prefix = target
// relation storage map             [sIdent][sId]   [tIdent][tId]
var RelationStorage       = make(map[int]map[int]map[int]map[int]Relation)

// and relation reverse storage map
// (for faster queries)              [tIdent][Tid]   [sIdent][sId]
var RelationRStorage      = make(map[int]map[int]map[int]map[int]string)

// relation index max id            [sIdent][sId] 
var RelationStorageMutex  = make(map[int]map[int]*sync.Mutex)





// - - - - - - - - - - - - - - - - - - - - - - - - - -
// + + + + + + FUNCTIONS + + + + + + 
// - - - - - - - - - - - - - - - - - - - - - - - - - -

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// init/construct function for storage package
func init() {
    // first we gonne create the mutex
    // maps inside the RelationStorage
    //RelationStorageMutex[1]   = make(map[string]*sync.Mutex)
    //RelationStorageMutex[2]   = make(map[string]*sync.Mutex)
}



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
    var newID                 = EntityIdentIDMax
    // finally create the new ident in our
    // EntityIdents index and reverse index
    EntityIdents[newID]       = name
    EntityRIdents[name]       = newID
    // and create mutex for EntityStorage ident+type
    EntityStorageMutex[newID] = &sync.Mutex{}
    // now we prepare the submaps in the entity
    // storage itseöf....
    EntityStorage[newID]      = make(map[int]Entity)
    // set the maxID for the new
    // ident type
    EntityIDMax[newID]        = 0
    // and create the basic submaps for
    // the relation storage
    // now we unlock the mutex
    // and return the new id
    EntityIdentMutex.Unlock()
    return newID
}

//func DeleteEntityIdent() {
//    
//}

func CreateEntity(entity Entity) (int, error){
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
        return -1, errors.New("Entity ident not existing");
    }
    // the ident seems to exist, now lets lock the
    // storage mutex before Unlocking the Entity
    // ident mutex to prevent the ident beeing
    // deleted before we start locking (small
    // timing still possible )
    EntityStorageMutex[entity.Ident].Lock()
    EntityIdentMutex.Unlock()
    // upcount our ID Max and copy it
    // into another variable so we can be sure
    // between unlock of the ressource and return
    // it doesnt get upcounted
    // lets upcount the entity id max fitting to
    //         [ident]  and  [type]
    EntityIDMax[entity.Ident]++
    var newID = EntityIDMax[entity.Ident]
    // and tell the entity its own id
    entity.ID = newID
    // now we store the entity element
    // in the EntityStorage
    EntityStorage[entity.Ident][newID] = entity
    // since we now stored the entity we can unlock
    // the storage ressource and return the ID
    EntityStorageMutex[entity.Ident].Unlock()
    // create the mutex for our ressource on
    // relation. we have to create the sub maps too
    // golang things....
    
    //var tmpMap1              = make(map[int]map[int]*sync.Mutex)
    //RelationStorageMutex[1]  = tmp1
    //RelationStorageMutex[2]  = tmp1
    //var tmpMap2              = make(map[int]*sync.Mutex)
    //RelationStorageMutex[1][entity.Ident] = tmp2
    //RelationStorageMutex[2][entity.Ident] = tmp2
    //RelationStorageMutex[1][entity.Ident][newID] = &sync.Mutex{}
    //RelationStorageMutex[2][entity.Ident][newID] = &sync.Mutex{}
    
    // finally we return the new id
    return newID, nil
}

func GetEntityByPath(ident int, id int) (Entity, error){
    // lets check if entity witrh the given path exists
    if entity, ok := EntityStorage[ident][id]; ok {
        // if yes we return the entity
        // and nil for error
        return entity, nil
    }
    // the path seems to result empty , so
    // we throw an error 
    return Entity{}, errors.New("Entity on given path does not exist.");
}



func CreateRelation(Type int, srcIdent int, srcID int, targetIdent int, targetID int) {
    //if _, ok := RelationStorage[Type][srcIdent][srcID][targetIdent][targetID]; !ok {
    //
    //}
    //if _, ok := RelationStorage[Type][srcIdent][srcID][targetIdent][targetID]; !ok {
    //
    //}
    //if _, ok := RelationStorage[Type][srcIdent][srcID][targetIdent][targetID]; !ok {
    //
    //}
    //if _, ok := RelationStorage[Type][srcIdent][srcID][targetIdent][targetID]; !ok {
    //
    //}
}




func GetEntityByIdentAndType() {
    
    
}

func GetEntityByIdent() {
    
}

func DeleteEntity() {
    
}

func UpdateEntity() {
    
}

func GetRelation() {
    
}

func DeleteRelation() {
    
}

func UpdateRelation() {
    
}


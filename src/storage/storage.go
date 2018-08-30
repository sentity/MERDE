package storage

// handle all the imports
import (
	//"fmt"
	//"os"
    //"math/rand"
    "sync"
    //"time"
    "errors"
    "strconv"
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

// relation index max id            [sIdent]
var RelationStorageMutex  = make(map[int]*sync.Mutex)





// - - - - - - - - - - - - - - - - - - - - - - - - - -
// + + + + + + FUNCTIONS + + + + + + 
// - - - - - - - - - - - - - - - - - - - - - - - - - -

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// init/construct function for storage package
func init() {
    // first we gonne create the mutex
    // maps inside the RelationStorage
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
    // create the base maps in relation storage
    RelationStorage[newID]       = make(map[int]map[int]map[int]Relation)
    RelationRStorage[newID]      = make(map[int]map[int]map[int]string)
    // and the mutex
    RelationStorageMutex[newID]  = &sync.Mutex{}
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

    // create the mutex for our ressource on
    // relation. we have to create the sub maps too
    // golang things....
    RelationStorageMutex[entity.Ident].Lock()
    RelationStorage[entity.Ident][newID]     = make(map[int]map[int]Relation)
    RelationRStorage[entity.Ident][newID]    = make(map[int]map[int]string)
    RelationStorageMutex[entity.Ident].Unlock()
    // since we now stored the entity and created all
    // needed ressources we can unlock
    // the storage ressource and return the ID (or err)
    EntityStorageMutex[entity.Ident].Unlock()
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



func CreateRelation(srcIdent int, srcID int, targetIdent int, targetID int, relation Relation) (bool,error) {
    
    // lets make sure the source ident exist
    if _,ok := EntityIdents[srcIdent] ; !ok {
        return false, errors.New("Source ident not existing")
    }
    // and the target ident exists too
    if _,ok := EntityIdents[targetIdent] ; !ok {
        return false, errors.New("Target ident not existing")
    }
    // Now we lock the to link entities to
    // make sure they dont get deletet meawhile
    EntityStorageMutex[srcIdent].Lock()
    // if srcIdent and targetIdent differ,
    // we lock targetIdent too else we
    // would create deadlock
    if srcIdent != targetIdent {
        EntityStorageMutex[targetIdent].Lock()
    }
    // now we lock the relation mutex
    RelationStorageMutex[srcIdent].Lock()
    // lets check if their exists a map for our
    // source entity to the target ident if not
    // create it.... golang things...
    if _, ok := RelationStorage[srcIdent][srcID][targetIdent]; !ok {
        RelationStorage[srcIdent][srcID][targetIdent] = make(map[int]Relation)
        // if the map doesnt exist in this direction
        // it wont exist in the other as in reverse
        // map either so we should create it too
        // but we will store a pointer to the other
        // maps Relation instead of the complete
        // relation twice - defunct, refactor later (may create more problems then help)
        //RelationStorage[targetIdent][targetID][srcIdent] = make(map[int]Relation)
    }
    // now we prepare the reverse storage if necessary
    if _,ok := RelationRStorage[targetIdent][targetID][srcIdent]; !ok {
        RelationRStorage[targetIdent][targetID][srcIdent] = make(map[int]string)
    }
    // now we store the relation 
    RelationStorage[srcIdent][srcID][targetIdent][targetID] = relation
    // and a pointer in the reverse index
    //RelationRStorage[targetIdent][targetID][srcIdent][srcID] = &RelationStorage[srcIdent][srcID][targetIdent][targetID]
    // since the solution above doesnt work atm we do the following workarround temporary
    a := strconv.Itoa(srcIdent)
    b := strconv.Itoa(srcID)
    c := strconv.Itoa(targetIdent)
    d := strconv.Itoa(targetID)
    RelationRStorage[targetIdent][targetID][srcIdent][srcID] = a + ":" + b + ":" + c + ":" + d
    // we are done now we can unlock the entity idents
    EntityStorageMutex[srcIdent].Unlock()
    // if we locked the targetIdent too (see upper)
    // than we have to unlock it too
    if srcIdent != targetIdent {
        EntityStorageMutex[targetIdent].Unlock()
    }
    // and finally unlock the relation ident and return
    RelationStorageMutex[srcIdent].Unlock()
    return true, nil
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


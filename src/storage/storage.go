package storage

// handle all the imports
import (
	"fmt"
	//"os"
    //"math/rand"
    //"time"
    //"strconv"
    //"strings"
    //"builtin"
    "encoding/json"
    "errors"
    "sync"
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
var EntityStorage            = make(map[int]map[int]Entity)
// entity storage mutex is a per ident mutex so write
// operations only block on ident + type
//                                  [ident]
var EntityStorageMutex       = make(map[int]*sync.RWMutex )
// entity storage master mutex
var EntityStorageMasterMutex = &sync.RWMutex{}
// - - - - - - - - - - - - - - - - - - - - - - - - - -
// entity storage id max         [Ident]
var EntityIDMax            = make(map[int]int)
// and rw mutex for the IDmax
var EntityIDMaxMutex       = make(map[int]*sync.RWMutex )
// master mutexd for EntityIdMax
var EntityIDMaxMasterMutex = &sync.RWMutex{}



// - - - - - - - - - - - - - - - - - - - - - - - - - -
// maps to translate Idents to their INT and reverse
var EntityIdents          = make(map[int]string)
var EntityRIdents         = make(map[string]int)
// entity ident mutex (for adding and deleting ident types)
var EntityIdentMutex      = &sync.RWMutex{}
// and a fitting max ID
var EntityIdentIDMax  int = 0
// and total IdentID mutex to protect that array itself
var EntityIdentIDMaxMutex = &sync.RWMutex{}


// - - - - - - - - - - - - - - - - - - - - - - - - - -
// relation struct
type Relation struct {
    SourceIdent  int
    SourceID     int
    TargetIdent  int
    TargetID     int
    Context      string
    Properties map[string]string
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// s prefix = source
// t prefix = target
// relation storage map             [sIdent][sId]   [tIdent][tId]
var RelationStorage             = make(map[int]map[int]map[int]map[int]Relation)
// and relation reverse storage map
// (for faster queries)              [tIdent][Tid]   [sIdent][sId]
var RelationRStorage            = make(map[int]map[int]map[int]map[int]bool) 
// relation storage mutex                 [sIdent]
var RelationStorageMutex        = make(map[int]*sync.RWMutex)
// relation storage master mutex     
var RelationStorageMasterMutex  = &sync.RWMutex{}



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
func CreateEntityIdent(name string) (int, error){
    // first of allw e lock
    fmt.Println("CreateEntityIdent.EntityIdentMutex.Lock");
    EntityIdentMutex.Lock()
    // lets check if the ident allready exists
    // if it does we just return the ID
    if id, ok := EntityRIdents[name]; ok {
        // dont forget to unlock
        fmt.Println("CreateEntityIdent.EntityIdentMutex.Unlock");
        EntityIdentMutex.Unlock()
        return id, nil
    }
    // ok entity doesnt exist yet, lets
    // upcount our ID Max and copy it
    // into another variable so we can be sure
    // between unlock of the ressource and return
    // it doesnt get upcounted
    fmt.Println("CreateEntityIdent.EntityIdentIDMaxMutex.Lock");
    EntityIdentIDMaxMutex.Lock()
    EntityIdentIDMax++
    var newID  = EntityIdentIDMax
    fmt.Println("CreateEntityIdent.EntityIdentIDMaxMutex.Unlock");
    EntityIdentIDMaxMutex.Unlock()
    // finally create the new ident in our
    // EntityIdents index and reverse index
    EntityIdents[newID]          = name
    EntityRIdents[name]          = newID
    // and create mutex for EntityStorage ident+type
    fmt.Println("CreateEntityIdent.EntityStorageMasterMutex.Lock");
    EntityStorageMasterMutex.Lock()
    EntityStorageMutex[newID]    = &sync.RWMutex{}
    // now we prepare the submaps in the entity
    // storage itseöf....
    EntityStorage[newID]         = make(map[int]Entity)
    fmt.Println("CreateEntityIdent.EntityStorageMasterMutex.Unlock");
    EntityStorageMasterMutex.Unlock()
    // set the maxID for the new
    // ident type
    fmt.Println("CreateEntityIdent.EntityIDMaxMasterMutex.Lock");
    EntityIDMaxMasterMutex.Lock()
    EntityIDMax[newID]           = 0
    EntityIDMaxMutex[newID]      = &sync.RWMutex{}
    fmt.Println("CreateEntityIdent.EntityIDMaxMasterMutex.Unlock");
    EntityIDMaxMasterMutex.Unlock()
    // create the base maps in relation storage
    RelationStorage[newID]       = make(map[int]map[int]map[int]Relation)
    RelationRStorage[newID]      = make(map[int]map[int]map[int]bool)
    // and the mutex
    fmt.Println("CreateEntityIdent.RelationStorageMasterMutex.Lock");
    RelationStorageMasterMutex.Lock()
    RelationStorageMutex[newID]  = &sync.RWMutex{}
    fmt.Println("CreateEntityIdent.RelationStorageMasterMutex.Unlock");
    RelationStorageMasterMutex.Unlock()
    // and create the basic submaps for
    // the relation storage
    // now we unlock the mutex
    // and return the new id
    fmt.Println("CreateEntityIdent.EntityIdentMutex.Unlock");
    EntityIdentMutex.Unlock()
    return newID, nil
}

func CreateEntity(entity Entity) (int, error){
    // first we lock the entity ident mutex
    // to make sure while we check for the
    // existence it doesnt get deletet, this
    // may sound like a very rare upcoming case,
    //but better be safe than sorry
    fmt.Println("CreateEntity.EntityIdentMutex.RLock");
    EntityIdentMutex.RLock()
    // now
    if _, ok := EntityIdents[entity.Ident]; !ok {
        // the ident doest exist, lets unlock
        // the ident mutex and return -1 for fail0r
        fmt.Println("CreateEntity.EntityIdentMutex.RUnlock");
        EntityIdentMutex.RUnlock()
        return -1, errors.New("CreateEntity.Entity ident not existing");
    }
    // the ident seems to exist, now lets lock the
    // storage mutex before Unlocking the Entity
    // ident mutex to prevent the ident beeing
    // deleted before we start locking (small
    // timing still possible )
    fmt.Println("CreateEntity.EntityIdentMutex.RUnlock");
    EntityIdentMutex.RUnlock()

    // upcount our ID Max and copy it
    // into another variable so we can be sure
    // between unlock of the ressource and return
    // it doesnt get upcounted
    // and set the IDMaxMutex on write Lock
    // lets upcount the entity id max fitting to
    //         [ident]
    fmt.Println("CreateEntity.EntityIDMaxMasterMutex.Lock");
    EntityIDMaxMasterMutex.Lock()
    fmt.Println("CreateEntity.EntityIDMaxMutex[entity.Ident].Lock ",entity.Ident);
    EntityIDMaxMutex[entity.Ident].Lock()
    //fmt.Println("CreateEntity.EntityIDMaxMasterMutex.Unlock");
    //EntityIDMaxMasterMutex.Unlock()
    EntityIDMax[entity.Ident]++
    var newID = EntityIDMax[entity.Ident]
    //fmt.Println("CreateEntity.EntityIDMaxMasterMutex.Lock");
    //EntityIDMaxMasterMutex.Lock()
    fmt.Println("CreateEntity.EntityIDMaxMutex[entity.Ident].Unlock ",entity.Ident);
    EntityIDMaxMutex[entity.Ident].Unlock()
    fmt.Println("CreateEntity.EntityIDMaxMasterMutex.Unlock");
    EntityIDMaxMasterMutex.Unlock()
    // and tell the entity its own id
    entity.ID = newID
    // now we store the entity element
    // in the EntityStorage
    fmt.Println("CreateEntity.EntityStorageMasterMutex.Lock");
    EntityStorageMasterMutex.Lock()
    fmt.Println("CreateEntity.EntityStorageMutex[entity.Ident].Lock ",entity.Ident);
    EntityStorageMutex[entity.Ident].Lock()
    fmt.Println("CreateEntity.EntityStorageMasterMutex.Unlock");
    EntityStorageMasterMutex.Unlock()
    EntityStorage[entity.Ident][newID] = entity
    fmt.Println("CreateEntity.EntityStorageMasterMutex.Lock");
    EntityStorageMasterMutex.Lock()
    fmt.Println("CreateEntity.EntityStorageMutex[entity.Ident].Unlock ",entity.Ident);
    EntityStorageMutex[entity.Ident].Unlock()
    fmt.Println("CreateEntity.EntityStorageMasterMutex.Unlock");
    EntityStorageMasterMutex.Unlock()
    // create the mutex for our ressource on
    // relation. we have to create the sub maps too
    // golang things....
    fmt.Println("CreateEntity.RelationStorageMasterMutex.Lock");
    RelationStorageMasterMutex.Lock()
    fmt.Println("CreateEntity.RelationStorageMutex[entity.Ident].Lock ",entity.Ident);
    RelationStorageMutex[entity.Ident].Lock()
    fmt.Println("CreateEntity.RelationStorageMasterMutex.Unlock");
    RelationStorageMasterMutex.Unlock()
    RelationStorage[entity.Ident][newID]     = make(map[int]map[int]Relation)
    RelationRStorage[entity.Ident][newID]    = make(map[int]map[int]bool)
    fmt.Println("CreateEntity.RelationStorageMasterMutex.Lock");
    RelationStorageMasterMutex.Lock()
    fmt.Println("CreateEntity.RelationStorageMutex[entity.Ident].Unlock ",entity.Ident);
    RelationStorageMutex[entity.Ident].Unlock()
    fmt.Println("CreateEntity.RelationStorageMasterMutex.Unlock");
    RelationStorageMasterMutex.Unlock()
    // since we now stored the entity and created all
    // needed ressources we can unlock
    // the storage ressource and return the ID (or err)
    return newID, nil
}

func GetEntityByPath(ident int, id int) (Entity, error){
    // lets check if entity witrh the given path exists
    fmt.Println("GetEntityByPath.EntityStorageMasterMutex.Lock");
    EntityStorageMasterMutex.Lock()
    fmt.Println("GetEntityByPath.EntityStorageMutex[ident].RLock ",ident);
    EntityStorageMutex[ident].RLock()
    fmt.Println("GetEntityByPath.EntityStorageMasterMutex.Unlock");
    EntityStorageMasterMutex.Unlock()
    if entity, ok := EntityStorage[ident][id]; ok {
        // if yes we return the entity
        // and nil for error
        fmt.Println("GetEntityByPath.EntityStorageMasterMutex.Lock");
        EntityStorageMasterMutex.Lock()
        fmt.Println("GetEntityByPath.EntityStorageMutex[ident].RUnlock ",ident);
        EntityStorageMutex[ident].RUnlock()
        fmt.Println("GetEntityByPath.EntityStorageMasterMutex.Unlock");
        EntityStorageMasterMutex.Unlock()
        return entity, nil
    }
    fmt.Println("GetEntityByPath.EntityStorageMasterMutex.Lock");
    EntityStorageMasterMutex.Lock()
    fmt.Println("GetEntityByPath.EntityStorageMutex[ident].RUnlock ",ident);
    EntityStorageMutex[ident].RUnlock()
    fmt.Println("GetEntityByPath.EntityStorageMasterMutex.Unlock");
    EntityStorageMasterMutex.Unlock()
    // the path seems to result empty , so
    // we throw an error
    return Entity{}, errors.New("Entity on given path does not exist.")
}



func CreateRelation(srcIdent int, srcID int, targetIdent int, targetID int, relation Relation) (bool,error) {
    // first we Readlock the EntityIdentMutex
    fmt.Println("CreateRelation.EntityIdentMutex.RLock");
    EntityIdentMutex.RLock()
    // lets make sure the source ident exist
    if _,ok := EntityIdents[srcIdent] ; !ok {
        fmt.Println("CreateRelation.EntityIdentMutex.RUnlock");
        EntityIdentMutex.RUnlock()
        return false, errors.New("Source ident not existing")
    }
    // and the target ident exists too
    if _,ok := EntityIdents[targetIdent] ; !ok {
        fmt.Println("CreateRelation.EntityIdentMutex.RUnlock");
        EntityIdentMutex.RUnlock()
        return false, errors.New("Target ident not existing")
    }
    //// - - - - - - - - - - - - - - - - - 
    //// ### test outcommented to review
    //// Identities seem to exist, we unlock the ident mutex
    //EntityIdentMutex.RUnlock()
    //// Now we lock the to link entities to
    //// make sure they dont get deletet meawhile
    //EntityStorageMasterMutex.Lock()
    //EntityStorageMutex[srcIdent].Lock() 
    //EntityStorageMasterMutex.Unlock()
    //// if srcIdent and targetIdent differ,
    //// we lock targetIdent too else we
    //// would create deadlock
    //if srcIdent != targetIdent {
    //    EntityStorageMasterMutex.Lock()
    //    EntityStorageMutex[targetIdent].Lock()
    //    EntityStorageMasterMutex.Unlock()
    //}
    //// ### test outcommented to review
    //// - - - - - - - - - - - - - - - - - 
    // now we lock the relation mutex
    fmt.Println("CreateRelation.RelationStorageMasterMutex.Lock");
    RelationStorageMasterMutex.Lock()
    fmt.Println("CreateRelation.RelationStorageMutex[srcIdent].Lock ",srcIdent);
    RelationStorageMutex[srcIdent].Lock()
    fmt.Println("CreateRelation.RelationStorageMasterMutex.Unlock");
    RelationStorageMasterMutex.Unlock()
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
        RelationRStorage[targetIdent][targetID][srcIdent] = make(map[int]bool)
    }
    // now we store the relation
    RelationStorage[srcIdent][srcID][targetIdent][targetID] = relation
    // and an entry into the reverse index, its existence
    // allows us to use the coords in the normal index to revtrieve
    // the Relation. We dont create a pointer because golang doesnt
    // allow pointer on submaps in nested maps
    RelationRStorage[targetIdent][targetID][srcIdent][srcID] = true
    // we are done now we can unlock the entity idents
    //// - - - - - - - - - - - - - - - - - 
    //// ### test outcommented to review
    //EntityStorageMasterMutex.Lock()
    //EntityStorageMutex[srcIdent].Unlock()
    //EntityStorageMasterMutex.Unlock()
    //// if we locked the targetIdent too (see upper)
    //// than we have to unlock it too
    //if srcIdent != targetIdent {
    //    EntityStorageMasterMutex.Lock()
    //    EntityStorageMutex[targetIdent].Unlock()
    //    EntityStorageMasterMutex.Unlock()
    //}
    //// ### test outcommented to review
    //// - - - - - - - - - - - - - - - - 
     //and finally unlock the relation ident and return
    fmt.Println("CreateRelation.RelationStorageMasterMutex.Lock");
    RelationStorageMasterMutex.Lock()
    fmt.Println("CreateRelation.RelationStorageMutex[srcIdent].Unlock", srcIdent);
    RelationStorageMutex[srcIdent].Unlock()
    fmt.Println("CreateRelation.RelationStorageMasterMutex.Unlock");
    RelationStorageMasterMutex.Unlock()
    return true, nil
}


func GetRelationsBySourceIdentAndSourceId(ident int, id int) (map[int]Relation , error) {
    // initialice the return map
    var mapRet = make(map[int]Relation)
    // set counter for the loop
    var cnt    = 0
    // copy the pool we have to search in
    // to prevent crashes on RW concurrency
    // we lock the RelationStorage mutex with
    // fitting ident. this allows us to proceed
    // faster since we just block to copy instead
    // of blocking for the whole process
    fmt.Println("GetRelationsBySourceIdentAndSourceId.RelationStorageMasterMutex.Lock");
    RelationStorageMasterMutex.Lock()
    fmt.Println("GetRelationsBySourceIdentAndSourceId.RelationStorageMutex[ident].RLock ",ident);
    RelationStorageMutex[ident].RLock()
    fmt.Println("GetRelationsBySourceIdentAndSourceId.RelationStorageMasterMutex.Unlock");
    RelationStorageMasterMutex.Unlock()
    var pool   = RelationStorage[ident][id];
    fmt.Println("GetRelationsBySourceIdentAndSourceId.RelationStorageMasterMutex.Lock");
    RelationStorageMasterMutex.Lock()
    fmt.Println("GetRelationsBySourceIdentAndSourceId.RelationStorageMutex[ident].RUnlock ",ident);
    RelationStorageMutex[ident].RUnlock()
    fmt.Println("GetRelationsBySourceIdentAndSourceId.RelationStorageMasterMutex.Unlock");
    RelationStorageMasterMutex.Unlock()
    // for each possible targtIdent
    for _,targetIdentMap := range pool {
        // for each possible targetId per targetIdent
        for _,relation := range targetIdentMap {
            // copy the relation into the return map
            // and upcount the int
            mapRet[cnt] = relation
            cnt++
        }
    }
    // + + + + + + + 
    //fmt.Println("Relations: ",cnt," - ",len(mapRet))
    // + + + + + + + 
    return mapRet, nil
}

func IdentExists(strIdent string) (bool){
    fmt.Println("IdentExists.EntityIdentMutex.RLock");
    EntityIdentMutex.RLock()
    // lets check if this ident exists
    if _,ok := EntityRIdents[strIdent]; ok {
        // it does lets return it
        fmt.Println("IdentExists.EntityIdentMutex.RUnlock");
        EntityIdentMutex.RUnlock()
        return true
    }
    fmt.Println("IdentExists.EntityIdentMutex.RUnlock");
    EntityIdentMutex.RUnlock()
    return false
}

func GetIdentIdByString(strIdent string)(int,error) {
    fmt.Println("IdentExists.EntityIdentMutex.RLock");
    EntityIdentMutex.RLock()
    // lets check if this ident exists
    if id,ok := EntityRIdents[strIdent]; ok {
        // it does lets return it
        fmt.Println("IdentExists.EntityIdentMutex.RUnlock");
        EntityIdentMutex.RUnlock()
        return id, nil
    }
    fmt.Println("IdentExists.EntityIdentMutex.RUnlock");
    EntityIdentMutex.RUnlock()
    return -1, errors.New("Entity ident string does not exist")
}


func debugPrint(param map[int]Relation) {
    fmt.Println("- - - - - - - - - - \n")
    out, _    := json.MarshalIndent(param, "", "  ")
    fmt.Print(string(out))  
    fmt.Println("- - - - - - - - - - \n")
}
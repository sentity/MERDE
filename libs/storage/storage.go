package storage

// handle all the imports
import (
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
var EntityStorage = make(map[int]map[int]Entity)

// entity storage master mutex
var EntityStorageMutex = &sync.RWMutex{}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// entity storage id max         [Ident]
var EntityIDMax = make(map[int]int)

// master mutexd for EntityIdMax
var EntityIDMaxMutex = &sync.RWMutex{}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// maps to translate Idents to their INT and reverse
var EntityIdents = make(map[int]string)
var EntityRIdents = make(map[string]int)

// and a fitting max ID
var EntityIdentIDMax int = 0

// entity ident mutex (for adding and deleting ident types)
var EntityIdentMutex = &sync.RWMutex{}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// relation struct
type Relation struct {
	SourceIdent int
	SourceID    int
	TargetIdent int
	TargetID    int
	Context     string
	Properties  map[string]string
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// s prefix = source
// t prefix = target
// relation storage map             [sIdent][sId]   [tIdent][tId]
var RelationStorage = make(map[int]map[int]map[int]map[int]Relation)

// and relation reverse storage map
// (for faster queries)              [tIdent][Tid]   [sIdent][sId]
var RelationRStorage = make(map[int]map[int]map[int]map[int]bool)

// relation storage master mutex
var RelationStorageMutex = &sync.RWMutex{}

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
func CreateEntityIdent(name string) (int, error) {
	// first of allw e lock
	//printMutexActions("CreateEntityIdent.EntityIdentMutex.Lock");
	EntityIdentMutex.Lock()
	// lets check if the ident allready exists
	// if it does we just return the ID
	if id, ok := EntityRIdents[name]; ok {
		// dont forget to unlock
		//printMutexActions("CreateEntityIdent.EntityIdentMutex.Unlock");
		EntityIdentMutex.Unlock()
		return id, nil
	}
	// ok entity doesnt exist yet, lets
	// upcount our ID Max and copy it
	// into another variable so we can be sure
	// between unlock of the ressource and return
	// it doesnt get upcounted
	EntityIdentIDMax++
	var newID = EntityIdentIDMax
	// finally create the new ident in our
	// EntityIdents index and reverse index
	EntityIdents[newID] = name
	EntityRIdents[name] = newID
	// and create mutex for EntityStorage ident+type
	//printMutexActions("CreateEntityIdent.EntityStorageMutex.Lock");
	EntityStorageMutex.Lock()
	// now we prepare the submaps in the entity
	// storage itseöf....
	EntityStorage[newID] = make(map[int]Entity)
	// set the maxID for the new
	// ident type
	EntityIDMax[newID] = 0
	//printMutexActions("CreateEntityIdent.EntityStorageMutex.Unlock");
	EntityStorageMutex.Unlock()
	// create the base maps in relation storage
	//printMutexActions("CreateEntityIdent.RelationStorageMutex.Lock");
	RelationStorageMutex.Lock()
	RelationStorage[newID] = make(map[int]map[int]map[int]Relation)
	RelationRStorage[newID] = make(map[int]map[int]map[int]bool)
	//printMutexActions("CreateEntityIdent.RelationStorageMutex.Unlock");
	RelationStorageMutex.Unlock()
	// and create the basic submaps for
	// the relation storage
	// now we unlock the mutex
	// and return the new id
	//printMutexActions("CreateEntityIdent.EntityIdentMutex.Unlock");
	EntityIdentMutex.Unlock()
	return newID, nil
}

func CreateEntity(entity Entity) (int, error) {
	// first we lock the entity ident mutex
	// to make sure while we check for the
	// existence it doesnt get deletet, this
	// may sound like a very rare upcoming case,
	//but better be safe than sorry
	//printMutexActions("CreateEntity.EntityIdentMutex.RLock");
	EntityIdentMutex.RLock()
	// now
	if _, ok := EntityIdents[entity.Ident]; !ok {
		// the ident doest exist, lets unlock
		// the ident mutex and return -1 for fail0r
		//printMutexActions("CreateEntity.EntityIdentMutex.RUnlock");
		EntityIdentMutex.RUnlock()
		return -1, errors.New("CreateEntity.Entity ident not existing")
	}
	// the ident seems to exist, now lets lock the
	// storage mutex before Unlocking the Entity
	// ident mutex to prevent the ident beeing
	// deleted before we start locking (small
	// timing still possible )
	//printMutexActions("CreateEntity.EntityIdentMutex.RUnlock");
	EntityIdentMutex.RUnlock()
	// upcount our ID Max and copy it
	// into another variable so we can be sure
	// between unlock of the ressource and return
	// it doesnt get upcounted
	// and set the IDMaxMutex on write Lock
	// lets upcount the entity id max fitting to
	//         [ident]
	//printMutexActions("CreateEntity.EntityStorageMutex.Lock");
	EntityStorageMutex.Lock()
	//fmt.Println("CreateEntity.EntityIDMaxMasterMutex.Unlock");
	//EntityIDMaxMasterMutex.Unlock()
	EntityIDMax[entity.Ident]++
	var newID = EntityIDMax[entity.Ident]
	//fmt.Println("CreateEntity.EntityIDMaxMasterMutex.Lock");
	//EntityIDMaxMasterMutex.Lock()
	// and tell the entity its own id
	entity.ID = newID
	// now we store the entity element
	// in the EntityStorage
	EntityStorage[entity.Ident][newID] = entity
	//printMutexActions("CreateEntity.EntityStorageMutex.Unlock");
	EntityStorageMutex.Unlock()
	// create the mutex for our ressource on
	// relation. we have to create the sub maps too
	// golang things....
	//printMutexActions("CreateEntity.RelationStorageMutex.Lock");
	RelationStorageMutex.Lock()
	RelationStorage[entity.Ident][newID] = make(map[int]map[int]Relation)
	RelationRStorage[entity.Ident][newID] = make(map[int]map[int]bool)
	//printMutexActions("CreateEntity.RelationStorageMutex.Unlock");
	RelationStorageMutex.Unlock()
	// since we now stored the entity and created all
	// needed ressources we can unlock
	// the storage ressource and return the ID (or err)
	return newID, nil
}

func GetEntityByPath(ident int, id int) (Entity, error) {
	// lets check if entity witrh the given path exists
	//printMutexActions("GetEntityByPath.EntityStorageMutex.Lock");
	EntityStorageMutex.Lock()
	if entity, ok := EntityStorage[ident][id]; ok {
		// if yes we return the entity
		// and nil for error
		//printMutexActions("GetEntityByPath.EntityStorageMutex.Unlock");
		EntityStorageMutex.Unlock()
		return entity, nil
	}
	//printMutexActions("GetEntityByPath.EntityStorageMutex.Unlock");
	EntityStorageMutex.Unlock()
	// the path seems to result empty , so
	// we throw an error
	return Entity{}, errors.New("Entity on given path does not exist.")
}

func CreateRelation(srcIdent int, srcID int, targetIdent int, targetID int, relation Relation) (bool, error) {
	// first we Readlock the EntityIdentMutex
	//printMutexActions("CreateRelation.EntityIdentMutex.RLock");
	EntityIdentMutex.RLock()
	// lets make sure the source ident exist
	if _, ok := EntityIdents[srcIdent]; !ok {
		//printMutexActions("CreateRelation.EntityIdentMutex.RUnlock");
		EntityIdentMutex.RUnlock()
		return false, errors.New("Source ident not existing")
	}
	// and the target ident exists too
	if _, ok := EntityIdents[targetIdent]; !ok {
		//printMutexActions("CreateRelation.EntityIdentMutex.RUnlock");
		EntityIdentMutex.RUnlock()
		return false, errors.New("Target ident not existing")
	}
	//// - - - - - - - - - - - - - - - - -
	// now we lock the relation mutex
	//printMutexActions("CreateRelation.RelationStorageMutex.Lock");
	RelationStorageMutex.Lock()
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
	if _, ok := RelationRStorage[targetIdent][targetID][srcIdent]; !ok {
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
	//// - - - - - - - - - - - - - - - -
	//and finally unlock the relation ident and return
	//printMutexActions("CreateRelation.RelationStorageMutex.Unlock");
	RelationStorageMutex.Unlock()
	return true, nil
}

func GetRelationsBySourceIdentAndSourceId(ident int, id int) (map[int]Relation, error) {
	// initialice the return map
	var mapRet = make(map[int]Relation)
	// set counter for the loop
	var cnt = 0
	// copy the pool we have to search in
	// to prevent crashes on RW concurrency
	// we lock the RelationStorage mutex with
	// fitting ident. this allows us to proceed
	// faster since we just block to copy instead
	// of blocking for the whole process
	//printMutexActions("GetRelationsBySourceIdentAndSourceId.RelationStorageMutex.Lock");
	RelationStorageMutex.Lock()
	var pool = RelationStorage[ident][id]
	//printMutexActions("GetRelationsBySourceIdentAndSourceId.RelationStorageMutex.Unlock");
	RelationStorageMutex.Unlock()
	// for each possible targtIdent
	for _, targetIdentMap := range pool {
		// for each possible targetId per targetIdent
		for _, relation := range targetIdentMap {
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

func IdentExists(strIdent string) bool {
	//printMutexActions("IdentExists.EntityIdentMutex.RLock");
	EntityIdentMutex.RLock()
	// lets check if this ident exists
	if _, ok := EntityRIdents[strIdent]; ok {
		// it does lets return it
		//printMutexActions("IdentExists.EntityIdentMutex.RUnlock");
		EntityIdentMutex.RUnlock()
		return true
	}
	//printMutexActions("IdentExists.EntityIdentMutex.RUnlock");
	EntityIdentMutex.RUnlock()
	return false
}

func GetIdentIdByString(strIdent string) (int, error) {
	//printMutexActions("IdentExists.EntityIdentMutex.RLock");
	EntityIdentMutex.RLock()
	// lets check if this ident exists
	if id, ok := EntityRIdents[strIdent]; ok {
		// it does lets return it
		//printMutexActions("IdentExists.EntityIdentMutex.RUnlock");
		EntityIdentMutex.RUnlock()
		return id, nil
	}
	//printMutexActions("IdentExists.EntityIdentMutex.RUnlock");
	EntityIdentMutex.RUnlock()
	return -1, errors.New("Entity ident string does not exist")
}

func printMutexActions(param string) {
	//fmt.Println(param)
}

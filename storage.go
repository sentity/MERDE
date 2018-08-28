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
// entity storage map          [Ident] [Type]  [ID]
var entityStorage    = make(map[int]map[int]map[int]Entity)

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// entity storage map          [Ident] [Type]  [ID]
var entityIDMax      = make(map[int]map[int]map[int]int)

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// entity storage mutex
var entityStoragemutex = &sync.Mutex{}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// maps to translate Idents to their INT and reverse
var entityIdents     = make(map[string]int)
var entityRIdents    = make(map[int]string)
// and a fitting max ID
var entityIdentIDMax = int

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
var relationStorage  = make(map[int]map[string]map[string]map[int]Relation)

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// relation reverse storage map
// (for faster queries)        [Type]  [Target]   [source]   [ID] [path]
var relationRStorage = make(map[int]map[string]map[string]map[int]string)

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// relation index max id       [Type]  [Source]   [Target]   [ID]
var relationIDMax    = make(map[int]map[string]map[string]map[int]int)


// - - - - - - - - - - - - - - - - - - - - - - - - - -
// + + + + + + FUNCTIONS + + + + + + 
// - - - - - - - - - - - - - - - - - - - - - - - - - -

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// 
func createEntityIdent() {
    
}

func deleteEntityIdent() {
    
}

func createEntity() {
    
}

func getEntity() {
    
}

func getEntityByPath() {
    
}

func getEntityByIdent() {
    
}

func deleteEntity() {
    
}

func updateEntity() {
    
}

func createRelation() {
    
}

func getRelation() {
    
}

func deleteRelation() {
    
}

func updateRelation() {
    
}


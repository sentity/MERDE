package mapper

// handle all the imports
import (
	//"fmt"
	//"os"
    //"math/rand"
    //"time"
    //"strconv"
    //"strings"
    //"builtin"
    "encoding/json"
    "errors"
    //"sync"
    "storage"
)


// - - - - - - - - - - - - - - - - - - - - - - - - - -
// rentity transport struct for json<->entity
type Entity struct {
    ID           int
    Ident        string
    Context      string
    Value        string
    Properties   map[string]string
    Children     map[int]Entity
}


// - - - - - - - - - - - - - - - - - - - - - - - - - -
// relation transport struct 
type Relatiion struct {
    SourceIdent  int
    SourceID     int
    TargetIdent  int
    TargetID     int
    Context      string
    Properties   map[string]string
}


func MapJson(data []byte) ( int,  error) {
    // define an entity struct as template
    // for json unmarshal and parse the
    // json byte array into the var
    var entity Entity
    if err := json.Unmarshal(data, &entity); err != nil {
        return -1 , errors.New("Invalid Json")
    }
    // + + + + + + +
    //fmt.Printf("%#v", entity.Children[1].Children[2]) 
    // + + + + + + +
    newID, err:=MapEntitiesRecursive(entity,-1,-1)
    if err != nil {
        return -1, errors.New("Couldnt map entities .... why tho?")
    }
    return newID, nil
}

func MapEntitiesRecursive(entity Entity,parentIdent int,parentID int ) ( int,  error) {
    // first we get the right identID
    var identID = HandleIdent(entity.Ident)
    // now we create the fitting entity
    tmpEntity := storage.Entity{
        ID         : -1,
        Ident      : identID,
        Value      : entity.Value,
        Context    : entity.Context,
        Properties : entity.Properties,
    }
    // now we create the entity
    var newID, _ = storage.CreateEntity(tmpEntity)
    // lets check if there are child elements
    if len(entity.Children) != 0 {
        // there are children lets iteater over
        // the map 
        for _, childEntity := range entity.Children {
            // pas the child entity and the parent coords to
            // create the relation after inserting the entity
            MapEntitiesRecursive(childEntity, identID, newID)
        }
    }
    // now lets check if ourparent ident and id
    // are not -1 , if so we need to create
    // a relation
    if parentIdent != -1 && parentID != -1 {
        // lets create the relation to our parent
        tmpRelation := storage.Relation{
            SourceIdent : parentIdent,
            SourceID    : parentID,
            TargetIdent : tmpEntity.Ident,
            TargetID    : newID,
        }
        storage.CreateRelation(parentIdent, parentID, tmpEntity.Ident, newID, tmpRelation)
    }
    // only the first return is interesting since it
    // returns the most parent id
    return newID, nil
}

func HandleIdent (strIdent string) (int){
    // lets check if this ident exists
    if id,ok := storage.EntityRIdents[strIdent]; ok {
        // it does lets return it
        return id
    }
    newID,_ := storage.CreateEntityIdent(strIdent)
    return newID
}

func GetEntityRecursive(entityIdent int, entityID int) (Entity, error){
    // first we get the entity
    entity, err := storage.GetEntityByPath(entityIdent, entityID)
    // if it doesnt exist we stop
    if err != nil {
        return Entity{}, errors.New("Entity with give path does not exist")
    }
    // now we define the data into
    // mapper entity struct
    returnEntity := Entity{
        ID      : entity.ID,
        Ident   : storage.EntityIdents[entity.Ident],
        Value   : entity.Value,
        Context : entity.Context,
        Properties : entity.Properties,
    }
    // now we retrieve all relations
    // from this entity
    tmpRelations, _ := storage.GetRelationsBySourceIdentAndSourceId(entity.Ident,entity.ID)
    // lets check if we got anyrelations
    if len(tmpRelations) != 0 {
        returnEntity.Children = make(map[int]Entity)
        // there are children lets iterate 
        // through the map 
        var i = 0
        for _, tmpRelation := range tmpRelations {
            i++
            // call the function recursive and add the object
            var tmpEntity,_ = GetEntityRecursive(tmpRelation.TargetIdent,
                                                 tmpRelation.TargetID)
            // store the subentity in child field
            returnEntity.Children[i] = tmpEntity
        }
    }
    // return the entity
    return returnEntity, nil
}



package query

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"goer/libs/mapper"
	"os"
	"strconv"
	"strings"
)

type Condition struct {
	Alpha    string
	Beta     string
	Operator string
}

type Query struct {
	Type      string
	Ident     string
	Mode      string
	Imply     int
	Direction string
	Set       map[string]string
	Where     []string
	From      []int
	To        []int
	Traverse  int
	Rejoin    int
}

type ResultSet struct {
	Data []ResultAddress
}

type ResultAddress struct {
	Ident int
	Id    int
}

type Filter struct {
	Id   int
	Sets []ResultSet
	Mode string
}

func HandleQuery(query string) {
	// unmarshal the json input
	var jsonData []Query
	if err := json.Unmarshal([]byte(query), &jsonData); err != nil {
		panic(err)
	}
	queryCount := len(jsonData)
	// if there really are queries
	if queryCount > 0 {
		results, err := ProcessQuery(jsonData)
	}

}

func ProcessQuery(multiQuery []Query) (ResultSet, error) {
	// prepare an empty result set
	ResultStorage := make(map[int]ResultSet)
	FilterStorage := make(map[int]Filter)
	// iterate through each of the multiquery to
	// handle them, always providing the last results
	for queryId, query := range multiQuery {
		// first we get the info what type of query
		// we are running | maybe validate later ###
		qtype := strings.Split(query.Type, ".")
		// if we are doin a system query we hard dispatch here
		if qtype[0] == "system" {
			SysMessage := SystemCall(query)
			return ResultSet{}, errors.New(SysMessage)
		}
		// now , if we got an entity
		if qtype[0] == "entity" {
			HandleEntityQuery()

		}

		results, err := ProcessQuery(query, results)
		if err != nil {
			return ResultSet{}, err
		}
		ret = results
	}
	// lets check if our final option is a find
	return results, nil
}

func HandleEntityQuery(action string, query Query, results ResultSet) {
	switch action {
	case "find":
	case "update":
	case "create":
	case "delete":
	default:
		// return unknown action fail
	}
}

func SystemCall(query Query) (Result string) {
	return ""
}

func ProcessQuery(query Query, results ResultSet) (ResultSet, error) {

	if len(qtype) == 2 {
		switch qtype[0] {
		case "entity":
			HandleEntityQuery(qtype[1], query, results)
		//case "relation":
		//	HandleRelationQuery(qtype[1], query, results)
		case "system":
			HandleSystemQuery(qtype[1], query, results)
		default:
			// return error unknown query type resource
		}
	} else {
		// return error invalid query type form
	}
}

func buildEntityReturn(results ResultSet, traverse int) []mapper.Entity {

}

func HandleRelationQuery(action string, query Query, results ResultSet) {
	switch action {
	case "find":
	case "update":
	case "create":
	case "delete":
	default:
		// return unknown action fail
	}
}

func HandleSystemQuery(action string, query Query, results ResultSet) {
	// no idea yet , will add later
}

func ParseConditions(condition string) (map[int]Condition, error) {
	arrCondition, valueMap := ExtractParamStrings(condition)
	//fmt.Print(parsedCondition)
	arrOperators := [7]string{"==", "!=", ">=", ">=", "%=", ">", "<"}
	//arrCondition      := strings.Split(parsedCondition,"&&")
	arrReturn := make(map[int]Condition)
	retCnt := int(0)
	cache := make(map[string]string)
	names := [2]string{"Alpha", "Beta"}
	tmp := ""
	for _, conditionPart := range arrCondition {
		for _, operator := range arrOperators {
			if strings.Contains(conditionPart, operator) {
				arrOperatableSub := strings.Split(conditionPart, operator)
				if len(arrOperatableSub) == 2 {
					for subId, operatableSub := range arrOperatableSub {
						//DebugPrint(valueMap)
						// check if var has to be resolved from valuemap
						if strings.Contains(operatableSub, "'") {
							tmp = strings.Replace(operatableSub, "'", "", 2)
							cnt, err := strconv.Atoi(tmp)
							if err != nil {
								return nil, errors.New("malformed condition")
							}
							// add check if value map entry exists... or not... im fkn sleepy
							cache[names[subId]] = UnescapeConditionValue(valueMap[cnt])
						} else {
							cache[names[subId]] = operatableSub
						}
					}
					tmpCondition := Condition{
						Alpha:    cache["Alpha"],
						Beta:     cache["Beta"],
						Operator: operator,
					}
					arrReturn[retCnt] = tmpCondition
					retCnt++
					break
				}
			}
		}
		//return nil, errors.New("condition with malformed or non existing operator")
	}
	if len(arrReturn) > 0 {
		return arrReturn, nil
	}
	return nil, errors.New("unknown error occurred parsing the condition")
}

func UnescapeConditionValue(value string) string {
	value = strings.Replace(value, "\\'", "'", 1024)
	return value
}

func ExtractParamStrings(condition string) (map[int]string, map[int]string) {
	valueMap := make(map[int]string)
	parsedCond := new(bytes.Buffer)
	param := new(bytes.Buffer)
	cnt := int(0)
	conditionMap := make(map[int]string)
	arrCondition := strings.Split(condition, "&&")
	for key, condition := range arrCondition {
		if strings.Contains(condition, "'") {
			cnt++
			arrCondition := strings.Split(condition, "'")
			splitAmount := len(arrCondition)
			for id, splitPart := range arrCondition {
				if id == 0 {
					(*parsedCond).WriteString(splitPart)
					(*parsedCond).WriteString("'")
					(*parsedCond).WriteString(strconv.Itoa(cnt))
				} else {
					if id == splitAmount {
						(*parsedCond).WriteString("'")
						(*parsedCond).WriteString(splitPart)
					} else {
						(*param).WriteString(splitPart)
					}
				}
			}
			conditionMap[key] = (*parsedCond).String()
			valueMap[cnt] = (*param).String()
			//fmt.Println("Extracted param: ",valueMap[cnt]," with key " , cnt)
			(*parsedCond).Reset()
			(*param).Reset()
		} else {
			conditionMap[key] = condition
		}
	}
	return conditionMap, valueMap
}

func MazeParseConditions(condition string) []Condition {

	// base consts
	comparisonOperators := []string{"==", "!=", ">=", ">=", "%=", ">", "<"}

	// result container
	conditionalStringMap := make(map[int]string)
	var logicalResults []string
	var finalResults []Condition

	// variables from current test
	charArr := strings.Split(condition, "")
	stringLength := len(condition)

	// temporary variables for ' parsing
	stringMode := false
	skipNext := false
	stringReplaceConst := "__TEXT__"
	intermediateString := new(bytes.Buffer)
	stringModeString := new(bytes.Buffer)

	//fmt.Println("Current Test: ", condition);

	// parse string char wise and replace '-strings by replacements
	for i, currentChar := range charArr {
		if skipNext {
			skipNext = false
			continue
		}
		hasLookaheadChar := (i + 1) < stringLength
		nextChar := ""
		if hasLookaheadChar {
			nextChar = string(condition[i+1])
		}
		if currentChar == "'" {
			if stringMode {
				currentMapIndex := len(conditionalStringMap)
				conditionalStringMap[currentMapIndex] = stringModeString.String()
				stringModeString.Reset()
				intermediateString.WriteString(stringReplaceConst)
				intermediateString.WriteString(strconv.Itoa(currentMapIndex))
			}
			stringMode = !stringMode
			continue
		}
		if currentChar == "\\" {
			if hasLookaheadChar && stringMode && nextChar == "'" {
				stringModeString.WriteString("'")
				skipNext = true
				continue
			}
			if hasLookaheadChar && stringMode && nextChar == "\\" {
				stringModeString.WriteString("\\")
				skipNext = true
				continue
			}
		}
		if stringMode {
			stringModeString.WriteString(currentChar)

		} else {
			// poor mans trim
			if currentChar == " " {
				continue
			} else if currentChar == "&" && nextChar == "&" {
				logicalResults = append(logicalResults, intermediateString.String())
				intermediateString.Reset()
				skipNext = true
			} else {
				intermediateString.WriteString(currentChar)
			}
		}
	}
	logicalResults = append(logicalResults, intermediateString.String())

	// fmt.Println(intermediateString);
	// fmt.Println("replacements");
	// for i, currentString := range conditionalStringMap {
	// 		fmt.Println(i, ": ",currentString );
	// }

	found := false
	for _, logicalResult := range logicalResults {
		for _, comparisonOperator := range comparisonOperators {
			tmp := strings.Split(logicalResult, comparisonOperator)
			if len(tmp) != 2 {
				continue
			}

			for i, tmpEntry := range tmp {
				if strings.HasPrefix(tmpEntry, stringReplaceConst) {
					replacementNumber, err := strconv.Atoi(tmpEntry[len(stringReplaceConst):len(tmpEntry)])
					if err != nil {
						// handle error
						fmt.Println(err)
						os.Exit(1)
					}
					tmp[i] = conditionalStringMap[replacementNumber]
				}
			}

			finalResults = append(finalResults, Condition{
				Alpha:    tmp[0],
				Beta:     tmp[1],
				Operator: comparisonOperator,
			})
			found = true
		}
		if !found {
			//fmt.Println("malformed result, no conditional");
			os.Exit(1)
		}
		found = false
	}
	return finalResults
}

func DebugPrint(param map[int]Condition) {
	//func DebugPrint(param []Condition) {
	fmt.Println("- - - - - - - - - - \n")
	out, _ := json.MarshalIndent(param, "", "  ")
	fmt.Print(string(out))
	fmt.Println("- - - - - - - - - - \n")
}

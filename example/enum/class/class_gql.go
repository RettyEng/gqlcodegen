// DO NOT EDIT. this file is generated by gqlcodegen.
package class

import (
	"errors"
	"strconv"
)

/*
   Description:
     """
     Driver class
     """
   Directives:
     @legend()
*/
type Class int

const (
	ROOKIE Class = iota
	ELITE

	/*
	   Directives:
	     @special()
	*/
	KING_OF_ROAD
	LEGEND
)

const _Class_Name = "ROOKIEELITEKING_OF_ROADLEGEND"

var _Class_Index = []int{0, 6, 11, 23, 29}

func (v Class) String() string {
	if v < 0 || v >= Class(len(_Class_Index)-1) {
		return "Class(" + strconv.FormatInt(int64(v), 10) + ")"
	}
	return _Class_Name[_Class_Index[v]:_Class_Index[v+1]]
}

func _ClassFromString(str string) (Class, error) {
	for i := 0; i < len(_Class_Index)-1; i++ {
		if v := Class(i); str == v.String() {
			return v, nil
		}
	}
	return -1, errors.New(str + " is not found")
}

func (Class) ImplementsGraphQLType(name string) bool {
	return name == "Class"
}

func (v *Class) UnmarshalGraphQL(input interface{}) error {
	switch input := input.(type) {
	case string:
		value, err := _ClassFromString(input)
		if err != nil {
			return err
		}
		*v = value
		return nil
	default:
		return errors.New("wrong type")
	}
}

func (v Class) MarshalJSON() ([]byte, error) {
	return []byte(`"` + v.String() + `"`), nil
}

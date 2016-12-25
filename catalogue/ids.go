package catalogue

import (
	"reflect"
	"strings"

	"github.com/satori/go.uuid"
)

type ID string

var (
	zeroID = ID("")
	idType = reflect.TypeOf(zeroID)
)

func GenerateID() ID {
	return ID(uuid.NewV4())
}

// EnsureID checks if v contains a field named "ID" of catalogue.ID type, and
// (if empty) sets its contents to a new ID. It does nothing otherwise, so it can
// be safely used on every type.
func EnsureID(v interface{}) ID {
	filterFunc := func(name string) bool {
		return strings.EqualFold(name, "id")
	}

	vType := reflect.TypeOf(v)
	vValue := reflect.ValueOf(v)
	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
		vValue = vValue.Elem()
	}

	if vType.Kind() == reflect.Struct {
		if structField, ok := vType.FieldByNameFunc(filterFunc); ok {
			if reflect.DeepEqual(structField.Type, idType) {
				fValue := vValue.FieldByNameFunc(filterFunc)
				if reflect.DeepEqual(fValue.Interface(), zeroID) {
					newID := GenerateID()

					fValue.SetString(newID)
					return newID
				}
			}
		}
	}

	return zeroID
}

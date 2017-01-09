package catalogue

import (
	"reflect"
	"strings"

	"github.com/satori/go.uuid"
)

func GenerateID() string {
	return uuid.NewV4().String()
}

// EnsureID checks if v contains a field named "ID" of string type, and
// (if empty) sets its contents to a new ID. It does nothing otherwise, so it can
// be safely used on every type.
func EnsureID(v interface{}) string {
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
			if structField.Type.Kind() == reflect.String {
				fValue := vValue.FieldByNameFunc(filterFunc)
				if fValue.Interface().(string) == "" {
					newID := GenerateID()

					fValue.SetString(newID)
					return newID
				}
			}
		}
	}

	return ""
}

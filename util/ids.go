/*
 *  Copyright (c) 2017 Open Baton (http://openbaton.org)
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package util

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

package recap

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func Parse(item interface{}, regex string, value string) (error, bool) {
	rx := regexp.MustCompile(regex)

	match := rx.FindStringSubmatch(value)
	m := map[string]string{}

	if match == nil {
		if rx.MatchString(value) {
			return nil, true
		}
		return nil, false
	}

	for i, name := range rx.SubexpNames() {
		if i > 0 && i < len(match) {
			m[name] = match[i]
		}
	}

	err := mapRegexGroupsToStruct(item, m, "")
	if err != nil {
		return err, false
	} else {
		return nil, true
	}
}

func mapRegexGroupsToStruct(item interface{}, values map[string]string, prefix string) error {
	v := reflect.ValueOf(item).Elem()

	var err error

	if !v.CanAddr() {
		return fmt.Errorf("cannot assign to the item passed, item must be a pointer in order to assign")
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)
		
		if field.Kind() == reflect.Struct {
			v := field.Addr()
			tag := fieldType.Tag
			subPrefix := prefix
			
			if tagValue, ok := tag.Lookup("recap"); ok {
				subPrefix = tagValue
			}

			err = mapRegexGroupsToStruct(v.Interface(), values, subPrefix)

			if err != nil {
				return err
			}
			continue
		}

		tag := fieldType.Tag

		if tagVal, ok := tag.Lookup("recap"); ok {
			parts := strings.Split(tagVal,";")
			valName := prefix+parts[0]

			if rxVal, ok := values[valName]; ok {
				if (len(parts) > 1) {
					rxVal = parseConditions(parts[1], rxVal)
				}

				typeVal, err := convertToType(rxVal, fieldType.Type.Name())
				if err != nil {
					return err
				} else {
					field.Set(reflect.ValueOf(typeVal))
				}
			}
		}
	}

	return err
}

func parseConditions(conds string, value string) string {
	parts := strings.Split(conds, ",")

	for _, v := range parts {
		method := strings.Split(v, "=")

		if len(method) < 2 {
			continue
		}

		if (method[0] == "contains") {
			return strconv.FormatBool(strings.Contains(value, method[1]))
		}
		if (method[0] == "default") {
			if len(value) == 0 {
				return method[1]
			} else {
				return value
			}
		}
	}

	return "false"
}

func convertToType(value string, valueType string) (interface{}, error) {
	var intVal int64
	var uintVal uint64
	var floatVal float64
	var boolVal bool
	var retVal interface{}
	var err error

	if (strings.HasPrefix(valueType, "int")) {
		intVal, err = strconv.ParseInt(value, 10, 64)
	} else if (strings.HasPrefix(valueType, "uint")) {
		uintVal, err = strconv.ParseUint(value, 10, 64)
	} else if (strings.HasPrefix(valueType, "float")) {
		floatVal, err = strconv.ParseFloat(value, 64)
	} else if valueType == "string" {
		return value, nil
	} else if valueType == "bool" {
		boolVal, err = strconv.ParseBool(value)
	}

	switch (valueType) {
	case "bool":
		retVal = boolVal
	case "int":
		retVal = int(intVal)
	case "int8":
		retVal = int8(intVal)
	case "int16":
		retVal = int16(intVal)
	case "int32":
		retVal = int32(intVal)
	case "int64":
		retVal = int64(intVal)
	case "uint":
		retVal = uint(uintVal)
	case "uint8":
		retVal = uint8(uintVal)
	case "uint16":
		retVal = uint16(uintVal)
	case "uint32":
		retVal = uint32(uintVal)
	case "uint64":
		retVal = uint64(uintVal)
	case "float32":
		retVal = float32(floatVal)
	case "float64":
		retVal = float64(floatVal)
	}

	return retVal, err
}
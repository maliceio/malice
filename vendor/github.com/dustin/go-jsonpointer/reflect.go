package jsonpointer

import (
	"reflect"
	"strconv"
	"strings"
)

// Reflect gets the value at the specified path from a struct.
func Reflect(o interface{}, path string) interface{} {
	if path == "" {
		return o
	}

	parts := parsePointer(path)
	var rv interface{} = o

OUTER:
	for _, p := range parts {
		val := reflect.ValueOf(rv)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		if val.Kind() == reflect.Struct {
			typ := val.Type()
			for i := 0; i < typ.NumField(); i++ {
				sf := typ.Field(i)
				tag := sf.Tag.Get("json")
				name := parseJSONTagName(tag)
				if (name != "" && name == p) || sf.Name == p {
					rv = val.Field(i).Interface()
					continue OUTER
				}
			}
			// Found no matching field.
			return nil
		} else if val.Kind() == reflect.Map {
			// our pointer always gives us a string key
			// here we try to convert it into the correct type
			mapKey, canConvert := makeMapKeyFromString(val.Type().Key(), p)
			if canConvert {
				field := val.MapIndex(mapKey)
				if field.IsValid() {
					rv = field.Interface()
				} else {
					return nil
				}
			} else {
				return nil
			}
		} else if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
			i, err := strconv.Atoi(p)
			if err == nil && i < val.Len() {
				rv = val.Index(i).Interface()
			} else {
				return nil
			}
		} else {
			return nil
		}
	}

	return rv
}

// ReflectListPointers lists all possible pointers from the given struct.
func ReflectListPointers(o interface{}) ([]string, error) {
	return reflectListPointersRecursive(o, ""), nil
}

func reflectListPointersRecursive(o interface{}, prefix string) []string {
	rv := []string{prefix + ""}

	val := reflect.ValueOf(o)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {

		typ := val.Type()
		for i := 0; i < typ.NumField(); i++ {
			child := val.Field(i).Interface()
			sf := typ.Field(i)
			tag := sf.Tag.Get("json")
			name := parseJSONTagName(tag)
			if name != "" {
				// use the tag name
				childReults := reflectListPointersRecursive(child, prefix+encodePointer([]string{name}))
				rv = append(rv, childReults...)
			} else {
				// use the original field name
				childResults := reflectListPointersRecursive(child, prefix+encodePointer([]string{sf.Name}))
				rv = append(rv, childResults...)
			}
		}

	} else if val.Kind() == reflect.Map {
		for _, k := range val.MapKeys() {
			child := val.MapIndex(k).Interface()
			mapKeyName := makeMapKeyName(k)
			childReults := reflectListPointersRecursive(child, prefix+encodePointer([]string{mapKeyName}))
			rv = append(rv, childReults...)
		}
	} else if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			child := val.Index(i).Interface()
			childResults := reflectListPointersRecursive(child, prefix+encodePointer([]string{strconv.Itoa(i)}))
			rv = append(rv, childResults...)
		}
	}
	return rv
}

// makeMapKeyName takes a map key value and creates a string representation
func makeMapKeyName(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		fv := v.Float()
		return strconv.FormatFloat(fv, 'f', -1, v.Type().Bits())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		iv := v.Int()
		return strconv.FormatInt(iv, 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		iv := v.Uint()
		return strconv.FormatUint(iv, 10)
	default:
		return v.String()
	}
}

// makeMapKeyFromString takes the key type for a map, and a string
// representing the key, it then tries to convert the string
// representation into a value of the correct type.
func makeMapKeyFromString(mapKeyType reflect.Type, pointer string) (reflect.Value, bool) {
	valp := reflect.New(mapKeyType)
	val := reflect.Indirect(valp)
	switch mapKeyType.Kind() {
	case reflect.String:
		return reflect.ValueOf(pointer), true
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		iv, err := strconv.ParseInt(pointer, 10, mapKeyType.Bits())
		if err == nil {
			val.SetInt(iv)
			return val, true
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		iv, err := strconv.ParseUint(pointer, 10, mapKeyType.Bits())
		if err == nil {
			val.SetUint(iv)
			return val, true
		}
	case reflect.Float32, reflect.Float64:
		fv, err := strconv.ParseFloat(pointer, mapKeyType.Bits())
		if err == nil {
			val.SetFloat(fv)
			return val, true
		}
	}

	return reflect.ValueOf(nil), false
}

// parseJSONTagName extracts the JSON field name from a struct tag
func parseJSONTagName(tag string) string {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag
}

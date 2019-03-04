package wall

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

func Parse(re *regexp.Regexp, input string, output interface{}) error {
	matches := re.FindStringSubmatch(input)
	names := normalizeNames(re.SubexpNames())

	if matches == nil {
		return errors.New("input does not match")
	}

	val := reflect.ValueOf(output)

	if val.Kind() != reflect.Ptr {
		return errors.New("output must be a pointer to a slice, map, or struct")
	}

	val = reflect.Indirect(val)

	switch val.Kind() {
	case reflect.Slice:
		return parseSlice(matches, val)
	case reflect.Map:
		return parseMap(matches, names, val)
	case reflect.Struct:
		return parseStruct(matches, names, val)
	default:
		return errors.New("output must be a pointer to a slice, map, or struct")
	}
}

func parseSlice(matches []string, val reflect.Value) error {
	if !(val.Kind() == reflect.Slice && val.Type().Elem().Kind() == reflect.String) {
		return errors.New("output slice must be []string")
	}

	if len(matches) > 1 {
		val.Set(reflect.ValueOf(matches[1:]))
	} else {
		val.Set(reflect.ValueOf([]string{}))
	}

	return nil
}

func parseMap(matches []string, names []string, val reflect.Value) error {
	if !(val.Kind() == reflect.Map &&
		val.Type().Elem().Kind() == reflect.String &&
		val.Type().Key().Kind() == reflect.String) {
		return errors.New("output map must be map[string]string")
	}

	outmap := make(map[string]string, len(matches))
	for i, v := range matches {
		if i == 0 {
			continue
		}

		outmap[names[i]] = v
	}

	val.Set(reflect.ValueOf(outmap))

	return nil
}

func parseStruct(matches []string, names []string, val reflect.Value) error {
	numFields := val.Type().NumField()

	for i := 0; i < numFields; i++ {
		sf := val.Type().Field(i)

		tag := sf.Tag.Get("wall")

		if tag != "" && strSliceContains(tag, names) {
			if sf.Type.Kind() != reflect.String {
				return fmt.Errorf(`field %s tagged with "wall" must be a string`, sf.Name)
			}

			f := val.Field(i)

			if !f.CanSet() {
				return fmt.Errorf(`field %s is not settable`, sf.Name)
			}

			f.Set(reflect.ValueOf(matchForName(tag, names, matches)))
		}
	}

	return nil
}

func normalizeNames(in []string) []string {
	out := make([]string, 0, len(in))

	for i, v := range in {
		if v == "" {
			out = append(out, strconv.FormatInt(int64(i), 10))
		} else {
			out = append(out, v)
		}
	}

	return out
}

func matchForName(name string, names []string, matches []string) string {
	for i, v := range names {
		if name == v {
			return matches[i]
		}
	}

	return ""
}

func strSliceContains(s string, sl []string) bool {
	for _, v := range sl {
		if s == v {
			return true
		}
	}

	return false
}

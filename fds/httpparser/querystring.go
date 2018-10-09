package httpparser

import (
	"fmt"
	"net/url"
	"reflect"
)

func QueryString(v interface{}) (url.Values, error) {
	values := make(url.Values)
	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return values, nil
		}
		val = val.Elem()
	}

	if v == nil {
		return values, nil
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("querystring: Values() expects struct input. Got %v", val.Kind())
	}

	err := reflectQueryString(values, val, "")
	return values, err
}

func reflectQueryString(values url.Values, val reflect.Value, scope string) error {
	var embedded []reflect.Value

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		if sf.PkgPath != "" && !sf.Anonymous { // unexported
			continue
		}

		sv := val.Field(i)
		tag := sf.Tag.Get(querystringTag)
		if tag == skipTag {
			continue
		}
		name, opts := parseTag(tag)
		if name == emptyTag {
			if sf.Anonymous && sv.Kind() == reflect.Struct {
				// save embedded struct for later processing
				embedded = append(embedded, sv)
				continue
			}

			name = sf.Name
		}

		if scope != "" {
			name = scope + "[" + name + "]"
		}

		if opts.Contains(omitemptyTag) && sv.Kind() == reflect.Bool {
			return fmt.Errorf("header: bool value can not be omitempty")
		}

		if opts.Contains(omitemptyTag) && isEmpty(sv) {
			continue
		}

		for sv.Kind() == reflect.Ptr {
			if sv.IsNil() {
				break
			}
			sv = sv.Elem()
		}

		if sv.Type() == timeType {
			values.Add(name, toString(sv, opts))
			continue
		}

		if sv.Kind() == reflect.Struct {
			reflectQueryString(values, sv, name)
			continue
		}

		values.Add(name, toString(sv, opts))
	}

	for _, f := range embedded {
		if err := reflectQueryString(values, f, scope); err != nil {
			return err
		}
	}

	return nil
}

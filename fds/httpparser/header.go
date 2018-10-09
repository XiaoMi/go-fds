package httpparser

import (
	"fmt"
	"net/http"
	"reflect"
)

var headerType = reflect.TypeOf(http.Header{})

func Header(v interface{}) (http.Header, error) {
	h := make(http.Header)
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return h, nil
		}
		val = val.Elem()
	}

	if v == nil {
		return h, nil
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("header: Header() expects struct input. Got %v", val.Kind())
	}

	err := reflectHeader(h, val)
	return h, err
}

func reflectHeader(header http.Header, val reflect.Value) error {
	var embedded []reflect.Value

	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		tf := typ.Field(i)
		if tf.PkgPath != "" && !tf.Anonymous { // unexported
			continue
		}

		vf := val.Field(i)
		tag := tf.Tag.Get(headerTag)

		if tag == skipTag {
			continue
		}

		name, opts := parseTag(tag)
		if name == emptyTag {
			if tf.Anonymous && vf.Kind() == reflect.Struct {
				embedded = append(embedded, vf)
				continue
			}

			name = tf.Name
		}

		if opts.Contains(omitemptyTag) && vf.Kind() == reflect.Bool {
			return fmt.Errorf("header: bool value can not be omitempty")
		}

		if opts.Contains(omitemptyTag) && isEmpty(vf) {
			continue
		}

		if vf.Kind() == reflect.Ptr {
			if vf.IsNil() {
				break
			}

			vf = vf.Elem()
		}

		if vf.Type() == timeType {
			header.Add(name, toString(vf, opts))
			continue
		}

		if vf.Type() == headerType {
			h := vf.Interface().(http.Header)
			for k, vs := range h {
				for _, v := range vs {
					header.Add(k, v)
				}
			}
			continue
		}

		if vf.Kind() == reflect.Struct {
			reflectHeader(header, vf)
			continue
		}

		header.Add(name, toString(vf, opts))
	}

	for _, v := range embedded {
		if err := reflectHeader(header, v); err != nil {
			return err
		}
	}
	return nil
}

package httpparser

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

var timeType = reflect.TypeOf(time.Time{})

const (
	skipTag        = "-"
	emptyTag       = ""
	omitemptyTag   = "omitempty"
	headerTag      = "header"
	querystringTag = "param"
)

type tags []string

func parseTag(tag string) (string, tags) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}

func (o tags) Contains(option string) bool {
	for _, s := range o {
		if s == option {
			return true
		}
	}
	return false
}

func isEmpty(v reflect.Value) bool {
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func toString(v reflect.Value, opts tags) string {

	if v.Type() == timeType {
		t := v.Interface().(time.Time)
		return t.Format(time.RFC822)
	}

	return fmt.Sprint(v.Interface())
}

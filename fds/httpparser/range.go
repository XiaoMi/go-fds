package httpparser

import (
	"fmt"
	"strconv"
	"strings"
)

type HTTPRange struct {
	Start int64
	End   int64
}

func Range(r string) ([]HTTPRange, error) {
	var ranges []HTTPRange

	if r == "" {
		return ranges, nil
	}

	index := strings.Index(r, "=")

	if index == -1 {
		return ranges, fmt.Errorf("fds: error range format")
	}

	if r[0:index] != "bytes" {
		return ranges, fmt.Errorf("fds: only support bytes range")
	}

	rangeBody := r[index+1:]
	splitRange := strings.Split(rangeBody, ",")

	for _, item := range splitRange {
		splitItem := strings.Split(item, "-")
		if len(splitItem) != 2 {
			return ranges, fmt.Errorf("fds: error range format")
		}

		var start int64
		var err error
		if splitItem[0] != "" {
			start, err = strconv.ParseInt(splitItem[0], 10, 0)
			if err != nil {
				return ranges, fmt.Errorf("fds: error range format")
			}
		}

		var end int64
		if splitItem[1] == "" {
			return ranges, fmt.Errorf("fds: error range format")
		}
		end, err = strconv.ParseInt(splitItem[1], 10, 0)
		if err != nil {
			return ranges, fmt.Errorf("fds: error range format")
		}

		ranges = append(ranges, HTTPRange{start, end})
	}

	return ranges, nil
}

package fds

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

var subResourceMap = map[string]string{
	"acl":                "",
	"quota":              "",
	"uploads":            "",
	"partNumber":         "",
	"uploadId":           "",
	"storageAccessToken": "",
	"metadata":           "",
}

func signature(sk string, method HTTPMethod, url string, header http.Header) (string, error) {
	var buf bytes.Buffer
	contentMd5 := header.Get(HTTPHeaderContentMD5)
	contentType := header.Get(HTTPHeaderContentType)
	date := expires(url)
	if len(date) == 0 {
		date = header.Get(HTTPHeaderDate)
	}
	buf.WriteString(string(method))
	buf.WriteString("\n")
	buf.WriteString(contentMd5)
	buf.WriteString("\n")
	buf.WriteString(contentType)
	buf.WriteString("\n")
	buf.WriteString(date)
	buf.WriteString("\n")

	ch, err := miHeader(header)
	if err != nil {
		return "", newServerError(err.Error(), -1)
	}
	buf.Write(ch)
	cr, err := resource(url)
	if err != nil {
		return "", newServerError(err.Error(), -1)
	}
	buf.Write(cr)
	h := hmac.New(sha1.New, []byte(sk))
	_, err = h.Write(buf.Bytes())
	if err != nil {
		return "", newServerError(err.Error(), -1)
	}
	b := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return b, nil
}

func resource(uri string) ([]byte, error) {
	uriParsed, err := url.Parse(uri)
	if err != nil {
		return nil, newServerError(err.Error(), -1)
	}
	var path bytes.Buffer
	path.Write([]byte(uriParsed.Path))

	param := uriParsed.Query()
	var filteredKey []string
	filteredMap := map[string]string{}
	for k, v := range param {
		_, ok := subResourceMap[k]
		if !ok {
			continue
		}
		filteredKey = append(filteredKey, k)
		if len(v) > 0 {
			filteredMap[k] = v[0]
		} else {
			filteredMap[k] = ""
		}
	}

	if len(filteredKey) == 0 {
		return path.Bytes(), nil
	}

	sort.Strings(filteredKey)

	for i, k := range filteredKey {
		if i == 0 {
			path.WriteString("?")
		} else {
			path.WriteString("&")
		}
		path.WriteString(k)
		if len(filteredMap[k]) > 0 {
			path.WriteString("=")
			path.WriteString(filteredMap[k])
		}
	}

	return path.Bytes(), nil
}

func expires(urlStr string) string {
	urlParsed, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	queryParams := urlParsed.Query()
	d, ok := queryParams["Expires"]

	if !ok || len(d) == 0 {
		return ""
	}
	return d[0]
}

func miHeader(h http.Header) ([]byte, error) {
	if len(h) == 0 {
		return nil, nil
	}

	var keyList []string
	filteredMap := map[string]string{}
	for k, v := range h {
		key := strings.ToLower(k)
		if !strings.HasPrefix(key, XiaomiPrefix) {
			continue
		}

		filteredMap[key] = strings.Join(v, ",")
		keyList = append(keyList, key)
	}
	sort.Strings(keyList)

	var r bytes.Buffer
	for _, k := range keyList {
		r.WriteString(k)
		r.WriteString(":")
		r.WriteString(filteredMap[k])
		r.WriteString("\n")
	}

	return r.Bytes(), nil
}

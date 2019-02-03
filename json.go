package superjson

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/tidwall/gjson"
)

type (
	// KeyConvert key convert function
	KeyConvert func(string) string
	// KeyFilter key filter function
	KeyFilter func(string) (omit bool, newKey string)
)

const (
	underscoreRune = '_'
	spaceRune      = ' '
)

func cutWords(str string) []string {
	result := make([]string, 0, 10)

	currentCutIndex := 0
	for index, ch := range str {
		if ch == underscoreRune ||
			ch == spaceRune {
			if currentCutIndex != index {
				result = append(result, str[currentCutIndex:index])
			}
			currentCutIndex = index + 1
		}
		if index == 0 || !unicode.IsUpper(ch) {
			continue
		}
		if currentCutIndex != index {
			result = append(result, str[currentCutIndex:index])
			currentCutIndex = index
		}
	}
	if currentCutIndex < len(str)-1 {
		result = append(result, str[currentCutIndex:])
	}
	return result
}

// Filter json filter
func Filter(buf []byte, filter KeyFilter) []byte {
	result := gjson.ParseBytes(buf)
	arr := make([][]byte, 0, 10)
	result.ForEach(func(key, value gjson.Result) bool {
		k := key.String()
		omit, newKey := filter(k)
		if omit {
			return true
		}
		if newKey != "" {
			k = newKey
		}
		if value.Type == gjson.Null {
			return true
		}
		arr = append(arr, []byte(`"`+k+`":`+value.Raw))
		return true
	})
	data := bytes.Join(arr, []byte(","))
	data = bytes.Join([][]byte{
		[]byte("{"),
		data,
		[]byte("}"),
	}, nil)
	return data
}

// Pick pick fields from json
func Pick(buf []byte, fields []string) []byte {
	pickKeys := make(map[string]bool)
	for _, key := range fields {
		pickKeys[key] = true
	}
	return Filter(buf, func(k string) (bool, string) {
		return !pickKeys[k], k
	})
}

// Omit omit fields from json
func Omit(buf []byte, fields []string) []byte {
	omitKeys := make(map[string]bool)
	for _, key := range fields {
		omitKeys[key] = true
	}
	return Filter(buf, func(k string) (bool, string) {
		return omitKeys[k], k
	})
}

// camelCase convert string to camel case
// https://github.com/lodash/lodash/blob/master/camelCase.js
func camelCase(str string) string {
	result := cutWords(str)
	for index, item := range result {
		if index == 0 {
			// 第一个单词首字母小写
			result[index] = strings.ToLower(item)
		} else {
			// 后续的单词首字母大写
			result[index] = strings.ToUpper(item[0:1]) + strings.ToLower(item[1:])
		}
	}
	return strings.Join(result, "")
}

// snakeCase convert string to snake case
// https://github.com/lodash/lodash/blob/master/snakeCase.js
func snakeCase(str string) string {
	result := cutWords(str)
	for index, item := range result {
		word := strings.ToLower(item)
		if index == 0 {
			result[index] = word
		} else {
			result[index] = "_" + word
		}
	}
	return strings.Join(result, "")
}

func convertJSON(t gjson.Result, fn KeyConvert, builder *strings.Builder) {
	isArray := t.IsArray()

	if isArray {
		builder.WriteString("[")
	} else {
		builder.WriteString("{")
	}

	index := 0
	t.ForEach(func(key, value gjson.Result) bool {
		if index != 0 {
			builder.WriteString(",")
		}
		k := fn(key.String())
		// 如果有key，则设置key
		if k != "" {
			builder.WriteString(`"`)
			builder.WriteString(k)
			builder.WriteString(`":`)
		}

		// 如果是array或者object，则递归
		if value.IsArray() || value.IsObject() {
			convertJSON(value, fn, builder)
		} else {
			builder.WriteString(value.Raw)
		}

		index++
		return true
	})
	if isArray {
		builder.WriteString("]")
	} else {
		builder.WriteString("}")
	}
}

func createBuilder() *strings.Builder {
	builder := new(strings.Builder)
	builder.Grow(4096)
	return builder
}

// CamelCase convert json to camel case
func CamelCase(buf []byte) []byte {
	builder := createBuilder()
	result := gjson.ParseBytes(buf)
	convertJSON(result, camelCase, builder)
	return []byte(builder.String())
}

// SnakeCase convert json to snake case
func SnakeCase(buf []byte) []byte {
	builder := createBuilder()
	result := gjson.ParseBytes(buf)
	convertJSON(result, snakeCase, builder)
	return []byte(builder.String())
}

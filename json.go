package superjson

import (
	"bytes"
	"strings"

	"github.com/tidwall/gjson"
)

type (
	// Convert key convert function
	Convert func(string) string
)

func cutWords(str string) []string {
	result := make([]string, 0, 10)
	indexList := (cutWordsReg.FindAllStringIndex(str, -1))
	currentIndex := 0
	prefix := ""
	for _, items := range indexList {
		start := items[0]
		if start == 0 {
			continue
		}
		end := items[1]

		result = append(result, prefix+str[currentIndex:start])
		tmp := str[start:end]
		prefix = omitReg.ReplaceAllString(tmp, "")

		currentIndex = end
	}
	if currentIndex != len(str) {
		result = append(result, prefix+str[currentIndex:])
	}
	return result
}

// Pick pick fields from json
func Pick(buf []byte, fields []string) []byte {
	result := gjson.GetManyBytes(buf, fields...)
	max := len(result)
	arr := make([][]byte, max)
	currentIndex := 0
	for index, item := range result {
		raw := item.Raw
		// nil的数据忽略
		if item.Type == gjson.Null {
			continue
		}
		arr[currentIndex] = []byte(`"` + fields[index] + `":` + raw)
		currentIndex++
	}
	// 如果部分数据跳过，则裁剪数组
	if currentIndex != max {
		arr = arr[0:currentIndex]
	}

	data := bytes.Join(arr, []byte(","))
	data = bytes.Join([][]byte{
		[]byte("{"),
		data,
		[]byte("}"),
	}, nil)
	return data
}

// Omit omit fields from json
func Omit(buf []byte, fields []string) []byte {
	newFields := make([]string, 0, 10)
	omitKeys := make(map[string]bool)
	for _, key := range fields {
		omitKeys[key] = true
	}
	result := gjson.ParseBytes(buf)
	result.ForEach(func(key, value gjson.Result) bool {
		k := key.String()
		if !omitKeys[k] {
			newFields = append(newFields, k)
		}
		return true
	})
	return Pick(buf, newFields)
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

func convertJSON(t gjson.Result, fn Convert) string {
	json := make([]string, 0)
	isArray := t.IsArray()
	iterator := func(key, value gjson.Result) bool {
		k := fn(key.String())
		var valueStr string
		if value.IsObject() || value.IsArray() {
			valueStr = convertJSON(value, fn)
		} else {
			valueStr = value.Raw
		}
		v := ""
		// 如果数组，则没有key
		if isArray {
			v = valueStr
		} else {
			v = `"` + k + `":` + valueStr
		}
		json = append(json, v)
		return true
	}
	t.ForEach(iterator)
	joinJSON := strings.Join(json, ",")
	if isArray {
		return "[" + joinJSON + "]"
	}
	return "{" + joinJSON + "}"
}

// CamelCase convert json to camel case
func CamelCase(buf []byte) []byte {
	result := gjson.ParseBytes(buf)
	json := convertJSON(result, camelCase)
	return []byte(json)
}

// SnakeCase convert json to snake case
func SnakeCase(buf []byte) []byte {
	result := gjson.ParseBytes(buf)
	json := convertJSON(result, snakeCase)
	return []byte(json)
}

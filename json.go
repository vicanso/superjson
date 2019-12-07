package superjson

import (
	"bytes"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/tidwall/gjson"
)

type (
	// KeyConvert key convert function
	KeyConvert func(string) string
	// KeyFilter key filter function
	KeyFilter func(key, value string) (omit bool, newKey string)
	// ValueMask value mask function
	ValueMask func(key, value string) (newValue string)
)

func doJSON(buf []byte, filter KeyFilter, mask ValueMask) []byte {
	result := gjson.ParseBytes(buf)
	arr := make([][]byte, 0, 10)
	result.ForEach(func(key, value gjson.Result) bool {
		// 对于数组内的元素
		if key.Type == gjson.Null && value.Type == gjson.JSON {
			arr = append(arr, doJSON([]byte(value.Raw), filter, mask))
			return true
		}
		k := key.String()
		v := value.Raw
		if filter != nil {
			omit, newKey := filter(k, v)
			if omit {
				return true
			}
			if newKey != "" {
				k = newKey
			}
		}
		if value.Type == gjson.Null {
			return true
		}
		if mask != nil {
			newValue := mask(k, v)
			if newValue != "" {
				v = newValue
			}
		}

		arr = append(arr, []byte(`"`+k+`":`+v))
		return true
	})
	data := bytes.Join(arr, []byte(","))
	if result.IsArray() {
		data = bytes.Join([][]byte{
			[]byte("["),
			data,
			[]byte("]"),
		}, nil)
	} else {
		data = bytes.Join([][]byte{
			[]byte("{"),
			data,
			[]byte("}"),
		}, nil)
	}

	return data
}

// Filter json filter
func Filter(buf []byte, filter KeyFilter) []byte {
	return doJSON(buf, filter, nil)
}

// Mask json mask
func Mask(buf []byte, mask ValueMask) []byte {
	return doJSON(buf, nil, mask)
}

// Pick pick fields from json
func Pick(buf []byte, fields []string) []byte {
	pickKeys := make(map[string]bool)
	for _, key := range fields {
		pickKeys[key] = true
	}
	return Filter(buf, func(k, _ string) (bool, string) {
		return !pickKeys[k], k
	})
}

// Omit omit fields from json
func Omit(buf []byte, fields []string) []byte {
	omitKeys := make(map[string]bool)
	for _, key := range fields {
		omitKeys[key] = true
	}
	return Filter(buf, func(k, _ string) (bool, string) {
		return omitKeys[k], k
	})
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
	convertJSON(result, strcase.ToLowerCamel, builder)
	return []byte(builder.String())
}

// SnakeCase convert json to snake case
func SnakeCase(buf []byte) []byte {
	builder := createBuilder()
	result := gjson.ParseBytes(buf)
	convertJSON(result, strcase.ToSnake, builder)
	return []byte(builder.String())
}

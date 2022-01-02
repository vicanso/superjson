package superjson

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPickOmit(t *testing.T) {
	b64 := base64.StdEncoding.EncodeToString(make([]byte, 1024))
	m := map[string]interface{}{
		"_x": b64,
		"_y": b64,
		"_z": b64,
		"i":  1,
		"f":  1.12,
		"s":  "\"abc",
		"b":  false,
		"arr": []interface{}{
			1,
			"2",
			true,
		},
		"m": map[string]interface{}{
			"a": 1,
			"b": "2",
			"c": false,
		},
		"null": nil,
		"中文":   "名称",
	}
	buf, _ := json.Marshal(m)
	t.Run("pick", func(t *testing.T) {
		assert := assert.New(t)
		pickData := Pick(buf, strings.Split("i,f,s,b,arr,m,null,中文", ","))
		assert.Equal(`{"arr":[1,"2",true],"b":false,"f":1.12,"i":1,"m":{"a":1,"b":"2","c":false},"s":"\"abc","中文":"名称"}`, string(pickData))
	})

	t.Run("pick array", func(t *testing.T) {
		assert := assert.New(t)
		data := []map[string]string{
			{
				"name":  "foo",
				"type":  "1",
				"price": "1.11",
			},
			{
				"name":  "bar",
				"type":  "2",
				"price": "2.22",
			},
		}
		buf, _ := json.Marshal(data)

		pcikData := Pick(buf, []string{
			"name",
			"price",
		})
		assert.Equal(`[{"name":"foo","price":"1.11"},{"name":"bar","price":"2.22"}]`, string(pcikData))
	})

	t.Run("omit", func(t *testing.T) {
		assert := assert.New(t)
		omitData := Omit(buf, strings.Split("_x,_y,_z", ","))
		assert.Equal(`{"arr":[1,"2",true],"b":false,"f":1.12,"i":1,"m":{"a":1,"b":"2","c":false},"s":"\"abc","中文":"名称"}`, string(omitData))
	})
}

func TestMask(t *testing.T) {
	assert := assert.New(t)
	m := map[string]interface{}{
		"i": 1,
		"f": 1.12,
		"s": "\"abc",
		"b": false,
		"arr": []interface{}{
			1,
			"2",
			true,
		},
		"m": map[string]interface{}{
			"a": 1,
			"b": "2",
			"c": false,
		},
		"null": nil,
		"中文":   "名称",
	}
	buf, _ := json.Marshal(m)
	data := Mask(buf, func(key, _ string) string {
		if key != "i" {
			return `"***"`
		}
		return ""
	})
	assert.Equal(`{"arr":"***","b":"***","f":"***","i":1,"m":"***","s":"***","中文":"***"}`, string(data))
}

func BenchmarkPick(b *testing.B) {
	b.ReportAllocs()
	m := map[string]interface{}{
		"i": 1,
		"f": 1.12,
		"s": "\"abc",
		"b": false,
		"arr": []interface{}{
			1,
			"2",
			true,
		},
		"m": map[string]interface{}{
			"a": 1,
			"b": "2",
			"c": false,
		},
		"null": nil,
		"中文":   "名称",
	}
	buf, _ := json.Marshal(m)
	fields := strings.Split("i,f,s,b,arr,m,null,中文", ",")
	for i := 0; i < b.N; i++ {
		Pick(buf, fields)
	}
}

func BenchmarkOmit(b *testing.B) {
	b.ReportAllocs()
	m := map[string]interface{}{
		"i": 1,
		"f": 1.12,
		"s": "\"abc",
		"b": false,
		"arr": []interface{}{
			1,
			"2",
			true,
		},
		"m": map[string]interface{}{
			"a": 1,
			"b": "2",
			"c": false,
		},
		"null": nil,
		"中文":   "名称",
	}
	buf, _ := json.Marshal(m)
	fields := strings.Split("i,f,s,b", ",")
	for i := 0; i < b.N; i++ {
		Omit(buf, fields)
	}
}

func BenchmarkCamelCase(b *testing.B) {
	b.ReportAllocs()
	json := []byte(`{
		"book_name": "测试",
		"book_price": 12,
		"book_on_sale": true,
		"book_author": {
			"author_name": "tree.xie",
			"author_age": 0,
			"author_salary": 10.1,
		},
		"book_category": ["vip", "hot-sale"],
		"book_infos": [
			{
				"word_count": 100
			}
		]
	}`)
	for i := 0; i < b.N; i++ {
		CamelCase(json)
	}
}

func TestConvertJSON(t *testing.T) {
	assert := assert.New(t)
	json := []byte(`{
		"book_name": "测试",
		"book_price": 12,
		"book_on_sale": true,
		"book_author": {
			"author_name": "tree.xie",
			"author_age": 0,
			"author_salary": 10.1,
		},
		"book_category": ["vip", "hot-sale"],
		"book_infos": [
			{
				"word_count": 100
			}
		]
	}`)
	camelCaseJSON := CamelCase(json)
	assert.Equal(`{"bookName":"测试","bookPrice":12,"bookOnSale":true,"bookAuthor":{"authorName":"tree.xie","authorAge":0,"authorSalary":10.1},"bookCategory":["vip","hot-sale"],"bookInfos":[{"wordCount":100}]}`, string(camelCaseJSON))
	assert.Equal(`{"book_name":"测试","book_price":12,"book_on_sale":true,"book_author":{"author_name":"tree.xie","author_age":0,"author_salary":10.1},"book_category":["vip","hot-sale"],"book_infos":[{"word_count":100}]}`, string(SnakeCase(camelCaseJSON)))
}

# super json

[![Build Status](https://img.shields.io/travis/vicanso/superjson.svg?label=linux+build)](https://travis-ci.org/vicanso/superjson)


JSON picker and converter.

# API

## Pick

```go
buf := []byte(`{
  "name": "tree.xie",
  "address": "GZ",
  "no": 123
}`)
data := superjson.Pick(buf, []string{
  "name",
  "no",
})
// {"name":"tree.xie","no":123}
fmt.Println(string(data))
```

## Omit

```go
buf := []byte(`{
  "name": "tree.xie",
  "address": "GZ",
  "no": 123
}`)
data := superjson.Omit(buf, []string{
  "address",
})
// {"name":"tree.xie","no":123}
fmt.Println(string(data))
```

## Filter

```go
buf := []byte(`{
  "name": "tree.xie",
  "address": "GZ",
  "no": 123
}`)
data := superjson.Filter(buf, func(key string) (omit bool, newKey string) {
  // omit the no
  if key == "no" {
    return true, ""
  }
  // convert the address to addr
  if key == "address" {
    return false, "addr"
  }
  // key original
  return false, key
})
// {"name":"tree.xie","addr":"GZ"}
fmt.Println(string(data))
```

## CamelCase

```go
buf := []byte(`{
	"book_author_name": "tree.xie",
	"book_no": 123
}`)
data := superjson.CamelCase(buf)
// {"bookAuthorName":"tree.xie","bookNo":123}
fmt.Println(string(data)
```

## SnakeCase

```go
buf := []byte(`{
  "bookAuthorName": "tree.xie",
  "bookNo": 123
}`)
data := superjson.SnakeCase(buf)
// {"book_author_name":"tree.xie","book_no":123}
fmt.Println(string(data))
```

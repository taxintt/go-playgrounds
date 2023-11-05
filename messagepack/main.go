package main

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
)

type Item struct {
	Foo string
	Bar string
}

func main() {
	item := &Item{Foo: "foo", Bar: "bar"}

	serialize(item)
	// fmt.Println(msgpack)
}

func serialize(item *Item) {
	b, err := msgpack.Marshal(item)
	if err != nil {
		panic(err)
	}

	fmt.Println(concatBytes(b))
	// err = msgpack.Unmarshal(b, &item)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(item.Foo)
	// Output: bar
}

func concatBytes(array []byte) string {
	var result string
	for _, v := range array {
		result += string(v)
	}
	return result
}

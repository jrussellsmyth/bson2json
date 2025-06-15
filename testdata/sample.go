package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"os"
)

func main() {
	docs := []interface{}{
		bson.M{"name": "Alice", "age": 30},
		bson.M{"name": "Bob", "age": 25},
	}
	f, err := os.Create("sample.bson")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for i := 0; i < 1000; i++ {
		for _, doc := range docs {
			b, err := bson.Marshal(doc)
			if err != nil {
				panic(err)
			}
			f.Write(b)
		}
	}
}

package main

import (
	"os"
	"go.mongodb.org/mongo-driver/bson"
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
	for _, doc := range docs {
		b, err := bson.Marshal(doc)
		if err != nil {
			panic(err)
		}
		f.Write(b)
	}
}

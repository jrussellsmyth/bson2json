package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func BenchmarkReadAllAndProcess(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		f, err := os.Open("../../testdata/sample.bson")
		if err != nil {
			b.Fatal(err)
		}
		data, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			b.Fatal(err)
		}
		buf := bytes.NewBuffer(data)
		var docs []map[string]interface{}
		for buf.Len() > 0 {
			if buf.Len() < 4 {
				break
			}
			lengthBytes := buf.Next(4)
			length := int(uint32(lengthBytes[0]) | uint32(lengthBytes[1])<<8 | uint32(lengthBytes[2])<<16 | uint32(lengthBytes[3])<<24)
			if length < 5 || length > buf.Len()+4 {
				break
			}
			docBytes := append(lengthBytes, buf.Next(length-4)...)
			var doc map[string]interface{}
			if err := bson.Unmarshal(docBytes, &doc); err != nil {
				b.Fatal(err)
			}
			docs = append(docs, doc)
		}
	}
}

func BenchmarkStreamProcess(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		f, err := os.Open("../../testdata/sample.bson")
		if err != nil {
			b.Fatal(err)
		}
		defer f.Close()
		for {
			lenBuf := make([]byte, 4)
			_, err := io.ReadFull(f, lenBuf)
			if err == io.EOF {
				break
			}
			if err != nil {
				if err == io.ErrUnexpectedEOF {
					break
				}
				b.Fatal(err)
			}
			length := int(uint32(lenBuf[0]) | uint32(lenBuf[1])<<8 | uint32(lenBuf[2])<<16 | uint32(lenBuf[3])<<24)
			if length < 5 {
				break
			}
			docBuf := make([]byte, length-4)
			_, err = io.ReadFull(f, docBuf)
			if err != nil {
				b.Fatal(err)
			}
			docBytes := append(lenBuf, docBuf...)
			var doc map[string]interface{}
			if err := bson.Unmarshal(docBytes, &doc); err != nil {
				b.Fatal(err)
			}
		}
		f.Close()
	}
}

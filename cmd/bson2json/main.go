package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	inputPath := flag.String("input", "", "Path to the BSON file (optionally .gz)")
	flag.Parse()

	if *inputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: --input flag is required")
		os.Exit(1)
	}

	file, err := os.Open(*inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var reader io.Reader = file
	if strings.HasSuffix(*inputPath, ".gz") {
		gz, err := gzip.NewReader(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create gzip reader: %v\n", err)
			os.Exit(1)
		}
		defer gz.Close()
		reader = gz
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file: %v\n", err)
		os.Exit(1)
	}

	var docs []map[string]interface{}
	buf := bytes.NewBuffer(data)
	for buf.Len() > 0 {
		var doc map[string]interface{}
		// Read the length of the next BSON document (first 4 bytes)
		if buf.Len() < 4 {
			break
		}
		lengthBytes := buf.Next(4)
		length := int(uint32(lengthBytes[0]) | uint32(lengthBytes[1])<<8 | uint32(lengthBytes[2])<<16 | uint32(lengthBytes[3])<<24)
		if length < 5 || length > buf.Len()+4 {
			break
		}
		// Reconstruct the full document
		docBytes := append(lengthBytes, buf.Next(length-4)...)
		if err := bson.Unmarshal(docBytes, &doc); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to decode BSON document: %v\n", err)
			os.Exit(1)
		}
		docs = append(docs, doc)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(docs); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to encode JSON: %v\n", err)
		os.Exit(1)
	}
}

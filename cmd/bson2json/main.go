package main

import (
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

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	fmt.Print("[")
	first := true
	for {
		// Read the length of the next BSON document (first 4 bytes)
		lenBuf := make([]byte, 4)
		_, err := io.ReadFull(reader, lenBuf)
		if err == io.EOF {
			break
		}
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				break
			}
			fmt.Fprintf(os.Stderr, "Failed to read BSON document length: %v\n", err)
			os.Exit(1)
		}
		length := int(uint32(lenBuf[0]) | uint32(lenBuf[1])<<8 | uint32(lenBuf[2])<<16 | uint32(lenBuf[3])<<24)
		if length < 5 {
			break
		}
		docBuf := make([]byte, length-4)
		_, err = io.ReadFull(reader, docBuf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read BSON document: %v\n", err)
			os.Exit(1)
		}
		docBytes := append(lenBuf, docBuf...)
		var doc map[string]interface{}
		if err := bson.Unmarshal(docBytes, &doc); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to decode BSON document: %v\n", err)
			os.Exit(1)
		}
		if !first {
			fmt.Print(",\n")
		}
		first = false
		jsonBytes, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to encode JSON: %v\n", err)
			os.Exit(1)
		}
		os.Stdout.Write(jsonBytes)
	}
	fmt.Println("]")
}

package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func isGzipHeader(header []byte) bool {
	return len(header) >= 2 && header[0] == 0x1f && header[1] == 0x8b
}

func main() {
	var input io.Reader
	var closer func()

	if len(os.Args) > 1 && os.Args[1] != "-" {
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
			os.Exit(1)
		}
		closer = func() { file.Close() }
		input = file
	} else {
		input = os.Stdin
	}

	bufr := bufio.NewReader(input)
	header, err := bufr.Peek(2)
	if err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "Failed to peek input: %v\n", err)
		os.Exit(1)
	}

	var reader io.Reader
	if isGzipHeader(header) {
		gz, err := gzip.NewReader(bufr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create gzip reader: %v\n", err)
			if closer != nil {
				closer()
			}
			os.Exit(1)
		}
		defer gz.Close()
		reader = bufio.NewReader(gz)
	} else {
		reader = bufr
	}
	if closer != nil {
		defer closer()
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")	

	fmt.Print("[")
	first := true
	for {
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

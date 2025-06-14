# Bson2Json

A simple command-line tool to convert BSON files (optionally gzip-compressed) to JSON.

## Usage

```
go run cmd/bson2json/main.go --input <input.bson|input.bson.gz>
```

## Build for Mac (Apple Silicon)

```
GOOS=darwin GOARCH=arm64 go build -o bson2json-macos-arm64 ./cmd/bson2json
```

## Features
- Supports BSON files and gzip-compressed BSON files
- Outputs JSON to stdout
- Handles single or multiple BSON documents (if implemented)

## License
MIT

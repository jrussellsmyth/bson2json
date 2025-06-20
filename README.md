# Bson2Json

A simple command-line tool to convert BSON files (optionally gzip-compressed) to JSON.

## Usage

```
# Convert a file
bson2json <input.bson|input.bson.gz>

# Or pipe data from stdin
cat input.bson | bson2json
cat input.bson.gz | bson2json
```

## Build Instructions

This project uses [go-task](https://taskfile.dev) for builds. Install it with:

```
curl -sL https://taskfile.dev/install.sh | sh
```

### Build for current platform
```
task build
```

### Build for all major platforms (outputs in `bin/`):
```
task build-all
```

### Clean build artifacts
```
task clean
```

## Features
- Supports BSON files and gzip-compressed BSON files
- Auto-detects gzip by file content, not extension
- Accepts input from a file or stdin
- Outputs JSON to stdout
- Multi-platform builds (Linux, macOS, Windows, Intel & ARM)

## License
MIT

# tsv

A fast, reflection-based TSV (Tab-Separated Values) encoder for Go.

## Features

- Encodes Go primitives, structs, maps, slices, arrays, pointers, and interfaces to TSV format
- Complex types (structs, maps, slices, arrays, pointers) inside slices/arrays are serialized as JSON to preserve TSV line integrity
- Nil pointers in slices produce empty strings to maintain column alignment
- Configurable time formatting (Unix epoch by default)
- CRLF line ending support for Windows compatibility
- String escaping for special characters (`\t`, `\n`, `\r`, `\\`)
- Zero external dependencies

## Installation

```bash
go get github.com/alex-cos/tsv
```

## Quick Start

```go
package main

import (
  "fmt"
  "github.com/alex-cos/tsv"
)

func main() {
  enc := tsv.NewTSVEncoder()

  // Primitives
  b, _ := enc.Encode("hello")
  fmt.Println(string(b)) // hello

  // Struct
  type User struct {
    Name string
    Age  int
  }
  b, _ = enc.Encode(User{Name: "Alice", Age: 30})
  fmt.Println(string(b)) // Alice  30
}
```

## Options

Options are passed as variadic arguments to `NewTSVEncoder`:

### WithTimeFormat

Sets a custom time format. By default, `time.Time` values are encoded as Unix epoch timestamps.

```go
enc := tsv.NewTSVEncoder(tsv.WithTimeFormat("2006/01/02 15:04:05"))
```

### WithCRLF

Uses `\r\n` line endings instead of `\n` for Windows compatibility.

```go
enc := tsv.NewTSVEncoder(tsv.WithCRLF())
```

### Combining options

```go
enc := tsv.NewTSVEncoder(
  tsv.WithTimeFormat(time.RFC3339),
  tsv.WithCRLF(),
)
```

## Supported Types

| Type | Output |
| ------ | -------- |
| `bool` | `true` / `false` |
| `int`, `int8`–`int64` | Decimal number |
| `uint`, `uint8`–`uint64`, `uintptr` | Decimal number |
| `float32`, `float64` | Decimal number |
| `string` | String (with escaping) |
| `time.Time` | Unix epoch or custom format |
| `struct` | Tab-separated fields |
| `map` | Key-value pairs, one per line |
| `slice` / `array` (primitives) | Tab-separated values |
| `slice` / `array` (complex) | JSON-serialized, tab-separated |
| `pointer` | Dereferenced value, or empty if nil |
| `interface{}` | Underlying type, or empty if nil |

## String Escaping

Special characters in strings are escaped to preserve TSV structure:

| Character | Escape |
| ----------- | -------- |
| `\` | `\\` |
| Tab | `\t` |
| Newline | `\n` |
| Carriage return | `\r` |

## Complex Types in Slices

When a slice or array contains structs, maps, or other slices, each element is JSON-serialized to keep the TSV on a single line:

```go
type Row struct {
  Name string `json:"name"`
  ID   int    `json:"id"`
}

data := []Row{
  {Name: "Alice", ID: 1},
  {Name: "Bob", ID: 2},
}

enc := tsv.NewTSVEncoder()
b, _ := enc.Encode(data)
// {"name":"Alice","id":1}  {"name":"Bob","id":2}
```

## Nil Pointers

Nil pointers inside slices produce an empty string to maintain column alignment:

```go
s1 := "hello"
s3 := "world"
data := []*string{&s1, nil, &s3}

b, _ := tsv.NewTSVEncoder().Encode(data)
// "hello"    "world"
```

## Maps

Maps are encoded as key-value pairs, one per line:

```go
data := map[string]int{
  "a": 1,
  "b": 2,
}

b, _ := tsv.NewTSVEncoder().Encode(data)
// a  1
// b  2
```

Note: Map iteration order is non-deterministic in Go.

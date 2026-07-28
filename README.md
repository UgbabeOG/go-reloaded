# go-reloaded

A Go command-line tool that transforms plain text using special inline markers. It applies number conversions, case modifiers, article fixes, and punctuation normalization.

## Features

- **Number conversions**
  - `(hex)` — convert the preceding hexadecimal number to decimal
  - `(bin)` — convert the preceding binary number to decimal

- **Case modifiers** (apply to the preceding word, or preceding N words if a count is given)
  - `(up)` — uppercase
  - `(low)` — lowercase
  - `(cap)` — capitalize

- **Article correction**
  - `a` → `an` before vowels or words starting with `h`

- **Punctuation & quotes**
  - Attach punctuation marks to the preceding word
  - Remove spaces inside quote pairs

## Installation

```bash
go build -o go-reloaded .
```

## Usage

```bash
./go-reloaded input.txt output.txt
```

The program reads `input.txt`, applies all transformations, and writes the result to `output.txt`.

## Examples

Input:
```
1 (hex) cats (up, 2) are better than a (cap) unicorn (hex, 2)
```

Output:
```
1 cats ARE BETTER THAN An unicorn 256, 100
```

Input:
```
'a man, a plan, a (low, 3) 'canal' , Panama 'wow' ...
```

Output:
```
a man, a plan, a canal, Panama wow...
```

Input:
```
(up) 1010 (bin) and (cap) 1F (hex)
```

Output:
```
1010 10 And 31
```

## Project Structure

- `main.go` — CLI entry point
- `parser.go` — Tokenizer and modifier parser
- `transformations.go` — Transformation implementations
- `*_test.go` — Unit tests

## Running Tests

```bash
go test ./...
```

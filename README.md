# xtracego

xtracego is a command-line tool to run Go source code injecting xtrace.

## Example

### FizzBuzz

Source code:

```go
// examples/fizzbuzz/main.go
package main

import "fmt"

const N = 20

func main() {
	for i := 1; i <= N; i++ {
		if i%15 == 0 {
			fmt.Println("FizzBuzz")
		} else if i%3 == 0 {
			fmt.Println("Fizz")
		} else if i%5 == 0 {
			fmt.Println("Buzz")
		} else {
			fmt.Println(i)
		}
	}
}
```

Run:

```sh
xtracego run ./examples/fizzbuzz
```

Got trace output from stderr:

```
2025-11-03T11:55:22Z [ 1] :const N = 20 ............................. [ /path/to/examples/fizzbuzz/main.go:6:7 ]
2025-11-03T11:55:22Z [ 1] :[VAR] N=20
2025-11-03T11:55:22Z [ 1] :[CALL] (main.main)
2025-11-03T11:55:22Z [ 1] :    for i := 1; i <= N; i++ { ............ [ /path/to/examples/fizzbuzz/main.go:9:2 ]
2025-11-03T11:55:22Z [ 1] :[VAR] i=1
2025-11-03T11:55:22Z [ 1] :        if i%15 == 0 { ................... [ /path/to/examples/fizzbuzz/main.go:10:3 ]
2025-11-03T11:55:22Z [ 1] :        } else if i%3 == 0 { ............. [ /path/to/examples/fizzbuzz/main.go:12:10 ]
2025-11-03T11:55:22Z [ 1] :        } else if i%5 == 0 { ............. [ /path/to/examples/fizzbuzz/main.go:14:10 ]
2025-11-03T11:55:22Z [ 1] :        } else { ......................... [ /path/to/examples/fizzbuzz/main.go:16:3 ]
2025-11-03T11:55:22Z [ 1] :            fmt.Println(i) ............... [ /path/to/examples/fizzbuzz/main.go:17:4 ]
2025-11-03T11:55:22Z [ 1] :[VAR] i=2
2025-11-03T11:55:22Z [ 1] :        if i%15 == 0 { ................... [ /path/to/examples/fizzbuzz/main.go:10:3 ]
2025-11-03T11:55:22Z [ 1] :        } else if i%3 == 0 { ............. [ /path/to/examples/fizzbuzz/main.go:12:10 ]
2025-11-03T11:55:22Z [ 1] :        } else if i%5 == 0 { ............. [ /path/to/examples/fizzbuzz/main.go:14:10 ]
2025-11-03T11:55:22Z [ 1] :        } else { ......................... [ /path/to/examples/fizzbuzz/main.go:16:3 ]
2025-11-03T11:55:22Z [ 1] :            fmt.Println(i) ............... [ /path/to/examples/fizzbuzz/main.go:17:4 ]
2025-11-03T11:55:22Z [ 1] :[VAR] i=3
2025-11-03T11:55:22Z [ 1] :        if i%15 == 0 { ................... [ /path/to/examples/fizzbuzz/main.go:10:3 ]
2025-11-03T11:55:22Z [ 1] :        } else if i%3 == 0 { ............. [ /path/to/examples/fizzbuzz/main.go:12:10 ]
2025-11-03T11:55:22Z [ 1] :            fmt.Println("Fizz") .......... [ /path/to/examples/fizzbuzz/main.go:13:4 ]
...
```


Got output from stdout:

```
1
2
Fizz
...
```

## Features

- Run Go source files directly with injected xtrace.
- Build an executable file from source files with injected xtrace.
- Rewrite source files to inject xtrace.

### Xtrace

Xtrace is execution trace information. The following traces are available:

- Traces of basic statements
- Traces of function and method calls
- Traces of variables and constants

## Installation

### Using Go

```shell
go install github.com/Jumpaku/xtracego/cmd/xtracego@latest
```

### Using Docker

```shell
docker run -i -v $(pwd):/workspace ghcr.io/jumpaku/xtracego:latest xtracego
```

### Downloading executable binary files

https://github.com/Jumpaku/xtracego/releases

Note that the downloaded executable binary file may require a security confirmation before it can be run.

### Building from source

```shell
git clone https://github.com/Jumpaku/xtracego.git
cd xtracego
go install ./cmd/xtracego
```

## Usage

### Run Go source files directly with injected xtrace

```sh
xtracego run ./path/to/package
```

### Build an executable file from source files with injected xtrace

```sh
xtracego build -o=build_dir ./path/to/package
```

### Rewrite source files to inject xtrace

```sh
xtracego rewrite -o=out_dir ./path/to/package
```

## Documentation

### Command-line interface

See detailed CLI documentation:

https://github.com/Jumpaku/xtracego/blob/main/docs/xtracego.md

## Limitation

- Comments are not handled; therefore, compiler directives (e.g., `//go:embed`) are ignored.

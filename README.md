# go-struct-optimiser
Optimise memory allocation for your Golang structs.

## Why
The order of which fields are declared in a Go struct impact on the way the memory is allocated.
This becomes tedious on larger structures, thus this tool.

Even though it's still in a very basic state, this tool can already help optimise Go structs on 64 bits systems.

## Usage
The tool takes in input a path to a file containing exactly one Go `struct`.

It, then, prints the optimal structure size and the optimal `struct` layout.
```
Usage:
  go-struct-optimiser [flags]

Flags:
  -h, --help                help for go-struct-optimiser
  -i, --input-file string   The path to the file containing ONLY your go struct

```

## Example

Our input file is `/users/home/myself/input.txt`, which contains:
```golang
type SpecialStruct struct {
	a,b,c int32
	ptr *error
	f float64
	flag,flag2 bool
}
```

We can then run the tool:
```
$ go-struct-optimiser -i "/users/home/myself/input.txt"

the best size for your struct is 32
type SpecialStruct struct {
        ptr     *error
        f       float64
        a       int32
        b       int32
        c       int32
        flag    bool
        flag2   bool
}
```

We can now go ahead and replace our old struct with the new and improve one, in our codebase.

## TODO
- Handle comments
- Handle multiple structs per file
- Handle nested structs
- Refactor into packages / segregate logic
- Unit tests
- Handle 32 bits architectures
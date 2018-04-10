bin2go
======

Converts binary files to go source code.

Installation
------------

	go get github.com/gonutz/bin2go/v2/bin2go

Usage
-----

```
Usage of bin2go:
  -export
        if true, will make the first letter of the variable name upper-case; if
        false, use the variable name as is
  -package string
        Package name. Empty string to omit package clause. (default "main")
  -var string
        Variable name to use. Must not be empty.
```

Example of converting a binary file to Go:

	bin2go -var=VarName < file > file.go

This creates a file of the following format:

```
package main

var VarName = []byte{
	0x62, 0x6C, 0x61, 0x68, 0x0D, 0x0A,
}

```
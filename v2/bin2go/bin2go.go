package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	packageName = flag.String("package", "main", "Package name. Empty string to omit package clause.")
	varName     = flag.String("var", "", "Variable name to use. Must not be empty.")
)

func main() {
	flag.Parse()

	if *varName == "" {
		flag.Usage()
		return
	}

	if *packageName != "" {
		fmt.Printf("package %s\n\n", *packageName)
	}

	fmt.Printf("var %s = []byte{", *varName)

	n, err := io.Copy(&generator{bytesInLine: maxBytesInLine}, os.Stdin)
	if err != nil {
		panic(err)
	}
	if n > 0 {
		fmt.Println()
	}

	fmt.Print("}")

	if *packageName != "" {
		fmt.Println()
	}
}

const maxBytesInLine = 12

type generator struct {
	bytesInLine int
}

func (g *generator) Write(p []byte) (n int, err error) {
	for _, b := range p {
		if g.bytesInLine >= maxBytesInLine {
			fmt.Print("\n\t")
			g.bytesInLine = 0
		} else {
			fmt.Print(" ")
		}
		fmt.Printf("0x%02X,", b)
		g.bytesInLine++
	}
	return len(p), nil
}

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

var (
	packageName = flag.String("package", "main", "Package name. Empty string to omit package clause.")
	varName     = flag.String("var", "", "Variable name to use. Must not be empty.")
	export      = flag.Bool("export", false, "if true, will make the first letter of the variable name upper-case; if false, use the variable name as is")
)

func main() {
	flag.Parse()
	err := generate(os.Stdin, os.Stdout, *varName, *packageName, *export)
	if err != nil {
		flag.Usage()
		panic(err)
	}
}

func generate(in io.Reader, out io.Writer, varName, packageName string, export bool) error {
	if varName == "" {
		return errors.New("variable name must not be empty")
	}

	if export {
		r, size := utf8.DecodeRuneInString(varName)
		r = unicode.ToUpper(r)
		varName = string(r) + varName[size:]
	}

	p := printer{out: out}

	if packageName != "" {
		p.printf("package %s\n\n", packageName)
	}

	p.printf("var %s = []byte{", varName)

	n, err := io.Copy(&generator{p: &p}, in)
	if err != nil {
		return err
	}
	if n > 0 {
		p.println()
	}

	p.print("}")

	if packageName != "" {
		p.println()
	}

	return p.err
}

type printer struct {
	out io.Writer
	err error
}

func (p *printer) println(a ...interface{}) {
	p.print(a...)
	p.print("\n")
}

func (p *printer) print(a ...interface{}) {
	if p.err == nil {
		_, p.err = fmt.Fprint(p.out, a...)
	}
}

func (p *printer) printf(format string, a ...interface{}) {
	if p.err == nil {
		_, p.err = fmt.Fprintf(p.out, format, a...)
	}
}

const maxBytesInLine = 12

type generator struct {
	p                *printer
	availBytesInLine int
}

func (g *generator) Write(p []byte) (n int, err error) {
	for _, b := range p {
		if g.availBytesInLine <= 0 {
			g.p.print("\n\t")
			g.availBytesInLine = maxBytesInLine
		} else {
			g.p.print(" ")
		}
		g.p.printf("0x%02X,", b)
		g.availBytesInLine--
	}
	return len(p), nil
}

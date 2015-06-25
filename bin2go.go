// Copyright (c) 2011, Christoph Schunk
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//     * Redistributions of source code must retain the above copyright
//       notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above copyright
//       notice, this list of conditions and the following disclaimer in the
//       documentation and/or other materials provided with the distribution.
//     * Neither the name of the author nor the
//       names of its contributors may be used to endorse or promote products
//       derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER BE LIABLE FOR ANY
// DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

var (
	packageName  = flag.String("p", "main", "Package name.")
	bytesPerLine = flag.Int("l", 8, "Number of bytes per line.")
	writeComment = flag.Bool("c", false, "Add the bytes as text as line comments.")
	useSlice     = flag.Bool("s", false, "Use slice ([]byte) instead of array ([...]byte)")
	variableName = flag.String("v", "", "Variable name to use. If empty, the name"+
		" is generated from the file name.")
	outputFileName = flag.String("o", "", "Output file path. This is only used"+
		" if there is just one file to be converted.")
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		// no files were given so tell the user how to use this program
		flag.Usage()
		return
	}

	for _, fileName := range flag.Args() {
		inputFile, err := os.Open(fileName)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer inputFile.Close()

		outputFile, err := os.Create(makeOutputFileName(fileName))
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer outputFile.Close()

		varName := *variableName
		if len(varName) == 0 {
			varName = camelCase(fileName)
		}
		if err = bin2go(inputFile, outputFile, varName); err != nil {
			fmt.Println(err)
		}
	}
}

func makeOutputFileName(fileName string) string {
	if len(*outputFileName) > 0 && len(flag.Args()) == 1 {
		name := *outputFileName
		if !strings.HasSuffix(name, ".go") {
			name += ".go"
		}
		return name
	}
	return fileName + ".go"
}

func bin2go(inputFile, outputFile *os.File, variableName string) error {
	buffer := make([]byte, *bytesPerLine)
	varType := "[...]byte"
	if *useSlice {
		varType = "[]byte"
	}
	fmt.Fprintf(outputFile, "package %s\n", *packageName)
	fmt.Fprintf(outputFile, "\nvar %s = %s{\n", variableName, varType)
	for {
		nRead, err := inputFile.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		fmt.Fprintf(outputFile, "\t%0#2x,", buffer[0])
		for _, c := range buffer[1:nRead] {
			fmt.Fprintf(outputFile, " %0#2x,", c)
		}
		if *writeComment {
			fmt.Fprintf(outputFile, "\t// %s",
				removeNonDisplayableRunes(string(buffer)))
		}
		fmt.Fprint(outputFile, "\n")
	}
	fmt.Fprint(outputFile, "}\n")
	return nil
}

func removeNonDisplayableRunes(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return '_'
	}, s)
}

func camelCase(s string) string {
	camel := false
	first := true
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if first { // always lower case first rune
				first = false
				return unicode.ToLower(r)
			}
			if camel {
				camel = false
				return unicode.ToTitle(r)
			}
			return r
		}

		if !first { // if first runes aren't letters or digits, don't capitalize
			camel = true // unknown rune type -> ignore and title case next rune
		}

		return -1
	}, s)
}

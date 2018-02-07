package main

import (
	"bytes"
	"testing"
)

func TestDoNotForgetTheVar(t *testing.T) {
	err := generate(nil, nil, "", "")
	if err == nil {
		t.Error("error expected when variable name is empty")
	}
}

func TestEmptyInputGeneratesEmptyByteSlice(t *testing.T) {
	checkGeneratedCode(
		t,
		"v", "",
		[]byte{},
		`var v = []byte{}`,
	)
}

func TestSingleByteSliceStartsOnNewLine(t *testing.T) {
	checkGeneratedCode(
		t,
		"v", "",
		[]byte{0},
		`var v = []byte{
	0x00,
}`,
	)
}

func TestLinesDoNotBecomeTooLongToRead(t *testing.T) {
	checkGeneratedCode(
		t,
		"v", "",
		[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		`var v = []byte{
	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B,
	0x0C,
}`,
	)
}

func TestPackageNameMeansPackageClauseAndNewLineAtEnd(t *testing.T) {
	checkGeneratedCode(
		t,
		"abc", "main",
		[]byte{0, 1, 2},
		`package main

var abc = []byte{
	0x00, 0x01, 0x02,
}
`,
	)
}

func checkGeneratedCode(t *testing.T, varName, packageName string, data []byte, expectedCode string) {
	var output bytes.Buffer
	err := generate(bytes.NewReader(data), &output, varName, packageName)
	if err != nil {
		t.Fatal(err)
	}
	code := string(output.Bytes())
	if code != expectedCode {
		t.Errorf("have code\n---\n%s\n---", code)
	}
}

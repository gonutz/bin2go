#!/bin/bash

# this Linux shell script takes all .txt files in this folder and generates
# a single Go file containing their content with the variable name being the
# file name without extension. The -export flag makes the variable names start
# upper-case so they are exported.

echo package data>gen.go
for x in ./*.txt; do
  echo "">>gen.go
  bin2go -var=$(basename $x .txt) -export -package="" < "$x" >> gen.go
  echo "">>gen.go
done
echo "">>gen.go

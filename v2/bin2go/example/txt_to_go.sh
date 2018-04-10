#!/bin/bash

echo package data>gen.go
for x in ./*.txt; do
  echo "">>gen.go
  bin2go -var=$(basename $x .txt) -export -package="" < "$x" >> gen.go
  echo "">>gen.go
done
echo "">>gen.go

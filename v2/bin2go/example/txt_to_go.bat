REM this Windows batch script takes all .txt files in this folder and generates
REM a single Go file containing their content with the variable name being the
REM file name without extension. The -export flag makes the variable names start
REM upper-case so they are exported.

echo package data>gen.go
for %%x in (*.txt) do (
  echo.>>gen.go
  bin2go -var=%%~nx -export -package="" < %%x >> gen.go
  echo.>>gen.go
)
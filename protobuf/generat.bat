@echo off
set a=%cd%

echo a=%a%

::set b=%a%\protos
set b=%a%\protos
echo proto path=%b%

set c=%a%\pb
echo outPath=%c%

echo "GOPATH= %GOPATH%"

set d=%GOPATH%\bin\protoc.exe
echo "d=%d%" 

set e=%GOPATH%\bin\protoc-gen-go.exe
echo "e=%e%"

if exist %c% (
    for /f "delims=" %%i in ('dir /b "%c%\*.go"') do (
            echo del file%%i
            del %c%\%%i
    )
) else (
    echo create dir%%i
    md pb %a%
)

for /f "delims=" %%i in ('dir /b "%b%\*.proto"') do (
    echo create %%i
    %d% --plugin=protoc-gen-go=%e% --proto_path=%b% --go_out=%c% %%i
)
pause

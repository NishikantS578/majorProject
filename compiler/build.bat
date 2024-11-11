gofmt -l -w src\
go -C .\src build -o ..\build
.\build\compiler.exe
# majorProject

# Requirements

## For development
- go
- gofmt

## Others
- nasm
- gcc (ld)

## How to build
```
cd <path to root of project>
go -C ./src mod tidy
./build.sh
```

## Usage
```
mnm <path to program file> // Compiles your mnm code
nasm -felf64 app.asm -o out.o && ld -o <app name> out.o && ./<app name> //for 64-bit
```

# Tasks
- [x] Basic Lexer.
- [x] Basic Parser.
- [x] Basic Generator.
- [ ] Add variables.
- [ ] Add most of data types.
- [ ] GUI Text/Code Editor.
- [ ] GUI Visual Editor.
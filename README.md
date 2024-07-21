# majorProject

# Requirements

## For development
- go

## Others
- nasm
- gcc

## How to build
```
cd <path to root of project>
go -C ./src mod tidy
go -C ./src build -o ../build
```

## Usage
```
mnm <path to program file> // Compiles your mnm code
```
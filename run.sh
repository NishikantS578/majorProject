build/compiler test.mnm
nasm -felf64 test.so
ld test.o -o test
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>

int main(int argc, char* argv[]){
    if(argc != 2){
        printf("Incorrect usage.\n");
        printf("Correct usage: compiler <input.mnm>\n");
        return 0;
    }

    char *fileBuffer;
    char *fileName = argv[1];
    uint64_t fileSize;

    FILE *programFile = fopen(fileName, "r");
    if(!programFile){
        printf("Couldn't open file: %s", fileName);
        return 0;
    }
    fseek(programFile, 0, SEEK_END);
    fileSize = ftell(programFile);
    fseek(programFile, 0, SEEK_SET);
    fileBuffer = (char *)malloc(fileSize);
    int fileRead = fread(fileBuffer, fileSize, 1, programFile);
    if(fileRead != 1){
        printf("Error while reading the file: %s", fileName);
        return 0;
    }
    fclose(programFile);

    for(int i = 0; i < fileSize; i++){
        printf("%c", fileBuffer[i]);
    }
    printf("\n");
    return 0;
}
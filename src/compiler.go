package main

import "fmt"

func main(){
	fmt.Println("Hello")
}

// #include <iostream>
// #include <fstream>
// #include <vector>
// #include <optional>
// #include <cstdint>

// enum TokenType{
//     _return,
//     integerLiteral,
//     semicolon
// };

// struct Token{
//     TokenType type;
//     std::optional<std::string> value;
// };

// std::vector<Token> tokenize(std::string str, int size){
//     std::vector<Token> tokenArr = {};
//     for(int i = 0; i < str)
//     return tokenArr;
// }

// int main(int argc, char* argv[]){
//     if(argc != 2){
//         std::cout<<"Incorrect usage.\n";
//         std::cout<<"Correct usage: compiler <input.mnm>\n";
//         return 0;
//     }

//     std::string fileName = argv[1];
//     int fileSize = 0;
//     std::string fileContent;

//     std::fstream inputFile(fileName, std::ios_base::in);
//     if(!inputFile.is_open()){
//         std::cout<<"Failed to open file: "<<fileName<<"\n";
//     }

//     inputFile.seekg(0, std::ios_base::end);
//     fileSize = inputFile.tellg();
//     inputFile.seekg(0, std::ios_base::beg);
//     fileContent.resize(fileSize);
//     inputFile.read(&fileContent[0], fileSize);
//     inputFile.close();

//     std::vector<Token> tokenArr;
//     tokenArr = tokenize(fileContent, fileSize);
//     for(Token t: tokenArr){
//         switch(t.type){
//             case _return:{
//                 std::cout<<"return ";
//                 break;
//             }
//             case integerLiteral:{
//                 std::cout<<"integer-literal ";
//                 break;
//             }
//             case semicolon:{
//                 std::cout<<"semicolon ";
//                 break;
//             }
//         }
//         std::string val = t.value.value_or("");
//         std::cout<<std::endl;
//     }
//     return 0;
// }

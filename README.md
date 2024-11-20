# Major Project: MNM Language Editor and Compiler  

This project introduces a **custom programming language** called **MNM**, complete with a **compiler** and an **editor** that emphasizes **code visualization**. Built using Go and Electron, the system is fully cross-platform, offering an integrated environment for coding, compiling, and visualizing program execution in real-time.  

## Features  

### Editor  
- **Technology**: Built using Electron.  
- **Standard Functionalities**: Save, edit, open files, and compile/run code.  
- **Code Visualization**: Unique "View Flow Chart" option that generates a real-time flowchart representation of code structure.  
- **Educational Focus**: Assists users in understanding program flow and debugging logic efficiently.  
- **Cross-Platform Compatibility**: Runs on Windows, macOS, and Linux.  

### Compiler  
- **Technology**: Developed in Go.  
- **Language**: Supports MNM files (`.mnm`).  
- **Compiler Features**:  
  - **Tokenization and Parsing**: Implements the **Pratt parsing algorithm** for efficient handling of operator precedence and expression parsing.  
  - **Code Generation**: Converts MNM code to run on a **stack-based virtual machine**.  
  - **Supported Features**: Algebraic expressions, logical conditions, if-else statements, variable declarations, and core data types (integers and booleans).  
- **Efficient Execution**: Compiles and executes code with detailed error handling.  

---

## Project Setup  

### Prerequisites  
- **Node.js** (v16.x or higher)  
- **Go** (v1.18 or higher)  
- **npm** or **yarn**  

### Installation  

1. **Clone the Repository**:  
   ```bash  
   git clone https://github.com/Moh1tsingh/majorProject.git  
   cd majorProject
   ```
2. **Install Dependencies**:
   Navigate to the editor directory and install dependencies:
   ```bash
   cd editor  
   npm install  
   ```
3. **Build the Compiler**:
   Navigate to the `compiler` directory and run the build script:  
   ```bash
   cd compiler  
   build.bat   
   ```
### Running the Project
1. **Start the Editor**:
   Navigate to the `compiler` directory and run the build script:  
   ```bash
   cd editor  
   npm start  
   ```
2. **Compile MNM Files**:
  Write and visualize `.mnm` files directly within the editor.

   

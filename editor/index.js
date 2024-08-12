const { app, BrowserWindow, ipcMain, dialog } = require("electron")
const path = require("node:path")
const fs = require("fs")
const { exec } = require("node:child_process")

let openedFilePath = ""

const createWindow = () => {
    const win = new BrowserWindow({
        width: 800, height: 600,
        webPreferences: {
            preload: path.join(__dirname, "preload.js")
        },
        autoHideMenuBar: true,
    })

    win.loadFile("index.html")
    return win
}

app.on("window-all-closed", () => {
    if (process.platform != "darwin") app.quit()
})

app.whenReady().then(() => {
    let win = createWindow()
    ipcMain.handle("saveFile", saveFile)
    ipcMain.handle("openFile", (e) => openFile(e, win))
    ipcMain.handle("compileAndRun", (e) => compileAndRun(e, win))
})

const saveFile = (e, fileContent) => {
    dialog.showSaveDialog().then((res) => {
        fs.writeFile(res.filePath, fileContent, (err) => {
            if(err != null){
                console.log(err)
            }
        })
    })
}

const openFile = (e, win) => {
    dialog.showOpenDialog().then((res) => {
        if(res.filePaths[0] == undefined){
            return
        }
        openedFilePath = res.filePaths[0]
        fs.readFile(res.filePaths[0], "utf-8", (err, fileContent) => {
            win.webContents.send("openedFile", fileContent)
        })
    })
}

const compileAndRun = (e, win) => {
    exec(" ../compiler/build/mnm " + openedFilePath + " && nasm -felf64 app.asm -o out.o && ld out.o && ./a.out", (err, stdout, stderr)=>{
        win.webContents.send("compiledAndRun", stdout + stderr + err.code)
    })
}
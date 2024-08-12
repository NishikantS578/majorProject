const {contextBridge, ipcRenderer} = require("electron")

contextBridge.exposeInMainWorld("versions", {
    saveFile: (fileContent) => ipcRenderer.invoke("saveFile", fileContent),
    openFile: () => ipcRenderer.invoke("openFile"),
    compileAndRun: () => ipcRenderer.invoke("compileAndRun"),
    openedFile: (callback) => ipcRenderer.on("openedFile", (e, fileContent) => callback(fileContent)),
    compiledAndRun: (callback) => ipcRenderer.on("compiledAndRun", (e, output) => callback(output))
})
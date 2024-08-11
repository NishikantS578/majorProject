const {contextBridge, ipcRenderer} = require("electron")

contextBridge.exposeInMainWorld("versions", {
    saveFile: (fileContent) => ipcRenderer.invoke("saveFile", fileContent),
    openFile: () => ipcRenderer.invoke("openFile"),
    openedFile: (callback) => ipcRenderer.on("openedFile", (e, fileContent) => callback(fileContent))
})
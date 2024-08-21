const saveBtn = document.querySelector("#save")
const openBtn = document.querySelector("#open")
const compileAndRunBtn = document.querySelector("#compileAndRun")
const inputEl = document.querySelector("textarea")
const outputWindowOutput = document.querySelector("#outputWindow #output")

saveBtn.addEventListener("mousedown", (e) => {
    versions.saveFile(inputEl.value)
})

export const save =  (e) => {
    versions.saveFile(inputEl.value)
}

openBtn.addEventListener("mousedown", (e) => {
    versions.openFile().then(() => {
    })
})

export const open = (e) => {
    versions.openFile().then(() => {
    })
}

compileAndRunBtn.addEventListener("mousedown", (e)=>{
    versions.compileAndRun().then(()=>{
    })
})

export const compileRun = (e)=>{
    versions.compileAndRun().then(()=>{
    })
}

window.versions.openedFile((data) => {
    inputEl.value = data
})

window.versions.compiledAndRun((data) => {
    outputWindowOutput.textContent = data
})
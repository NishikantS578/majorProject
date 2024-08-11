const inputEl = document.querySelector("textarea")
const saveBtn = document.querySelector("#save")
const openBtn = document.querySelector("#open")

saveBtn.addEventListener("mousedown", (e) => {
    versions.saveFile(inputEl.value)
})

openBtn.addEventListener("mousedown", (e) => {
    versions.openFile().then(() => {
    })
})

window.versions.openedFile((data) => {
    inputEl.value = data
})
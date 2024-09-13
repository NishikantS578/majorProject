import { useState } from 'react'
import CodeBlock from './components/CodeBlock'


const save = (editor_content_state: any) => {
  console.log(editor_content_state)
  window.ipcRenderer.invoke("saveFile", editor_content_state)
}

const open = () => {
  window.ipcRenderer.invoke("openFile")
}

const compileRun = () => {
  window.ipcRenderer.invoke("compileAndRun")
}

function App() {
  const [editor_content, set_editor_content_state] = useState("")
  const [output_content, set_output_content] = useState("")
  const [window_type, set_window_type] = useState("text_editor")

  window.ipcRenderer.on("openedFile", (_, data) => {
    set_editor_content_state(data)
  })

  window.ipcRenderer.on("compiledAndRun", (_, data) => {
    set_output_content(data)
  })

  return (
    <>
      <div className="menubar">
        <button onMouseDown={open}>Open</button>
        <button onMouseDown={() => save(editor_content)}>
          Save
        </button>
        <button onMouseDown={compileRun}>
          Compile and Run
        </button>
        <button onClick={() => {
          if (window_type == "text_editor") {
            set_window_type("visual_editor")
          } else {
            set_window_type("text_editor")
          }
        }}>
          {window_type}
        </button>
      </div>

      <div className='editor min-w-fit w-full flex-auto'>
        {
          window_type == "text_editor" ?
            (
              <textarea className="editorTab min-w-fit text-nowrap w-full h-full"
                onChange={
                  (e) => set_editor_content_state(e.target.value)}
                value={editor_content}>
              </textarea>
            )
            :
            (
              <div className='visualEditor w-full h-full'>
                <CodeBlock arrow_start_x={0} arrow_start_y={0}></CodeBlock>
              </div>
            )
        }
      </div>

      <div className='outputWindow'>
        <h3>OUTPUT</h3>
        <textarea className='output'
          disabled value={output_content}>
        </textarea>
      </div>
    </>
  )
}

export default App